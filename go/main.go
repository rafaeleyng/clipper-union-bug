package main

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/rafaeleyng/clipper-union-bug/clipper"
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
			{X: 44, Y: 170},
			{X: 68, Y: 200},
			{X: 44, Y: 200},
		},
		{
			{X: 65, Y: 160},
			{X: 58, Y: 189},
			{X: 30, Y: 189},
		},
		{
			{X: 61, Y: 189},
			{X: 50, Y: 195},
			{X: 46, Y: 187},
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
	page := pdfcpu.NewPage(pdfcpu.RectForDim(250, 250))
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
