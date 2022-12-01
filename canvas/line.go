package canvas

import (
	"math"

	"github.com/leonsal/gux/gb"
)

func (c *Canvas) AddPolyLine(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	//c.polyLineBasic(points, col, flags, thickness)
}

func (c *Canvas) polyLineAntiAliased(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	// Anti-aliased stroke
	//const AA_SIZE = 1.0
	//colTrans := uint32(col) & ^gb.ColorMaskA
	//var closed bool
	//if (flags & Flag_Closed) != 0 {
	//	closed = true
	//}
}

func (c *Canvas) polyLineTextured(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	// Checks if 'flags' specifies closed line path (last point == first point)
	var closed bool
	if (flags & Flag_Closed) != 0 {
		closed = true
	}

	// Adjusts line thickness
	if thickness < 1.0 {
		thickness = 1.0
	}
	//_, frac := math.Modf(float64(thickness))
	//fracThickness := float32(frac)

	// Number of line segments to draw
	pointCount := len(points) - 1
	segCount := pointCount - 1
	if closed {
		segCount = pointCount
	}

	// Calculates the number of indices and vertices needed and reserve command
	idxCount := segCount * 6
	vtxCount := pointCount * 2
	_, bufIdx, bufVtx := c.DrawList.ReserveCmd(idxCount, vtxCount)

	// Calculate normals for each line segment: 2 points for each line point.
	tempNormals := c.ReserveVec2(pointCount * 2)
	for i1 := 0; i1 < segCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 > pointCount {
			i2 = 0
		}

		dx := points[i2].X - points[i1].X
		dy := points[i2].Y - points[i1].Y
		tempNormals[i1].X = dy
		tempNormals[i1].Y = -dx
	}
	if closed {
		tempNormals[pointCount-1] = tempNormals[pointCount-2]
	}

	// Generates
	tempPoints := c.ReserveVec2(pointCount)
	halfDrawSize := (thickness * 0.5) + 1
	// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
	if !closed {
		tempPoints[0] = *points[0].Add(tempNormals[0].MultScalar(halfDrawSize))
		tempPoints[1] = *points[0].Sub(tempNormals[0].MultScalar(halfDrawSize))
		tempPoints[pointCount-1] = *points[0].Add(tempNormals[pointCount-1].MultScalar(halfDrawSize))
		tempPoints[pointCount] = *points[0].Sub(tempNormals[pointCount-1].MultScalar(halfDrawSize))

	}

	// Generate the indices to form 2 triangles for each line segment, and the vertices for the line edges
	// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
	idx1 := uint32(0) // Vertex index for start of line segment
	idxPos := 0       // Start index for indices buffer
	for i1 := 0; i1 < pointCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 == pointCount {
			i2 = 0
		}

		// Calculates vertex index for end of segment
		var idx2 uint32
		if i1+1 == pointCount {
			idx2 = idx1
		} else {
			idx2 = idx1 + 2
		}

		// Average normals
		dmX := (tempNormals[i1].X + tempNormals[i2].X) * 0.5
		dmY := (tempNormals[i1].Y + tempNormals[i2].Y) * 0.5
		dmX, dmY = fixNormal2f(dmX, dmY)

		// Add temporary vertexes for the outer edges
		outVtx := i2 * 2
		tempPoints[outVtx].X = points[i2].X + dmX
		tempPoints[outVtx].Y = points[i2].X + dmY
		tempPoints[outVtx+1].X = points[i2].X - dmY
		tempPoints[outVtx+1].Y = points[i2].Y - dmY

		// Add indices for two triangles
		bufIdx[idxPos] = idx2 // Right triangle
		bufIdx[idxPos+1] = idx1
		bufIdx[idxPos+2] = idx1 + 2
		bufIdx[idxPos+3] = idx2 + 1 // Left triangle
		bufIdx[idxPos+4] = idx1 + 1
		bufIdx[idxPos+5] = idx2
		idxPos += 6
	}

	// Add vertexes for each point on the line
	vtxPos := 0
	for i := 0; i < pointCount; i++ {
		bufVtx[vtxPos+0].Pos = tempPoints[i*2+0]
		bufVtx[vtxPos+0].Col = col
		bufVtx[vtxPos+1].Pos = tempPoints[i*2+1]
		bufVtx[vtxPos+1].Col = col
	}
}

func (c *Canvas) AddPolyLineBasic(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	pointCount := len(points) - 1
	cmd, indices, vertices := c.DrawList.ReserveCmd(pointCount*6, pointCount*4)
	uv := gb.Vec2{0, 0}

	vtxCurrent := uint32(0)
	vtx := 0
	idx := 0
	for i1 := 0; i1 < pointCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 > pointCount {
			i2 = 0
		}

		// Calculate normals
		p1 := points[i1]
		p2 := points[i2]
		dx := p2.X - p1.X
		dy := p2.Y - p1.Y
		dx, dy = normalize2f(dx, dy)
		dx *= thickness * 0.5
		dy *= thickness * 0.5

		vertices[vtx].Pos.X = p1.X + dy
		vertices[vtx].Pos.Y = p1.Y - dx
		vertices[vtx].UV = uv
		vertices[vtx].Col = col
		vtx++

		vertices[vtx].Pos.X = p2.X + dy
		vertices[vtx].Pos.Y = p2.Y - dx
		vertices[vtx].UV = uv
		vertices[vtx].Col = col
		vtx++

		vertices[vtx].Pos.X = p2.X - dy
		vertices[vtx].Pos.Y = p2.Y + dx
		vertices[vtx].UV = uv
		vertices[vtx].Col = col
		vtx++

		vertices[vtx].Pos.X = p1.X - dy
		vertices[vtx].Pos.Y = p1.Y + dx
		vertices[vtx].UV = uv
		vertices[vtx].Col = col
		vtx++

		indices[idx] = vtxCurrent
		idx++
		indices[idx] = vtxCurrent + 1
		idx++
		indices[idx] = vtxCurrent + 2
		idx++
		indices[idx] = vtxCurrent
		idx++
		indices[idx] = vtxCurrent + 2
		idx++
		indices[idx] = vtxCurrent + 3
		idx++

		vtxCurrent += 4
	}
	c.DrawList.AdjustIdx(cmd)
}

// #define IM_NORMALIZE2F_OVER_ZERO(VX,VY)     { float d2 = VX*VX + VY*VY; if (d2 > 0.0f) { float inv_len = ImRsqrt(d2); VX *= inv_len; VY *= inv_len; } } (void)0
func normalize2f(vx, vy float32) (float32, float32) {

	d2 := vx*vx + vy*vy
	if d2 > 0 {
		invLen := 1.0 / math.Sqrt(float64(d2))
		return vx * float32(invLen), vy * float32(invLen)
	}
	return vx, vy
}

// #define IM_FIXNORMAL2F(VX,VY)
// float d2 = VX*VX + VY*VY;
//
//	if (d2 > 0.000001f) {
//	   float inv_len2 = 1.0f / d2;
//	      if (inv_len2 > IM_FIXNORMAL2F_MAX_INVLEN2) {
//	         inv_len2 = IM_FIXNORMAL2F_MAX_INVLEN2;
//	      }
//	      VX *= inv_len2; VY *= inv_len2;
//	}
func fixNormal2f(vx, vy float32) (float32, float32) {

	d2 := vx*vx + vy*vy
	if d2 > 0.000001 {
		invLen2 := 1.0 / d2
		const maxINVLEN2 = float32(100.0) / float32(500.0)
		if invLen2 > maxINVLEN2 {
			invLen2 = maxINVLEN2
		}
		return vx * invLen2, vy * invLen2
	}
	return vx, vy
}

//// On AddPolyline() and AddConvexPolyFilled() we intentionally avoid using ImVec2 and superfluous function calls to optimize debug/non-inlined builds.
//// - Those macros expects l-values and need to be used as their own statement.
//// - Those macros are intentionally not surrounded by the 'do {} while (0)' idiom because even that translates to runtime with debug compilers.
//#define IM_NORMALIZE2F_OVER_ZERO(VX,VY)     { float d2 = VX*VX + VY*VY; if (d2 > 0.0f) { float inv_len = ImRsqrt(d2); VX *= inv_len; VY *= inv_len; } } (void)0
//#define IM_FIXNORMAL2F_MAX_INVLEN2          100.0f // 500.0f (see #4053, #3366)
//#define IM_FIXNORMAL2F(VX,VY)               { float d2 = VX*VX + VY*VY; if (d2 > 0.000001f) { float inv_len2 = 1.0f / d2; if (inv_len2 > IM_FIXNORMAL2F_MAX_INVLEN2) inv_len2 = IM_FIXNORMAL2F_MAX_INVLEN2; VX *= inv_len2; VY *= inv_len2; } } (void)0
//
//// TODO: Thickness anti-aliased lines cap are missing their AA fringe.
//// We avoid using the ImVec2 math operators here to reduce cost to a minimum for debug/non-inlined builds.
//void ImDrawList::AddPolyline(const ImVec2* points, const int points_count, ImU32 col, ImDrawFlags flags, float thickness)
//{
//    if (points_count < 2)
//        return;
//
//    const bool closed = (flags & ImDrawFlags_Closed) != 0;
//    const ImVec2 opaque_uv = _Data->TexUvWhitePixel;
//    const int count = closed ? points_count : points_count - 1; // The number of line segments we need to draw
//    const bool thick_line = (thickness > _FringeScale);
//
//    if (Flags & ImDrawListFlags_AntiAliasedLines)
//    {
//        // Anti-aliased stroke
//        const float AA_SIZE = _FringeScale;
//        const ImU32 col_trans = col & ~IM_COL32_A_MASK;
//
//        // Thicknesses <1.0 should behave like thickness 1.0
//        thickness = ImMax(thickness, 1.0f);
//        const int integer_thickness = (int)thickness;
//        const float fractional_thickness = thickness - integer_thickness;
//
//        // Do we want to draw this line using a texture?
//        // - For now, only draw integer-width lines using textures to avoid issues with the way scaling occurs, could be improved.
//        // - If AA_SIZE is not 1.0f we cannot use the texture path.
//        const bool use_texture = (Flags & ImDrawListFlags_AntiAliasedLinesUseTex) && (integer_thickness < IM_DRAWLIST_TEX_LINES_WIDTH_MAX) && (fractional_thickness <= 0.00001f) && (AA_SIZE == 1.0f);
//
//        // We should never hit this, because NewFrame() doesn't set ImDrawListFlags_AntiAliasedLinesUseTex unless ImFontAtlasFlags_NoBakedLines is off
//        IM_ASSERT_PARANOID(!use_texture || !(_Data->Font->ContainerAtlas->Flags & ImFontAtlasFlags_NoBakedLines));
//
//        const int idx_count = use_texture ? (count * 6) : (thick_line ? count * 18 : count * 12);
//        const int vtx_count = use_texture ? (points_count * 2) : (thick_line ? points_count * 4 : points_count * 3);
//        PrimReserve(idx_count, vtx_count);
//
//        // Temporary buffer
//        // The first <points_count> items are normals at each line point, then after that there are either 2 or 4 temp points for each line point
//        _Data->TempBuffer.reserve_discard(points_count * ((use_texture || !thick_line) ? 3 : 5));
//        ImVec2* temp_normals = _Data->TempBuffer.Data;
//        ImVec2* temp_points = temp_normals + points_count;
//
//        // Calculate normals (tangents) for each line segment
//        for (int i1 = 0; i1 < count; i1++)
//        {
//            const int i2 = (i1 + 1) == points_count ? 0 : i1 + 1;
//            float dx = points[i2].x - points[i1].x;
//            float dy = points[i2].y - points[i1].y;
//            IM_NORMALIZE2F_OVER_ZERO(dx, dy);
//            temp_normals[i1].x = dy;
//            temp_normals[i1].y = -dx;
//        }
//        if (!closed)
//            temp_normals[points_count - 1] = temp_normals[points_count - 2];
//
//        // If we are drawing a one-pixel-wide line without a texture, or a textured line of any width, we only need 2 or 3 vertices per point
//        if (use_texture || !thick_line)
//        {
//            // [PATH 1] Texture-based lines (thick or non-thick)
//            // [PATH 2] Non texture-based lines (non-thick)
//
//            // The width of the geometry we need to draw - this is essentially <thickness> pixels for the line itself, plus "one pixel" for AA.
//            // - In the texture-based path, we don't use AA_SIZE here because the +1 is tied to the generated texture
//            //   (see ImFontAtlasBuildRenderLinesTexData() function), and so alternate values won't work without changes to that code.
//            // - In the non texture-based paths, we would allow AA_SIZE to potentially be != 1.0f with a patch (e.g. fringe_scale patch to
//            //   allow scaling geometry while preserving one-screen-pixel AA fringe).
//            const float half_draw_size = use_texture ? ((thickness * 0.5f) + 1) : AA_SIZE;
//
//            // If line is not closed, the first and last points need to be generated differently as there are no normals to blend
//            if (!closed)
//            {
//                temp_points[0] = points[0] + temp_normals[0] * half_draw_size;
//                temp_points[1] = points[0] - temp_normals[0] * half_draw_size;
//                temp_points[(points_count-1)*2+0] = points[points_count-1] + temp_normals[points_count-1] * half_draw_size;
//                temp_points[(points_count-1)*2+1] = points[points_count-1] - temp_normals[points_count-1] * half_draw_size;
//            }
//
//            // Generate the indices to form a number of triangles for each line segment, and the vertices for the line edges
//            // This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
//            // FIXME-OPT: Merge the different loops, possibly remove the temporary buffer.
//            unsigned int idx1 = _VtxCurrentIdx; // Vertex index for start of line segment
//            for (int i1 = 0; i1 < count; i1++) // i1 is the first point of the line segment
//            {
//                const int i2 = (i1 + 1) == points_count ? 0 : i1 + 1; // i2 is the second point of the line segment
//                const unsigned int idx2 = ((i1 + 1) == points_count) ? _VtxCurrentIdx : (idx1 + (use_texture ? 2 : 3)); // Vertex index for end of segment
//
//                // Average normals
//                float dm_x = (temp_normals[i1].x + temp_normals[i2].x) * 0.5f;
//                float dm_y = (temp_normals[i1].y + temp_normals[i2].y) * 0.5f;
//                IM_FIXNORMAL2F(dm_x, dm_y);
//                dm_x *= half_draw_size; // dm_x, dm_y are offset to the outer edge of the AA area
//                dm_y *= half_draw_size;
//
//                // Add temporary vertexes for the outer edges
//                ImVec2* out_vtx = &temp_points[i2 * 2];
//                out_vtx[0].x = points[i2].x + dm_x;
//                out_vtx[0].y = points[i2].y + dm_y;
//                out_vtx[1].x = points[i2].x - dm_x;
//                out_vtx[1].y = points[i2].y - dm_y;
//
//                if (use_texture)
//                {
//                    // Add indices for two triangles
//                    _IdxWritePtr[0] = (ImDrawIdx)(idx2 + 0); _IdxWritePtr[1] = (ImDrawIdx)(idx1 + 0); _IdxWritePtr[2] = (ImDrawIdx)(idx1 + 1); // Right tri
//                    _IdxWritePtr[3] = (ImDrawIdx)(idx2 + 1); _IdxWritePtr[4] = (ImDrawIdx)(idx1 + 1); _IdxWritePtr[5] = (ImDrawIdx)(idx2 + 0); // Left tri
//                    _IdxWritePtr += 6;
//                }
//                else
//                {
//                    // Add indexes for four triangles
//                    _IdxWritePtr[0] = (ImDrawIdx)(idx2 + 0); _IdxWritePtr[1] = (ImDrawIdx)(idx1 + 0); _IdxWritePtr[2] = (ImDrawIdx)(idx1 + 2); // Right tri 1
//                    _IdxWritePtr[3] = (ImDrawIdx)(idx1 + 2); _IdxWritePtr[4] = (ImDrawIdx)(idx2 + 2); _IdxWritePtr[5] = (ImDrawIdx)(idx2 + 0); // Right tri 2
//                    _IdxWritePtr[6] = (ImDrawIdx)(idx2 + 1); _IdxWritePtr[7] = (ImDrawIdx)(idx1 + 1); _IdxWritePtr[8] = (ImDrawIdx)(idx1 + 0); // Left tri 1
//                    _IdxWritePtr[9] = (ImDrawIdx)(idx1 + 0); _IdxWritePtr[10] = (ImDrawIdx)(idx2 + 0); _IdxWritePtr[11] = (ImDrawIdx)(idx2 + 1); // Left tri 2
//                    _IdxWritePtr += 12;
//                }
//
//                idx1 = idx2;
//            }
//
//            // Add vertexes for each point on the line
//            if (use_texture)
//            {
//                // If we're using textures we only need to emit the left/right edge vertices
//                ImVec4 tex_uvs = _Data->TexUvLines[integer_thickness];
//                /*if (fractional_thickness != 0.0f) // Currently always zero when use_texture==false!
//                {
//                    const ImVec4 tex_uvs_1 = _Data->TexUvLines[integer_thickness + 1];
//                    tex_uvs.x = tex_uvs.x + (tex_uvs_1.x - tex_uvs.x) * fractional_thickness; // inlined ImLerp()
//                    tex_uvs.y = tex_uvs.y + (tex_uvs_1.y - tex_uvs.y) * fractional_thickness;
//                    tex_uvs.z = tex_uvs.z + (tex_uvs_1.z - tex_uvs.z) * fractional_thickness;
//                    tex_uvs.w = tex_uvs.w + (tex_uvs_1.w - tex_uvs.w) * fractional_thickness;
//                }*/
//                ImVec2 tex_uv0(tex_uvs.x, tex_uvs.y);
//                ImVec2 tex_uv1(tex_uvs.z, tex_uvs.w);
//                for (int i = 0; i < points_count; i++)
//                {
//                    _VtxWritePtr[0].pos = temp_points[i * 2 + 0]; _VtxWritePtr[0].uv = tex_uv0; _VtxWritePtr[0].col = col; // Left-side outer edge
//                    _VtxWritePtr[1].pos = temp_points[i * 2 + 1]; _VtxWritePtr[1].uv = tex_uv1; _VtxWritePtr[1].col = col; // Right-side outer edge
//                    _VtxWritePtr += 2;
//                }
//            }
//            else
//            {
//                // If we're not using a texture, we need the center vertex as well
//                for (int i = 0; i < points_count; i++)
//                {
//                    _VtxWritePtr[0].pos = points[i];              _VtxWritePtr[0].uv = opaque_uv; _VtxWritePtr[0].col = col;       // Center of line
//                    _VtxWritePtr[1].pos = temp_points[i * 2 + 0]; _VtxWritePtr[1].uv = opaque_uv; _VtxWritePtr[1].col = col_trans; // Left-side outer edge
//                    _VtxWritePtr[2].pos = temp_points[i * 2 + 1]; _VtxWritePtr[2].uv = opaque_uv; _VtxWritePtr[2].col = col_trans; // Right-side outer edge
//                    _VtxWritePtr += 3;
//                }
//            }
//        }
//        else
//        {
//            // [PATH 2] Non texture-based lines (thick): we need to draw the solid line core and thus require four vertices per point
//            const float half_inner_thickness = (thickness - AA_SIZE) * 0.5f;
//
//            // If line is not closed, the first and last points need to be generated differently as there are no normals to blend
//            if (!closed)
//            {
//                const int points_last = points_count - 1;
//                temp_points[0] = points[0] + temp_normals[0] * (half_inner_thickness + AA_SIZE);
//                temp_points[1] = points[0] + temp_normals[0] * (half_inner_thickness);
//                temp_points[2] = points[0] - temp_normals[0] * (half_inner_thickness);
//                temp_points[3] = points[0] - temp_normals[0] * (half_inner_thickness + AA_SIZE);
//                temp_points[points_last * 4 + 0] = points[points_last] + temp_normals[points_last] * (half_inner_thickness + AA_SIZE);
//                temp_points[points_last * 4 + 1] = points[points_last] + temp_normals[points_last] * (half_inner_thickness);
//                temp_points[points_last * 4 + 2] = points[points_last] - temp_normals[points_last] * (half_inner_thickness);
//                temp_points[points_last * 4 + 3] = points[points_last] - temp_normals[points_last] * (half_inner_thickness + AA_SIZE);
//            }
//
//            // Generate the indices to form a number of triangles for each line segment, and the vertices for the line edges
//            // This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
//            // FIXME-OPT: Merge the different loops, possibly remove the temporary buffer.
//            unsigned int idx1 = _VtxCurrentIdx; // Vertex index for start of line segment
//            for (int i1 = 0; i1 < count; i1++) // i1 is the first point of the line segment
//            {
//                const int i2 = (i1 + 1) == points_count ? 0 : (i1 + 1); // i2 is the second point of the line segment
//                const unsigned int idx2 = (i1 + 1) == points_count ? _VtxCurrentIdx : (idx1 + 4); // Vertex index for end of segment
//
//                // Average normals
//                float dm_x = (temp_normals[i1].x + temp_normals[i2].x) * 0.5f;
//                float dm_y = (temp_normals[i1].y + temp_normals[i2].y) * 0.5f;
//                IM_FIXNORMAL2F(dm_x, dm_y);
//                float dm_out_x = dm_x * (half_inner_thickness + AA_SIZE);
//                float dm_out_y = dm_y * (half_inner_thickness + AA_SIZE);
//                float dm_in_x = dm_x * half_inner_thickness;
//                float dm_in_y = dm_y * half_inner_thickness;
//
//                // Add temporary vertices
//                ImVec2* out_vtx = &temp_points[i2 * 4];
//                out_vtx[0].x = points[i2].x + dm_out_x;
//                out_vtx[0].y = points[i2].y + dm_out_y;
//                out_vtx[1].x = points[i2].x + dm_in_x;
//                out_vtx[1].y = points[i2].y + dm_in_y;
//                out_vtx[2].x = points[i2].x - dm_in_x;
//                out_vtx[2].y = points[i2].y - dm_in_y;
//                out_vtx[3].x = points[i2].x - dm_out_x;
//                out_vtx[3].y = points[i2].y - dm_out_y;
//
//                // Add indexes
//                _IdxWritePtr[0]  = (ImDrawIdx)(idx2 + 1); _IdxWritePtr[1]  = (ImDrawIdx)(idx1 + 1); _IdxWritePtr[2]  = (ImDrawIdx)(idx1 + 2);
//                _IdxWritePtr[3]  = (ImDrawIdx)(idx1 + 2); _IdxWritePtr[4]  = (ImDrawIdx)(idx2 + 2); _IdxWritePtr[5]  = (ImDrawIdx)(idx2 + 1);
//                _IdxWritePtr[6]  = (ImDrawIdx)(idx2 + 1); _IdxWritePtr[7]  = (ImDrawIdx)(idx1 + 1); _IdxWritePtr[8]  = (ImDrawIdx)(idx1 + 0);
//                _IdxWritePtr[9]  = (ImDrawIdx)(idx1 + 0); _IdxWritePtr[10] = (ImDrawIdx)(idx2 + 0); _IdxWritePtr[11] = (ImDrawIdx)(idx2 + 1);
//                _IdxWritePtr[12] = (ImDrawIdx)(idx2 + 2); _IdxWritePtr[13] = (ImDrawIdx)(idx1 + 2); _IdxWritePtr[14] = (ImDrawIdx)(idx1 + 3);
//                _IdxWritePtr[15] = (ImDrawIdx)(idx1 + 3); _IdxWritePtr[16] = (ImDrawIdx)(idx2 + 3); _IdxWritePtr[17] = (ImDrawIdx)(idx2 + 2);
//                _IdxWritePtr += 18;
//
//                idx1 = idx2;
//            }
//
//            // Add vertices
//            for (int i = 0; i < points_count; i++)
//            {
//                _VtxWritePtr[0].pos = temp_points[i * 4 + 0]; _VtxWritePtr[0].uv = opaque_uv; _VtxWritePtr[0].col = col_trans;
//                _VtxWritePtr[1].pos = temp_points[i * 4 + 1]; _VtxWritePtr[1].uv = opaque_uv; _VtxWritePtr[1].col = col;
//                _VtxWritePtr[2].pos = temp_points[i * 4 + 2]; _VtxWritePtr[2].uv = opaque_uv; _VtxWritePtr[2].col = col;
//                _VtxWritePtr[3].pos = temp_points[i * 4 + 3]; _VtxWritePtr[3].uv = opaque_uv; _VtxWritePtr[3].col = col_trans;
//                _VtxWritePtr += 4;
//            }
//        }
//        _VtxCurrentIdx += (ImDrawIdx)vtx_count;
//    }
//    else
//    {
//        // [PATH 4] Non texture-based, Non anti-aliased lines
//        const int idx_count = count * 6;
//        const int vtx_count = count * 4;    // FIXME-OPT: Not sharing edges
//        PrimReserve(idx_count, vtx_count);
//
//        for (int i1 = 0; i1 < count; i1++)
//        {
//            const int i2 = (i1 + 1) == points_count ? 0 : i1 + 1;
//            const ImVec2& p1 = points[i1];
//            const ImVec2& p2 = points[i2];
//
//            float dx = p2.x - p1.x;
//            float dy = p2.y - p1.y;
//            IM_NORMALIZE2F_OVER_ZERO(dx, dy);
//            dx *= (thickness * 0.5f);
//            dy *= (thickness * 0.5f);
//
//            _VtxWritePtr[0].pos.x = p1.x + dy; _VtxWritePtr[0].pos.y = p1.y - dx; _VtxWritePtr[0].uv = opaque_uv; _VtxWritePtr[0].col = col;
//            _VtxWritePtr[1].pos.x = p2.x + dy; _VtxWritePtr[1].pos.y = p2.y - dx; _VtxWritePtr[1].uv = opaque_uv; _VtxWritePtr[1].col = col;
//            _VtxWritePtr[2].pos.x = p2.x - dy; _VtxWritePtr[2].pos.y = p2.y + dx; _VtxWritePtr[2].uv = opaque_uv; _VtxWritePtr[2].col = col;
//            _VtxWritePtr[3].pos.x = p1.x - dy; _VtxWritePtr[3].pos.y = p1.y + dx; _VtxWritePtr[3].uv = opaque_uv; _VtxWritePtr[3].col = col;
//            _VtxWritePtr += 4;
//
//            _IdxWritePtr[0] = (ImDrawIdx)(_VtxCurrentIdx); _IdxWritePtr[1] = (ImDrawIdx)(_VtxCurrentIdx + 1); _IdxWritePtr[2] = (ImDrawIdx)(_VtxCurrentIdx + 2);
//            _IdxWritePtr[3] = (ImDrawIdx)(_VtxCurrentIdx); _IdxWritePtr[4] = (ImDrawIdx)(_VtxCurrentIdx + 2); _IdxWritePtr[5] = (ImDrawIdx)(_VtxCurrentIdx + 3);
//            _IdxWritePtr += 6;
//            _VtxCurrentIdx += 4;
//        }
//    }
//}
