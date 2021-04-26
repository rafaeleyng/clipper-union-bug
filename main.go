package main

import (
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
	combinedPaths, ok :=  clip.Execute1(clipper.CtUnion, clipper.PftNonZero, clipper.PftNonZero)
	if !ok {
		panic("failed to execute")
	}

	if len(combinedPaths) > 1 {
		panic("unexpected multiple paths in combined paths")
	}

	output(pathsToPolygons(paths), "paths.pdf")
	output(pathsToPolygons(combinedPaths), "combinedPaths.pdf")
}

/*
	coordinates for the test case
*/
func getPolygons() Polygons {
	return Polygons{
		{
			{X: 58.485000000000014, Y: 346.63700000000006},
			{X: 25.673000000000002, Y: 346.63700000000006},
			{X: 0, Y: 340},
			{X: -2.5960000000000036, Y: 323.154},
			{X: 13.808999999999997, Y: 294.73900000000003},
			{X: 27.034999999999997, Y: 289.023},
			{X: 44.412000000000006, Y: 291.09900000000005},
			{X: 60.81900000000002, Y: 319.514},
		},
		{
			{X: 58.484999999999985, Y: 318.22200000000004},
			{X: 25.672999999999973, Y: 318.22200000000004},
			{X: -0.00000000000002842170943040401, Y: 311.585},
			{X: -2.596000000000032, Y: 294.739},
			{X: 13.808999999999969, Y: 266.324},
			{X: 27.034999999999968, Y: 260.608},
			{X: 44.41199999999998, Y: 262.684},
			{X: 60.81899999999999, Y: 291.099},
		},
		{
			{X: 58.484999999999985, Y: 289.8070000000001},
			{X: 25.672999999999973, Y: 289.8070000000001},
			{X: -0.00000000000002842170943040401, Y: 283.17},
			{X: -2.596000000000032, Y: 266.324},
			{X: 13.808999999999969, Y: 237.90900000000005},
			{X: 27.034999999999968, Y: 232.19300000000004},
			{X: 44.41199999999998, Y: 234.26900000000006},
			{X: 60.81899999999999, Y: 262.684},
		},
		{
			{X: 53.287999998999965, Y: 272.76193988800003},
			{X: 178.24399999899998, Y: 266.01993988800007},
			{X: 180.57799999900004, Y: 238.89693988800002},
			{X: 102.66699999800005, Y: 213.99793988800002},
			{X: 59.572999999000075, Y: 195.325939888},
			{X: 42.19599999900001, Y: 193.24993988799997},
			{X: 28.969999999000038, Y: 198.96593988799998},
			{X: -2.5960000009999646, Y: 221.26993988799995},
			{X: -0.0000000009999894245993346, Y: 238.11593988799996},
			{X: 27.61499999899999, Y: 266.124939888},
		},
		{
			{X: 73.73396682299999, Y: 345.02749608500005},
			{X: 56.35696682299998, Y: 342.95149608500003},
			{X: 30.68396682299999, Y: 336.314496085},
			{X: 28.087966822999988, Y: 319.46849608499997},
			{X: 30.42196682299999, Y: 292.3454960849999},
			{X: 43.647966822999976, Y: 286.6294960849999},
			{X: 61.024966822999986, Y: 288.70549608499994},
			{X: 86.69796682299999, Y: 295.342496085},
			{X: 89.29396682299999, Y: 312.188496085},
			{X: 86.95996682299999, Y: 339.31149608500004},
		},
	}
}

/*
	conversion between clipper types and my custom types
*/
type (
	Point struct { X, Y float64 }
	Polygon []Point
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
			if j < len(polygon) - 1 {
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
