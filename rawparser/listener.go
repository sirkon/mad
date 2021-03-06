package rawparser

import (
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// TokenStorage is an abstraction over token storage
type TokenStorage interface {
	Header(lin, col, xcol int, value string)
	ValueNumber(lin, col, xcol int, value string)
	ValueString(lin, col, xcol int, value string)
	Boolean(lin, col, xcol int, value string)
}

// Listener implementation of RawListener to scan over raw fenced code blocks
type Listener struct {
	lin     int
	col     int
	storage TokenStorage
}

// NewListener constructor
func NewListener(lin, col int, storage TokenStorage) *Listener {
	return &Listener{
		lin:     lin,
		col:     col,
		storage: storage,
	}
}

// VisitTerminal ...
func (l *Listener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode ...
func (l *Listener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule ...
func (l *Listener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule ...
func (l *Listener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSet ...
func (l *Listener) EnterSet(c *SetContext) {}

// EnterEoc ...
func (l *Listener) EnterEoc(c *EocContext) {}

// EnterLine adds identifier as a header
func (l *Listener) EnterLine(c *LineContext) {
	lin, col := l.pos(c.IDENTIFIER().GetSymbol())
	l.storage.Header(
		lin,
		col,
		col+len(c.IDENTIFIER().GetText()),
		c.IDENTIFIER().GetText(),
	)
}

func (l *Listener) pos(token antlr.Token) (lin int, col int) {
	return l.lin + token.GetLine() - 1, l.col + token.GetColumn()
}

// EnterValue ...
func (l *Listener) EnterValue(c *ValueContext) {
	switch {
	case c.STRING_LITERAL() != nil:
		lin, col := l.pos(c.STRING_LITERAL().GetSymbol())
		xcol := col + len(c.STRING_LITERAL().GetText())
		text := c.STRING_LITERAL().GetText()
		text = strings.Replace(text, "\\\"", "\"", -1)
		text = strings.Replace(text, "\\'", "'", -1)
		text = strings.Replace(text, "\\t", "\t", -1)
		text = strings.Replace(text, "\\r", "\r", -1)
		text = strings.Replace(text, "\\n", "\n", -1)
		text = strings.Replace(text, "\\b", "\b", -1)
		text = strings.Replace(text, "\\f", "\f", -1)
		l.storage.ValueString(lin, col, xcol, text[1:len(text)-1])
	case c.NUMBER() != nil:
		lin, col := l.pos(c.NUMBER().GetSymbol())
		l.storage.ValueNumber(lin, col, col+len(c.NUMBER().GetText()), c.NUMBER().GetText())
	case c.INLINE_STRING() != nil:
		lin, col := l.pos(c.INLINE_STRING().GetSymbol())
		l.storage.ValueString(lin, col, col+len(c.INLINE_STRING().GetText()), c.INLINE_STRING().GetText())
	case c.IDENTIFIER() != nil:
		lin, col := l.pos(c.IDENTIFIER().GetSymbol())
		l.storage.ValueString(lin, col, col+len(c.IDENTIFIER().GetText()), c.IDENTIFIER().GetText())
	}

}

// ExitSet ...
func (l *Listener) ExitSet(c *SetContext) {}

// ExitEoc ...
func (l *Listener) ExitEoc(c *EocContext) {}

// ExitLine ...
func (l *Listener) ExitLine(c *LineContext) {}

// ExitValue ...
func (l *Listener) ExitValue(c *ValueContext) {}

// EnterBoolean ...
func (l *Listener) EnterBoolean(c *BooleanContext) {
	lin, col := l.pos(c.GetStart())
	l.storage.Boolean(lin, col, col+len(c.GetText()), c.GetText())
}

// ExitBoolean ...
func (l *Listener) ExitBoolean(c *BooleanContext) {}
