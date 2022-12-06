package canvas

import (
	"github.com/leonsal/gux/gb"
)

// AddConvexPolyFilled adds a filled convex polygon to this canvas
func (c *Canvas) AddConvexPolyFilled(points []gb.Vec2, col gb.Color, flags Flags) {

	pointsCount := len(points)
	if pointsCount < 3 {
		return
	}

	if (flags & Flag_AntiAliasedFill) != 0 {

		AA_SIZE := c.w.FringeScale
		colTrans := gb.Color(uint32(col) & ^gb.ColorMaskA)

		// Allocates command
		idxCount := (pointsCount-2)*3 + pointsCount*6
		vtxCount := pointsCount * 2
		cmd, bufIdx, bufVtx := c.dl.ReserveCmd(idxCount, vtxCount)
		cmd.TexId = c.w.TexWhiteId

		// Add indexes for fill
		vtxInnerIdx := uint32(0)
		vtxOuterIdx := vtxInnerIdx + 1
		idxPos := 0
		for i := uint32(2); i < uint32(pointsCount); i++ {
			bufIdx[idxPos+0] = vtxInnerIdx
			bufIdx[idxPos+1] = vtxInnerIdx + (i-1)<<1
			bufIdx[idxPos+2] = vtxInnerIdx + (i << 2)
			idxPos++
		}

		// Calculate normals
		tempNormals := c.ReserveVec2(pointsCount)
		i0 := uint32(pointsCount - 1)
		for i1 := uint32(0); i1 < uint32(pointsCount-1); i1++ {
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
		i0 = uint32(pointsCount - 1)
		vtxPos := 0
		idxPos = 0
		for i1 := uint32(0); i1 < uint32(pointsCount-1); i1++ {

			// Average normals
			n0 := tempNormals[i0]
			n1 := tempNormals[i1]
			dmX := (n0.X + n1.X) * 0.5
			dmY := (n0.Y + n1.Y) * 0.5
			dmX, dmY = fixNormal2f(dmX, dmY)
			dmX *= AA_SIZE * 0.5
			dmY *= AA_SIZE * 0.5

			// Add vertices
			// Inner
			bufVtx[vtxPos+0].Pos.X = points[i1].X - dmX
			bufVtx[vtxPos+0].Pos.Y = points[i1].Y - dmY
			bufVtx[vtxPos+0].Col = col
			// Outer
			bufVtx[vtxPos+1].Pos.X = points[i1].X + dmX
			bufVtx[vtxPos+1].Pos.Y = points[i1].Y + dmY
			bufVtx[vtxPos+1].Col = colTrans
			vtxPos += 2

			// Add indexes for fringes
			bufIdx[idxPos+0] = vtxInnerIdx + (i1 << 1)
			bufIdx[idxPos+1] = vtxInnerIdx + (i0 << 1)
			bufIdx[idxPos+2] = vtxOuterIdx + (i0 << 1)
			bufIdx[idxPos+3] = vtxOuterIdx + (i0 << 1)
			bufIdx[idxPos+4] = vtxOuterIdx + (i1 << 1)
			bufIdx[idxPos+5] = vtxInnerIdx + (i1 << 1)
			idxPos += 6
			i0 = i1
		}
		return
	}

	// Non-Anti aliased filled
	// Allocate command
	idxCount := (pointsCount - 2) * 3
	vtxCount := pointsCount
	cmd, bufIdx, bufVtx := c.dl.ReserveCmd(idxCount, vtxCount)
	cmd.TexId = c.w.TexWhiteId

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

}

/***
    if (points_count < 3)
        return;

    const ImVec2 uv = _Data->TexUvWhitePixel;

    if (Flags & ImDrawListFlags_AntiAliasedFill)
    {
        // Anti-aliased Fill
        const float AA_SIZE = _FringeScale;
        const ImU32 col_trans = col & ~IM_COL32_A_MASK;
        const int idx_count = (points_count - 2)*3 + points_count * 6;
        const int vtx_count = (points_count * 2);
        PrimReserve(idx_count, vtx_count);

        // Add indexes for fill
        unsigned int vtx_inner_idx = _VtxCurrentIdx;
        unsigned int vtx_outer_idx = _VtxCurrentIdx + 1;
        for (int i = 2; i < points_count; i++)
        {
            _IdxWritePtr[0] = (ImDrawIdx)(vtx_inner_idx); _IdxWritePtr[1] = (ImDrawIdx)(vtx_inner_idx + ((i - 1) << 1)); _IdxWritePtr[2] = (ImDrawIdx)(vtx_inner_idx + (i << 1));
            _IdxWritePtr += 3;
        }

        // Compute normals
        _Data->TempBuffer.reserve_discard(points_count);
        ImVec2* temp_normals = _Data->TempBuffer.Data;
        for (int i0 = points_count - 1, i1 = 0; i1 < points_count; i0 = i1++)
        {
            const ImVec2& p0 = points[i0];
            const ImVec2& p1 = points[i1];
            float dx = p1.x - p0.x;
            float dy = p1.y - p0.y;
            IM_NORMALIZE2F_OVER_ZERO(dx, dy);
            temp_normals[i0].x = dy;
            temp_normals[i0].y = -dx;
        }

        for (int i0 = points_count - 1, i1 = 0; i1 < points_count; i0 = i1++)
        {
            // Average normals
            const ImVec2& n0 = temp_normals[i0];
            const ImVec2& n1 = temp_normals[i1];
            float dm_x = (n0.x + n1.x) * 0.5f;
            float dm_y = (n0.y + n1.y) * 0.5f;
            IM_FIXNORMAL2F(dm_x, dm_y);
            dm_x *= AA_SIZE * 0.5f;
            dm_y *= AA_SIZE * 0.5f;

            // Add vertices
            _VtxWritePtr[0].pos.x = (points[i1].x - dm_x); _VtxWritePtr[0].pos.y = (points[i1].y - dm_y); _VtxWritePtr[0].uv = uv; _VtxWritePtr[0].col = col;        // Inner
            _VtxWritePtr[1].pos.x = (points[i1].x + dm_x); _VtxWritePtr[1].pos.y = (points[i1].y + dm_y); _VtxWritePtr[1].uv = uv; _VtxWritePtr[1].col = col_trans;  // Outer
            _VtxWritePtr += 2;

            // Add indexes for fringes
            _IdxWritePtr[0] = (ImDrawIdx)(vtx_inner_idx + (i1 << 1)); _IdxWritePtr[1] = (ImDrawIdx)(vtx_inner_idx + (i0 << 1)); _IdxWritePtr[2] = (ImDrawIdx)(vtx_outer_idx + (i0 << 1));
            _IdxWritePtr[3] = (ImDrawIdx)(vtx_outer_idx + (i0 << 1)); _IdxWritePtr[4] = (ImDrawIdx)(vtx_outer_idx + (i1 << 1)); _IdxWritePtr[5] = (ImDrawIdx)(vtx_inner_idx + (i1 << 1));
            _IdxWritePtr += 6;
        }
        _VtxCurrentIdx += (ImDrawIdx)vtx_count;
    }
    else
    {
        // Non Anti-aliased Fill
        const int idx_count = (points_count - 2)*3;
        const int vtx_count = points_count;
        PrimReserve(idx_count, vtx_count);
        for (int i = 0; i < vtx_count; i++)
        {
            _VtxWritePtr[0].pos = points[i]; _VtxWritePtr[0].uv = uv; _VtxWritePtr[0].col = col;
            _VtxWritePtr++;
        }
        for (int i = 2; i < points_count; i++)
        {
            _IdxWritePtr[0] = (ImDrawIdx)(_VtxCurrentIdx); _IdxWritePtr[1] = (ImDrawIdx)(_VtxCurrentIdx + i - 1); _IdxWritePtr[2] = (ImDrawIdx)(_VtxCurrentIdx + i);
            _IdxWritePtr += 3;
        }
        _VtxCurrentIdx += (ImDrawIdx)vtx_count;
    }
}
***/
