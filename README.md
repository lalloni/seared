Seared
======

*Seared* is a Go library targeted at allowing easy implementation of text parsers based on powerful [Parsing Expression Grammars](https://en.wikipedia.org/wiki/Parsing_expression_grammar) without the hassle of parser generation steps, while trying to be easy to use, lightweight and appropriate for high-performance parsing needs.

So, how does a grammar definition look in practice? The well-known example "calculator" grammar can be defined like this:

```go
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
```

And this is how that parser could be directly used for syntax recognizing:

```go
parser := Calculator()
success := parser.Recognize("2+1*3+4*(2-1)")
```

License
=======

*Seared* is released under the [Simplified BSD License](./LICENSE) which can be found at the root of this project.
