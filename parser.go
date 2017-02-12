package seared

import "github.com/lalloni/seared/buffer"

type Parser struct {
	name  string
	main  Rule
	debug bool
	log   Log
}

func NewParser(main func(*Rules) Rule) *Parser {
	_, name := callerKeyName()
	parser := &Parser{name: name, log: StandardLog()}
	parser.main = main(rules(parser))
	return parser
}

func (p *Parser) Name() string {
	return p.name
}

func (p *Parser) Recognize(input string) (success bool) {
	success, _ = p.main.Apply(buffer.NewStringBuffer(input), 0)
	return
}

func (p *Parser) SetLog(log Log) {
	p.log = log
}

func (p *Parser) SetDebug(debug bool) {
	p.debug = debug
}
