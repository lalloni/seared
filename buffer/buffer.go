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
