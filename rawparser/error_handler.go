package rawparser

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Error ...
type Error struct {
	Lin int
	Col int
	Err error
}

// ErrorHandler handle errors of generated parser
type ErrorHandler struct {
	lin    int
	col    int
	errors []Error
}

// SyntaxError ...
func (eh *ErrorHandler) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	eh.errors = append(eh.errors, Error{
		Lin: line,
		Col: column,
		Err: fmt.Errorf(msg),
	})
}

// ReportAmbiguity TODO
func (eh *ErrorHandler) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	panic("implement me")
}

// ReportAttemptingFullContext dubious
func (eh *ErrorHandler) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

// ReportContextSensitivity TODO
func (eh *ErrorHandler) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

// NewErrorHandler constructor
func NewErrorHandler(lin, col int) *ErrorHandler {
	return &ErrorHandler{
		lin: lin,
		col: col,
	}
}
