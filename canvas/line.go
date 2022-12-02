package canvas

import (
	"math"

	"github.com/leonsal/gux/gb"
)

func (c *Canvas) AddPolyLine(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	//c.polyLineBasic(points, col, flags, thickness)
}

func (c *Canvas) polyLineAntiAliased(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

	/*
		// Anti-aliased stroke
		const AA_SIZE = 1.0
		colTrans := uint32(col) & ^gb.ColorMaskA
		var closed bool
		if (flags & Flag_Closed) != 0 {
			closed = true
		}

		const FringeScale = 1.0
		thickLine := false
		if thickness > FringScale {
			thickLine = true
		}

		// Adjusts line thickness
		if thickness < 1.0 {
			thickness = 1.0
		}
		//_, frac := math.Modf(float64(thickness))
		//fracThickness := float32(frac)

		// Number of line segments to draw
		pointCount := len(points)
		segCount := pointCount - 1
		if closed {
			segCount = pointCount
		}

		// Calculates the number of indices and vertices needed and reserve command
		var idxCount int
		if thickLine {
			idxCount = 18
		} else {
			idxCount = 12
		}
		var vtxCount int
		if thickLine {
			idxCount = pointCount * 4
		} else {
			idxCount = pointCount * 3
		}
		cmd, bufIdx, bufVtx := c.DrawList.ReserveCmd(idxCount, vtxCount)

		// Calculate normals for each line segment: 2 points for each line point.
		tempNormals := c.ReserveVec2(pointCount)
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
	*/
}

func (c *Canvas) AddPolyLineTextured(points []gb.Vec2, col gb.Color, flags Flags, thickness float32) {

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
	pointCount := len(points)
	segCount := pointCount - 1
	if closed {
		segCount = pointCount
	}

	// Calculates the number of indices and vertices needed and reserve command
	idxCount := segCount * 6
	vtxCount := pointCount * 2
	cmd, bufIdx, bufVtx := c.DrawList.ReserveCmd(idxCount, vtxCount)

	// Calculate normals for each line segment: 2 points for each line point.
	tempNormals := c.ReserveVec2(pointCount)
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
	tempPoints := c.ReserveVec2(pointCount * 2)
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
		idx2 := idx1 + 1
		if i1+1 == pointCount {
			idx2 = 0
		} else {
			idx2 = idx1 + 2
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

	// Add vertexes for each point on the line
	vtxPos := 0

	for i := 0; i < pointCount; i++ {
		bufVtx[vtxPos+0].Pos = tempPoints[i*2+0]
		bufVtx[vtxPos+0].Col = col
		bufVtx[vtxPos+1].Pos = tempPoints[i*2+1]
		bufVtx[vtxPos+1].Col = col
		vtxPos += 2
	}
	c.DrawList.AdjustIdx(cmd)
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

//	#define IM_NORMALIZE2F_OVER_ZERO(VX,VY) {
//		   float d2 = VX*VX + VY*VY;
//		   if (d2 > 0.0f) {
//		      float inv_len = ImRsqrt(d2);
//		      VX *= inv_len; VY *= inv_len;
//		   }
//	} (void)0
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
		const maxINVLEN2 = 100.0
		if invLen2 > maxINVLEN2 {
			invLen2 = maxINVLEN2
		}
		return vx * invLen2, vy * invLen2
	}
	return vx, vy
}
