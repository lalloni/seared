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
	"strings"

	"github.com/lalloni/seared/buffer"
	"github.com/lalloni/seared/node"
)

type Result struct {
	Expression Expression
	Success    bool
	Input      buffer.Buffer
	Start      int
	End        int
	// Parent is the parent PEG expression Result
	Parent *Result
	// Results are the children PEG expressions Results
	Results []*Result
	// Nodes are the parse trees produced
	Nodes []*node.Node
}

func (r *Result) Match() string {
	return r.Input.String(r.Start, r.End)
}

func (r *Result) Length() int {
	return r.End - r.Start
}

func (r *Result) Error() string {
	return "Invalid input '" + r.Input.String(r.Start, r.Start+1) + "' at " + r.Input.Location(r.Start).String() + ", expected " + r.Expression.Expectation()
}

func (r *Result) WithNodes(nodes ...*node.Node) *Result {
	for _, node := range nodes {
		r.Nodes = append(r.Nodes, node)
	}
	return r
}

func (r *Result) WithResults(results ...*Result) *Result {
	for _, result := range results {
		r.Results = append(r.Results, result)
	}
	return r
}

func (r *Result) HasChildren() bool {
	return len(r.Results) > 0
}

func (r *Result) ChildlessResults() []*Result {
	if !r.HasChildren() {
		return []*Result{r}
	}
	rs := []*Result{}
	for _, child := range r.Results {
		rs = append(rs, child.ChildlessResults()...)
	}
	return rs
}

func (r *Result) FailedChildlessResults() []*Result {
	cls := r.ChildlessResults()
	rs := []*Result{}
	for _, cl := range cls {
		if !cl.Success {
			rs = append(rs, cl)
		}
	}
	return rs
}

func formatResultTree(d int, r *Result) string {
	s := strings.Repeat(" ", 2*d) + r.Expression.Name()
	if r.Success {
		s += ": matched '" + r.Match() + "'"
	} else {
		s += ": " + r.Error()
	}
	s += "\n"
	for _, child := range r.Results {
		s += formatResultTree(d+1, child)
	}
	return s
}

func (r *Result) FormatResultTree() string {
	return formatResultTree(0, r)
}

func (r *Result) FormatNodeTree() string {
	ss := []string{}
	for _, node := range r.Nodes {
		ss = append(ss, node.Format())
	}
	return strings.Join(ss, "\n")
}

func (r *Result) FirstRuleAncestor() *Result {
	if r.Expression == nil {
		return nil
	}
	if _, ok := r.Expression.(*rule); ok {
		return r
	}
	if r.Parent == nil {
		return nil
	}
	return r.Parent.FirstRuleAncestor()
}

func deepestFunc(r *Result, f func(r *Result) bool) (result *Result, depth int) {
	result = r
	depth = 0
	for _, child := range r.Results {
		if f(child) {
			rr, dd := deepestFunc(child, f)
			if depth < dd+1 {
				result = rr
				depth = dd + 1
			}
		}
	}
	return
}

func (r *Result) DeepestFailedResult() (result *Result) {
	result, _ = deepestFunc(r, func(r *Result) bool { return r.HasChildren() || !r.Success })
	return
}

func (r *Result) FarthestFailedResult() (result *Result) {
	p := 0
	for _, f := range r.FailedChildlessResults() {
		if f.Start > p {
			p = f.Start
			result = f
		}
	}
	return
}

func (r *Result) BetterError() string {
	ffr := r.FarthestFailedResult()
	if ffr == nil {
		return ""
	}
	frs := make([]*Result, 0)
	for _, fr := range r.FailedChildlessResults() {
		if fr.Start == ffr.Start {
			frs = append(frs, fr)
		}
	}
	ss := make([]string, 0)
	for _, fr := range frs {
		ss = append(ss, fr.Expression.Expectation())
	}
	return "Invalid input '" + r.Input.String(ffr.Start, ffr.Start+1) + "' at " + r.Input.Location(ffr.Start).String() + ", expected " + strings.Join(ss, " or ")
}

func (r *Result) Depth() int {
	if r.Parent == nil {
		return 0
	}
	return 1 + r.Parent.Depth()
}

func Success(expression Expression, input buffer.Buffer, start, end int) *Result {
	return &Result{
		Expression: expression,
		Success:    true,
		Input:      input,
		Start:      start,
		End:        end,
	}
}

func Failure(expression Expression, input buffer.Buffer, start, end int) *Result {
	return &Result{
		Expression: expression,
		Success:    false,
		Input:      input,
		Start:      start,
		End:        end,
	}
}

func ResultsNodes(results []*Result) []*node.Node {
	nodes := []*node.Node{}
	for _, result := range results {
		nodes = append(nodes, result.Nodes...)
	}
	return nodes
}
