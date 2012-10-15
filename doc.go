/*
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
*/
package matcher
