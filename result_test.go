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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lalloni/seared/buffer"
)

func TestErrorDepth(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("input")

	var e *Result

	e = Failure(ruleE("l0"), i, 0, 0).DeepestFailedResult()
	a.NotNil(e)
	a.Equal("l0", e.Expression.Expectation())

	e = Failure(ruleE("l0"), i, 0, 0).WithResults(Failure(ruleE("l1"), i, 0, 0)).DeepestFailedResult()
	a.NotNil(e)
	a.Equal("l1", e.Expression.Expectation())

	e = Failure(ruleE("l0"), i, 0, 0).WithResults(Failure(ruleE("l1"), i, 0, 0).WithResults(Failure(ruleE("l2"), i, 0, 0))).DeepestFailedResult()
	a.NotNil(e)
	a.Equal("l2", e.Expression.Expectation())

	e = Failure(ruleE("l0"), i, 0, 0).WithResults(Failure(ruleE("l1a"), i, 0, 0), Failure(ruleE("l1b"), i, 0, 0)).DeepestFailedResult()
	a.NotNil(e)
	a.Equal("l1a", e.Expression.Expectation())
}

func TestErrorLeaves(t *testing.T) {
	a := assert.New(t)
	i := buffer.StringBuffer("input")

	var e []*Result

	e = Failure(ruleE("l0"), i, 0, 0).ChildlessResults()
	a.NotNil(e)
	a.Equal(1, len(e))
	a.Equal("l0", e[0].Expression.Expectation())

	e =
		Failure(ruleE("l0"), i, 0, 0).WithResults(
			Failure(ruleE("l1"), i, 0, 0)).ChildlessResults()
	a.NotNil(e)
	a.Equal(1, len(e))
	a.Equal("l1", e[0].Expression.Expectation())

	e =
		Failure(ruleE("l0"), i, 0, 0).WithResults(
			Failure(ruleE("l1"), i, 0, 0).WithResults(
				Failure(ruleE("l2"), i, 0, 0))).ChildlessResults()
	a.NotNil(e)
	a.Equal(1, len(e))
	a.Equal("l2", e[0].Expression.Expectation())

	e =
		Failure(ruleE("1"), i, 0, 0).WithResults(
			Failure(ruleE("1.1"), i, 0, 0),
			Failure(ruleE("1.2"), i, 0, 0)).ChildlessResults()
	a.NotNil(e)
	a.Equal(2, len(e))
	a.Equal("1.1", e[0].Expression.Expectation())
	a.Equal("1.2", e[1].Expression.Expectation())

	e =
		Failure(ruleE("1"), i, 0, 0).WithResults(
			Failure(ruleE("1.1"), i, 0, 0).WithResults(
				Failure(ruleE("1.1.1"), i, 0, 0).WithResults(
					Failure(ruleE("1.1.1.1"), i, 0, 0),
					Failure(ruleE("1.1.1.2"), i, 0, 0))),
			Failure(ruleE("1.2"), i, 0, 0).WithResults(
				Failure(ruleE("1.2.1"), i, 0, 0)),
			Failure(ruleE("1.3"), i, 0, 0)).ChildlessResults()
	a.NotNil(e)
	a.Equal(4, len(e))
	a.Equal("1.1.1.1", e[0].Expression.Expectation())
	a.Equal("1.1.1.2", e[1].Expression.Expectation())
	a.Equal("1.2.1", e[2].Expression.Expectation())
	a.Equal("1.3", e[3].Expression.Expectation())
}
