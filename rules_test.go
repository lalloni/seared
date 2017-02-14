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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/seared/buffer"
)

func rule(matcher Matcher) Rule {
	return &testRule{matcher}
}

type testRule struct {
	apply Matcher
}

func (r *testRule) Apply(input buffer.Buffer, pos int) (success bool, nextpos int) {
	return r.apply(input, pos)
}

func successful(consume int) Rule {
	return rule(func(input buffer.Buffer, pos int) (success bool, nextpos int) {
		return true, pos + consume
	})
}

func failed() Rule {
	return rule(func(input buffer.Buffer, pos int) (success bool, nextpos int) {
		return false, pos
	})
}

func successfulTimes(consume int, times int) Rule {
	return rule(func(input buffer.Buffer, pos int) (success bool, nextpos int) {
		times--
		success = times > -1
		nextpos = pos
		if success {
			nextpos++
		}
		return
	})
}

func Foo(r *Rules) Rule {
	return r.Rule(func() Rule {
		return r.Choice(r.Rune('f'), Bar(r))
	})
}

func Bar(r *Rules) Rule {
	return r.Rule(func() Rule {
		return r.Choice(Baz(r), r.Rune('b'), Foo(r))
	})
}

func Baz(r *Rules) Rule {
	return r.Rule(func() Rule {
		return r.Rune('z')
	})
}

func TestRecurse(t *testing.T) {
	a := assert.New(t)
	p := NewParser(Foo)
	a.NotNil(p)
	p.SetDebug(true)
	p.SetLog(TestingLog(t))
	a.True(p.Recognize("b"))
}

func TestRange(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("blanco")
	b := rules(nil)

	var (
		s bool
		p int
	)

	s, p = b.Range('a', 'c').Apply(i, 0)
	a.True(s)
	a.Equal(1, p)

	s, p = b.Range('x', 'z').Apply(i, 2)
	a.False(s)
	a.Equal(2, p)
}

func TestAny(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("blanco")
	b := rules(nil)

	var (
		s bool
		p int
	)

	s, p = b.Any("abc").Apply(i, 0)
	a.True(s)
	a.Equal(1, p)

	s, p = b.Any("123").Apply(i, 2)
	a.False(s)
	a.Equal(2, p)
}

func TestRune(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("blanco")
	b := rules(nil)

	var (
		s bool
		p int
	)

	s, p = b.Rune('a').Apply(i, 0)
	a.False(s)
	a.Equal(0, p)

	s, p = b.Rune('a').Apply(i, 2)
	a.True(s)
	a.Equal(3, p)
}

func TestLiteral(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("blanco")
	b := rules(nil)

	var (
		s bool
		p int
	)

	s, p = b.Literal("bla").Apply(i, 0)
	a.True(s)
	a.EqualValues(3, p)

	s, p = b.Literal("nop").Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = b.Literal("anco").Apply(i, 2)
	a.True(s)
	a.EqualValues(6, p)

	s, p = b.Literal("a").Apply(i, 6)
	a.False(s)
	a.EqualValues(6, p)
}

func TestSequence(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.Sequence(failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Sequence(successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Sequence(successful(1), failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Sequence(successful(1), successful(1), successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(3, p)

	s, p = r.Sequence(successful(1), successful(1), successful(1)).Apply(i, 5)
	a.True(s)
	a.EqualValues(8, p)

	s, p = r.Sequence(successful(1), failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Sequence(successful(1), failed(), successful(1)).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Sequence(failed(), successful(1)).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)
}

func TestChoice(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.Choice(successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Choice(successful(1), successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Choice(successful(1), failed()).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Choice(failed(), successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Choice(failed(), failed(), successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Choice(failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Choice(failed(), failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)
}

func TestZeroOrMore(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("_abababababzzz")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.ZeroOrMore(successfulTimes(1, 3)).Apply(i, 0)
	a.True(s)
	a.EqualValues(3, p)

	s, p = r.ZeroOrMore(successfulTimes(1, 1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.ZeroOrMore(successfulTimes(1, 1)).Apply(i, 5)
	a.True(s)
	a.EqualValues(6, p)

	s, p = r.ZeroOrMore(failed()).Apply(i, 5)
	a.True(s)
	a.EqualValues(5, p)

	s, p = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("ab"))).Apply(i, 0)
	a.True(s)
	a.EqualValues(11, p)

	s, p = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("x"))).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("ab"))).Apply(i, 1)
	a.False(s)
	a.EqualValues(1, p)

	s, p = r.ZeroOrMore(r.Literal("ab")).Apply(i, 1)
	a.True(s)
	a.EqualValues(11, p)
}

func TestOneOrMore(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("_abababababzzz")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.OneOrMore(successfulTimes(1, 3)).Apply(i, 0)
	a.True(s)
	a.EqualValues(3, p)

	s, p = r.OneOrMore(successfulTimes(1, 1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.OneOrMore(successfulTimes(1, 1)).Apply(i, 5)
	a.True(s)
	a.EqualValues(6, p)

	s, p = r.OneOrMore(failed()).Apply(i, 5)
	a.False(s)
	a.EqualValues(5, p)

	s, p = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("ab"))).Apply(i, 0)
	a.True(s)
	a.EqualValues(11, p)

	s, p = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("x"))).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("ab"))).Apply(i, 1)
	a.False(s)
	a.EqualValues(1, p)

	s, p = r.OneOrMore(r.Literal("ab")).Apply(i, 1)
	a.True(s)
	a.EqualValues(11, p)

	s, p = r.OneOrMore(r.Literal("xx")).Apply(i, 5)
	a.False(s)
	a.EqualValues(5, p)

	s, p = r.OneOrMore(r.Literal("z")).Apply(i, 11)
	a.True(s)
	a.EqualValues(14, p)

	s, p = r.OneOrMore(r.Literal("x")).Apply(i, 11)
	a.False(s)
	a.EqualValues(11, p)

	s, p = r.OneOrMore(r.Literal("_")).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)
}

func TestOptional(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("aaa")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.Optional(successful(1)).Apply(i, 0)
	a.True(s)
	a.EqualValues(1, p)

	s, p = r.Optional(failed()).Apply(i, 0)
	a.True(s)
	a.EqualValues(0, p)
}

func TestAnd(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.And(successful(10)).Apply(i, 0)
	a.True(s)
	a.EqualValues(0, p)

	s, p = r.And(failed()).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)
}

func TestNot(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.Not(successful(10)).Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.Not(failed()).Apply(i, 0)
	a.True(s)
	a.EqualValues(0, p)
}

func TestEnd(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("123")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.End().Apply(i, 0)
	a.False(s)
	a.EqualValues(0, p)

	s, p = r.End().Apply(i, 1)
	a.False(s)
	a.EqualValues(1, p)

	s, p = r.End().Apply(i, 2)
	a.False(s)
	a.EqualValues(2, p)

	s, p = r.End().Apply(i, 3)
	a.True(s)
	a.EqualValues(3, p)
}

func TestEmpty(t *testing.T) {
	a := assert.New(t)
	i := buffer.NewStringBuffer("123")
	r := rules(nil)

	var (
		s bool
		p int
	)

	s, p = r.Empty().Apply(i, 0)
	a.True(s)
	a.EqualValues(0, p)

	s, p = r.Empty().Apply(i, 2)
	a.True(s)
	a.EqualValues(2, p)
}
