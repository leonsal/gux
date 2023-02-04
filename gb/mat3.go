package gb

import (
	"errors"
	"math"
)

// Mat3 is 3x3 matrix organized internally as column matrix
type Mat3 [9]float32

// Set sets all the elements of the matrix row by row starting at row1, column1,
// row1, column2, row1, column3 and so forth.
// Returns the pointer to this updated Matrix.
func (m *Mat3) Set(n11, n12, n13, n21, n22, n23, n31, n32, n33 float32) *Mat3 {

	m[0] = n11
	m[3] = n12
	m[6] = n13
	m[1] = n21
	m[4] = n22
	m[7] = n23
	m[2] = n31
	m[5] = n32
	m[8] = n33
	return m
}

// Identity sets this matrix as the identity matrix.
// Returns the pointer to this updated matrix.
func (m *Mat3) Identity() *Mat3 {

	m.Set(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
	return m
}

// Zero sets this matrix as the zero matrix.
// Returns the pointer to this updated matrix.
func (m *Mat3) Zero() *Mat3 {

	m.Set(
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	)
	return m
}

// Multiply multiply this matrix by the other matrix
// Returns pointer to this updated matrix.
func (m *Mat3) Mult(other *Mat3) *Mat3 {

	return m.MultMat(m, other)
}

// MultiplyMatrices multiply matrix a by b storing the result in this matrix.
// Returns pointer to this updated matrix.
func (m *Mat3) MultMat(a, b *Mat3) *Mat3 {

	a11 := a[0]
	a12 := a[3]
	a13 := a[6]
	a21 := a[1]
	a22 := a[4]
	a23 := a[7]
	a31 := a[2]
	a32 := a[5]
	a33 := a[8]

	b11 := b[0]
	b12 := b[3]
	b13 := b[6]
	b21 := b[1]
	b22 := b[4]
	b23 := b[7]
	b31 := b[2]
	b32 := b[5]
	b33 := b[8]

	m[0] = a11*b11 + a12*b21 + a13*b31
	m[3] = a11*b12 + a12*b22 + a13*b32
	m[6] = a11*b13 + a12*b23 + a13*b33

	m[1] = a21*b11 + a22*b21 + a23*b31
	m[4] = a21*b12 + a22*b22 + a23*b32
	m[7] = a21*b13 + a22*b23 + a23*b33

	m[2] = a31*b11 + a32*b21 + a33*b31
	m[5] = a31*b12 + a32*b22 + a33*b32
	m[8] = a31*b13 + a32*b23 + a33*b33

	return m
}

// MultiplyScalar multiplies each of this matrix's components by the specified scalar.
// Returns pointer to this updated matrix.
func (m *Mat3) MultScalar(s float32) *Mat3 {

	m[0] *= s
	m[3] *= s
	m[6] *= s
	m[1] *= s
	m[4] *= s
	m[7] *= s
	m[2] *= s
	m[5] *= s
	m[8] *= s
	return m
}

// Determinant calculates and returns the determinant of this matrix.
func (m *Mat3) Determinant() float32 {

	return m[0]*m[4]*m[8] -
		m[0]*m[5]*m[7] -
		m[1]*m[3]*m[8] +
		m[1]*m[5]*m[6] +
		m[2]*m[3]*m[7] -
		m[2]*m[4]*m[6]
}

// GetInverse sets this matrix to the inverse of the src matrix.
// If the src matrix cannot be inverted returns error and
// sets this matrix to the identity matrix.
func (m *Mat3) GetInverse(src *Mat3) error {

	n11 := src[0]
	n21 := src[1]
	n31 := src[2]
	n12 := src[3]
	n22 := src[4]
	n32 := src[5]
	n13 := src[6]
	n23 := src[7]
	n33 := src[8]

	t11 := n33*n22 - n32*n23
	t12 := n32*n13 - n33*n12
	t13 := n23*n12 - n22*n13

	det := n11*t11 + n21*t12 + n31*t13

	// no inverse
	if det == 0 {
		m.Identity()
		return errors.New("cannot invert matrix")
	}

	detInv := 1 / det

	m[0] = t11 * detInv
	m[1] = (n31*n23 - n33*n21) * detInv
	m[2] = (n32*n21 - n31*n22) * detInv
	m[3] = t12 * detInv
	m[4] = (n33*n11 - n31*n13) * detInv
	m[5] = (n31*n12 - n32*n11) * detInv
	m[6] = t13 * detInv
	m[7] = (n21*n13 - n23*n11) * detInv
	m[8] = (n22*n11 - n21*n12) * detInv

	return nil
}

// Transpose transposes this matrix.
// Returns pointer to this updated matrix.
func (m *Mat3) Transpose() *Mat3 {

	m[1], m[3] = m[3], m[1]
	m[2], m[6] = m[6], m[2]
	m[5], m[7] = m[7], m[5]
	return m
}

// SetTranslation sets this matrix to a translation matrix for the specified x and y values.
// Returns pointer to this updated matrix.
func (m *Mat3) SetTranslation(x, y float32) *Mat3 {

	m.Set(
		1, 0, x,
		0, 1, y,
		0, 0, 1,
	)
	return m
}

// SetTranslationVec sets this matrix to a translation matrix for the specified Vec2.
// Returns pointer to this updated matrix.
func (m *Mat3) SetTranslationVec(t Vec2) *Mat3 {

	return m.SetTranslation(t.X, t.Y)
}

// SetRotation set this matrix to a rotation matrix around the origin clockwise direction
// by the amount of theta radians.
func (m *Mat3) SetRotation(theta float32) *Mat3 {

	cosf := float32(math.Cos(float64(theta)))
	sinf := float32(math.Sin(float64(theta)))
	m.Set(
		cosf, -sinf, 0,
		sinf, cosf, 0,
		0, 0, 1,
	)
	return m
}

// SetScale set this matrix to a scale matrix for the specified x and y values
func (m *Mat3) SetScale(x, y float32) *Mat3 {

	m.Set(
		x, 0, 0,
		0, y, 0,
		0, 0, 1,
	)
	return m
}

// SetScaleVec set this matrix to a scale matrix for the specified Vec2
func (m *Mat3) SetScaleVec(s Vec2) *Mat3 {

	m.SetScale(s.X, s.Y)
	return m
}

// Rotate applies a rotation of theta radians around the origin
// in the clockwise direction to this matrix.
func (m *Mat3) Rotate(theta float32) *Mat3 {

	var sm Mat3
	sm.SetRotation(theta)
	m.Mult(&sm)
	return m
}

func (m *Mat3) Translate(x, y float32) *Mat3 {

	var tm Mat3
	tm.SetTranslation(x, y)
	m.Mult(&tm)
	return m
}

func (m *Mat3) TranslateVec(t Vec2) *Mat3 {

	m.Translate(t.X, t.Y)
	return m
}

func (m *Mat3) Scale(x, y float32) *Mat3 {

	var sm Mat3
	sm.SetScale(x, y)
	m.Mult(&sm)
	return m
}

func (m *Mat3) ScaleVec(s Vec2) *Mat3 {

	m.Scale(s.X, s.Y)
	return m
}
