package main

// Standard library imports
import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// iNamik imports
import (
	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_lexer/rangeutil"
	"github.com/iNamik/go_lexer_matcher"
	"github.com/iNamik/go_parser"
)

// We define our lexer tokens starting from the pre-defined EOF token
const (
	T_UNKNOWN lexer.TokenType = lexer.TokenTypeUnknown
	T_EOF                     = lexer.TokenTypeEOF
	T_NIL                     = lexer.TokenTypeEOF + iota
	T_OPEN_BRACE
	T_CLOSE_BRACE
	T_OPEN_BRACKET
	T_CLOSE_BRACKET
	T_COLON
	T_COMMA
	T_MINUS
	T_BACK_SLASH
	T_NUMBER
	T_UNQUOTED_STRING
	T_OPEN_QUOTE
	T_CLOSE_QUOTE
	T_TEXT
	T_ESCAPE_CHAR
	T_CHAR_QUOTE
	T_CHAR_BACK_SLASH
	T_CHAR_SLASH
	T_CHAR_BACK_SPACE
	T_CHAR_FORM_FEED
	T_CHAR_LINE_FEED
	T_CHAR_CARRIAGE_RETURN
	T_CHAR_TAB
	T_CHAR_LOWER_U
	T_CHAR_HEX_WORD
	T_CHAR_CONTROL
)

// Single-character tokens
var singleChars = []byte{'{', '}', '[', ']', ':', ','}

var singleTokens = []lexer.TokenType{T_OPEN_BRACE, T_CLOSE_BRACE, T_OPEN_BRACKET, T_CLOSE_BRACKET, T_COLON, T_COMMA}

// Multi-character tokens
var bytesWhitespace = []byte{' ', '\t', '\r', '\n'}

var bytesDigits = rangeutil.RangeToBytes("0-9")

var bytes1to9 = rangeutil.RangeToBytes("1-9")

var bytesAlpha = rangeutil.RangeToBytes("a-zA-Z")

var bytesAlphaNum = rangeutil.RangeToBytes("0-9a-zA-Z")

var bytesHex = rangeutil.RangeToBytes("0-9a-fA-F")

var bytesNonText = rangeutil.RangeToBytes("\u0000-\u001f\\\"")

var escapeChars = []byte{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'}

var escapeTokens = []lexer.TokenType{T_CHAR_QUOTE, T_CHAR_BACK_SLASH, T_CHAR_SLASH, T_CHAR_BACK_SPACE, T_CHAR_FORM_FEED, T_CHAR_LINE_FEED, T_CHAR_CARRIAGE_RETURN, T_CHAR_TAB}

type jsonValue struct {
	value interface{}
	err   error
}

/**
 * tokenTypeAsString returns the string name of the specified token.
 * Used for debugging and error messages.
 */
func tokenTypeAsString(t lexer.TokenType) string {
	var typeString string

	switch t {
	case T_UNKNOWN:
		typeString = "T_UNKNOWN"
	case T_EOF:
		typeString = "T_EOF"
	case T_NIL:
		typeString = "T_NIL"
	case T_OPEN_BRACE:
		typeString = "T_OPEN_BRACE"
	case T_CLOSE_BRACE:
		typeString = "T_CLOSE_BRACE"
	case T_OPEN_BRACKET:
		typeString = "T_OPEN_BRACKET"
	case T_CLOSE_BRACKET:
		typeString = "T_CLOSE_BRACKET"
	case T_COLON:
		typeString = "T_COLON"
	case T_COMMA:
		typeString = "T_COMMA"
	case T_MINUS:
		typeString = "T_MINUS"
	case T_BACK_SLASH:
		typeString = "T_BACK_SLASH"
	case T_NUMBER:
		typeString = "T_NUMBER"
	case T_UNQUOTED_STRING:
		typeString = "T_UNQUOTED_STRING"
	case T_OPEN_QUOTE:
		typeString = "T_OPEN_QUOTE"
	case T_CLOSE_QUOTE:
		typeString = "T_CLOSE_QUOTE"
	case T_TEXT:
		typeString = "T_TEXT"
	case T_ESCAPE_CHAR:
		typeString = "T_ESCAPE_CHAR"
	case T_CHAR_QUOTE:
		typeString = "T_CHAR_QUOTE"
	case T_CHAR_BACK_SLASH:
		typeString = "T_CHAR_BACK_SLASH"
	case T_CHAR_SLASH:
		typeString = "T_CHAR_SLASH"
	case T_CHAR_BACK_SPACE:
		typeString = "T_CHAR_BACK_SPACE"
	case T_CHAR_FORM_FEED:
		typeString = "T_CHAR_FORM_FEED"
	case T_CHAR_LINE_FEED:
		typeString = "T_CHAR_LINE_FEED"
	case T_CHAR_CARRIAGE_RETURN:
		typeString = "T_CHAR_CARRIAGE_RETURN"
	case T_CHAR_TAB:
		typeString = "T_CHAR_TAB"
	case T_CHAR_LOWER_U:
		typeString = "T_CHAR_LOWER_U"
	case T_CHAR_HEX_WORD:
		typeString = "T_CHAR_HEX_WORD"
	case T_CHAR_CONTROL:
		typeString = "T_CHAR_CONTROL"

	default:
		typeString = strconv.Itoa(int(t))
	}

	return typeString
}

/**
 * main
 */
func main() {

	jsonText, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Printf("%v\n", err)
		usage()
	}

	if jsonText == nil || len(jsonText) == 0 {
		usage()
	}

	// Create a new lexer to turn the input text into tokens
	// NOTE : Yes, it seems redundant to create a new reader after going out of
	//        our way to read all the input from io.Stdin, but since we can't
	//        know the longest quoted string on the input, we have to use the
	//        length of the entire string as our buffer len.
	l := lexer.New(lex, strings.NewReader(string(jsonText)), len(jsonText), 3)

	// Create a new parser that feeds off the lexer and generates expression values
	p := parser.New(parse, l, 1)

	value := p.Next() // Pull a JSON value off the parser

	fmt.Printf("%v\n", value)
}

/**
 * usage
 */
func usage() {
	fmt.Println("usage:")
	fmt.Println("  echo '<JSON_TEXT>' | json_decode")
	os.Exit(1)
}

/**
 * lex is the starting StateFn for lexing the input into tokens
 */
func lex(l lexer.Lexer) lexer.StateFn {
	// EOF
	if l.MatchEOF() {
		l.EmitEOF()
		return nil // We're done here
	}

	// Skip whitespace
	if l.MatchOneOrMoreBytes(bytesWhitespace) {
		l.IgnoreToken()
	}

	// Single-char token?
	if i := bytes.IndexRune(singleChars, l.PeekRune(0)); i >= 0 {
		l.NextRune()
		l.EmitToken(singleTokens[i])
		return lex
	}

	// Quote
	if l.MatchOneRune('"') {
		l.EmitToken(T_OPEN_QUOTE)
		return lexQuotedString

		// Unquoted String
	} else if l.MatchOneBytes(bytesAlpha) && l.MatchZeroOrMoreBytes(bytesAlphaNum) {
		l.EmitTokenWithBytes(T_UNQUOTED_STRING)

		// Number:  /-?(0|([1-9][0-9]*))(\.[0-9]+)?([eE][-+]?[0-9]+)?/
	} else if m := matcher.New(l); m.
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
		l.EmitTokenWithBytes(T_NUMBER)
		return lex

		// Unknown
	} else {
		l.NextRune() // Consume unknown rune
		l.EmitTokenWithBytes(T_UNKNOWN)
	}

	// See you again soon!
	return lex
}

/**
 * lexQuotedString parses tokens from within a quoted string, including the
 * final close quote.
 */
func lexQuotedString(l lexer.Lexer) lexer.StateFn {
	r := l.PeekRune(0)

	// If EOF, return
	if r == lexer.RuneEOF {
		return lex
	}

	// Normal text
	if l.NonMatchOneOrMoreBytes(bytesNonText) {
		l.EmitTokenWithBytes(T_TEXT)
		return lexQuotedString
	}

	r = l.PeekRune(0)

	// If control character
	if r < ' ' {
		l.NextRune() // Consume rune
		l.EmitTokenWithBytes(T_CHAR_CONTROL)

		// Close-quote
	} else if r == '"' {
		l.NextRune() // Consume the '"'
		l.EmitToken(T_CLOSE_QUOTE)
		return lex

		// Escape sequence
	} else if r == '\\' {
		l.NextRune() // Consume '\'
		l.EmitToken(T_ESCAPE_CHAR)

		r := l.PeekRune(0)

		// if EOF, return
		if r == lexer.RuneEOF {
			return lex

			// If control character
		} else if r < ' ' {
			l.NextRune() // Consume  rune
			l.EmitTokenWithBytes(T_CHAR_CONTROL)

			// Single-char escape?
		} else if i := bytes.IndexRune(escapeChars, r); i >= 0 {
			l.NextRune() // Consume escape char
			l.EmitToken(escapeTokens[i])

			// Unicode escape sequence
		} else if r == 'u' {
			l.NextRune() // Consume 'u'
			l.EmitToken(T_CHAR_LOWER_U)

			m := matcher.New(l)

			// Match a 4-char hex string
			if m.
				MatchOneBytes(bytesHex).
				And().
				MatchOneBytes(bytesHex).
				And().
				MatchOneBytes(bytesHex).
				And().
				MatchOneBytes(bytesHex).
				Result() {
				l.EmitTokenWithBytes(T_CHAR_HEX_WORD)
			}

			// Unknown escape sequence, emit character after '\' as TEXT
		} else {
			l.NextRune() // consume unknown char
			l.EmitTokenWithBytes(T_TEXT)
		}
	}

	return lexQuotedString
}

/**
 * parse is the starting point for the parser.
 */
func parse(p parser.Parser) parser.StateFn {
	value, err := parseValue(p)

	p.Emit(&jsonValue{value, err})

	return nil // One-pass only
}

/**
 * parseValue parses a json value.  This function may be called recursively
 * if values contain other values.
 */
func parseValue(p parser.Parser) (value interface{}, err error) {
	m := p.Marker()

	err = nil

	t := p.NextToken()

	switch t.Type() {

	// Object (map) '{'
	case T_OPEN_BRACE:
		valueMap := make(map[string]interface{})
		for err == nil && p.PeekTokenType(0) != T_CLOSE_BRACE {
			var entryName string
			var entryValue interface{}

			entryName, entryValue, err = parseNameValuePair(p)

			if err == nil {
				valueMap[entryName] = entryValue

				if p.PeekTokenType(0) == T_COMMA {
					t = p.NextToken() // consume ','

				} else if p.PeekTokenType(0) != T_CLOSE_BRACE {
					t = p.NextToken()

					err = errors.New(fmt.Sprintf("Expecing T_COMMA OR T_CLOSE_BRACE, found %s", tokenTypeAsString(t.Type())))
				}
			}
		}
		if err == nil {
			t = p.NextToken()

			// Expecting '}'
			if t.Type() == T_CLOSE_BRACE {
				value = valueMap
			} else {
				err = wrongTokenError(T_CLOSE_BRACE, t.Type())
			}
		}

	// Array '['
	case T_OPEN_BRACKET:
		valueArray := make([]interface{}, 0)

		for err == nil && p.PeekTokenType(0) != T_CLOSE_BRACKET {
			var entryValue interface{}

			entryValue, err = parseValue(p)

			if err == nil {
				valueArray = append(valueArray, entryValue)

				if p.PeekTokenType(0) == T_COMMA {
					t = p.NextToken() // consume ','

				} else if p.PeekTokenType(0) != T_CLOSE_BRACKET {
					t = p.NextToken()

					err = errors.New(fmt.Sprintf("Expecing T_COMMA OR T_CLOSE_BRACKET, found %s", tokenTypeAsString(t.Type())))
				}
			}
		}
		if err == nil {
			t = p.NextToken()

			// Expecting ']'
			if t.Type() == T_CLOSE_BRACKET {
				value = valueArray
			} else {
				err = wrongTokenError(T_CLOSE_BRACKET, t.Type())
			}
		}

	// null | true | false
	case T_UNQUOTED_STRING:
		switch tString := string(t.Bytes()); tString {
		case "null":
			value = nil
		case "true":
			value = true
		case "false":
			value = false

		default:
			err = errors.New(fmt.Sprintf("Unknown literal: '%s'", tString))
		}

	// number
	case T_NUMBER:
		sNumber := string(t.Bytes())

		value, err = strconv.ParseFloat(sNumber, 64)

		if err != nil {
			err = errors.New(fmt.Sprintf("Error parsing number '%s': %s", sNumber, err))
		}

	// Quoted String '\"'
	case T_OPEN_QUOTE:
		value, err = parseQuotedString(p)

	default:
		err = unexpectedTokenError(t.Type())
	}

	if err != nil {
		p.Reset(m)
		value = nil
	}

	return value, err
}

/**
 * parseQuotedString parses tokens generated by lexQuotedString, returning
 * a decoded string.
 */
func parseQuotedString(p parser.Parser) (value string, err error) {
	var t *lexer.Token

	buffer := new(bytes.Buffer)

	for err == nil && p.PeekTokenType(0) != T_CLOSE_QUOTE {
		t = p.NextToken()

		switch t.Type() {

		// Text
		case T_TEXT:
			buffer.Write(t.Bytes())

		// Escape-char
		case T_ESCAPE_CHAR:
			t = p.NextToken()
			switch t.Type() {
			// Quote '"'
			case T_CHAR_QUOTE:
				buffer.WriteRune('"')

			// Back-slash '\'
			case T_CHAR_BACK_SLASH:
				buffer.WriteRune('\\')

			// Slash '/'
			case T_CHAR_SLASH:
				buffer.WriteRune('/')

			// Back-space '\b'
			case T_CHAR_BACK_SPACE:
				buffer.WriteRune('\b')

			// Form-feed '\f'
			case T_CHAR_FORM_FEED:
				buffer.WriteRune('\f')

			// Line-feed '\n'
			case T_CHAR_LINE_FEED:
				buffer.WriteRune('\n')

			// Carriage-return '\r'
			case T_CHAR_CARRIAGE_RETURN:
				buffer.WriteRune('\r')

			// Tab '\t'
			case T_CHAR_TAB:
				buffer.WriteRune('\t')

			// Hex-word '\uXXXX'
			case T_CHAR_LOWER_U:
				t = p.NextToken()
				if t.Type() == T_CHAR_HEX_WORD {
					hexBytes, err := hex.DecodeString(string(t.Bytes()))
					if err == nil {
						buffer.Write(hexBytes)
					}
				} else {
					err = wrongTokenError(T_CHAR_HEX_WORD, t.Type())
				}

			default:
				err = unexpectedTokenError(t.Type())
			}
		default:
			err = unexpectedTokenError(t.Type())
		}
	}

	if err == nil {
		t = p.NextToken()

		// Expecting '"'
		if t.Type() == T_CLOSE_QUOTE {
			value = buffer.String()
		} else {
			err = wrongTokenError(T_CLOSE_QUOTE, t.Type())
		}
	}

	return value, err
}

/**
 * parseNameValuePair parses a name/value pair from a JSON object
 */
func parseNameValuePair(p parser.Parser) (name string, value interface{}, err error) {
	m := p.Marker()

	err = nil

	if p.PeekTokenType(0) == T_UNQUOTED_STRING && p.PeekTokenType(1) == T_COLON {
		name = string(p.NextToken().Bytes())

		p.NextToken() // Ignore ':'

		value, err = parseValue(p)
	} else if p.PeekTokenType(0) == T_OPEN_QUOTE {
		p.NextToken() // consume '"'

		name, err = parseQuotedString(p)

		if err == nil {
			t := p.NextToken()

			if t.Type() == T_COLON {
				value, err = parseValue(p)
			} else {
				err = wrongTokenError(T_COLON, t.Type())
			}
		}
	}

	if err != nil {
		p.Reset(m)
		name = ""
		value = nil
	}

	return name, value, err
}

/**
 * wrongTokenError
 */
func wrongTokenError(expected, actual lexer.TokenType) error {
	return errors.New(fmt.Sprintf("Expecting %s, found %s", tokenTypeAsString(expected), tokenTypeAsString(actual)))
}

/**
 * unexpectedTokenError
 */
func unexpectedTokenError(t lexer.TokenType) error {
	return errors.New(fmt.Sprintf("Unexpected token: '%s'", tokenTypeAsString(t)))
}
