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
