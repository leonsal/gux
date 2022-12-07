package gux

import "math"

type Number interface {
	~int | ~float32 | ~float64
}

func Min[T Number](lhs, rhs T) T {

	if lhs < rhs {
		return lhs
	}
	return rhs
}

func Max[T Number](lhs, rhs T) T {

	if lhs >= rhs {
		return lhs
	}
	return rhs
}

func Clamp[T Number](v, min, max T) T {

	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func Lerp[T Number](a, b T, t float32) T {

	return (a + (b-a)*T(t))
}

type Float interface {
	~float32 | ~float64
}

func Cos[T Float](a T) T {

	return T(math.Cos(float64(a)))
}

func Sin[T Float](a T) T {

	return T(math.Sin(float64(a)))
}
