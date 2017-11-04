package mdl

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"
)

// Location is a Token position in an input stream
type Location struct {
	Lin int
	Col int

	// XLin and XCol for token's end
	XLin int
	XCol int
}

// String ...
type String struct {
	Location
	Value string
}

// Header represents header
type Header struct {
	Location
	Content String
	Level   int
}

// Code represents fenced code block
type Code struct {
	Location
	Syntax  String
	Content String
}

// Comment represents everything else
type Comment struct {
	Location
	Content String
}

// Tokenizer serializes input stream into the set of tokens
type Tokenizer struct {
	input   []byte
	scan    *bufio.Scanner
	curLine []byte
	token   interface{}

	err   error
	warns []error

	// location
	lin int
	col int
}

// NewTokenizer constructor
func NewTokenizer(input []byte) *Tokenizer {
	return &Tokenizer{
		input: input,
		scan:  bufio.NewScanner(bytes.NewReader(input)),
		lin:   -1,
	}
}

// Err returns error status. It is guaranteed Err can only return error if the previous Next call returned false
func (t *Tokenizer) Err() error {
	return t.err
}

// Warnings return list of warnings
func (t *Tokenizer) Warnings() []error {
	return t.warns
}

// Next checks if something to be extracted left
func (t *Tokenizer) Next() bool {
	return t.nextHeader() //|| t.nextCode() || t.nextComment()
}

// Token returns token extracted with
func (t *Tokenizer) Token() interface{} {
	return t.token
}

func (t *Tokenizer) commitLine() {
	t.curLine = nil
}

func passHeadSpaces(src []rune) (spaces []rune, rest []rune) {
	if src == nil {
		return
	}

	var beg int
	for i, r := range src {
		if !unicode.IsSpace(r) {
			beg = i
			break
		}
	}
	spaces = src[:beg]
	rest = src[beg:]
	return
}

func throwTrailingSpaces(src []rune) (res []rune) {
	if src == nil {
		return
	}

	var beg int
	for i, r := range src {
		if !unicode.IsSpace(r) {
			beg = i + 1
		}
	}
	return src[:beg]
}

// passing whitespaces
func (t *Tokenizer) passWhitespaces() bool {
	if t.curLine != nil {
		return true
	}
	// If the previous line was read out and processed t.curLine must be set to nil
	for t.scan.Scan() {
		t.lin++
		t.col = 0
		for _, r := range []rune(t.scan.Text()) {
			if !unicode.IsSpace(r) {
				t.curLine = t.scan.Bytes()
				return true
			}
		}
	}
	t.err = t.scan.Err()
	return false
}

func (t *Tokenizer) locErr(lin int, col int, err error) error {
	return fmt.Errorf("%d:%d: %s", lin, col, err)
}

func (t *Tokenizer) locReport(lin int, col int, format string, a ...interface{}) error {
	return t.locErr(lin, col, fmt.Errorf(format, a...))
}

func (t *Tokenizer) appendWarnErr(lin, col int, err error) {
	t.warns = append(t.warns, t.locErr(lin, col, err))
}

func (t *Tokenizer) appendWarn(lin, col int, format string, a ...interface{}) {
	t.warns = append(t.warns, t.locReport(lin, col, format, a...))
}

func (t *Tokenizer) nextHeader() bool {
	if !t.passWhitespaces() {
		return false
	}

	pos := bytes.IndexByte(t.curLine, '#')
	if pos < 0 {
		return false
	}
	if pos > 3 {
		return false
	}

	nextPos := bytes.IndexFunc(
		t.curLine[pos:],
		func(r rune) bool {
			return r != '#'
		},
	)
	if nextPos > 6 {
		t.err = t.locReport(t.lin, pos, "header level limit exceeded: %d, cannot be greater than 6", nextPos)
		return false
	}
	if pos > 0 {
		t.appendWarn(t.lin, pos, `please align this line to the left border`)
		return false
	}

	rest := []rune(string(t.curLine[pos+nextPos:]))
	spaces, tail := passHeadSpaces(rest)
	body := throwTrailingSpaces(tail)

	t.token = Header{
		Location: Location{
			Lin:  t.lin,
			Col:  pos + t.col,
			XLin: t.lin,
			XCol: t.col + pos + nextPos + len(spaces) + len(body),
		},
		Content: String{
			Location: Location{
				Lin:  t.lin,
				Col:  t.col + pos + nextPos + len(spaces),
				XLin: t.lin,
				XCol: t.col + pos + nextPos + len(spaces) + len(body),
			},
			Value: string(body),
		},
		Level: nextPos,
	}
	t.commitLine()

	return true
}
