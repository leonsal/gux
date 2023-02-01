package util

import (
	"math"
	"unicode"
)

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

func Assert(cond bool, msg string) {

	if !cond {
		panic("Gux Assertion Failed:" + msg)
	}
}

// AsciiSet returns slice of ASCII runes
func AsciiSet() []rune {

	runes := make([]rune, unicode.MaxASCII-32)
	for i := range runes {
		runes[i] = rune(32 + i)
	}
	return runes
}

// RangeTableSet returns a slice of runes from the specified unicode.RangeTable
func RangeTableSet(table *unicode.RangeTable) []rune {

	var runes []rune
	for _, rng := range table.R16 {
		for r := rng.Lo; r <= rng.Hi; r += rng.Stride {
			runes = append(runes, rune(r))
		}
	}
	for _, rng := range table.R32 {
		for r := rng.Lo; r <= rng.Hi; r += rng.Stride {
			runes = append(runes, rune(r))
		}
	}
	return runes
}
