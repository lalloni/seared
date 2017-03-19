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

package examples

import (
	"github.com/lalloni/seared"
)

func BooleanExpressionParser() *seared.Parser {
	return seared.NewParser(BooleanExpression)
}

func BooleanExpression(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Disjunction(b), b.End())
	})
}

func Disjunction(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Conjunction(b), b.ZeroOrMore(b.Rune('|'), b.Optional(Space(b)), Conjunction(b)))
	})
}

func Conjunction(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Value(b), b.ZeroOrMore(b.Rune('&'), b.Optional(Space(b)), Value(b)))
	})
}

func Value(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(b.Choice(b.Rune('T'), b.Rune('F')), b.Optional(Space(b)))
	})
}

func Space(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.OneOrMore(b.AnyOf(" \t\n"))
	}, b.DropNode())
}
