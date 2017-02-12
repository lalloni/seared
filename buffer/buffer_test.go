package buffer

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/seared/location"
)

func TestStringBuffer(t *testing.T) {
	b := NewStringBuffer("input")
	assert.NotNil(t, b)
	assert.EqualValues(t, "input", b.Input())
	assert.EqualValues(t, 'p', b.Rune(2))
	assert.EqualValues(t, "np", b.String(1, 3))
	assert.EqualValues(t, "nput", b.String(1, 10))
	assert.EqualValues(t, "", b.String(10, 10))
}

func TestStringBufferUnicode(t *testing.T) {
	b := NewStringBuffer("aaŧ←↓ŋħ5ł")
	assert.NotNil(t, b)
	assert.EqualValues(t, "aaŧ←↓ŋħ5ł", b.Input())
	assert.EqualValues(t, "←↓", b.String(3, 5))
	assert.EqualValues(t, 'ħ', b.Rune(6))
}

func TestStringBufferLocation(t *testing.T) {
	s := "lots\nof text\nin multiple lines\nand more,\nmore, much more"
	b := NewStringBuffer(s)
	assert.EqualValues(t, location.NewLocation(1, 1, 0), b.Location(0))
	assert.EqualValues(t, location.NewLocation(1, 2, 1), b.Location(1))
	assert.EqualValues(t, location.NewLocation(2, 6, 10), b.Location(10))
	assert.EqualValues(t, location.NewLocation(3, 8, 20), b.Location(20))
	assert.EqualValues(t, location.NewLocation(2, 7, 11), b.Location(11))
	assert.EqualValues(t, location.NewLocation(2, 8, 12), b.Location(12))
	assert.EqualValues(t, location.NewLocation(3, 1, 13), b.Location(13))
}

func TestStringBufferRuneReader(t *testing.T) {
	a := assert.New(t)
	s := "lħs\nof¶↓n€ß"
	b := NewStringBuffer(s)
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
	b := NewStringBuffer("fsjdsljdsfl")
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
