package seared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	a := assert.New(t)
	p := NewParser(func(r *Rules) Rule {
		return r.Rule(func() Rule {
			return r.Rune('a')
		})
	})
	a.Equal("TestNewParser", p.name)
	a.NotNil(p.log)
	a.NotNil(p.main)
}
