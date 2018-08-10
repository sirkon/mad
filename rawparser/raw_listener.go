// Code generated from Raw.g4 by ANTLR 4.7.1. DO NOT EDIT.

package rawparser // Raw
import "github.com/antlr/antlr4/runtime/Go/antlr"

// RawListener is a complete listener for a parse tree produced by RawParser.
type RawListener interface {
	antlr.ParseTreeListener

	// EnterSet is called when entering the set production.
	EnterSet(c *SetContext)

	// EnterEoc is called when entering the eoc production.
	EnterEoc(c *EocContext)

	// EnterLine is called when entering the line production.
	EnterLine(c *LineContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// EnterBoolean is called when entering the boolean production.
	EnterBoolean(c *BooleanContext)

	// ExitSet is called when exiting the set production.
	ExitSet(c *SetContext)

	// ExitEoc is called when exiting the eoc production.
	ExitEoc(c *EocContext)

	// ExitLine is called when exiting the line production.
	ExitLine(c *LineContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)

	// ExitBoolean is called when exiting the boolean production.
	ExitBoolean(c *BooleanContext)
}
