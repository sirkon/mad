// Generated from Tag.g4 by ANTLR 4.7.

package tagparser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 6, 29, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 3, 2, 3, 2, 3, 3, 3,
	3, 3, 4, 6, 4, 17, 10, 4, 13, 4, 14, 4, 18, 3, 5, 3, 5, 7, 5, 23, 10, 5,
	12, 5, 14, 5, 26, 11, 5, 3, 5, 3, 5, 2, 2, 6, 3, 3, 5, 4, 7, 5, 9, 6, 3,
	2, 4, 6, 2, 47, 47, 50, 59, 67, 92, 99, 124, 3, 2, 36, 36, 2, 30, 2, 3,
	3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 3, 11,
	3, 2, 2, 2, 5, 13, 3, 2, 2, 2, 7, 16, 3, 2, 2, 2, 9, 20, 3, 2, 2, 2, 11,
	12, 7, 60, 2, 2, 12, 4, 3, 2, 2, 2, 13, 14, 7, 34, 2, 2, 14, 6, 3, 2, 2,
	2, 15, 17, 9, 2, 2, 2, 16, 15, 3, 2, 2, 2, 17, 18, 3, 2, 2, 2, 18, 16,
	3, 2, 2, 2, 18, 19, 3, 2, 2, 2, 19, 8, 3, 2, 2, 2, 20, 24, 7, 36, 2, 2,
	21, 23, 10, 3, 2, 2, 22, 21, 3, 2, 2, 2, 23, 26, 3, 2, 2, 2, 24, 22, 3,
	2, 2, 2, 24, 25, 3, 2, 2, 2, 25, 27, 3, 2, 2, 2, 26, 24, 3, 2, 2, 2, 27,
	28, 7, 36, 2, 2, 28, 10, 3, 2, 2, 2, 5, 2, 18, 24, 2,
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
	"", "':'", "' '",
}

var lexerSymbolicNames = []string{
	"", "", "", "IDENTIFIER", "STRING_LITERAL",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "IDENTIFIER", "STRING_LITERAL",
}

type TagLexer struct {
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

func NewTagLexer(input antlr.CharStream) *TagLexer {

	l := new(TagLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Tag.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// TagLexer tokens.
const (
	TagLexerT__0           = 1
	TagLexerT__1           = 2
	TagLexerIDENTIFIER     = 3
	TagLexerSTRING_LITERAL = 4
)
