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
