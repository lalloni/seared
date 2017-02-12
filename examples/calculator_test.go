package examples

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/seared"
)

func TestCalculatorRecognizing(t *testing.T) {
	a := assert.New(t)
	parser := Calculator()
	parser.SetDebug(true)
	parser.SetLog(seared.TestingLog(t))
	cases := []struct {
		expression string
		expected   bool
	}{
		{"1", true},
		{"10", true},
		{"1+1", true},
		{"1*10+1", true},
		{"10*2", true},
		{"10*(2+1)", true},
		{"a20", false},
		{"", false},
		{"1*10+a", false},
		{"329842498274982", true},
	}
	for _, c := range cases {
		t.Logf("testing %q...", c.expression)
		actual := parser.Recognize(c.expression)
		s := "recognized"
		if !actual {
			s = "not " + s
		}
		a.EqualValues(c.expected, actual, fmt.Sprintf("expression %q was %s", c.expression, s))
	}
}
