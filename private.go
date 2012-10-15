package matcher

import (
	"github.com/iNamik/go_container/queue"
	"github.com/iNamik/go_lexer"
)

type matcherCallback func() bool

type matcherFn func(bool, matcherCallback) bool

type matcherEndFn func(bool) bool

type matcherState struct {
	skipAll  bool
	skipNext bool
	result   bool
	fn       matcherFn
	marker   *lexer.Marker
}

type matcher struct {
	lexer     lexer.Lexer
	stack     queue.Interface
	hasResult bool
	state     *matcherState
}

// (matcherFn) matcherNil
func matcherNil(b1 bool, f matcherCallback) bool {
	return f()
}

// matcherAnd
func matcherAnd(b1 bool, f matcherCallback) bool {
	return b1 && f()
}

// matcherOr
func matcherOr(b1 bool, f matcherCallback) bool {
	return b1 || f()
}

// endMatchZeroOrOne
func endMatchZeroOrOne(b bool) bool {
	return true
}

// endMatchOne
func endMatchOne(b bool) bool {
	return b
}

// matcher::doMatch
func (m *matcher) doMatch(f matcherCallback) {
	if m.state.skipNext == false {
		m.state.result = m.state.fn(m.state.result, f)

		m.hasResult = true
	}

	m.state.skipNext = m.state.skipAll

	m.state.fn = matcherNil
}

// matcher::clearState
func (m *matcher) clearState() {
	m.state.result = false

	m.state.skipAll = false

	m.state.skipNext = false

	m.state.fn = matcherNil

	m.state.marker = m.lexer.Marker()
}

// matcher::pushState
func (m *matcher) pushState() {
	m.stack.Add(m.state)

	m.state = &matcherState{}

	m.clearState()
}

// matcher::popState
func (m *matcher) popState() {
	i := m.stack.Remove()

	m.state = i.(*matcherState)
}

// matcher::begin
func (m *matcher) begin() Matcher {
	tmpSkipAll := m.state.skipAll || m.state.skipNext

	m.pushState()

	m.state.skipAll = tmpSkipAll

	m.state.skipNext = tmpSkipAll

	return m
}

// matcher::end provides the cleanup and call-back for the End* functions
func (m *matcher) end(endFn matcherEndFn) {
	if m.state.result == false {
		m.lexer.Reset(m.state.marker)
	}

	b := endFn(m.state.result)

	m.popState()

	m.doMatch(func() bool { return b })
}
