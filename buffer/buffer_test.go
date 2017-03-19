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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/seared/location"
)

func TestStringBuffer(t *testing.T) {
	b := StringBuffer("input")
	assert.NotNil(t, b)
	assert.EqualValues(t, "input", b.Input())
	assert.EqualValues(t, 'p', b.Rune(2))
	assert.EqualValues(t, "np", b.String(1, 3))
	assert.EqualValues(t, "nput", b.String(1, 10))
	assert.EqualValues(t, "", b.String(10, 10))
}

func TestStringBufferUnicode(t *testing.T) {
	b := StringBuffer("aaŧ←↓ŋħ5ł")
	assert.NotNil(t, b)
	assert.EqualValues(t, "aaŧ←↓ŋħ5ł", b.Input())
	assert.EqualValues(t, "←↓", b.String(3, 5))
	assert.EqualValues(t, 'ħ', b.Rune(6))
}

func TestStringBufferLocation(t *testing.T) {
	s := "lots\nof text\nin multiple lines\nand more,\nmore, much more"
	b := StringBuffer(s)
	assert.EqualValues(t, location.New(1, 1, 0), b.Location(0))
	assert.EqualValues(t, location.New(1, 2, 1), b.Location(1))
	assert.EqualValues(t, location.New(2, 6, 10), b.Location(10))
	assert.EqualValues(t, location.New(3, 8, 20), b.Location(20))
	assert.EqualValues(t, location.New(2, 7, 11), b.Location(11))
	assert.EqualValues(t, location.New(2, 8, 12), b.Location(12))
	assert.EqualValues(t, location.New(3, 1, 13), b.Location(13))
}

func TestStringBufferLine(t *testing.T) {
	a := assert.New(t)
	s := "lots\nof text\nin multiple lines\nand more,\nmore, much more"
	b := StringBuffer(s)
	a.Equal("of text", b.Line(2))
	a.Equal("lots", b.Line(1))
	a.Equal("and more,", b.Line(4))
}

func TestStringBufferRuneReader(t *testing.T) {
	a := assert.New(t)
	s := "lħs\nof¶↓n€ß"
	b := StringBuffer(s)
	r := b.Reader(0)
	c := 0
	for _, ru := range []rune(s) {
		rr, _, err := r.ReadRune()
		a.Nil(err)
		a.EqualValues(ru, rr)
		c++
	}
	_, _, err := r.ReadRune()
	a.EqualValues(io.EOF, err)
	a.EqualValues(len([]rune(s)), c)
}

func TestStringBufferRunes(t *testing.T) {
	a := assert.New(t)
	b := StringBuffer("fsjdsljdsfl")
	a.EqualValues([]rune("f"), b.Runes(0, 1))
	a.EqualValues([]rune("fs"), b.Runes(0, 2))
	a.EqualValues([]rune("ds"), b.Runes(3, 5))
	a.EqualValues([]rune("fl"), b.Runes(9, 11))
	a.EqualValues([]rune("fl"), b.Runes(9, 12))
	a.EqualValues([]rune("fl"), b.Runes(9, 15))
	a.EqualValues([]rune("l"), b.Runes(10, 15))
	a.EqualValues([]rune("l"), b.Runes(10, 11))
	a.EqualValues([]rune(""), b.Runes(10, 10))
	a.Nil(b.Runes(11, 15))
}
