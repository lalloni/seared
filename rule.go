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
	"fmt"

	"github.com/lalloni/seared/buffer"
	"github.com/lalloni/seared/location"
	"github.com/lalloni/seared/node"
)

// Rule is a PEG rule
type Rule interface {
	Expression
	SetExpression(expression Expression)
	SetDropNode(b bool)
	SetOmitNode(b bool)
}

type rule struct {
	Rule
	name       string
	parser     *Parser
	expression Expression
	dropNode   bool
	omitNode   bool
}

func newRule(name string, p *Parser, expression Expression) *rule {
	return &rule{
		name:       name,
		parser:     p,
		expression: expression,
	}
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Expectation() string {
	return r.Name()
}

func (r *rule) SetExpression(e Expression) {
	r.expression = e
}

func (r *rule) SetDropNode(b bool) {
	r.dropNode = b
}

func (r *rule) SetOmitNode(b bool) {
	r.omitNode = b
}

func (r *rule) Apply(input buffer.Buffer, pos int) (result *Result) {
	var loc location.Location
	if r.parser.debug {
		loc = input.Location(pos)
		r.parser.log.Debugf("Trying %q at %s of %q", r.Name(), loc, input.Input())
	}
	inner := r.expression.Apply(input, pos)
	if inner.Success {
		result = Success(r, input, inner.Start, inner.End).WithResults(inner)
		if !r.dropNode {
			if r.omitNode {
				result.WithNodes(inner.Nodes...)
			} else {
				result.WithNodes(node.NewNonTerminal(r.Name(), inner.Nodes))
			}
		}
	} else {
		result = Failure(r, input, inner.Start, inner.End).WithResults(inner)
	}
	if r.parser.debug {
		var s string
		if result.Success {
			s = fmt.Sprintf("succeed consuming %q", input.String(pos, result.End))
		} else {
			s = fmt.Sprintf("failed to consume: %+v", result.Expression.Expectation())
		}
		r.parser.log.Debugf("Result of %q at %s of %q: %s\n", r.Name(), loc, input.Input(), s)
	}
	return
}
