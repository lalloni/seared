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

func Number(r *seared.Rules) seared.Rule {
	return r.Rule(func() seared.Rule {
		return r.OneOrMore(r.Range('0', '9'))
	})
}

func Factor(r *seared.Rules) seared.Rule {
	return r.Rule(func() seared.Rule {
		return r.Choice(Number(r), r.Sequence(r.Rune('('), Sum(r), r.Rune(')')))
	})
}

func Term(r *seared.Rules) seared.Rule {
	return r.Rule(func() seared.Rule {
		return r.Sequence(Factor(r), r.ZeroOrMore(r.Any("*/"), Factor(r)))
	})
}

func Sum(r *seared.Rules) seared.Rule {
	return r.Rule(func() seared.Rule {
		return r.Sequence(Term(r), r.ZeroOrMore(r.Any("+-"), Term(r)))
	})
}

func Operation(r *seared.Rules) seared.Rule {
	return r.Rule(func() seared.Rule {
		return r.Sequence(Sum(r), r.End())
	})
}

func Calculator() *seared.Parser {
	return seared.NewParser(Operation)
}
