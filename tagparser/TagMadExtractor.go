package tagparser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// MadExtractor is specifically designed to deal with complex tags
type MadExtractor struct {
	madValue string
}

// VisitTerminal ...
func (m *MadExtractor) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode ...
func (m *MadExtractor) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule ...
func (m *MadExtractor) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule ...
func (m *MadExtractor) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSet ...
func (m *MadExtractor) EnterSet(c *SetContext) {}

// EnterTag ...
func (m *MadExtractor) EnterTag(c *TagContext) {
	if c.IDENTIFIER() == nil {
		return
	}
	if c.IDENTIFIER().GetText() == "mad" {
		m.madValue = c.STRING_LITERAL().GetText()
	}
}

// ExitSet ...
func (m *MadExtractor) ExitSet(c *SetContext) {}

// ExitTag ...
func (m *MadExtractor) ExitTag(c *TagContext) {}

type errorPasser struct{}

// SyntaxError ...
func (errorPasser) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
}

// ReportAmbiguity ...
func (errorPasser) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

// ReportAttemptingFullContext ...
func (errorPasser) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
}

// ReportContextSensitivity ...
func (errorPasser) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
}

// ExtractMad extracts mad tag
func ExtractMad(tag string) string {
	m := &MadExtractor{}
	input := antlr.NewInputStream(tag)
	lexer := NewTagLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := NewTagParser(stream)

	p.RemoveErrorListeners()
	p.AddErrorListener(errorPasser{})

	tree := p.Set()
	walker := antlr.NewParseTreeWalker()
	walker.Walk(m, tree)
	if len(m.madValue) > 0 {
		return m.madValue[1 : len(m.madValue)-1]
	}
	return ""
}
