go_lexer_matcher
================

**Fluent Interface to the iNamik go_lexer API with support for concise, regex-like matching**


ABOUT
-----

Package matcher provides a fluent interface to the Match* functions of
iNamik/go_lexer, along with a couple of extra features to support concise,
regex-like matching.

You can now group multiple lexer Match* calls into sub-expressions,
similar to how parens '()' work in standard regular expressions.  You can even
make entire sub-expressions optional (i.e. '()?')

You can concisely code your expressions to the token you want to match, without
any regard to failed matching conditions or cleanup of a partially matched token.
Matcher leverages iNamik/go_lexer's Marker/Reset functionality to automatically
cleanup the lexer state if your expression fails.


MATCHER INTERFACE
-----------------

Below are the interfaces for the main Matcher types:

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


EXAMPLE
-------

Below is an (slightly modified) excerpt from examples/json_decoder that matches a
JSON integer as outlined in RFC 4627 ( see http://www.ietf.org/rfc/rfc4627.txt ).
Notice how the regex shown maps directly to Matcher funtion calls.

	var bytesDigits = rangeutil.RangeToBytes("0-9")

	var bytes1to9 = rangeutil.RangeToBytes("1-9")

	var myLexer = lexer.New(...)

	matcher := matcher.New(myLexer)

	// Regex:  /-?(0|([1-9][0-9]*))(\.[0-9]+)?([eE][-+]?[0-9]+)?/
	if matcher.
		MatchZeroOrOneRune('-').                     // -?      // Leading '-' (optional)
		And().Begin().                               // (       // Begin Integer (required)
		MatchOneRune('0').                           // 0       //   0 is stand alone, unsigned
		Or().Begin().                                // |(      //   Begin non-zero
		MatchOneBytes(bytes1to9).                    // [1-9]   //     First digit - No leading 0's (required)
		And().MatchZeroOrMoreBytes(bytesDigits).     // [0-9]*  //     Extra digits (optional)
		End().MatchOne().                            // )       //   End non-zero
		End().MatchOne().                            // )       // End Integer
		And().Begin().                               // (       // Begin Fraction (optional)
		MatchOneRune('.').                           // \.      //   Decimal '.' (required)
		And().MatchOneOrMoreBytes(bytesDigits).      // [0-9]+  //   Digits (required)
		End().MatchZeroOrOne().                      // )?      // End Fraction
		And().Begin().                               // (       // Begin Exponent (optional)
		MatchOneBytes([]byte{'e', 'E'}).             // [eE]    //   'e' | 'E' (required)
		And().MatchZeroOrOneBytes([]byte{'-', '+'}). // [-+]?   //   Sign (optional)
		And().MatchOneOrMoreBytes(bytesDigits).      // [0-9]+  //   Digits (required)
		End().MatchZeroOrOne().                      //  )?     // End Exponent
		Result() {
			myLexer.EmitTokenWithBytes(T_NUMBER)
	}


INSTALL
-------

The package is built using the Go tool.  Assuming you have correctly set the
$GOPATH variable, you can run the folloing command:

	go get github.com/iNamik/go_lexer_matcher


DEPENDENCIES
------------

go_lexer_matcher depends on the iNamik go_lexer and go_container packages.  Additionally the
'json_decoder' example program requires the iNamik go_parser package:

* https://github.com/iNamik/go_lexer
* https://github.com/iNamik/go_parser
* https://github.com/iNamik/go_container


AUTHORS
-------

 * David Farrell

