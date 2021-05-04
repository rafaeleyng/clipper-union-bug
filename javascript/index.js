const ClipperLib = require('./clipper')

const run = () => {
    const clip = new ClipperLib.Clipper()
    const combinedPaths = new ClipperLib.Paths()

    const polygons = getPolygons()
    let paths = polygonsToPaths(polygons)
    // paths = ClipperLib.Clipper.SimplifyPolygons(paths, ClipperLib.PolyFillType.pftNonZero)

    clip.AddPaths(paths, ClipperLib.PolyType.ptSubject, true)

    const ok = clip.Execute(ClipperLib.ClipType.ctUnion, combinedPaths, ClipperLib.PolyFillType.pftNonZero, ClipperLib.PolyFillType.pftNonZero)
    if (!ok) {
        throw new Error('failed to execute')
    }

    if (combinedPaths.length > 1) {
        throw new ERror('unexpected multiple paths in combined paths')
    }

    const combinedPolygons = pathsToPolygons(combinedPaths)
    console.log(combinedPolygons)
}

const main = () => {
    if (process.env.INFINITE === 'true') {
        while (true) {
            run()
        }
    } else {
        run()
    }
}

/*
    coordinates for the test case
*/
const getPolygons = () => {
    return [
        [
            { X: 53, Y: 180 },
            { X: 68, Y: 200 },
            { X: 44, Y: 199 },
        ],
        [
            { X: 65, Y: 160 },
            { X: 58, Y: 189 },
            { X: 30, Y: 190 },
        ],
        [
            { X: 61, Y: 189 },
            { X: 52, Y: 195 },
            { X: 48, Y: 187 },
        ],
    ]
}

/*
    conversion between clipper types and my custom types
*/
const scaleFactor = 1e9

const scaleUp = (value) => (value * scaleFactor)
const scaleDown = (value) => (value / scaleFactor)

const pointToIntPoint = (point) => ({ X: scaleUp(point.X), Y: scaleUp(point.Y) })
const polygonToPath = (polygon) => polygon.map(p => pointToIntPoint(p))
const polygonsToPaths = (polygons) => polygons.map(p => polygonToPath(p))

const intPointToPoint = (intPoint) => ({ X: scaleDown(intPoint.X), Y: scaleDown(intPoint.Y) })
const pathToPolygon = (path) => path.map(p => intPointToPoint(p))
const pathsToPolygons = (paths) => paths.map(p => pathToPolygon(p))

main()
