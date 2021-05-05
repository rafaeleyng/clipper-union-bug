Investigation with 3 triangles:

- -9007199254740992 in JS becomes -Inf in Go, in `c.m_edges.Dx`
- Go doesn't have `m_IntersectNodeComparer`, doesn't seem to be a problem though, it solves it another way
- Go doesn't have `m_Maxima`, `InsertMaxima`, doesn't seem to be a problem. TODO: should investigate more
- JS calls `c.InsertLocalMinimaIntoAEL(botY)` differently
- in JS `ExecuteInternal`, topY.v ends with `0`, while in Go ends with `160000000000`
- after the main loop in `ExecuteInternal`, `m_PolyOuts[0].Pts.Prev.Pt` is different between Go and JS


---

// cpp
[
    [
        { X: 58.4261, Y: 187.235 },
        { X: 59.587, Y: 188.783 },
        { X: 61, Y: 189 },
        { X: 60.1667, Y: 189.556 },
        { X: 68, Y: 200 },
        { X: 44, Y: 199 },
        { X: 48.5774, Y: 189.337 },
        { X: 30, Y: 190 },
        { X: 65, Y: 160 }
    ]
]

// javascript
[
    [
        { X: 58.426086957, Y: 187.234782609 },
        { X: 59.586956522, Y: 188.782608696 },
        { X: 61, Y: 189 },
        { X: 60.166666667, Y: 189.555555556 },
        { X: 68, Y: 200 },
        { X: 44, Y: 199 },
        { X: 48.577437859, Y: 189.336520076 },
        { X: 30, Y: 190 },
        { X: 65, Y: 160 }
    ]
]

// go
[
    [
        { X: 58.10741688, Y: 188.554987212 },
        { X: 61, Y: 189 },
        { X: 60.166666667, Y: 189.555555556 },
        { X: 68, Y: 200 },
        { X: 44, Y: 199 },
        { X: 48.577437859, Y: 189.336520076 },
        { X: 30, Y: 190 },
        { X: 65, Y: 160 }
    ]
]

---

Calls to AddOutPt:

go:
### AddOutPt CountOuts 0 {68000000000, 200000000000}
### AddOutPt CountOuts 2 {44000000000, 199000000000}
### AddOutPt CountOuts 2 {30000000000, 190000000000}
### AddOutPt CountOuts 4 {48577437859, 189336520076}
### AddOutPt CountOuts 4 {60166666667, 189555555556}
### AddOutPt CountOuts 5 {61000000000, 189000000000}
### AddOutPt CountOuts 6 {58107416880, 188554987212}
### AddOutPt CountOuts 7 {65000000000, 160000000000}

this is unique to go:
### AddOutPt CountOuts 6 {58107416880, 188554987212}

js:
### AddOutPt CountOuts 0 { X: 68000000000, Y: 200000000000 }
### AddOutPt CountOuts 2 { X: 44000000000, Y: 199000000000 }
### AddOutPt CountOuts 2 { X: 30000000000, Y: 190000000000 }
### AddOutPt CountOuts 4 { X: 60166666667, Y: 189555555556 }
### AddOutPt CountOuts 5 { X: 48577437859, Y: 189336520076 }
### AddOutPt CountOuts 5 { X: 61000000000, Y: 189000000000 }
### AddOutPt CountOuts 6 { X: 59586956522, Y: 188782608696 }
### AddOutPt CountOuts 7 { X: 58426086957, Y: 187234782609 }
### AddOutPt CountOuts 8 { X: 65000000000, Y: 160000000000 }

this is unique to js:
### AddOutPt CountOuts 6 {59586956522, 188782608696}
### AddOutPt CountOuts 7 {58426086957, 187234782609}

this is the order in go:
{48577437859, 189336520076}
{60166666667, 189555555556}

this is the order in js:
{60166666667, 189555555556}
{48577437859, 189336520076}


---

the difference happens in:

```
for i := 0; i < len(c.m_PolyOuts); i++ {
    outRec := c.m_PolyOuts[i]
    if outRec.Pts != nil && !outRec.IsOpen {
        c.FixupOutPolygon(outRec)
    }
}
```
