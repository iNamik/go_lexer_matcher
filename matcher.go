package matcher

import (
	"github.com/iNamik/go_container/stack"
	"github.com/iNamik/go_lexer"
)

type Matcher interface {

	// MatchZeroOrOneBytes consumes the next rune if it matches, always returning true
	MatchZeroOrOneBytes([]byte) MatcherOperator

	// MatchZeroOrOneRunes consumes the next rune if it matches, always returning true
	MatchZeroOrOneRunes([]rune) MatcherOperator

	// MatchZeroOrOneRune consumes the next rune if it matches, always returning true
	MatchZeroOrOneRune(rune) MatcherOperator

	// MatchZeroOrMoreBytes consumes a run of matching runes, always returning true
	MatchZeroOrMoreBytes([]byte) MatcherOperator

	// MatchZeroOrMoreRunes consumes a run of matching runes, always returning true
	MatchZeroOrMoreRunes([]rune) MatcherOperator

	// MatchOneBytes consumes the next rune if its in the list of bytes
	MatchOneBytes([]byte) MatcherOperator

	// MatchOneRunes consumes the next rune if its in the list of bytes
	MatchOneRunes([]rune) MatcherOperator

	// MatchOneRune consumes the next rune if it matches
	MatchOneRune(rune) MatcherOperator

	// MatchOneOrMoreBytes consumes a run of matching runes
	MatchOneOrMoreBytes([]byte) MatcherOperator

	// MatchOneOrMoreRunes consumes a run of matching runes
	MatchOneOrMoreRunes([]rune) MatcherOperator

	// NonMatchZeroOrOneBytes consumes the next rune if it does not match, always returning true
	NonMatchZeroOrOneBytes([]byte) MatcherOperator

	// NonMatchZeroOrOneRuness consumes the next rune if it does not match, always returning true
	NonMatchZeroOrOneRunes([]rune) MatcherOperator

	// NonMatchZeroOrMoreBytes consumes a run of non-matching runes, always returning true
	NonMatchZeroOrMoreBytes([]byte) MatcherOperator

	// NonMatchZeroOrMoreRunes consumes a run of non-matching runes, always returning true
	NonMatchZeroOrMoreRunes([]rune) MatcherOperator

	// NonMatchOneBytes consumes the next rune if its NOT in the list of bytes
	NonMatchOneBytes([]byte) MatcherOperator

	// NonMatchOneRuness consumes the next rune if its NOT in the list of bytes
	NonMatchOneRunes([]rune) MatcherOperator

	// NonMatchOneOrMoreBytes consumes a run of non-matching runes
	NonMatchOneOrMoreBytes([]byte) MatcherOperator

	// NonMatchOneOrMoreRunes consumes a run of non-matching runes
	NonMatchOneOrMoreRunes([]rune) MatcherOperator

	// MatchEOF tries to match the next rune against RuneEOF
	MatchEOF() MatcherOperator

	// BeginOne begins a new grouping that is expected to match (i.e required)
	Begin() Matcher

	// End ends a grouping. NOTE You are expected to call one of the MatcherEnd
	// functions in order to apply the result of the grouping to your current result.
	End() MatcherEnd

	// EndMatchOne performs End(), followed by MatchOne()
	EndMatchOne() MatcherOperator

	// EndMatchZeroOrOne performs End(), followed by MatchZeroOrOne
	EndMatchZeroOrOne() MatcherOperator

	// Result returns the final result of the matcher, resetting the
	// matcher state if the result is false.
	Result() bool
}

type MatcherEnd interface {
	// MatchOne
	MatchOne() MatcherOperator

	// MatchZeroOrOne
	MatchZeroOrOne() MatcherOperator
}

type MatcherOperator interface {

	// And Performs a logical 'and' between the current matcher state and the
	// next operand.  Short-circuit logic is performed, whereby the next operand
	// will not actually be executed if the current matcher state is already
	// false
	And() Matcher

	// Or Performs a logical 'or' between the current matcher result and the
	// next operand.  Short-circuit logic is performed, whereby the next operand
	// will not actually be executed if the current matcher state is already
	// true
	Or() Matcher

	// AndBegin performs an And(), followed by a Begin()
	AndBegin() Matcher

	// OrBegin performs an Or(), followed by a Begin()
	OrBegin() Matcher

	// End ends a grouping. NOTE You are expected to call one of the MatcherEnd
	// functions in order to apply the result of the grouping to your current result.
	End() MatcherEnd

	// EndMatchOne performs End(), followed by MatchOne()
	EndMatchOne() MatcherOperator

	// EndMatchZeroOrOne performs End(), followed by MatchZeroOrOne
	EndMatchZeroOrOne() MatcherOperator

	// Result returns the final result of the matcher, resetting the
	// matcher state if the result is false.
	Result() bool

	// Reset resets the state of the matcher
	Reset() Matcher
}

// New createas a new Matcher against the specifid Lexer
func New(l lexer.Lexer) Matcher {
	m := &matcher{
		lexer: l,
		stack: stack.New(4), // 4 is just a nice number that seems appropriate
		state: &matcherState{},
	}

	m.Reset()

	return m
}
