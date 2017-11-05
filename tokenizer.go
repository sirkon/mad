package mdl

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
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
	res := &Tokenizer{
		input: input,
		scan:  bufio.NewScanner(bytes.NewReader(input)),
		lin:   -1,
	}
	in := []rune(string(input))
	lin := 0
	col := 0
	for _, r := range in {
		switch r {
		case '\n':
			lin++
			col = 0
		case '\t':
			res.locReport(lin, col, "tabulations are not allowed, convert them to spaces")
		case '\r':
			res.locReport(lin, col, "\\r symbols are not allowed, remove them")
		default:
			col++
		}
	}
	return res
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
	if t.err != nil {
		return false
	}
	return t.nextHeader() || t.nextCode() // || t.nextComment()
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
		if unicode.IsSpace(r) {
			beg = i + 1
		} else {
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

var (
	codeBound = []byte("```")
)

// checkCodeBound returns -1 if bound is not found, otherwise returns index of symbol just after ```
func checkCodeBound(line []byte) int {
	pos := bytes.Index(line, codeBound)
	if pos < 0 || pos > 3 {
		return -1
	}
	pos += 3
	rest := []rune(string(line[pos:]))
	for _, r := range rest {
		if !unicode.IsSpace(r) {
			return -1
		}
	}
	return pos
}

func (t *Tokenizer) nextCode() bool {
	if !t.passWhitespaces() {
		return false
	}

	bbb := []rune(string(t.curLine))
	spaces, rest := passHeadSpaces(bbb)
	body := throwTrailingSpaces(rest)
	if len(spaces) > 3 {
		return false
	}
	if !strings.HasPrefix(string(rest), "```") {
		return false
	}

	// OK, it is looks like the fenced code block, getting syntax
	tail := body[3:]
	sss, syntax := passHeadSpaces(tail)
	if len(syntax) == 0 {
		t.appendWarn(t.lin, len(spaces)+3, "code block syntax (language) name required")
	}

	buf := &bytes.Buffer{}
	lin := t.lin
	var col int
	for {
		if !t.scan.Scan() {
			t.err = t.locReport(t.lin, len(spaces), "unclosed code block")
			return false
		}
		lin++
		col = checkCodeBound(t.scan.Bytes())
		if col > 0 {
			break
		}
		buf.Write(t.scan.Bytes())
		buf.WriteByte('\n')
	}

	t.token = Code{
		Location: Location{
			Lin:  t.lin,
			Col:  len(spaces),
			XLin: lin,
			XCol: col,
		},
		Syntax: String{
			Location: Location{
				Lin:  t.lin,
				Col:  len(spaces) + 3 + len(sss),
				XLin: t.lin,
				XCol: len(spaces) + len(body),
			},
			Value: string(syntax),
		},
		Content: String{
			Location: Location{
				Lin:  t.lin + 1,
				Col:  0,
				XLin: lin,
				XCol: 0,
			},
			Value: buf.String(),
		},
	}
	t.lin = lin
	t.col = col
	t.commitLine()
	return true
}
