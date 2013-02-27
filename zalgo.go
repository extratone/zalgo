// Copyright ©2013 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zalgo implements a zalgo text io.Writer.
package zalgo

import (
	"fmt"
	"io"
	"math/rand"
	"unicode/utf8"
)

var (
	up = []rune{
		0x030d, 0x030e, 0x0304, 0x0305,
		0x033f, 0x0311, 0x0306, 0x0310,
		0x0352, 0x0357, 0x0351, 0x0307,
		0x0308, 0x030a, 0x0342, 0x0343,
		0x0344, 0x034a, 0x034b, 0x034c,
		0x0303, 0x0302, 0x030c, 0x0350,
		0x0300, 0x0301, 0x030b, 0x030f,
		0x0312, 0x0313, 0x0314, 0x033d,
		0x0309, 0x0363, 0x0364, 0x0365,
		0x0366, 0x0367, 0x0368, 0x0369,
		0x036a, 0x036b, 0x036c, 0x036d,
		0x036e, 0x036f, 0x033e, 0x035b,
		0x0346, 0x031a,
	}

	down = []rune{
		0x0316, 0x0317, 0x0318, 0x0319,
		0x031c, 0x031d, 0x031e, 0x031f,
		0x0320, 0x0324, 0x0325, 0x0326,
		0x0329, 0x032a, 0x032b, 0x032c,
		0x032d, 0x032e, 0x032f, 0x0330,
		0x0331, 0x0332, 0x0333, 0x0339,
		0x033a, 0x033b, 0x033c, 0x0345,
		0x0347, 0x0348, 0x0349, 0x034d,
		0x034e, 0x0353, 0x0354, 0x0355,
		0x0356, 0x0359, 0x035a, 0x0323,
	}

	middle = []rune{
		0x0315, 0x031b, 0x0340, 0x0341,
		0x0358, 0x0321, 0x0322, 0x0327,
		0x0328, 0x0334, 0x0335, 0x0336,
		0x034f, 0x035c, 0x035d, 0x035e,
		0x035f, 0x0360, 0x0362, 0x0338,
		0x0337, 0x0361, 0x0489,
	}

	zalgoChars = func() map[rune]struct{} {
		zc := make(map[rune]struct{})
		for _, z := range up {
			zc[z] = struct{}{}
		}
		for _, z := range down {
			zc[z] = struct{}{}
		}
		for _, z := range middle {
			zc[z] = struct{}{}
		}
		return zc
	}
)

// Zalgo alters a Corrupter based in the number of bytes written by the Corrupter.
type Zalgo func(int, *Corrupter)

// Corrupter implements io.Writer transforming plain text to zalgo text.
type Corrupter struct {
	Up     complex128
	Middle complex128
	Down   complex128
	Zalgo
	w io.Writer
	n int
	b []byte
}

// NewCorrupter returns a new Corrupter that writes to w.
func NewCorrupter(w io.Writer) *Corrupter {
	return &Corrupter{w: w, b: make([]byte, utf8.MaxRune)}
}

// Write writes the byte slice p to the Corrupter's underlying io.Writer
// returning the number of bytes written and any error that occurred during
// the write operation.
func (z *Corrupter) Write(p []byte) (n int, err error) {
	for _, r := range string(p) {
		z.b = z.b[:utf8.RuneLen(r)]
		utf8.EncodeRune(z.b, r)
		n, err = z.w.Write(z.b)
		z.n += n
		if err != nil {
			return
		}
		if z.Zalgo != nil {
			z.Zalgo(z.n, z)
		}
		for i := real(z.Up); i > 0; i-- {
			if rand.Float64() < imag(z.Up) {
				n, err = fmt.Fprintf(z.w, "%c", up[rand.Intn(len(up))])
				z.n += n
				if err != nil {
					return
				}
			}
		}
		for i := real(z.Middle); i > 0; i-- {
			if rand.Float64() < imag(z.Middle) {
				n, err = fmt.Fprintf(z.w, "%c", middle[rand.Intn(len(middle))])
				z.n += n
				if err != nil {
					return
				}
			}
		}
		for i := real(z.Down); i > 0; i-- {
			if rand.Float64() < imag(z.Down) {
				n, err = fmt.Fprintf(z.w, "%c", down[rand.Intn(len(down))])
				z.n += n
				if err != nil {
					return
				}
			}
		}
	}
	return
}
