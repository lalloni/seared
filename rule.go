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
)

type Matcher func(input buffer.Buffer, pos int) (success bool, nextpos int)

type Rule interface {
	Apply(input buffer.Buffer, pos int) (success bool, nextpos int)
}

func newProxyRule(name string, parser *Parser, rule Rule) *proxyRule {
	return &proxyRule{
		name:   name,
		parser: parser,
		rule:   rule,
	}
}

type proxyRule struct {
	name   string
	parser *Parser
	rule   Rule
}

func (r *proxyRule) SetRule(rule Rule) {
	r.rule = rule
}

func (r *proxyRule) Name() string {
	return r.name
}

func (r *proxyRule) Apply(input buffer.Buffer, pos int) (success bool, nextpos int) {
	var loc location.Location
	if r.parser.debug {
		loc = input.Location(pos)
		r.parser.log.Debugf("Trying %q at %s of %q", r.Name(), loc, input.Input())
	}
	success, nextpos = r.rule.Apply(input, pos)
	if r.parser.debug {
		var s string
		if success {
			s = fmt.Sprintf("succeed consuming %q", input.String(pos, nextpos))
		} else {
			s = "failed"
		}
		r.parser.log.Debugf("Result of %q at %s of %q: %s\n", r.Name(), loc, input.Input(), s)
	}
	return
}

func newMatcherRule(name string, parser *Parser, matcher Matcher) *matcherRule {
	return &matcherRule{
		name:    name,
		matcher: matcher,
		parser:  parser,
	}
}

type matcherRule struct {
	name    string
	parser  *Parser
	matcher Matcher
}

func (r *matcherRule) Name() string {
	return r.name
}

func (r *matcherRule) Apply(input buffer.Buffer, pos int) (success bool, nextpos int) {
	success, nextpos = r.matcher(input, pos)
	return
}

func (r *matcherRule) SetMatcher(matcher Matcher) {
	r.matcher = matcher
}
