// Generated from Raw.g4 by ANTLR 4.7.

package rawparser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 12, 102,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 7, 5, 43, 10, 5, 12, 5, 14, 5, 46, 11,
	5, 3, 6, 3, 6, 3, 6, 7, 6, 51, 10, 6, 12, 6, 14, 6, 54, 11, 6, 3, 6, 3,
	6, 3, 7, 5, 7, 59, 10, 7, 3, 7, 3, 7, 3, 7, 6, 7, 64, 10, 7, 13, 7, 14,
	7, 65, 5, 7, 68, 10, 7, 3, 7, 5, 7, 71, 10, 7, 3, 8, 3, 8, 3, 8, 7, 8,
	76, 10, 8, 12, 8, 14, 8, 79, 11, 8, 5, 8, 81, 10, 8, 3, 9, 3, 9, 5, 9,
	85, 10, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10, 3, 11, 6, 11, 93, 10, 11, 13,
	11, 14, 11, 94, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 3, 13, 2, 2, 14, 3,
	3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 2, 19, 2, 21, 10, 23, 11,
	25, 12, 3, 2, 12, 5, 2, 67, 92, 97, 97, 99, 124, 6, 2, 50, 59, 67, 92,
	97, 97, 99, 124, 6, 2, 12, 12, 15, 15, 36, 36, 94, 94, 3, 2, 50, 59, 3,
	2, 51, 59, 4, 2, 71, 71, 103, 103, 4, 2, 45, 45, 47, 47, 10, 2, 36, 36,
	41, 41, 94, 94, 100, 100, 104, 104, 112, 112, 116, 116, 118, 118, 6, 2,
	11, 12, 15, 15, 34, 34, 63, 63, 5, 2, 11, 11, 15, 15, 34, 34, 2, 110, 2,
	3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2,
	11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 21, 3, 2, 2, 2,
	2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2, 3, 27, 3, 2, 2, 2, 5, 29, 3, 2, 2,
	2, 7, 34, 3, 2, 2, 2, 9, 40, 3, 2, 2, 2, 11, 47, 3, 2, 2, 2, 13, 58, 3,
	2, 2, 2, 15, 80, 3, 2, 2, 2, 17, 82, 3, 2, 2, 2, 19, 88, 3, 2, 2, 2, 21,
	92, 3, 2, 2, 2, 23, 96, 3, 2, 2, 2, 25, 98, 3, 2, 2, 2, 27, 28, 7, 63,
	2, 2, 28, 4, 3, 2, 2, 2, 29, 30, 7, 118, 2, 2, 30, 31, 7, 116, 2, 2, 31,
	32, 7, 119, 2, 2, 32, 33, 7, 103, 2, 2, 33, 6, 3, 2, 2, 2, 34, 35, 7, 104,
	2, 2, 35, 36, 7, 99, 2, 2, 36, 37, 7, 110, 2, 2, 37, 38, 7, 117, 2, 2,
	38, 39, 7, 103, 2, 2, 39, 8, 3, 2, 2, 2, 40, 44, 9, 2, 2, 2, 41, 43, 9,
	3, 2, 2, 42, 41, 3, 2, 2, 2, 43, 46, 3, 2, 2, 2, 44, 42, 3, 2, 2, 2, 44,
	45, 3, 2, 2, 2, 45, 10, 3, 2, 2, 2, 46, 44, 3, 2, 2, 2, 47, 52, 7, 36,
	2, 2, 48, 51, 10, 4, 2, 2, 49, 51, 5, 19, 10, 2, 50, 48, 3, 2, 2, 2, 50,
	49, 3, 2, 2, 2, 51, 54, 3, 2, 2, 2, 52, 50, 3, 2, 2, 2, 52, 53, 3, 2, 2,
	2, 53, 55, 3, 2, 2, 2, 54, 52, 3, 2, 2, 2, 55, 56, 7, 36, 2, 2, 56, 12,
	3, 2, 2, 2, 57, 59, 7, 47, 2, 2, 58, 57, 3, 2, 2, 2, 58, 59, 3, 2, 2, 2,
	59, 60, 3, 2, 2, 2, 60, 67, 5, 15, 8, 2, 61, 63, 7, 48, 2, 2, 62, 64, 9,
	5, 2, 2, 63, 62, 3, 2, 2, 2, 64, 65, 3, 2, 2, 2, 65, 63, 3, 2, 2, 2, 65,
	66, 3, 2, 2, 2, 66, 68, 3, 2, 2, 2, 67, 61, 3, 2, 2, 2, 67, 68, 3, 2, 2,
	2, 68, 70, 3, 2, 2, 2, 69, 71, 5, 17, 9, 2, 70, 69, 3, 2, 2, 2, 70, 71,
	3, 2, 2, 2, 71, 14, 3, 2, 2, 2, 72, 81, 7, 50, 2, 2, 73, 77, 9, 6, 2, 2,
	74, 76, 9, 5, 2, 2, 75, 74, 3, 2, 2, 2, 76, 79, 3, 2, 2, 2, 77, 75, 3,
	2, 2, 2, 77, 78, 3, 2, 2, 2, 78, 81, 3, 2, 2, 2, 79, 77, 3, 2, 2, 2, 80,
	72, 3, 2, 2, 2, 80, 73, 3, 2, 2, 2, 81, 16, 3, 2, 2, 2, 82, 84, 9, 7, 2,
	2, 83, 85, 9, 8, 2, 2, 84, 83, 3, 2, 2, 2, 84, 85, 3, 2, 2, 2, 85, 86,
	3, 2, 2, 2, 86, 87, 5, 15, 8, 2, 87, 18, 3, 2, 2, 2, 88, 89, 7, 94, 2,
	2, 89, 90, 9, 9, 2, 2, 90, 20, 3, 2, 2, 2, 91, 93, 10, 10, 2, 2, 92, 91,
	3, 2, 2, 2, 93, 94, 3, 2, 2, 2, 94, 92, 3, 2, 2, 2, 94, 95, 3, 2, 2, 2,
	95, 22, 3, 2, 2, 2, 96, 97, 7, 12, 2, 2, 97, 24, 3, 2, 2, 2, 98, 99, 9,
	11, 2, 2, 99, 100, 3, 2, 2, 2, 100, 101, 8, 13, 2, 2, 101, 26, 3, 2, 2,
	2, 14, 2, 44, 50, 52, 58, 65, 67, 70, 77, 80, 84, 94, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'='", "'true'", "'false'", "", "", "", "", "", "'\n'",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "IDENTIFIER", "STRING_LITERAL", "NUMBER", "INT", "INLINE_STRING",
	"NEWLINE", "WS",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "IDENTIFIER", "STRING_LITERAL", "NUMBER", "INT",
	"EXP", "EscapeSequence", "INLINE_STRING", "NEWLINE", "WS",
}

type RawLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewRawLexer(input antlr.CharStream) *RawLexer {

	l := new(RawLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Raw.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// RawLexer tokens.
const (
	RawLexerT__0           = 1
	RawLexerT__1           = 2
	RawLexerT__2           = 3
	RawLexerIDENTIFIER     = 4
	RawLexerSTRING_LITERAL = 5
	RawLexerNUMBER         = 6
	RawLexerINT            = 7
	RawLexerINLINE_STRING  = 8
	RawLexerNEWLINE        = 9
	RawLexerWS             = 10
)
