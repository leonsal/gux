package gb

import "math"

// Vec2Add sums vectors v1 with v2 returning the resulting vector.
// Vectors v1 and v2 are not changed.
func Vec2Add(v1, v2 Vec2) Vec2 {

	return Vec2{v1.X + v2.X, v1.Y + v2.Y}
}

// Vec2Sub subtracts vector v2 from vector v1 returning the resulting vector.
// Vectors v1 and v2 are not changed.
func Vec2Sub(v1, v2 Vec2) Vec2 {

	return Vec2{v1.X - v2.X, v1.Y - v2.Y}
}

// Vec2Mult multiplies each component of v1 with the corresponding component of v2 returning a new vector.
// Vectors v1 and v2 are not changed.
func Vec2Mult(v1, v2 Vec2) Vec2 {

	return Vec2{v1.X * v2.X, v1.Y * v2.Y}

}

// Vec2MultScalar multiplies each component of vector v with the scalar s returning the resulting vector.
// Vector v is not changed.
func Vec2MultScalar(v Vec2, s float32) Vec2 {

	return Vec2{v.X * s, v.Y * s}
}

// Add adds other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vec2) Add(other Vec2) *Vec2 {

	v.X += other.X
	v.Y += other.Y
	return v
}

// AddScalar adds scalar s to each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vec2) AddScalar(s float32) *Vec2 {

	v.X += s
	v.Y += s
	return v
}

// Sub subtracts other vector from this one.
// Returns the pointer to this updated vector.
func (v *Vec2) Sub(other Vec2) *Vec2 {

	v.X -= other.X
	v.Y -= other.Y
	return v
}

// SubScalar subtracts scalar s from each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vec2) SubScalar(s float32) *Vec2 {

	v.X -= s
	v.Y -= s
	return v
}

// MultScalar multiplies each component of this vector by the scalar s.
// Returns the pointer to this updated vector.
func (v *Vec2) MultScalar(s float32) *Vec2 {

	v.X *= s
	v.Y *= s
	return v
}

// DivScalar divides each component of this vector by the scalar s.
// If scalar is zero, sets this vector to zero.
// Returns the pointer to this updated vector.
func (v *Vec2) DivScalar(scalar float32) *Vec2 {

	if scalar != 0 {
		invScalar := 1 / scalar
		v.X *= invScalar
		v.Y *= invScalar
	} else {
		v.X = 0
		v.Y = 0
	}
	return v
}

// Negate negates each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vec2) Negate() *Vec2 {

	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Dot returns the dot product of this vector with other.
// None of the vectors are changed.
func (v *Vec2) Dot(other Vec2) float32 {

	return v.X*other.X + v.Y*other.Y
}

// LengthSq returns the length squared of this vector.
// LengthSq can be used to compare vectors' lengths without the need to perform a square root.
func (v *Vec2) LengthSq() float32 {

	return v.X*v.X + v.Y*v.Y
}

// Length returns the length of this vector.
func (v *Vec2) Length() float32 {

	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

// Normalize normalizes this vector so its length will be 1.
// Returns the pointer to this updated vector.
func (v *Vec2) Normalize() *Vec2 {

	return v.DivScalar(v.Length())
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vec2) Min(other Vec2) *Vec2 {

	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vec2) Max(other Vec2) *Vec2 {

	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	return v
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding components of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
// Returns the pointer to this updated vector.
func (v *Vec2) Clamp(min, max Vec2) *Vec2 {

	if v.X < min.X {
		v.X = min.X
	} else if v.X > max.X {
		v.X = max.X
	}

	if v.Y < min.Y {
		v.Y = min.Y
	} else if v.Y > max.Y {
		v.Y = max.Y
	}
	return v
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
// Returns the pointer to this updated vector.
func (v *Vec2) ClampScalar(minVal, maxVal float32) *Vec2 {

	if v.X < minVal {
		v.X = minVal
	} else if v.X > maxVal {
		v.X = maxVal
	}

	if v.Y < minVal {
		v.Y = minVal
	} else if v.Y > maxVal {
		v.Y = maxVal
	}
	return v
}
