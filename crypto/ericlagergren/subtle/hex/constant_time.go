// https://github.com/jedisct1/libsodium/blob/d4ee08ab8a1c674203796161af6d013283b33d69/src/libsodium/sodium/codecs.c
// https://github.com/jedisct1/libsodium/blob/561e556dad078af581f338fe3de9ee6362d28b16/LICENSE
//
//  Copyright (c) 2013-2022 Frank Denis <j at pureftpd dot org>
//  Portions Copyright (c) 2022 Eric Lagergren
//
// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package hex

import "crypto/subtle"

// Encode encodes src into EncodedLen(len(src)) bytes of dst.
// As a convenience, it returns the number of bytes written to
// dst, but this value is always EncodedLen(len(src)).
//
// Encode runs in constant time for the length of src.
func Encode(dst, src []byte) int {
	j := 0
	for _, v := range src {
		b := uint(v >> 4)
		c := uint(v & 0x0f)

		const (
			mask = ^uint(38)
		)
		dst[j+1] = byte(87 + c + (((c - 10) >> 8) & mask))
		dst[j] = byte(87 + b + (((b - 10) >> 8) & mask))
		j += 2
	}
	return len(src) * 2
}

// Decode decodes src into DecodedLen(len(src)) bytes, returning
// the actual number of bytes written to dst.
//
// Decode expects that src contains only hexadecimal characters
// and that src has even length. If the input is malformed,
// Decode returns the number of bytes decoded before the error.
//
// Decode runs in constant time for the length of src.
func Decode(dst, src []byte) (int, error) {
	// failed is set to 1 if the input is malformed, 0 otherwise.
	var failed int
	// badIdx is the number of bytes written to dst when
	// malformed data was found.
	//
	// Only has value if bad != 0.
	var badIdx int
	// badChar is the malformed character.
	var badChar int
	// acc is the accumulator between halves of a hexadecimal
	// character pair (04, e4, fe, ...).
	var acc byte
	// i is the index into dst.
	var i int

	for j := 0; j < len(src); j++ {
		c := uint(src[j])

		// Is c in '0' ... '9'?
		//
		// This is equivalent to
		//
		//    if n := c^'0'; n < 10 {
		//        val = n
		//    }
		//
		// which is true because
		//     y^(16*i) < 10 ∀ y ∈ [y, y+10)
		// and '0' == 48.
		//
		// If num < 10, subtracting 10 produces the two's
		// complement which flips the bits in [63:4] (which are
		// all zero because num < 10) to all one. Shifting by
		// 8 then ensures that bits [7:0] are all set to one,
		// resulting in 0xff.
		//
		// If num >= 10, subtracting 10 doesn't set any bits in
		// [63:8] (which are all zero because c < 256) and
		// shifting by 8 shifts off any set bits, resulting in
		// 0x00.
		num := c ^ '0'
		num0 := (num - 10) >> 8

		// Is c in 'a' ... 'f' or 'A' ... 'F'?
		//
		// This is equivalent to
		//
		//    const mask = ^uint(1<<5) // 0b11011111
		//    if a := c&mask; a >= 'A' && a < 'F' {
		//        val = a-55
		//    }
		//
		// The only difference between each uppercase and
		// lowercase ASCII pair ('a'-'A', 'e'-'E', etc.) is 32,
		// or bit #5. Masking that bit off folds the lowercase
		// letters into uppercase. The the range check should
		// then be obvious. Subtracting 55 converts the
		// hexadecimal character to binary by making 'A' = 10,
		// 'B' = 11, etc.
		//
		// If alpha is in [10, 15], subtracting 10 results in the
		// correct binary number, less 10. Notably, the bits in
		// [63:4] are all zero.
		//
		// If alpha is in [10, 15], subtracting 16 returns the
		// two's complement, flipping the bits in [63:4] (which
		// are all zero because alpha <= 15) to one.
		//
		// If alpha is in [10, 15], (alpha-10)^(alpha-16) sets
		// the bits in [63:4] to one. The bits in [3:0] are
		// irrelevant. Otherwise, if alpha == 9 or alpha >= 16,
		// both halves of the XOR have the same bits in [63:4],
		// so the XOR sets them to zero.
		//
		// Shifting by 8 clears the irrelevant bits in [3:0],
		// leaving only the interesting bits from the XOR. Thus,
		// if alpha is in [10, 15] bits [7:0] are all one,
		// resulting in 0xff. Otherwise, bits [7:0] are all zero,
		// resulting in 0x00.
		alpha := (c & ^uint(32)) - 55
		alpha0 := ((alpha - 10) ^ (alpha - 16)) >> 8

		// If both num0 and alpha0 are 0x00 then the character is
		// invalid.
		//
		// This is the constant-time equivalent of
		//
		//    if num0|alph0 == 0 {
		//        bad = 1
		//     } else {
		//        bad = 0
		//     }
		//
		bad := subtle.ConstantTimeByteEq(byte(num0|alpha0), 0)

		// If we haven't encountered an invalid character yet,
		// check whether the most recent character is invalid. If
		// so, record the invalid character and the number of
		// bytes we've written to dst.
		//
		// This is the constant-time equivalent of
		//
		//    if failed == 0 {
		//        if bad != 0 {
		//            badIdx = i
		//            badChar = c
		//        }
		//    }
		//
		// where both "if" cases have an implicit "else" that
		// sets
		//
		//    badIdx = badIdx
		//    badChar = badChar
		//
		badIdx = subtle.ConstantTimeSelect(failed, badIdx,
			subtle.ConstantTimeSelect(bad, i, badIdx))
		badChar = subtle.ConstantTimeSelect(failed, badChar,
			subtle.ConstantTimeSelect(bad, int(c), badChar))

		failed |= bad

		// Since num0 is either 0xff or 0x00, the bitwise AND
		// either leaves num unchanged or sets it to zero.
		// Similarly so for alpha0 and alpha.
		//
		// Only num or alpha can be non-zero here, so both OR and
		// XOR would work.
		val := byte(num0&num | alpha0&alpha)
		if j%2 == 0 {
			acc = val * 16
		} else {
			dst[i] = acc | val
			i++
		}
	}

	// Go checks for invalid length after checking for an invalid
	// character, so we do that too.
	if failed != 0 {
		return badIdx, InvalidByteError(badChar)
	}
	if len(src)%2 == 1 {
		return i, ErrLength
	}
	return i, nil
}

// validHexChar reports, in constant time, whether c is a valid
// hexadecimal character.
func validHexChar(c byte) bool {
	num := uint(c) ^ '0'
	num0 := (num - 10) >> 8
	alpha := (uint(c) & ^uint(32)) - 55
	alpha0 := ((alpha - 10) ^ (alpha - 16)) >> 8
	return subtle.ConstantTimeByteEq(byte(num0|alpha0), 0) == 0
}
