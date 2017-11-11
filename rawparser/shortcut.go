package rawparser

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Parse shortcut for parsing raw config format
func Parse(lin, col int, data string, storage TokenStorage) (errors []error) {
	input := antlr.NewInputStream(data)
	lexer := NewRawLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewRawParser(stream)
	p.RemoveErrorListeners()
	el := NewErrorHandler(lin, col)
	p.AddErrorListener(el)
	tree := p.Set()

	listener := NewListener(lin, col, storage)

	walker := antlr.NewParseTreeWalker()
	walker.Walk(listener, tree)

	for _, err := range el.errors {
		errors = append(errors, fmt.Errorf("%d:%d: %s", err.Lin+1, err.Col+1, err.Err))
	}
	return
}
