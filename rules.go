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

	"github.com/lalloni/seared/buffer"
)

func rules(parser *Parser) *Rules {
	return &Rules{
		parser: parser,
		rules:  map[string]Rule{},
	}
}

type Rules struct {
	parser *Parser
	rules  map[string]Rule
}

func (p *Rules) Rule(rule func() Rule) Rule {
	key, name := callerKeyName()
	r, ok := p.rules[key]
	if ok {
		return r
	}
	nr := newProxyRule(name, p.parser, nil)
	p.rules[key] = nr
	nr.SetRule(rule())
	return nr
}

func callerKeyName() (key string, label string) {
	pc, _, _, _ := runtime.Caller(2)
	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	return f.Function, f.Function[1+strings.LastIndex(f.Function, "."):]
}

// ------------------------- MATCHERS ------------------------------------------

func (b *Rules) Empty() Rule {
	return newMatcherRule("Empty", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			return true, pos
		})
}

func (b *Rules) End() Rule {
	return newMatcherRule("End", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			return pos >= input.Length(), pos
		})
}

func (b *Rules) Rune(r rune) Rule {
	return newMatcherRule("Terminal", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			nextpos = pos
			success = input.Rune(pos) == r
			if success {
				nextpos++
			}
			return
		})
}

func (b *Rules) Literal(literal string) Rule {
	return newMatcherRule("Literal", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			nextpos = pos + len([]rune(literal))
			success = input.String(pos, nextpos) == literal
			if !success {
				nextpos = pos
			}
			return
		})
}

func (b *Rules) Range(start, end rune) Rule {
	return newMatcherRule("Range", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			r := input.Rune(pos)
			nextpos = pos
			success = r >= start && r <= end
			if success {
				nextpos++
			}
			return
		})
}

func (b *Rules) Any(runes string) Rule {
	return newMatcherRule("Any", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			r := input.Rune(pos)
			nextpos = pos
			for _, rr := range runes {
				success = r == rr
				if success {
					nextpos++
					return
				}
			}
			return
		})
}

// -------------------- OPERATORS ----------------------------------------------

func (b *Rules) Sequence(rules ...Rule) Rule {
	return newMatcherRule("Sequence", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			nextpos = pos
			for _, rule := range rules {
				success, nextpos = rule.Apply(input, nextpos)
				if !success {
					return false, pos
				}
			}
			return
		})
}

func (b *Rules) Choice(rules ...Rule) Rule {
	return newMatcherRule("Choice", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			for _, rule := range rules {
				success, nextpos = rule.Apply(input, pos)
				if success {
					return
				}
			}
			return
		})
}

func (b *Rules) ZeroOrMore(rules ...Rule) Rule {
	rule := b.Sequence(rules...)
	return newMatcherRule("ZeroOrMore", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			success = true
			nextpos = pos
			matching := true
			for matching {
				matching, nextpos = rule.Apply(input, nextpos)
			}
			return
		})
}

func (b *Rules) OneOrMore(rule Rule) Rule {
	return newMatcherRule("OneOrMore", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			nextpos = pos
			matching := true
			for matching {
				matching, nextpos = rule.Apply(input, nextpos)
				if matching {
					success = true
				}
			}
			return
		})
}

func (b *Rules) Optional(rule Rule) Rule {
	return newMatcherRule("Optional", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			_, nextpos = rule.Apply(input, pos)
			return true, nextpos
		})
}

func (b *Rules) And(rule Rule) Rule {
	return newMatcherRule("And", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			success, _ = rule.Apply(input, pos)
			return success, pos
		})
}

func (b *Rules) Not(rule Rule) Rule {
	return newMatcherRule("Not", b.parser,
		func(input buffer.Buffer, pos int) (success bool, nextpos int) {
			success, _ = rule.Apply(input, pos)
			return !success, pos
		})
}
