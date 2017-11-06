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
type Comment String

// Tokenizer serializes input stream into the set of tokens
type Tokenizer struct {
	input   []byte
	scan    *bufio.Scanner
	curLine []byte
	token   interface{}

	err   error
	warns []error

	comment struct {
		ready bool
		data  *bytes.Buffer
		lin   int
		col   int
		xlin  int
		xcol  int
	}

	// location
	lin int
	col int
}

// careLen computes an offset in visible letters, i.e. \t → 4 characters
func careLen(src []rune) int {
	res := 0
	for _, r := range src {
		switch r {
		case '\t':
			res += 4
		case '\r':
			res = 0
		default:
			res++
		}
	}
	return res
}

// NewTokenizer constructor
func NewTokenizer(input []byte) *Tokenizer {
	input = bytes.Replace(input, []byte("\r"), []byte{}, -1)
	input = bytes.Replace(input, []byte("\t"), []byte("    "), -1)
	res := &Tokenizer{
		input: input,
		scan:  bufio.NewScanner(bytes.NewReader(input)),
		lin:   -1,
	}
	res.comment.data = &bytes.Buffer{}
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
func (t *Tokenizer) Next() (ok bool) {
	if t.err != nil {
		return false
	}
	if t.token != nil {
		return true
	}
	var lineBreak bool
	for {
		if t.curLine == nil {
			lineBreak, ok = t.passWhitespaces()
			if !ok {
				if t.comment.ready {
					t.commitComment()
					return true
				}
				return false
			}
		}
		res := t.nextHeader() || t.nextCode()
		if res {
			return res
		}
		if !t.comment.ready {
			t.comment.ready = true
			t.comment.lin = t.lin
			t.comment.col = t.col
		}
		if lineBreak {
			t.comment.data.WriteByte('\n')
		}
		t.comment.data.Write(t.curLine)
		t.comment.data.WriteByte('\n')
		t.comment.xlin = t.lin
		t.comment.xcol = careLen([]rune(string(t.curLine)))
		t.commitLine()
	}
	return true
}

func (t *Tokenizer) commitComment() {
	t.comment.ready = false
	t.token = Comment{
		Location: Location{
			Lin:  t.comment.lin,
			Col:  t.comment.col,
			XLin: t.lin,
			XCol: t.col,
		},
		Value: t.comment.data.String(),
	}
	t.comment.data.Reset()
}

// Token returns token extracted with
func (t *Tokenizer) Token() interface{} {
	if t.comment.ready {
		res := Comment{
			Location: Location{
				Lin:  t.comment.lin,
				Col:  t.comment.col,
				XLin: t.comment.xlin,
				XCol: t.comment.xcol,
			},
			Value: t.comment.data.String(),
		}
		t.comment.data.Reset()
		t.comment.ready = false
		return res
	}
	res := t.token
	t.token = nil
	return res
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

// passing whitespaces:
//
// 1. is only called when curLine == nil
// 2. passing empty lines (signal if passed any)
// 3. return false if no line was read
func (t *Tokenizer) passWhitespaces() (passedEmpty bool, ok bool) {
	// If the previous line was read out and processed t.curLine must be set to nil
	for t.scan.Scan() {
		t.lin++
		t.col = 0
		for _, r := range []rune(t.scan.Text()) {
			if !unicode.IsSpace(r) {
				t.curLine = t.scan.Bytes()
				ok = true
				return
			}
		}
		passedEmpty = true
	}
	t.err = t.scan.Err()
	return
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
	veryHead, rest := passHeadSpaces([]rune(string(t.curLine)))
	pos := careLen(veryHead)
	if pos < 0 {
		return false
	}
	if pos > 3 {
		return false
	}
	if len(rest) == 0 || rest[0] != '#' {
		return false
	}

	nextPos := -1
	for i, r := range rest {
		if r != '#' {
			nextPos = i
			break
		}
	}
	if nextPos > 6 {
		t.err = t.locReport(t.lin, pos, "header level limit exceeded: %d, cannot be greater than 6", nextPos)
		return false
	}
	if pos > 0 {
		t.appendWarn(t.lin, pos, `please align this line to the left border`)
	}

	spaces, tail := passHeadSpaces(rest[nextPos:])
	body := throwTrailingSpaces(tail)

	t.token = Header{
		Location: Location{
			Lin:  t.lin,
			Col:  pos + t.col,
			XLin: t.lin,
			XCol: t.col + pos + nextPos + careLen(spaces) + careLen(body),
		},
		Content: String{
			Location: Location{
				Lin:  t.lin,
				Col:  t.col + pos + nextPos + careLen(spaces),
				XLin: t.lin,
				XCol: t.col + pos + nextPos + careLen(spaces) + careLen(body),
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
	if pos < 0 {
		return -1
	}
	pos = careLen([]rune(string([]byte(line[:pos]))))
	if pos > 3 {
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
	bbb := []rune(string(t.curLine))
	spaces, rest := passHeadSpaces(bbb)
	body := throwTrailingSpaces(rest)
	if careLen(spaces) > 3 {
		return false
	}
	if !strings.HasPrefix(string(rest), "```") {
		return false
	}

	// OK, it is looks like the fenced code block, getting syntax
	tail := body[3:]
	sss, syntax := passHeadSpaces(tail)
	if careLen(syntax) == 0 {
		t.appendWarn(t.lin, careLen(spaces)+3, "code block syntax (language) name required")
	}

	buf := &bytes.Buffer{}
	codeLin := t.lin
	var codeCol int
	lin := t.lin
	var col int
	for {
		if !t.scan.Scan() {
			t.err = t.locReport(t.lin, careLen(spaces), "unclosed code block")
			return false
		}
		lin++
		col = checkCodeBound(t.scan.Bytes())
		if col >= 0 {
			break
		}
		codeCol = careLen(throwTrailingSpaces([]rune(t.scan.Text())))
		codeLin++
		buf.Write(t.scan.Bytes())
		buf.WriteByte('\n')
	}

	t.token = Code{
		Location: Location{
			Lin:  t.lin,
			Col:  careLen(spaces),
			XLin: lin,
			XCol: col,
		},
		Syntax: String{
			Location: Location{
				Lin:  t.lin,
				Col:  careLen(spaces) + 3 + careLen(sss),
				XLin: t.lin,
				XCol: careLen(spaces) + careLen(body),
			},
			Value: string(syntax),
		},
		Content: String{
			Location: Location{
				Lin:  t.lin + 1,
				Col:  0,
				XLin: codeLin,
				XCol: codeCol,
			},
			Value: buf.String(),
		},
	}
	t.lin = codeLin
	t.col = codeCol
	t.commitLine()
	return true
}
