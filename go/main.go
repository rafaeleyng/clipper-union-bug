package main

import (
	"fmt"

	clipper "github.com/ctessum/go.clipper"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	clip := clipper.NewClipper(clipper.IoNone)

	polygons := getPolygons()
	paths := polygonsToPaths(polygons)

	for i, path := range paths {
		result := clip.SimplifyPolygon(path, clipper.PftNonZero)
		if len(result) > 1 {
			panic("unexpected multiple paths polygon")
		}
		paths[i] = result[0]
	}

	clip.AddPaths(paths, clipper.PtSubject, true)
	combinedPaths, ok := clip.Execute1(clipper.CtUnion, clipper.PftNonZero, clipper.PftNonZero)
	if !ok {
		panic("failed to execute")
	}

	if len(combinedPaths) > 1 {
		panic("unexpected multiple paths in combined paths")
	}

	combinedPolygons := pathsToPolygons(combinedPaths)
	fmt.Printf("%+v\n", combinedPolygons)

	output(pathsToPolygons(paths), "output-paths.pdf")
	output(combinedPolygons, "output-combinedPaths.pdf")
}

/*
	coordinates for the test case
*/
func getPolygons() Polygons {
	return Polygons{
		{
			{X: 44.412, Y: 291.099},
			{X: 60.819, Y: 319.514},
			{X: 58.485, Y: 346.637},
			{X: 25.673, Y: 346.637},
			{X: 1, Y: 340},
			{X: 1, Y: 323.154},
			{X: 13.809, Y: 294.739},
			{X: 27.035, Y: 289.023},
		},
		{
			{X: 44.412, Y: 262.684},
			{X: 60.819, Y: 291.099},
			{X: 58.485, Y: 318.222},
			{X: 25.673, Y: 318.222},
			{X: 1, Y: 311.585},
			{X: 1, Y: 294.739},
			{X: 13.809, Y: 266.324},
			{X: 27.035, Y: 260.608},
		},
		{
			{X: 44.412, Y: 234.269},
			{X: 60.819, Y: 262.684},
			{X: 58.485, Y: 289.807},
			{X: 25.673, Y: 289.807},
			{X: 1, Y: 283.17},
			{X: 1, Y: 266.324},
			{X: 13.809, Y: 237.909},
			{X: 27.035, Y: 232.193},
		},
		{
			{X: 59.573, Y: 195.326},
			{X: 102.667, Y: 213.998},
			{X: 180.578, Y: 238.897},
			{X: 178.244, Y: 266.02},
			{X: 53.288, Y: 272.762},
			{X: 27.615, Y: 266.125},
			{X: 1, Y: 238.116},
			{X: 1, Y: 221.27},
			{X: 28.97, Y: 198.966},
			{X: 42.196, Y: 193.25},
		},
		{
			{X: 61.025, Y: 288.7054},
			{X: 86.698, Y: 295.3424},
			{X: 89.294, Y: 312.1884},
			{X: 86.96, Y: 339.3114},
			{X: 73.734, Y: 345.0274},
			{X: 56.357, Y: 342.9514},
			{X: 30.684, Y: 336.3144},
			{X: 28.088, Y: 319.4684},
			{X: 30.422, Y: 292.3454},
			{X: 43.648, Y: 286.6294},
		},
	}
}

/*
	conversion between clipper types and my custom types
*/
type (
	Point    struct{ X, Y float64 }
	Polygon  []Point
	Polygons []Polygon
)

func pointToIntPoint(point Point) *clipper.IntPoint {
	return &clipper.IntPoint{X: scaleUp(point.X), Y: scaleUp(point.Y)}
}

func polygonToPath(polygon Polygon) clipper.Path {
	path := clipper.Path(make([]*clipper.IntPoint, len(polygon)))
	for i, point := range polygon {
		path[i] = pointToIntPoint(point)
	}
	return path
}

func polygonsToPaths(polygons Polygons) clipper.Paths {
	paths := clipper.Paths(make([]clipper.Path, len(polygons)))
	for i, polygon := range polygons {
		paths[i] = polygonToPath(polygon)
	}
	return paths
}

func intPointToPoint(intPoint *clipper.IntPoint) Point {
	return Point{X: scaleDown(intPoint.X), Y: scaleDown(intPoint.Y)}
}

func pathToPolygon(path clipper.Path) Polygon {
	polygon := Polygon(make([]Point, len(path)))
	for i, intPoint := range path {
		polygon[i] = intPointToPoint(intPoint)
	}
	return polygon
}

func pathsToPolygons(paths clipper.Paths) Polygons {
	polygons := Polygons(make([]Polygon, len(paths)))
	for i, path := range paths {
		polygons[i] = pathToPolygon(path)
	}
	return polygons
}

const scaleFactor = 1e9

func scaleUp(value float64) clipper.CInt {
	return clipper.CInt(value * scaleFactor)
}

func scaleDown(value clipper.CInt) float64 {
	// Using:
	//  return float64(value / scaleFactor)
	// can seem to fix the problem, but it actually just hides it more.
	return float64(value) / scaleFactor
}

/*
	pdf output
*/
func output(polygons Polygons, outFile string) {
	page := pdfcpu.NewPage(pdfcpu.RectForDim(400, 400))
	pdfcpu.SetLineWidth(page.Buf, 0)

	for i, polygon := range polygons {
		pdfcpu.SetStrokeColor(page.Buf, colors[i%len(colors)])
		for j, curr := range polygon {
			nextJ := 0
			if j < len(polygon)-1 {
				nextJ = j + 1
			}
			next := polygon[nextJ]
			pdfcpu.DrawLine(page.Buf, curr.X, curr.Y, next.X, next.Y)
		}
	}

	xRefTable, err := pdfcpu.CreateDemoXRef(page)
	if err != nil {
		panic(err)
	}

	err = api.CreatePDFFile(xRefTable, outFile, nil)
	if err != nil {
		panic(err)
	}
}

var colors = []pdfcpu.SimpleColor{
	{R: 0.8984375, G: 0.1484375, B: 0.12109375},
	{R: 0.91796875, G: 0.45703125, B: 0.1953125},
	{R: 0.96484375, G: 0.8125, B: 0.21875},
	{R: 0.63671875, G: 0.875, B: 0.28125},
	{R: 0.28515625, G: 0.8515625, B: 0.6015625},
}
