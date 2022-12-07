package gux

import (
	"github.com/leonsal/gux/gb"
)

// AddConvexPolyFilled adds a filled convex polygon to draw list
func (w *Window) AddConvexPolyFilled(dl *gb.DrawList, points []gb.Vec2, col gb.Color) {

	pointsCount := len(points)
	if pointsCount < 3 {
		return
	}

	if (w.drawFlags & DrawListFlags_AntiAliasedFill) != 0 {

		AA_SIZE := w.FringeScale
		colTrans := gb.Color(col & ^gb.ColorMaskA)

		// Allocates command
		idxCount := (pointsCount-2)*3 + pointsCount*6
		vtxCount := pointsCount * 2
		cmd, bufIdx, bufVtx := dl.ReserveCmd(idxCount, vtxCount)
		cmd.TexId = w.TexWhiteId

		// Add indexes for inner triangles
		idxPos := uint32(0)
		for i := uint32(2); i < uint32(pointsCount); i++ {
			bufIdx[idxPos+0] = 0
			bufIdx[idxPos+1] = (i - 1) << 1
			bufIdx[idxPos+2] = (i << 1)
			idxPos += 3
		}

		// Calculate normals
		tempNormals := w.ReserveVec2(pointsCount)
		i0 := uint32(pointsCount - 1)
		for i1 := uint32(0); i1 < uint32(pointsCount); i1++ {
			p0 := points[i0]
			p1 := points[i1]
			dx := p1.X - p0.X
			dy := p1.Y - p0.Y
			dx, dy = normalize2f(dx, dy)
			tempNormals[i0].X = dy
			tempNormals[i0].Y = -dx
			i0 = i1
		}

		// Set indices and vertices
		vtxPos := 0
		i0 = uint32(pointsCount - 1)
		for i1 := uint32(0); i1 < uint32(pointsCount); i1++ {

			// Average normals
			n0 := tempNormals[i0]
			n1 := tempNormals[i1]
			dmX := (n0.X + n1.X) * 0.5
			dmY := (n0.Y + n1.Y) * 0.5
			dmX, dmY = fixNormal2f(dmX, dmY)
			dmX *= AA_SIZE * 0.5
			dmY *= AA_SIZE * 0.5

			// TODO WRONG IN IMGUI -> ISSUE ????
			// Add inner vertex
			bufVtx[vtxPos+0].Pos.X = points[i1].X + dmX
			bufVtx[vtxPos+0].Pos.Y = points[i1].Y + dmY
			bufVtx[vtxPos+0].Col = col
			// Add outer vertex
			bufVtx[vtxPos+1].Pos.X = points[i1].X - dmX
			bufVtx[vtxPos+1].Pos.Y = points[i1].Y - dmY
			bufVtx[vtxPos+1].Col = colTrans
			vtxPos += 2

			// Add indexes for fringes
			bufIdx[idxPos+0] = 0 + (i1 << 1)
			bufIdx[idxPos+1] = 0 + (i0 << 1)
			bufIdx[idxPos+2] = 1 + (i0 << 1)
			bufIdx[idxPos+3] = 1 + (i0 << 1)
			bufIdx[idxPos+4] = 1 + (i1 << 1)
			bufIdx[idxPos+5] = 0 + (i1 << 1)
			idxPos += 6
			i0 = i1
		}
		dl.AdjustIdx(cmd)
		return
	}

	// Non-Anti aliased filled
	// Allocate command
	idxCount := (pointsCount - 2) * 3
	vtxCount := pointsCount
	cmd, bufIdx, bufVtx := dl.ReserveCmd(idxCount, vtxCount)
	cmd.TexId = w.TexWhiteId

	// Set vertices
	for i := 0; i < vtxCount; i++ {
		bufVtx[i].Pos = points[i]
		bufVtx[i].Col = col
	}

	// Set indices
	idxPos := 0
	for i := 2; i < pointsCount; i++ {
		bufIdx[idxPos+0] = 0
		bufIdx[idxPos+1] = uint32(i - 1)
		bufIdx[idxPos+2] = uint32(i)
		idxPos += 3
	}
	dl.AdjustIdx(cmd)
}
