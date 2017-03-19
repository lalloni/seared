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

package seared

import (
	"strings"

	"github.com/lalloni/seared/buffer"
	"github.com/lalloni/seared/node"
)

func (b *Builder) Empty() (this Expression) {
	this = newExpression("Empty", "EMPTY", b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			return Success(this, input, start, start)
		})
	return
}

func (b *Builder) End() (this Expression) {
	this = newExpression("End", "END", b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			if start >= input.Length() {
				return Success(this, input, start, start)
			}
			return Failure(this, input, start, start)
		})
	return
}

func (b *Builder) Rune(r rune) (this Expression) {
	e := "'" + string(r) + "'"
	this = newExpression("Rune", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			if input.Rune(start) == r {
				return Success(this, input, start, start+1).WithNodes(node.NewTerminal(string(r)))
			}
			return Failure(this, input, start, start)
		})
	return
}

func (b *Builder) Literal(literal string) (this Expression) {
	e := "'" + literal + "'"
	this = newExpression("Literal", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			end := start + len([]rune(literal))
			if input.String(start, end) == literal {
				return Success(this, input, start, end).WithNodes(node.NewTerminal(literal))
			}
			return Failure(this, input, start, start)
		})
	return
}

func (b *Builder) Range(first, last rune) (this Expression) {
	e := "[" + string(first) + "-" + string(last) + "]"
	this = newExpression("Range", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			r := input.Rune(start)
			if r >= first && r <= last {
				return Success(this, input, start, start+1).WithNodes(node.NewTerminal(string(r)))
			}
			return Failure(this, input, start, start)
		})
	return
}

func (b *Builder) Any() (this Expression) {
	this = newExpression("Any", ".", b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			if start < input.Length() {
				return Success(this, input, start, start+1).WithNodes(node.NewTerminal(string(input.Rune(start))))
			}
			return Failure(this, input, start, start)
		})
	return
}

func (b *Builder) AnyOf(runes string) (this Expression) {
	rs := strings.Split(runes, "")
	if len(rs) == 0 {
		panic("AnyOf rules does not allow the empty string")
	}
	e := "[" + runes + "]"
	this = newExpression("AnyOf", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			r := input.Rune(start)
			for _, rr := range runes {
				if r == rr {
					return Success(this, input, start, start+1).WithNodes(node.NewTerminal(string(r)))
				}
			}
			return Failure(this, input, start, start)
		})
	return
}
