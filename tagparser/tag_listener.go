// Generated from Tag.g4 by ANTLR 4.7.

package tagparser // Tag
import "github.com/antlr/antlr4/runtime/Go/antlr"

// TagListener is a complete listener for a parse tree produced by TagParser.
type TagListener interface {
	antlr.ParseTreeListener

	// EnterSet is called when entering the set production.
	EnterSet(c *SetContext)

	// EnterTag is called when entering the tag production.
	EnterTag(c *TagContext)

	// ExitSet is called when exiting the set production.
	ExitSet(c *SetContext)

	// ExitTag is called when exiting the tag production.
	ExitTag(c *TagContext)
}
