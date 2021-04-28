const ClipperLib = require('./clipper')

const main = () => {
    const clip = new ClipperLib.Clipper()
    const combinedNfp = new ClipperLib.Paths()

    const polygons = getPolygons()
    const paths = polygonsToPaths(polygons)
        .map(p => ClipperLib.Clipper.SimplifyPolygon(p, ClipperLib.PolyFillType.pftNonZero)[0])

    clip.AddPaths(paths, ClipperLib.PolyType.ptSubject, true)

    const ok = clip.Execute(ClipperLib.ClipType.ctUnion, combinedNfp, ClipperLib.PolyFillType.pftNonZero, ClipperLib.PolyFillType.pftNonZero)
    if (!ok) {
        throw new Error('failed to execute')
    }

    if (combinedNfp.length > 1) {
        throw new ERror('unexpected multiple paths in combined paths')
    }

    const combinedPolygons = pathsToPolygons(combinedNfp)
    console.log(combinedPolygons)
}

/*
    coordinates for the test case
*/
const getPolygons = () => {
    return [
        [
            { X: 44.412, Y: 291.099 },
            { X: 60.819, Y: 319.514 },
            { X: 58.485, Y: 346.637 },
            { X: 25.673, Y: 346.637 },
            { X: 1, Y: 340 },
            { X: 1, Y: 323.154 },
            { X: 13.809, Y: 294.739 },
            { X: 27.035, Y: 289.023 },
        ],
        [
            { X: 44.412, Y: 262.684 },
            { X: 60.819, Y: 291.099 },
            { X: 58.485, Y: 318.222 },
            { X: 25.673, Y: 318.222 },
            { X: 1, Y: 311.585 },
            { X: 1, Y: 294.739 },
            { X: 13.809, Y: 266.324 },
            { X: 27.035, Y: 260.608 },
        ],
        [
            { X: 44.412, Y: 234.269 },
            { X: 60.819, Y: 262.684 },
            { X: 58.485, Y: 289.807 },
            { X: 25.673, Y: 289.807 },
            { X: 1, Y: 283.17 },
            { X: 1, Y: 266.324 },
            { X: 13.809, Y: 237.909 },
            { X: 27.035, Y: 232.193 },
        ],
        [
            { X: 59.573, Y: 195.326 },
            { X: 102.667, Y: 213.998 },
            { X: 180.578, Y: 238.897 },
            { X: 178.244, Y: 266.02 },
            { X: 53.288, Y: 272.762 },
            { X: 27.615, Y: 266.125 },
            { X: 1, Y: 238.116 },
            { X: 1, Y: 221.27 },
            { X: 28.97, Y: 198.966 },
            { X: 42.196, Y: 193.25 },
        ],
        [
            { X: 61.025, Y: 288.7054 },
            { X: 86.698, Y: 295.3424 },
            { X: 89.294, Y: 312.1884 },
            { X: 86.96, Y: 339.3114 },
            { X: 73.734, Y: 345.0274 },
            { X: 56.357, Y: 342.9514 },
            { X: 30.684, Y: 336.3144 },
            { X: 28.088, Y: 319.4684 },
            { X: 30.422, Y: 292.3454 },
            { X: 43.648, Y: 286.6294 },
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