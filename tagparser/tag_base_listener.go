// Generated from Tag.g4 by ANTLR 4.7.

package tagparser // Tag
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseTagListener is a complete listener for a parse tree produced by TagParser.
type BaseTagListener struct{}

var _ TagListener = &BaseTagListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTagListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTagListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTagListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTagListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSet is called when production set is entered.
func (s *BaseTagListener) EnterSet(ctx *SetContext) {}

// ExitSet is called when production set is exited.
func (s *BaseTagListener) ExitSet(ctx *SetContext) {}

// EnterTag is called when production tag is entered.
func (s *BaseTagListener) EnterTag(ctx *TagContext) {}

// ExitTag is called when production tag is exited.
func (s *BaseTagListener) ExitTag(ctx *TagContext) {}
