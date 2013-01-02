package matcher

import (
	"github.com/iNamik/go_lexer"
)

/*****************************************************************************
 * Matcher
 *****************************************************************************/

// Matcher::Reset
func (m *matcher) Reset() Matcher {
	m.stack.Clear()

	m.clearState()

	m.hasResult = false

	return m
}

// Matcher::MatchZeroOrOneBytes
func (m *matcher) MatchZeroOrOneBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrOneBytes(match) })
	return m
}

// Matcher::MatchZeroOrOneRunes
func (m *matcher) MatchZeroOrOneRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrOneRunes(match) })
	return m
}

// Matcher::MatchZeroOrOneRune
func (m *matcher) MatchZeroOrOneRune(match rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrOneRune(match) })
	return m
}

// Matcher::MatchZeroOrOneFunc
func (m *matcher) MatchZeroOrOneFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrOneFunc(match) })
	return m
}

// Matcher::MatchZeroOrMoreBytes
func (m *matcher) MatchZeroOrMoreBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrMoreBytes(match) })
	return m
}

// Matcher::MatchZeroOrMoreRunes
func (m *matcher) MatchZeroOrMoreRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrMoreRunes(match) })
	return m
}

// Matcher::MatchZeroOrMoreFunc
func (m *matcher) MatchZeroOrMoreFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchZeroOrMoreFunc(match) })
	return m
}

// Matcher::MatchOneBytes
func (m *matcher) MatchOneBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneBytes(match) })
	return m
}

// Matcher::MatchOneRunes
func (m *matcher) MatchOneRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneRunes(match) })
	return m
}

// Matcher::MatchOneRune
func (m *matcher) MatchOneRune(match rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneRune(match) })
	return m
}

// Matcher::MatchOneFunc
func (m *matcher) MatchOneFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneFunc(match) })
	return m
}

// Matcher::MatchOneOrMoreBytes
func (m *matcher) MatchOneOrMoreBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneOrMoreBytes(match) })
	return m
}

// Matcher::MatchOneOrMoreRuness
func (m *matcher) MatchOneOrMoreRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneOrMoreRunes(match) })
	return m
}

// Matcher::MatchOneOrMoreFunc
func (m *matcher) MatchOneOrMoreFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchOneOrMoreFunc(match) })
	return m
}

// MatchMinMaxBytes consumes a specified run of matching runes
func (m *matcher) MatchMinMaxBytes(match []byte, min int, max int) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchMinMaxBytes(match, min, max) })
	return m
}

// MatchMinMaxRunes consumes a specified run of matching runes
func (m *matcher) MatchMinMaxRunes(match []rune, min int, max int) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchMinMaxRunes(match, min, max) })
	return m
}

// MatchMinMaxFunc consumes a specified run of matching runes
func (m *matcher) MatchMinMaxFunc(match lexer.MatchFn, min int, max int) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchMinMaxFunc(match, min, max) })
	return m
}

// Matcher::NonMatchZeroOrOneBytes
func (m *matcher) NonMatchZeroOrOneBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrOneBytes(match) })
	return m
}

// Matcher::NonMatchZeroOrOneRunes
func (m *matcher) NonMatchZeroOrOneRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrOneRunes(match) })
	return m
}

// Matcher::NonMatchZeroOrOneFunc
func (m *matcher) NonMatchZeroOrOneFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrOneFunc(match) })
	return m
}

// Matcher::NonMatchZeroOrMoreBytes
func (m *matcher) NonMatchZeroOrMoreBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrMoreBytes(match) })
	return m
}

// Matcher::NonMatchZeroOrMoreRunes
func (m *matcher) NonMatchZeroOrMoreRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrMoreRunes(match) })
	return m
}

// Matcher::NonMatchZeroOrMoreFunc
func (m *matcher) NonMatchZeroOrMoreFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchZeroOrMoreFunc(match) })
	return m
}

// Matcher::NonMatchOneBytes
func (m *matcher) NonMatchOneBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneBytes(match) })
	return m
}

// Matcher::NonMatchOneRunes
func (m *matcher) NonMatchOneRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneRunes(match) })
	return m
}

// Matcher::NonMatchOneFunc
func (m *matcher) NonMatchOneFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneFunc(match) })
	return m
}

// Matcher::NonMatchOneOrMoreBytes
func (m *matcher) NonMatchOneOrMoreBytes(match []byte) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneOrMoreBytes(match) })
	return m
}

// Matcher::NonMatchOneOrMoreRunes
func (m *matcher) NonMatchOneOrMoreRunes(match []rune) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneOrMoreRunes(match) })
	return m
}

// Matcher::NonMatchOneOrMoreFunc
func (m *matcher) NonMatchOneOrMoreFunc(match lexer.MatchFn) MatcherOperator {
	m.doMatch(func() bool { return m.lexer.NonMatchOneOrMoreFunc(match) })
	return m
}

// Matcher::MatchEOF
func (m *matcher) MatchEOF() MatcherOperator {
	m.doMatch(func() bool { return m.lexer.MatchEOF() })
	return m
}

// Matcher::Begin
func (m *matcher) Begin() Matcher {
	return m.begin()
}

// Matcher::End
func (m *matcher) End() MatcherEnd {
	return m
}

// Matcher::EndMatchZeroOrOne
func (m *matcher) EndMatchZeroOrOne() MatcherOperator {

	return m.End().MatchZeroOrOne()
}

// Matcher::EndMatchOne
func (m *matcher) EndMatchOne() MatcherOperator {

	return m.End().MatchZeroOrOne()
}

// Matcher::Result
func (m *matcher) Result() bool {
	if !m.hasResult {
		panic("Calling Result() without trying to match anything")
	}

	result := m.state.result

	if result == false {
		m.lexer.Reset(m.state.marker)
	}

	m.Reset()

	return result
}

/*****************************************************************************
 * Matcher End
 *****************************************************************************/

// MatcherEnd::MatchZeroOrOne
func (m *matcher) MatchZeroOrOne() MatcherOperator {
	m.end(endMatchZeroOrOne)
	return m
}

// MatcherEnd::MatchOne
func (m *matcher) MatchOne() MatcherOperator {
	m.end(endMatchOne)
	return m
}

/*****************************************************************************
 * Matcher Operator
 *****************************************************************************/

// MatcherOperator::And
func (m *matcher) And() Matcher {
	if m.hasResult == false {
		panic("No operator executed before operand")
	}
	m.state.skipNext = m.state.skipAll == true || m.state.result == false
	m.state.fn = matcherAnd
	return m
}

// MatcherOperator::Or
func (m *matcher) Or() Matcher {
	if m.hasResult == false {
		panic("No operator executed before operand")
	}
	m.state.skipNext = m.state.skipAll == true || m.state.result == true
	m.state.fn = matcherOr
	return m
}

// MatcherOperator::AndBegin
func (m *matcher) AndBegin() Matcher {
	m.And()
	m.Begin()
	return m
}

// MatcherOperator::OrBegin
func (m *matcher) OrBegin() Matcher {
	m.Or()
	m.Begin()
	return m
}
