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

import "github.com/lalloni/seared"

// ----------------- Rules -----------------------------------------------------

func Digit(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Range('0', '9')
	})
}

func Number(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.OneOrMore(Digit(b))
	})
}

func Parenthesis(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(b.Rune('('), Sum(b), b.Rune(')'))
	})
}

func Factor(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Choice(Number(b), Parenthesis(b))
	})
}

func Term(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Factor(b), b.ZeroOrMore(b.AnyOf("*/"), Factor(b)))
	})
}

func Sum(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Term(b), b.ZeroOrMore(b.AnyOf("+-"), Term(b)))
	})
}

func Calculator(b *seared.Builder) seared.Expression {
	return b.Rule(func() seared.Expression {
		return b.Sequence(Sum(b), b.End())
	})
}

func CalculatorParser() *seared.Parser {
	return seared.NewParser(Calculator)
}
