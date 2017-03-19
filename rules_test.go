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

func ruleM(matcher Matcher) Expression {
	return &testRule{apply: matcher}
}

func ruleEM(expectation string, matcher Matcher) Expression {
	return &testRule{apply: matcher, expectation: expectation}
}

func ruleE(expectation string) Expression {
	return &testRule{expectation: expectation}
}

func ruleNE(name, expectation string) Expression {
	return &testRule{name: name, expectation: expectation}
}

type testRule struct {
	apply       Matcher
	name        string
	expectation string
}

func (r *testRule) Apply(input buffer.Buffer, pos int) (result *Result) {
	return r.apply(input, pos)
}

func (r *testRule) Name() string {
	return r.name
}

func (r *testRule) Expectation() string {
	return r.expectation
}

func successful(consume int) (this Expression) {
	this = ruleEM("successful", func(input buffer.Buffer, pos int) (result *Result) {
		return Success(this, input, pos, pos+consume)
	})
	return
}

func failed() (this Expression) {
	this = ruleEM("failed", func(input buffer.Buffer, pos int) (result *Result) {
		return Failure(this, input, pos, pos)
	})
	return
}

func successfulTimes(consume int, times int) (this Expression) {
	this = ruleM(func(input buffer.Buffer, pos int) (result *Result) {
		times--
		if times > -1 {
			return Success(this, input, pos, pos+consume)
		}
		return Failure(this, input, pos, pos)
	})
	return
}

func Foo(r *Builder) Expression {
	return r.Rule(func() Expression {
		return r.Choice(r.Rune('f'), Bar(r))
	})
}

func Bar(r *Builder) Expression {
	return r.Rule(func() Expression {
		return r.Choice(Baz(r), r.Rune('b'), Foo(r))
	})
}

func Baz(r *Builder) Expression {
	return r.Rule(func() Expression {
		return r.Rune('z')
	})
}

func TestRecurse(t *testing.T) {
	a := assert.New(t)
	p := NewParser(Foo)
	a.NotNil(p)
	p.SetDebug(true)
	p.SetLog(TestingLog(t))
	result := p.ParseString("b")
	a.True(result.Success)
}

func TestRange(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("blanco")
	b := newBuilder(nil)

	var result *Result

	result = b.Range('a', 'c').Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(0, len(result.Results))

	result = b.Range('x', 'z').Apply(i, 2)
	a.False(result.Success)
	a.Equal(2, result.End)
	a.Equal("[x-z]", result.Expression.Expectation())
	a.Equal(0, len(result.Results))
}

func TestAny(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("blanco")
	b := newBuilder(nil)

	var result *Result

	result = b.AnyOf("abc").Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(0, len(result.Results))

	result = b.AnyOf("123").Apply(i, 2)
	a.False(result.Success)
	a.Equal(2, result.End)
	a.Equal("[123]", result.Expression.Expectation())
	a.Equal(0, len(result.Results))
}

func TestRune(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("blanco")
	b := newBuilder(nil)

	var result *Result

	result = b.Rune('a').Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal(0, len(result.Results))

	result = b.Rune('a').Apply(i, 2)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(0, len(result.Results))
}

func TestLiteral(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("blanco")
	b := newBuilder(nil)

	var result *Result

	result = b.Literal("bla").Apply(i, 0)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(0, len(result.Results))

	result = b.Literal("nop").Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("'nop'", result.Expression.Expectation())
	a.Equal(0, len(result.Results))

	result = b.Literal("anco").Apply(i, 2)
	a.True(result.Success)
	a.Equal(6, result.End)
	a.Equal(0, len(result.Results))

	result = b.Literal("a").Apply(i, 6)
	a.False(result.Success)
	a.Equal(6, result.End)
	a.Equal("'a'", result.Expression.Expectation())
	a.Equal(0, len(result.Results))
}

func TestSequence(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("")
	r := newBuilder(nil)

	var result *Result

	result = r.Sequence(failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("failed", result.Expression.Expectation())
	a.Equal(1, len(result.Results))

	result = r.Sequence(successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.Sequence(successful(1), failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("successful failed", result.Expression.Expectation())
	a.Equal(2, len(result.Results))

	result = r.Sequence(successful(1), successful(1), successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(3, len(result.Results))

	result = r.Sequence(successful(1), successful(1), successful(1)).Apply(i, 5)
	a.True(result.Success)
	a.Equal(8, result.End)
	a.Equal(3, len(result.Results))

	result = r.Sequence(successful(1), failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("successful failed", result.Expression.Expectation())
	a.Equal(2, len(result.Results))

	result = r.Sequence(successful(1), failed(), successful(1)).Apply(i, 0)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("successful failed successful", result.Expression.Expectation())
	a.Equal(2, len(result.Results))

	result = r.Sequence(failed(), successful(1)).Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("failed successful", result.Expression.Expectation())
	a.Equal(1, len(result.Results))

}

func TestChoice(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("")
	r := newBuilder(nil)

	var result *Result

	result = r.Choice(successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.Choice(successful(1), successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.Choice(successful(1), failed()).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.Choice(failed(), successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(2, len(result.Results))

	result = r.Choice(failed(), failed(), successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(3, len(result.Results))

	result = r.Choice(failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("failed", result.Expression.Expectation())
	a.Equal(1, len(result.Results))

	result = r.Choice(failed(), failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("failed/failed", result.Expression.Expectation())
	a.Equal(2, len(result.Results))
}

func TestZeroOrMore(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("_abababababzzz")
	r := newBuilder(nil)

	var result *Result

	result = r.ZeroOrMore(successfulTimes(1, 3)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(3, len(result.Results))

	result = r.ZeroOrMore(successfulTimes(1, 1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.ZeroOrMore(successfulTimes(1, 1)).Apply(i, 5)
	a.True(result.Success)
	a.Equal(6, result.End)
	a.Equal(1, len(result.Results))

	result = r.ZeroOrMore(failed()).Apply(i, 5)
	a.True(result.Success)
	a.Equal(5, result.End)
	a.Equal(1, len(result.Results))

	result = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("ab"))).Apply(i, 0)
	a.True(result.Success)
	a.Equal(11, result.End)
	a.Equal(2, len(result.Results))
	a.Equal(0, len(result.Results[0].Results))
	a.Equal(5, len(result.Results[1].Results))
	a.Equal(0, len(result.Results[1].Results[0].Results))
	a.Equal(0, len(result.Results[1].Results[1].Results))

	result = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("x"))).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)

	result = r.Sequence(r.Literal("_"), r.ZeroOrMore(r.Literal("ab"))).Apply(i, 1)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("'_' 'ab'*", result.Expression.Expectation())
	a.Equal(1, len(result.Results))
	a.Equal(0, len(result.Results[0].Results))

	result = r.ZeroOrMore(r.Literal("ab")).Apply(i, 1)
	a.True(result.Success)
	a.Equal(11, result.End)
	a.Equal(5, len(result.Results))
}

func TestOneOrMore(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("_abababababzzz")
	r := newBuilder(nil)

	var result *Result

	result = r.OneOrMore(successfulTimes(1, 3)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(3, len(result.Results))
	a.Equal(0, len(result.Results[0].Results))
	a.Equal(0, len(result.Results[1].Results))

	result = r.OneOrMore(successfulTimes(1, 1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.OneOrMore(successfulTimes(1, 1)).Apply(i, 5)
	a.True(result.Success)
	a.Equal(6, result.End)
	a.Equal(1, len(result.Results))

	result = r.OneOrMore(failed()).Apply(i, 5)
	a.False(result.Success)
	a.Equal(5, result.End)
	a.Equal("failed+", result.Expression.Expectation())
	a.Equal(1, len(result.Results))

	result = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("ab"))).Apply(i, 0)
	a.True(result.Success)
	a.Equal(11, result.End)
	a.Equal(2, len(result.Results))
	a.Equal(5, len(result.Results[1].Results))

	result = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("x"))).Apply(i, 0)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("'_' 'x'+", result.Expression.Expectation())

	result = r.Sequence(r.Literal("_"), r.OneOrMore(r.Literal("ab"))).Apply(i, 1)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("'_' 'ab'+", result.Expression.Expectation())

	result = r.OneOrMore(r.Literal("ab")).Apply(i, 1)
	a.True(result.Success)
	a.Equal(11, result.End)

	result = r.OneOrMore(r.Literal("xx")).Apply(i, 5)
	a.False(result.Success)
	a.Equal(5, result.End)
	a.Equal("'xx'+", result.Expression.Expectation())

	result = r.OneOrMore(r.Literal("z")).Apply(i, 11)
	a.True(result.Success)
	a.Equal(14, result.End)

	result = r.OneOrMore(r.Literal("x")).Apply(i, 11)
	a.False(result.Success)
	a.Equal(11, result.End)
	a.Equal("'x'+", result.Expression.Expectation())

	result = r.OneOrMore(r.Literal("_")).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
}

func TestOptional(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("aaa")
	r := newBuilder(nil)

	var result *Result

	result = r.Optional(successful(1)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(1, result.End)
	a.Equal(1, len(result.Results))

	result = r.Optional(failed()).Apply(i, 0)
	a.True(result.Success)
	a.Equal(0, result.End)
	a.Equal(1, len(result.Results))
}

func TestAnd(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("")
	r := newBuilder(nil)

	var result *Result

	result = r.Test(successful(10)).Apply(i, 0)
	a.True(result.Success)
	a.Equal(0, result.End)
	a.Equal(1, len(result.Results))

	result = r.Test(failed()).Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("&failed", result.Expression.Expectation())
	a.Equal(1, len(result.Results))
}

func TestNot(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("")
	r := newBuilder(nil)

	var result *Result

	result = r.TestNot(successful(10)).Apply(i, 0)
	a.False(result.Success)
	a.Equal(10, result.End)
	a.Equal("!successful", result.Expression.Expectation())
	a.Equal(1, len(result.Results))

	result = r.TestNot(failed()).Apply(i, 0)
	a.True(result.Success)
	a.Equal(0, result.End)
	a.Equal(1, len(result.Results))
}

func TestEnd(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("123")
	r := newBuilder(nil)

	var result *Result

	result = r.End().Apply(i, 0)
	a.False(result.Success)
	a.Equal(0, result.End)
	a.Equal("END", result.Expression.Expectation())
	a.Equal(0, len(result.Results))

	result = r.End().Apply(i, 1)
	a.False(result.Success)
	a.Equal(1, result.End)
	a.Equal("END", result.Expression.Expectation())
	a.Equal(0, len(result.Results))

	result = r.End().Apply(i, 2)
	a.False(result.Success)
	a.Equal(2, result.End)
	a.Equal("END", result.Expression.Expectation())

	result = r.End().Apply(i, 3)
	a.True(result.Success)
	a.Equal(3, result.End)
	a.Equal(0, len(result.Results))
}

func TestEmpty(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("123")
	r := newBuilder(nil)

	var result *Result

	result = r.Empty().Apply(i, 0)
	a.True(result.Success)
	a.Equal(0, result.End)
	a.Equal(0, len(result.Results))

	result = r.Empty().Apply(i, 2)
	a.True(result.Success)
	a.Equal(2, result.End)
	a.Equal(0, len(result.Results))
}
