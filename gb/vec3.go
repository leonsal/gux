package gb

import "math"

// Set sets this vector X, Y and Z components.
// Returns the pointer to this updated vector.
func (v *Vec3) Set(x, y, z float32) *Vec3 {

	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// SetX sets this vector X component.
// Returns the pointer to this updated Vector.
func (v *Vec3) SetX(x float32) *Vec3 {

	v.X = x
	return v
}

// SetY sets this vector Y component.
// Returns the pointer to this updated vector.
func (v *Vec3) SetY(y float32) *Vec3 {

	v.Y = y
	return v
}

// SetZ sets this vector Z component.
// Returns the pointer to this updated vector.
func (v *Vec3) SetZ(z float32) *Vec3 {

	v.Z = z
	return v
}

// Add adds other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vec3) Add(other Vec3) *Vec3 {

	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// AddScalar adds scalar s to each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vec3) AddScalar(s float32) *Vec3 {

	v.X += s
	v.Y += s
	v.Z += s
	return v
}

// AddVec adds vectors a and b to this one.
// Returns the pointer to this updated vector.
func (v *Vec3) AddVecs(a, b *Vec3) *Vec3 {

	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	v.Z = a.Z + b.Z
	return v
}

// Sub subtracts other vector from this one.
// Returns the pointer to this updated vector.
func (v *Vec3) Sub(other *Vec3) *Vec3 {

	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// SubScalar subtracts scalar s from each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vec3) SubScalar(s float32) *Vec3 {

	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

// SubVectors sets this vector to a - b.
// Returns the pointer to this updated vector.
func (v *Vec3) SubVectors(a, b *Vec3) *Vec3 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}

// Multiply multiplies each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector.
func (v *Vec3) Multiply(other *Vec3) *Vec3 {

	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

// MultiplyScalar multiplies each component of this vector by the scalar s.
// Returns the pointer to this updated vector.
func (v *Vec3) MultiplyScalar(s float32) *Vec3 {

	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// Divide divides each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector
func (v *Vec3) Divide(other *Vec3) *Vec3 {

	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

// DivideScalar divides each component of this vector by the scalar s.
// If scalar is zero, sets this vector to zero.
// Returns the pointer to this updated vector.
func (v *Vec3) DivideScalar(scalar float32) *Vec3 {

	if scalar != 0 {
		invScalar := 1 / scalar
		v.X *= invScalar
		v.Y *= invScalar
		v.Z *= invScalar
	} else {
		v.X = 0
		v.Y = 0
		v.Z = 0
	}
	return v
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vec3) Min(other *Vec3) *Vec3 {

	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	if v.Z > other.Z {
		v.Z = other.Z
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vec3) Max(other *Vec3) *Vec3 {

	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	if v.Z < other.Z {
		v.Z = other.Z
	}
	return v
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding component of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
// Returns the pointer to this updated vector.
func (v *Vec3) Clamp(min, max *Vec3) *Vec3 {

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

	if v.Z < min.Z {
		v.Z = min.Z
	} else if v.Z > max.Z {
		v.Z = max.Z
	}
	return v
}

// Negate negates each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vec3) Negate() *Vec3 {

	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Dot returns the dot product of this vector with other.
// None of the vectors are changed.
func (v *Vec3) Dot(other *Vec3) float32 {

	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// LengthSq returns the length squared of this vector.
// LengthSq can be used to compare vectors' lengths without the need to perform a square root.
func (v *Vec3) LengthSq() float32 {

	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length returns the length of this vector.
func (v *Vec3) Length() float32 {

	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize normalizes this vector so its length will be 1.
// Returns the pointer to this updated vector.
func (v *Vec3) Normalize() *Vec3 {

	return v.DivideScalar(v.Length())
}

// DistanceTo returns the distance of this point to other.
func (v *Vec3) DistanceTo(other *Vec3) float32 {

	return float32(math.Sqrt(float64(v.DistanceToSquared(other))))
}

// DistanceToSquared returns the distance squared of this point to other.
func (v *Vec3) DistanceToSquared(other *Vec3) float32 {

	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// SetLength sets this vector to have the specified length.
// If the current length is zero, does nothing.
// Returns the pointer to this updated vector.
func (v *Vec3) SetLength(l float32) *Vec3 {

	oldLength := v.Length()
	if oldLength != 0 && l != oldLength {
		v.MultiplyScalar(l / oldLength)
	}
	return v
}

// Lerp sets each of this vector's components to the linear interpolated value of
// alpha between itself and the corresponding other component.
// Returns the pointer to this updated vector.
func (v *Vec3) Lerp(other *Vec3, alpha float32) *Vec3 {

	v.X += (other.X - v.X) * alpha
	v.Y += (other.Y - v.Y) * alpha
	v.Z += (other.Z - v.Z) * alpha
	return v
}

// Equals returns if this vector is equal to other.
func (v *Vec3) Equals(other *Vec3) bool {

	return (other.X == v.X) && (other.Y == v.Y) && (other.Z == v.Z)
}

// FromArray sets this vector's components from the specified array and offset
// Returns the pointer to this updated vector.
func (v *Vec3) FromArray(array []float32, offset int) *Vec3 {

	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	return v
}

// ToArray copies this vector's components to array starting at offset.
// Returns the array.
func (v *Vec3) ToArray(array []float32, offset int) []float32 {

	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	return array
}

// MultiplyVectors multiply vectors a and b storing the result in this vector.
// Returns the pointer to this updated vector.
func (v *Vec3) MultiplyVectors(a, b *Vec3) *Vec3 {

	v.X = a.X * b.X
	v.Y = a.Y * b.Y
	v.Z = a.Z * b.Z
	return v
}

// ApplyMat3 multiplies the specified 3x3 matrix by this vector.
// Returns the pointer to this updated vector.
func (v *Vec3) ApplyMat3(m *Mat3) *Vec3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[3]*y + m[6]*z
	v.Y = m[1]*x + m[4]*y + m[7]*z
	v.Z = m[2]*x + m[5]*y + m[8]*z
	return v
}

// Cross calculates the cross product of this vector with other and returns the result vector.
func (v *Vec3) Cross(other *Vec3) *Vec3 {

	cx := v.Y*other.Z - v.Z*other.Y
	cy := v.Z*other.X - v.X*other.Z
	cz := v.X*other.Y - v.Y*other.X
	v.X = cx
	v.Y = cy
	v.Z = cz
	return v
}

// CrossVectors calculates the cross product of a and b storing the result in this vector.
// Returns the pointer to this updated vector.
func (v *Vec3) CrossVectors(a, b *Vec3) *Vec3 {

	cx := a.Y*b.Z - a.Z*b.Y
	cy := a.Z*b.X - a.X*b.Z
	cz := a.X*b.Y - a.Y*b.X
	v.X = cx
	v.Y = cy
	v.Z = cz
	return v
}
