// Code generated from Raw.g4 by ANTLR 4.7.1. DO NOT EDIT.

package rawparser // Raw
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 12, 48, 4,
	2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 3, 2, 3, 2, 3,
	2, 7, 2, 16, 10, 2, 12, 2, 14, 2, 19, 11, 2, 3, 2, 3, 2, 3, 2, 5, 2, 24,
	10, 2, 3, 2, 5, 2, 27, 10, 2, 3, 3, 7, 3, 30, 10, 3, 12, 3, 14, 3, 33,
	11, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 5, 5, 44,
	10, 5, 3, 6, 3, 6, 3, 6, 2, 2, 7, 2, 4, 6, 8, 10, 2, 3, 3, 2, 4, 5, 2,
	50, 2, 26, 3, 2, 2, 2, 4, 31, 3, 2, 2, 2, 6, 34, 3, 2, 2, 2, 8, 43, 3,
	2, 2, 2, 10, 45, 3, 2, 2, 2, 12, 13, 5, 6, 4, 2, 13, 14, 5, 4, 3, 2, 14,
	16, 3, 2, 2, 2, 15, 12, 3, 2, 2, 2, 16, 19, 3, 2, 2, 2, 17, 15, 3, 2, 2,
	2, 17, 18, 3, 2, 2, 2, 18, 23, 3, 2, 2, 2, 19, 17, 3, 2, 2, 2, 20, 21,
	5, 6, 4, 2, 21, 22, 7, 2, 2, 3, 22, 24, 3, 2, 2, 2, 23, 20, 3, 2, 2, 2,
	23, 24, 3, 2, 2, 2, 24, 27, 3, 2, 2, 2, 25, 27, 3, 2, 2, 2, 26, 17, 3,
	2, 2, 2, 26, 25, 3, 2, 2, 2, 27, 3, 3, 2, 2, 2, 28, 30, 7, 11, 2, 2, 29,
	28, 3, 2, 2, 2, 30, 33, 3, 2, 2, 2, 31, 29, 3, 2, 2, 2, 31, 32, 3, 2, 2,
	2, 32, 5, 3, 2, 2, 2, 33, 31, 3, 2, 2, 2, 34, 35, 7, 6, 2, 2, 35, 36, 7,
	3, 2, 2, 36, 37, 5, 8, 5, 2, 37, 7, 3, 2, 2, 2, 38, 44, 7, 7, 2, 2, 39,
	44, 7, 8, 2, 2, 40, 44, 5, 10, 6, 2, 41, 44, 7, 10, 2, 2, 42, 44, 7, 6,
	2, 2, 43, 38, 3, 2, 2, 2, 43, 39, 3, 2, 2, 2, 43, 40, 3, 2, 2, 2, 43, 41,
	3, 2, 2, 2, 43, 42, 3, 2, 2, 2, 44, 9, 3, 2, 2, 2, 45, 46, 9, 2, 2, 2,
	46, 11, 3, 2, 2, 2, 7, 17, 23, 26, 31, 43,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'='", "'true'", "'false'", "", "", "", "", "", "'\n'",
}
var symbolicNames = []string{
	"", "", "", "", "IDENTIFIER", "STRING_LITERAL", "NUMBER", "INT", "INLINE_STRING",
	"NEWLINE", "WS",
}

var ruleNames = []string{
	"set", "eoc", "line", "value", "boolean",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type RawParser struct {
	*antlr.BaseParser
}

func NewRawParser(input antlr.TokenStream) *RawParser {
	this := new(RawParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Raw.g4"

	return this
}

// RawParser tokens.
const (
	RawParserEOF            = antlr.TokenEOF
	RawParserT__0           = 1
	RawParserT__1           = 2
	RawParserT__2           = 3
	RawParserIDENTIFIER     = 4
	RawParserSTRING_LITERAL = 5
	RawParserNUMBER         = 6
	RawParserINT            = 7
	RawParserINLINE_STRING  = 8
	RawParserNEWLINE        = 9
	RawParserWS             = 10
)

// RawParser rules.
const (
	RawParserRULE_set     = 0
	RawParserRULE_eoc     = 1
	RawParserRULE_line    = 2
	RawParserRULE_value   = 3
	RawParserRULE_boolean = 4
)

// ISetContext is an interface to support dynamic dispatch.
type ISetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSetContext differentiates from other interfaces.
	IsSetContext()
}

type SetContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetContext() *SetContext {
	var p = new(SetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = RawParserRULE_set
	return p
}

func (*SetContext) IsSetContext() {}

func NewSetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetContext {
	var p = new(SetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = RawParserRULE_set

	return p
}

func (s *SetContext) GetParser() antlr.Parser { return s.parser }

func (s *SetContext) AllLine() []ILineContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILineContext)(nil)).Elem())
	var tst = make([]ILineContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILineContext)
		}
	}

	return tst
}

func (s *SetContext) Line(i int) ILineContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILineContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILineContext)
}

func (s *SetContext) AllEoc() []IEocContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEocContext)(nil)).Elem())
	var tst = make([]IEocContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEocContext)
		}
	}

	return tst
}

func (s *SetContext) Eoc(i int) IEocContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEocContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEocContext)
}

func (s *SetContext) EOF() antlr.TerminalNode {
	return s.GetToken(RawParserEOF, 0)
}

func (s *SetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.EnterSet(s)
	}
}

func (s *SetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.ExitSet(s)
	}
}

func (p *RawParser) Set() (localctx ISetContext) {
	localctx = NewSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, RawParserRULE_set)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(24)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(15)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				{
					p.SetState(10)
					p.Line()
				}
				{
					p.SetState(11)
					p.Eoc()
				}

			}
			p.SetState(17)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
		}
		p.SetState(21)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == RawParserIDENTIFIER {
			{
				p.SetState(18)
				p.Line()
			}
			{
				p.SetState(19)
				p.Match(RawParserEOF)
			}

		}

	case 2:
		p.EnterOuterAlt(localctx, 2)

	}

	return localctx
}

// IEocContext is an interface to support dynamic dispatch.
type IEocContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEocContext differentiates from other interfaces.
	IsEocContext()
}

type EocContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEocContext() *EocContext {
	var p = new(EocContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = RawParserRULE_eoc
	return p
}

func (*EocContext) IsEocContext() {}

func NewEocContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EocContext {
	var p = new(EocContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = RawParserRULE_eoc

	return p
}

func (s *EocContext) GetParser() antlr.Parser { return s.parser }

func (s *EocContext) AllNEWLINE() []antlr.TerminalNode {
	return s.GetTokens(RawParserNEWLINE)
}

func (s *EocContext) NEWLINE(i int) antlr.TerminalNode {
	return s.GetToken(RawParserNEWLINE, i)
}

func (s *EocContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EocContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EocContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.EnterEoc(s)
	}
}

func (s *EocContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.ExitEoc(s)
	}
}

func (p *RawParser) Eoc() (localctx IEocContext) {
	localctx = NewEocContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, RawParserRULE_eoc)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(29)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == RawParserNEWLINE {
		{
			p.SetState(26)
			p.Match(RawParserNEWLINE)
		}

		p.SetState(31)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ILineContext is an interface to support dynamic dispatch.
type ILineContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLineContext differentiates from other interfaces.
	IsLineContext()
}

type LineContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLineContext() *LineContext {
	var p = new(LineContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = RawParserRULE_line
	return p
}

func (*LineContext) IsLineContext() {}

func NewLineContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LineContext {
	var p = new(LineContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = RawParserRULE_line

	return p
}

func (s *LineContext) GetParser() antlr.Parser { return s.parser }

func (s *LineContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(RawParserIDENTIFIER, 0)
}

func (s *LineContext) Value() IValueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IValueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *LineContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LineContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LineContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.EnterLine(s)
	}
}

func (s *LineContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.ExitLine(s)
	}
}

func (p *RawParser) Line() (localctx ILineContext) {
	localctx = NewLineContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, RawParserRULE_line)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(32)
		p.Match(RawParserIDENTIFIER)
	}
	{
		p.SetState(33)
		p.Match(RawParserT__0)
	}
	{
		p.SetState(34)
		p.Value()
	}

	return localctx
}

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = RawParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = RawParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(RawParserSTRING_LITERAL, 0)
}

func (s *ValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(RawParserNUMBER, 0)
}

func (s *ValueContext) Boolean() IBooleanContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBooleanContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBooleanContext)
}

func (s *ValueContext) INLINE_STRING() antlr.TerminalNode {
	return s.GetToken(RawParserINLINE_STRING, 0)
}

func (s *ValueContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(RawParserIDENTIFIER, 0)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.ExitValue(s)
	}
}

func (p *RawParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, RawParserRULE_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(41)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case RawParserSTRING_LITERAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(36)
			p.Match(RawParserSTRING_LITERAL)
		}

	case RawParserNUMBER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(37)
			p.Match(RawParserNUMBER)
		}

	case RawParserT__1, RawParserT__2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(38)
			p.Boolean()
		}

	case RawParserINLINE_STRING:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(39)
			p.Match(RawParserINLINE_STRING)
		}

	case RawParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(40)
			p.Match(RawParserIDENTIFIER)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IBooleanContext is an interface to support dynamic dispatch.
type IBooleanContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBooleanContext differentiates from other interfaces.
	IsBooleanContext()
}

type BooleanContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanContext() *BooleanContext {
	var p = new(BooleanContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = RawParserRULE_boolean
	return p
}

func (*BooleanContext) IsBooleanContext() {}

func NewBooleanContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanContext {
	var p = new(BooleanContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = RawParserRULE_boolean

	return p
}

func (s *BooleanContext) GetParser() antlr.Parser { return s.parser }
func (s *BooleanContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.EnterBoolean(s)
	}
}

func (s *BooleanContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(RawListener); ok {
		listenerT.ExitBoolean(s)
	}
}

func (p *RawParser) Boolean() (localctx IBooleanContext) {
	localctx = NewBooleanContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, RawParserRULE_boolean)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(43)
		_la = p.GetTokenStream().LA(1)

		if !(_la == RawParserT__1 || _la == RawParserT__2) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
