package subtle

import "crypto/subtle"

// ConstantTimeByteEq returns 1 if x == y and 0 otherwise.
func ConstantTimeByteEq(x, y uint8) int {
	return subtle.ConstantTimeByteEq(x, y)
}

// ConstantTimeCompare returns 1 if the two slices, x and y, have
// equal contents and 0 otherwise.
//
// The time taken is a function of the length of the slices and
// is independent of the contents.
func ConstantTimeCompare(x, y []byte) int {
	return subtle.ConstantTimeCompare(x, y)
}

// ConstantTimeCopy copies the contents of y into x (a slice of
// equal length) if v == 1. If v == 0, x is left unchanged. Its
// behavior is undefined if v takes any other value.
func ConstantTimeCopy(v int, x, y []byte) {
	subtle.ConstantTimeCopy(v, x, y)
}

// ConstantTimeEq returns 1 if x == y and 0 otherwise.
func ConstantTimeEq(x, y int32) int {
	return subtle.ConstantTimeEq(x, y)
}

// ConstantTimeLessOrEq returns 1 if x <= y and 0 otherwise.
// Its behavior is undefined if x or y are negative or > 2**31 - 1.
func ConstantTimeLessOrEq(x, y int) int {
	return subtle.ConstantTimeLessOrEq(x, y)
}

// ConstantTimeSelect returns x if v == 1 and y if v == 0.
// Its behavior is undefined if v takes any other value.
func ConstantTimeSelect(v, x, y int) int {
	return subtle.ConstantTimeSelect(v, x, y)
}

// ConstantTimeBigEndianZero reports, in constant time, whether
// the big-endian integer x is zero.
//
// It returns 1 if x <= y and 0 otherwise.
func ConstantTimeBigEndianZero(x []byte) int {
	var v byte
	for i := 0; i < len(x); i++ {
		v |= x[i]
	}
	return ConstantTimeByteEq(v, 0)
}

// ConstantTimeBigEndianLessOrEq compares x and y, which must
// have the same length, as big-endian integers in constant time.
//
// It returns 1 if x <= y and 0 otherwise.
func ConstantTimeBigEndianLessOrEq(x, y []byte) int {
	if len(x) != len(y) {
		panic("subtle: slices have different lengths")
	}
	var neq int
	var gt int
	for i := 0; i < len(x); i++ {
		// if neq == 0 {
		//     gt = ConstantTimeByteGreater(x[i], y[i])
		// }
		gt |= ConstantTimeSelect(neq, 0,
			ConstantTimeByteGreater(x[i], y[i]))
		// if gt == 0 {
		//     neq = ConstantTimeNeq(x[i], y[i])
		// }
		neq |= ConstantTimeSelect(gt, 0,
			ConstantTimeByteEq(x[i], y[i])^1)
	}
	return gt ^ 1
}

// ConstantTimeByteGreater returns 1 if x > y and 0 otherwise.
func ConstantTimeByteGreater(x, y uint8) int {
	return ConstantTimeByteLessOrEq(x, y) ^ 1
}

// ConstantTimeByteLessOrEq returns 1 if x <= y and 0 otherwise.
func ConstantTimeByteLessOrEq(x, y uint8) int {
	return ConstantTimeLessOrEq(int(x), int(y))
}
