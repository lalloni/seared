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
	"runtime"
	"strings"
)

func newBuilder(parser *Parser) *Builder {
	return &Builder{
		parser: parser,
		rules:  map[string]Expression{},
	}
}

type RuleOption func(Rule)

type Builder struct {
	parser *Parser
	rules  map[string]Expression
}

func (b *Builder) DropNode() RuleOption {
	return func(r Rule) {
		r.SetDropNode(true)
	}
}

func (b *Builder) OmitNode() RuleOption {
	return func(r Rule) {
		r.SetOmitNode(true)
	}
}

func (b *Builder) Rule(rule func() Expression, options ...RuleOption) Expression {
	key, name := callerKeyName()
	r, ok := b.rules[key]
	if ok {
		return r
	}
	this := newRule(name, b.parser, nil)
	b.rules[key] = this
	this.SetExpression(rule())
	for _, option := range options {
		option(this)
	}
	return this
}

func callerKeyName() (key string, label string) {
	pc, _, _, _ := runtime.Caller(2)
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return f.Function, f.Function[1+strings.LastIndex(f.Function, "."):]
}

func expectations(rules []Expression) []string {
	result := make([]string, len(rules))
	for i, rule := range rules {
		result[i] = rule.Expectation()
	}
	return result
}
