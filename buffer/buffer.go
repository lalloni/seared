// Copyright (C) 2017, Pablo Lalloni <plalloni@gmail.com>.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package buffer

import (
	"io"
	"sort"

	"github.com/lalloni/seared/location"
)

type Buffer interface {
	Length() int
	Input() string
	Rune(pos int) rune
	Runes(start, end int) []rune
	String(start, end int) string
	Reader(pos int) io.RuneReader
	Location(pos int) location.Location
}

type bufferReader struct {
	buf Buffer
	pos int
}

func (r *bufferReader) ReadRune() (rune, int, error) {
	if r.pos < r.buf.Length() {
		ru := r.buf.Rune(r.pos)
		r.pos++
		return ru, len(string(ru)), nil
	}
	return 0, 0, io.EOF
}

type buffer struct {
	input []rune
	nls   []int
}

func NewStringBuffer(input string) Buffer {
	return &buffer{input: []rune(input)}
}

func NewRunesBuffer(input []rune) Buffer {
	return &buffer{input: input}
}

func NewBytesBuffer(input []byte) Buffer {
	return &buffer{input: []rune(string(input))}
}

func (b *buffer) Rune(pos int) rune {
	if pos < b.Length() {
		return b.input[pos]
	}
	return 0
}

func (b *buffer) Runes(start, end int) []rune {
	l := b.Length()
	if start >= l {
		return nil
	}
	if end > l {
		end = l
	}
	return b.input[start:end]
}

func (b *buffer) String(start, end int) string {
	return string(b.Runes(start, end))
}

func (b *buffer) Input() string {
	return string(b.input)
}

func (b *buffer) Length() int {
	return len(b.input)
}

func (b *buffer) Reader(pos int) io.RuneReader {
	return &bufferReader{b, pos}
}

func (b *buffer) Location(pos int) location.Location {
	nls := b.newlines()
	l := sort.SearchInts(nls, pos)
	d := -1
	if l > 0 {
		d = nls[l-1]
	}
	return location.Location{
		Line:     l + 1,
		Column:   pos - d,
		Position: pos,
	}
}

func (b *buffer) newlines() []int {
	if b.nls == nil {
		n := []int{}
		for p, r := range b.input {
			if r == '\n' {
				n = append(n, p)
			}
		}
		b.nls = n
	}
	return b.nls
}
