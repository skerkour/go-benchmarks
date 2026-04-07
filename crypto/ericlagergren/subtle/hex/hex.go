// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hex

import (
	"encoding/hex"
	"io"
)

var ErrLength = hex.ErrLength

type InvalidByteError = hex.InvalidByteError

// bufferSize is the number of hexadecimal characters to buffer
// in encoder and decoder.
//
// It's taken from encoding/hex and seemingly completely
// arbitrary.
const bufferSize = 1024

// EncodedLen returns the length of an encoding of n source
// bytes.
// Specifically, it returns n * 2.
func EncodedLen(n int) int {
	return hex.EncodedLen(n)
}

// EncodeToString returns the hexadecimal encoding of src.
//
// Encode runs in constant time for the length of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// NewEncoder returns an io.Writer that writes lowercase
// hexadecimal characters to w.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / 2
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		var written int
		encoded := Encode(e.out[:], p[:chunkSize])
		written, e.err = e.w.Write(e.out[:encoded])
		n += written / 2
		p = p[chunkSize:]
	}
	return n, e.err
}

func DecodedLen(n int) int {
	return hex.DecodedLen(n)
}

// DecodeString returns the bytes represented by the hexadecimal
// string s.
//
// DecodeString expects that src contains only hexadecimal
// characters and that src has even length. If the input is
// malformed, DecodeString returns the bytes decoded before the
// error.
//
// DecodeString runs in constant time for the length of s.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	n, err := Decode(src, src)
	return src[:n], err
}

// NewDecoder returns an io.Reader that decodes hexadecimal
// characters from r.
//
// NewDecoder expects that r contain only an even number of
// hexadecimal characters.
//
// The first call to Read that encounters malformed hexadecimal
// characters will return a non-nil error. This means that the
// io.Reader does not operate in constant time over the entire
// stream, but rather for each chunk read from r.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

var _ io.Reader = (*decoder)(nil)

func (d *decoder) Read(p []byte) (n int, err error) {
	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < 2 && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in) // Copies either 0 or 1 bytes
		numRead, d.err = d.r.Read(d.arr[numCopy:])
		d.in = d.arr[:numCopy+numRead]
		if d.err == io.EOF && len(d.in)%2 != 0 {
			if !validHexChar(d.in[len(d.in)-1]) {
				d.err = InvalidByteError(d.in[len(d.in)-1])
			} else {
				d.err = io.ErrUnexpectedEOF
			}
		}
	}

	// Decode internal buffer into output buffer
	if numAvail := len(d.in) / 2; len(p) > numAvail {
		p = p[:numAvail]
	}
	numDec, err := Decode(p, d.in[:len(p)*2])
	d.in = d.in[2*numDec:]
	if err != nil {
		d.in, d.err = nil, err // Decode error; discard input remainder
	}

	if len(d.in) < 2 {
		return numDec, d.err // Only expose errors when buffer fully consumed
	}
	return numDec, nil
}
