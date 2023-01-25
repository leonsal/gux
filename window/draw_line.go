package window

import (
	"math"

	"github.com/leonsal/gux/gb"
)

func (w *Window) AddPolyLine(dl *gb.DrawList, points []gb.Vec2, col gb.RGBA, flags DrawFlags, thickness float32) {

	if thickness <= float32(TexLinesWidthMax) {
		w.AddPolyLineTextured(dl, points, col, flags, thickness)
	} else {
		w.AddPolyLineAntiAliased(dl, points, col, flags, thickness)
	}
}

func (w *Window) AddPolyLineAntiAliased(dl *gb.DrawList, points []gb.Vec2, col gb.RGBA, flags DrawFlags, thickness float32) {

	// Anti-aliased stroke
	AA_SIZE := w.FringeScale
	colTrans := gb.RGBA(col & ^gb.RGBAMaskA)
	var closed bool
	if (flags & DrawFlags_Closed) != 0 {
		closed = true
	}

	// Adjusts line thickness
	if thickness < 1.0 {
		thickness = 1.0
	}

	// Checks if the line is thick or not
	thickLine := false
	if thickness > w.FringeScale {
		thickLine = true
	}

	// Number of points and line segments to draw
	pointCount := len(points)
	segCount := pointCount - 1
	if closed {
		segCount = pointCount
	}

	// Calculates the number of indices and vertices needed and reserve draw command
	var idxCount int
	var vtxCount int
	if thickLine {
		idxCount = segCount * 18
		vtxCount = pointCount * 4
	} else {
		idxCount = segCount * 12
		vtxCount = pointCount * 3
	}
	_, bufIdx, bufVtx := w.NewDrawCmd(dl, idxCount, vtxCount)

	// Calculate normals for each line segment: 2 points for each line point.
	tempNormals := w.ReserveVec2(pointCount)
	for i1 := 0; i1 < segCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 == pointCount {
			i2 = 0
		}

		// Calculates the normal vector for segment point i1
		dx := points[i2].X - points[i1].X
		dy := points[i2].Y - points[i1].Y
		dx, dy = normalize2f(dx, dy)
		tempNormals[i1].X = dy
		tempNormals[i1].Y = -dx
	}
	if !closed {
		tempNormals[pointCount-1] = tempNormals[pointCount-2]
	}

	// Allocates temporary buffer for points
	tempCount := pointCount * 2
	if thickLine {
		tempCount = pointCount * 4
	}
	tempPoints := w.ReserveVec2(tempCount)

	// One pixel wide line
	if !thickLine {

		/*
			One pixel AA line
			- 3 vertices per point
			- 4 triangles per segment
			- 12 indices per segment
			+-------------------------------------+
			|                                     |	AA fringe
			|                                     |
			X-------------------------------------X Line segment
			|                                     |
			|                                     | AA fringe
			+-------------------------------------+
		*/

		halfDrawSize := float32(AA_SIZE)
		// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
		if !closed {
			tempPoints[0] = gb.Vec2Add(points[0], gb.Vec2MultScalar(tempNormals[0], halfDrawSize))
			tempPoints[1] = gb.Vec2Sub(points[0], gb.Vec2MultScalar(tempNormals[0], halfDrawSize))
			tempPoints[(pointCount-1)*2] = gb.Vec2Add(points[pointCount-1], gb.Vec2MultScalar(tempNormals[pointCount-1], halfDrawSize))
			tempPoints[(pointCount-1)*2+1] = gb.Vec2Sub(points[pointCount-1], gb.Vec2MultScalar(tempNormals[pointCount-1], halfDrawSize))
		}

		// Generate the indices to form 2 triangles for each line segment, and the vertices for the line edges
		// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
		idx1 := uint32(0) // Vertex index for start of line segment
		idxPos := 0       // Start index for indices buffer
		for i1 := 0; i1 < segCount; i1++ {

			// Calculates the index of the next point in the segment
			i2 := i1 + 1
			if i2 == pointCount {
				i2 = 0
			}

			// Calculates vertex index for end of segment
			idx2 := idx1 + 3
			if i1+1 == pointCount {
				idx2 = 0
			}

			// Average normals
			dmX := (tempNormals[i1].X + tempNormals[i2].X) * 0.5
			dmY := (tempNormals[i1].Y + tempNormals[i2].Y) * 0.5
			dmX, dmY = fixNormal2f(dmX, dmY)
			dmX *= halfDrawSize
			dmY *= halfDrawSize

			// Add temporary vertexes for the outer edges
			outVtx := i2 * 2
			tempPoints[outVtx].X = points[i2].X + dmX
			tempPoints[outVtx].Y = points[i2].Y + dmY
			tempPoints[outVtx+1].X = points[i2].X - dmX
			tempPoints[outVtx+1].Y = points[i2].Y - dmY

			// Add indices for four triangles:
			bufIdx[idxPos+0] = idx2
			bufIdx[idxPos+1] = idx1
			bufIdx[idxPos+2] = idx1 + 2
			bufIdx[idxPos+3] = idx1 + 2
			bufIdx[idxPos+4] = idx2 + 2
			bufIdx[idxPos+5] = idx2
			bufIdx[idxPos+6] = idx2 + 1
			bufIdx[idxPos+7] = idx1 + 1
			bufIdx[idxPos+8] = idx1
			bufIdx[idxPos+9] = idx1
			bufIdx[idxPos+10] = idx2
			bufIdx[idxPos+11] = idx2 + 1
			idxPos += 12
			idx1 = idx2
		}

		// Add vertices for each point on the line and the center vertex as well
		vtxPos := 0
		for i := 0; i < pointCount; i++ {
			bufVtx[vtxPos+0].Pos = points[i]
			bufVtx[vtxPos+0].Col = col
			bufVtx[vtxPos+1].Pos = tempPoints[i*2+0]
			bufVtx[vtxPos+1].Col = colTrans
			bufVtx[vtxPos+2].Pos = tempPoints[i*2+1]
			bufVtx[vtxPos+2].Col = colTrans
			vtxPos += 3
		}
		return
	}

	/*
		Non texture-based thick lines:
		- 4 vertices per point
		- 8 triangles per segment
		- 18 indices per segment
		+-------------------------------------+
		|                                     |	AA fringe
		+-------------------------------------+
		|                                     |
		|                                     |
		X-------------------------------------X Line segment
		|                                     |
		|                                     |
		+-------------------------------------+
		|                                     | AA fringe
		+-------------------------------------+
	*/

	// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
	halfInnerThickness := (thickness - AA_SIZE) * 0.5
	pointLast := pointCount - 1
	if !closed {
		tempPoints[0] = gb.Vec2Add(points[0], gb.Vec2MultScalar(tempNormals[0], halfInnerThickness+AA_SIZE))
		tempPoints[1] = gb.Vec2Add(points[0], gb.Vec2MultScalar(tempNormals[0], halfInnerThickness))
		tempPoints[2] = gb.Vec2Sub(points[0], gb.Vec2MultScalar(tempNormals[0], halfInnerThickness))
		tempPoints[3] = gb.Vec2Sub(points[0], gb.Vec2MultScalar(tempNormals[0], halfInnerThickness+AA_SIZE))
		tempPoints[pointLast*4+0] = gb.Vec2Add(points[pointLast], gb.Vec2MultScalar(tempNormals[pointLast], halfInnerThickness+AA_SIZE))
		tempPoints[pointLast*4+1] = gb.Vec2Add(points[pointLast], gb.Vec2MultScalar(tempNormals[pointLast], halfInnerThickness))
		tempPoints[pointLast*4+2] = gb.Vec2Sub(points[pointLast], gb.Vec2MultScalar(tempNormals[pointLast], halfInnerThickness))
		tempPoints[pointLast*4+3] = gb.Vec2Sub(points[pointLast], gb.Vec2MultScalar(tempNormals[pointLast], halfInnerThickness+AA_SIZE))
	}

	// Generate the indices to form 2 triangles for each line segment, and the vertices for the line edges
	// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
	idx1 := uint32(0) // Vertex index for start of line segment
	idxPos := 0       // Start index for indices buffer
	for i1 := 0; i1 < segCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 == pointCount {
			i2 = 0
		}

		// Calculates vertex index for end of segment
		idx2 := idx1 + 4
		if i1+1 == pointCount {
			idx2 = 0
		}

		// Average normals
		dmX := (tempNormals[i1].X + tempNormals[i2].X) * 0.5
		dmY := (tempNormals[i1].Y + tempNormals[i2].Y) * 0.5
		dmX, dmY = fixNormal2f(dmX, dmY)
		dmOutX := dmX * (halfInnerThickness + AA_SIZE)
		dmOutY := dmY * (halfInnerThickness + AA_SIZE)
		dmInX := dmX * halfInnerThickness
		dmInY := dmY * halfInnerThickness

		// Add temporary vertexes for the outer edges
		outVtx := i2 * 4
		tempPoints[outVtx+0].X = points[i2].X + dmOutX
		tempPoints[outVtx+0].Y = points[i2].Y + dmOutY
		tempPoints[outVtx+1].X = points[i2].X + dmInX
		tempPoints[outVtx+1].Y = points[i2].Y + dmInY
		tempPoints[outVtx+2].X = points[i2].X - dmInX
		tempPoints[outVtx+2].Y = points[i2].Y - dmInY
		tempPoints[outVtx+3].X = points[i2].X - dmOutX
		tempPoints[outVtx+3].Y = points[i2].Y - dmOutY

		// Add indices for 8 triangles:
		bufIdx[idxPos+0] = idx2 + 1
		bufIdx[idxPos+1] = idx1 + 1
		bufIdx[idxPos+2] = idx1 + 2
		bufIdx[idxPos+3] = idx1 + 2
		bufIdx[idxPos+4] = idx2 + 2
		bufIdx[idxPos+5] = idx2 + 1
		bufIdx[idxPos+6] = idx2 + 1
		bufIdx[idxPos+7] = idx1 + 1
		bufIdx[idxPos+8] = idx1
		bufIdx[idxPos+9] = idx1
		bufIdx[idxPos+10] = idx2
		bufIdx[idxPos+11] = idx2 + 1
		bufIdx[idxPos+12] = idx2 + 2
		bufIdx[idxPos+13] = idx1 + 2
		bufIdx[idxPos+14] = idx1 + 3
		bufIdx[idxPos+15] = idx1 + 3
		bufIdx[idxPos+16] = idx2 + 3
		bufIdx[idxPos+17] = idx2 + 2
		idxPos += 18
		idx1 = idx2

	}
	// Add vertices
	vtxPos := 0
	for i := 0; i < pointCount; i++ {
		bufVtx[vtxPos+0].Pos = tempPoints[i*4+0]
		bufVtx[vtxPos+0].Col = colTrans
		bufVtx[vtxPos+1].Pos = tempPoints[i*4+1]
		bufVtx[vtxPos+1].Col = col
		bufVtx[vtxPos+2].Pos = tempPoints[i*4+2]
		bufVtx[vtxPos+2].Col = col
		bufVtx[vtxPos+3].Pos = tempPoints[i*4+3]
		bufVtx[vtxPos+3].Col = colTrans
		vtxPos += 4
	}
}

func (w *Window) AddPolyLineTextured(dl *gb.DrawList, points []gb.Vec2, col gb.RGBA, flags DrawFlags, thickness float32) {

	/*
		- 2 vertices per point
		- 2 triangles per segment
		- 6 indices per segment
		+-------------------------------------+
		|                                     |
		|                                     |
		X-------------------------------------X Line segment
		|                                     |
		|                                     |
		+-------------------------------------+
	*/

	// Checks if 'flags' specifies closed line path (last point == first point)
	var closed bool
	if (flags & DrawFlags_Closed) != 0 {
		closed = true
	}

	// Adjusts line thickness
	if thickness < 1.0 {
		thickness = 1.0
	}

	// Number of line segments to draw
	pointCount := len(points)
	segCount := pointCount - 1
	if closed {
		segCount = pointCount
	}

	// Calculates the number of indices and vertices needed and reserve command
	idxCount := segCount * 6
	vtxCount := pointCount * 2
	cmd, bufIdx, bufVtx := w.NewDrawCmd(dl, idxCount, vtxCount)

	// Calculate normals for each line segment: 2 points for each line point.
	tempNormals := w.ReserveVec2(pointCount)
	for i1 := 0; i1 < segCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 == pointCount {
			i2 = 0
		}

		// Calculates the normal vector for segment point i1
		dx := points[i2].X - points[i1].X
		dy := points[i2].Y - points[i1].Y
		dx, dy = normalize2f(dx, dy)
		tempNormals[i1].X = dy
		tempNormals[i1].Y = -dx
	}
	if !closed {
		tempNormals[pointCount-1] = tempNormals[pointCount-2]
	}

	// Generates
	tempPoints := w.ReserveVec2(pointCount * 2)
	halfDrawSize := (thickness * 0.5) + 1
	// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
	if !closed {
		tempPoints[0] = gb.Vec2Add(points[0], gb.Vec2MultScalar(tempNormals[0], halfDrawSize))
		tempPoints[1] = gb.Vec2Sub(points[0], gb.Vec2MultScalar(tempNormals[0], halfDrawSize))
		tempPoints[(pointCount-1)*2] = gb.Vec2Add(points[pointCount-1], gb.Vec2MultScalar(tempNormals[pointCount-1], halfDrawSize))
		tempPoints[(pointCount-1)*2+1] = gb.Vec2Sub(points[pointCount-1], gb.Vec2MultScalar(tempNormals[pointCount-1], halfDrawSize))
	}

	// Generate the indices to form 2 triangles for each line segment, and the vertices for the line edges
	// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
	idx1 := uint32(0) // Vertex index for start of line segment
	idxPos := 0       // Start index for indices buffer
	for i1 := 0; i1 < segCount; i1++ {

		// Calculates the index of the next point in the segment
		i2 := i1 + 1
		if i2 == pointCount {
			i2 = 0
		}

		// Calculates vertex index for end of segment
		idx2 := idx1 + 2
		if i1+1 == pointCount {
			idx2 = 0
		}

		// Average normals
		dmX := (tempNormals[i1].X + tempNormals[i2].X) * 0.5
		dmY := (tempNormals[i1].Y + tempNormals[i2].Y) * 0.5
		dmX, dmY = fixNormal2f(dmX, dmY)
		dmX *= halfDrawSize
		dmY *= halfDrawSize

		// Add temporary vertexes for the outer edges
		outVtx := i2 * 2
		tempPoints[outVtx].X = points[i2].X + dmX
		tempPoints[outVtx].Y = points[i2].Y + dmY
		tempPoints[outVtx+1].X = points[i2].X - dmX
		tempPoints[outVtx+1].Y = points[i2].Y - dmY

		// Add indices for two triangles
		bufIdx[idxPos] = idx2 // Right triangle
		bufIdx[idxPos+1] = idx1
		bufIdx[idxPos+2] = idx1 + 1
		bufIdx[idxPos+3] = idx2 + 1 // Left triangle
		bufIdx[idxPos+4] = idx1 + 1
		bufIdx[idxPos+5] = idx2
		idxPos += 6
		idx1 = idx2
	}

	// If line width less than maximum for textured lines
	// sets the command texture id and UV coordinates
	intThickness := int(thickness)
	texUv0 := gb.Vec2{}
	texUv1 := gb.Vec2{}
	if intThickness < TexLinesWidthMax {
		cmd.TexID = w.TexLinesId
		texUvs := w.TexUvLines[intThickness]
		texUv0 = gb.Vec2{texUvs.X, texUvs.Y}
		texUv1 = gb.Vec2{texUvs.Z, texUvs.W}
	}

	// Add vertexes for each point on the line
	vtxPos := 0
	for i := 0; i < pointCount; i++ {
		bufVtx[vtxPos+0].Pos = tempPoints[i*2+0]
		bufVtx[vtxPos+0].UV = texUv0
		bufVtx[vtxPos+0].Col = col
		bufVtx[vtxPos+1].Pos = tempPoints[i*2+1]
		bufVtx[vtxPos+1].UV = texUv1
		bufVtx[vtxPos+1].Col = col
		vtxPos += 2
	}
}

func normalize2f(vx, vy float32) (float32, float32) {

	d2 := vx*vx + vy*vy
	if d2 > 0 {
		invLen := 1.0 / math.Sqrt(float64(d2))
		return vx * float32(invLen), vy * float32(invLen)
	}
	return vx, vy
}

func fixNormal2f(vx, vy float32) (float32, float32) {

	d2 := vx*vx + vy*vy
	if d2 > 0.000001 {
		invLen2 := 1.0 / d2
		const maxINVLEN2 = 100.0
		if invLen2 > maxINVLEN2 {
			invLen2 = maxINVLEN2
		}
		return vx * invLen2, vy * invLen2
	}
	return vx, vy
}
