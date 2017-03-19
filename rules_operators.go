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
)

func (b *Builder) Sequence(expressions ...Expression) (this Expression) {
	if len(expressions) < 1 {
		panic("Sequence rules must have inner rules")
	}
	e := strings.Join(expectations(expressions), " ")
	this = newExpression("Sequence", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			children := []*Result{}
			next := start
			for _, expression := range expressions {
				result = expression.Apply(input, next)
				children = append(children, result)
				if !result.Success {
					return Failure(this, input, start, result.End).WithResults(children...)
				}
				next = result.End
			}
			return Success(this, input, start, next).WithResults(children...).WithNodes(ResultsNodes(children)...)
		})
	return
}

func (b *Builder) Choice(expressions ...Expression) (this Expression) {
	if len(expressions) == 0 {
		panic("Choice rules must have inner rules")
	}
	e := strings.Join(expectations(expressions), "/")
	this = newExpression("Choice", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			children := []*Result{}
			for _, expression := range expressions {
				result = expression.Apply(input, start)
				children = append(children, result)
				if result.Success {
					return Success(this, input, start, result.End).WithResults(children...).WithNodes(result.Nodes...)
				}
			}
			return Failure(this, input, start, result.End).WithResults(children...)
		})
	return
}

func (b *Builder) ZeroOrMore(expressions ...Expression) (this Expression) {
	var (
		expression Expression
		e          string
	)
	switch len(expressions) {
	case 0:
		panic("ZeroOrMore rules must have inner rules")
	case 1:
		expression = expressions[0]
		e = expression.Expectation() + "*"
	default:
		expression = b.Sequence(expressions...)
		e = "(" + expression.Expectation() + ")*"
	}
	this = newExpression("ZeroOrMore", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			children := []*Result{}
			next := start
			for {
				result = expression.Apply(input, next)
				children = append(children, result)
				if !result.Success {
					if len(children) > 1 {
						children = children[0 : len(children)-1]
					}
					return Success(this, input, start, next).WithResults(children...).WithNodes(ResultsNodes(children)...)
				}
				next = result.End
			}
		})
	return
}

func (b *Builder) OneOrMore(expressions ...Expression) (this Expression) {
	var (
		expression Expression
		e          string
	)
	switch len(expressions) {
	case 0:
		panic("OneOrMore rules must have inner rules")
	case 1:
		expression = expressions[0]
		e = expression.Expectation() + "+"
	default:
		expression = b.Sequence(expressions...)
		e = "(" + expression.Expectation() + ")+"
	}
	this = newExpression("OneOrMore", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			children := []*Result{}
			next := start
			matched := false
			for {
				result = expression.Apply(input, next)
				children = append(children, result)
				if !result.Success {
					if matched {
						c := children[0 : len(children)-1]
						return Success(this, input, start, next).WithResults(c...).WithNodes(ResultsNodes(c)...)
					}
					return Failure(this, input, start, result.End).WithResults(children...)
				}
				next = result.End
				matched = true
			}
		})
	return
}

func (b *Builder) Optional(expressions ...Expression) (this Expression) {
	var (
		expression Expression
		e          string
	)
	switch len(expressions) {
	case 0:
		panic("Optional rules must have inner rules")
	case 1:
		expression = expressions[0]
		e = expression.Expectation() + "?"
	default:
		expression = b.Sequence(expressions...)
		e = "(" + expression.Expectation() + ")?"
	}
	this = newExpression("Optional", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			inner := expression.Apply(input, start)
			result = Success(this, input, start, inner.End).WithResults(inner).WithNodes(inner.Nodes...)
			return
		})
	return
}

func (b *Builder) Test(expressions ...Expression) (this Expression) {
	var (
		expression Expression
		e          string
	)
	switch len(expressions) {
	case 0:
		panic("Test rules must have inner rules")
	case 1:
		expression = expressions[0]
		e = "&" + expression.Expectation()
	default:
		expression = b.Sequence(expressions...)
		e = "&(" + expression.Expectation() + ")"
	}
	this = newExpression("Test", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			result = expression.Apply(input, start)
			if result.Success {
				return Success(this, input, start, start).WithResults(result)
			}
			return Failure(this, input, start, result.End).WithResults(result)
		})
	return
}

func (b *Builder) TestNot(expressions ...Expression) (this Expression) {
	var (
		expression Expression
		e          string
	)
	switch len(expressions) {
	case 0:
		panic("TestNot rules must have inner rules")
	case 1:
		expression = expressions[0]
		e = "!" + expression.Expectation()
	default:
		expression = b.Sequence(expressions...)
		e = "!(" + expression.Expectation() + ")"
	}
	this = newExpression("TestNot", e, b.parser,
		func(input buffer.Buffer, start int) (result *Result) {
			result = expression.Apply(input, start)
			if !result.Success {
				return Success(this, input, start, start).WithResults(result)
			}
			return Failure(this, input, start, result.End).WithResults(result)
		})
	return
}
