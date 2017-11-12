package mad

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/sirkon/mad/rawparser"
)

// Location is a getToken position in an input stream
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

// Integer integer number
type Integer struct {
	Location
	Value string
	Real  int64
}

// Unsigned unsigned integer number
type Unsigned struct {
	Location
	Value string
	Real  uint64
}

// Float represents floating point number
type Float struct {
	Location
	Value string
	Real  float64
}

// Boolean represents
type Boolean struct {
	Location
	Value string
	Real  bool
}

// tokenizer serializes input stream into the set of tokens
type tokenizer struct {
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

// newTokenizer constructor
func newTokenizer(input []byte) *tokenizer {
	input = bytes.Replace(input, []byte("\r"), []byte{}, -1)
	input = bytes.Replace(input, []byte("\t"), []byte("    "), -1)
	res := &tokenizer{
		input: input,
		scan:  bufio.NewScanner(bytes.NewReader(input)),
		lin:   -1,
	}
	res.comment.data = &bytes.Buffer{}
	return res
}

// Err returns error status. It is guaranteed Err can only return error if the previous next call returned false
func (t *tokenizer) Err() error {
	return t.err
}

// Warnings return list of warnings
func (t *tokenizer) Warnings() []error {
	return t.warns
}

// next checks if something to be extracted left
func (t *tokenizer) next() (ok bool) {
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

func (t *tokenizer) commitComment() {
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

// getToken returns token extracted with. next readout will return nil
func (t *tokenizer) getToken() interface{} {
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

func (t *tokenizer) commitLine() {
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
func (t *tokenizer) passWhitespaces() (passedEmpty bool, ok bool) {
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

func (t *tokenizer) locErr(lin int, col int, err error) error {
	return fmt.Errorf("%d:%d: %s", lin, col, err)
}

func (t *tokenizer) locReport(lin int, col int, format string, a ...interface{}) error {
	return t.locErr(lin, col, fmt.Errorf(format, a...))
}

func (t *tokenizer) appendWarnErr(lin, col int, err error) {
	t.warns = append(t.warns, t.locErr(lin, col, err))
}

func (t *tokenizer) appendWarn(lin, col int, format string, a ...interface{}) {
	t.warns = append(t.warns, t.locReport(lin, col, format, a...))
}

func (t *tokenizer) nextHeader() bool {
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

func (t *tokenizer) nextCode() bool {
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
	t.lin = codeLin + 1
	t.col = codeCol
	t.commitLine()
	return true
}

// Tokenizer abstraction
type Tokenizer interface {
	Next() bool
	Token() interface{}
	Confirm()
}

// MDTokenizer is a markdown level tokenizer
type MDTokenizer struct {
	t         *tokenizer
	confirmed bool
	token     interface{}
}

// NewTokenizer ...
func NewTokenizer(input []byte) *MDTokenizer {
	t := newTokenizer(input)
	return &MDTokenizer{
		t:         t,
		confirmed: true,
	}
}

// next moves underlying tokenizer to its next
func (ts *MDTokenizer) Next() bool {
	if !ts.confirmed {
		return true
	}
	ts.confirmed = false
	return ts.t.next()
}

// getToken returns token from underlying tokenizer
func (ts *MDTokenizer) Token() interface{} {
	if ts.token == nil {
		ts.token = ts.t.getToken()
	}
	return ts.token
}

// Confirm confirms token read out
func (ts *MDTokenizer) Confirm() {
	ts.confirmed = true
	ts.token = nil
}

// RawStorage storage for raw parser output
type RawStorage struct {
	level  int
	items  []interface{}
	errors []error
}

// NewRawStorage constructor
func NewRawStorage(level int) *RawStorage {
	return &RawStorage{
		level: level,
	}
}

func (rs *RawStorage) append(v interface{}) {
	rs.items = append(rs.items, v)
}

// Header consumes header
func (rs *RawStorage) Header(lin, col, xcol int, value string) {
	rs.append(Header{
		Location: Location{
			Lin:  lin,
			Col:  col,
			XLin: lin,
			XCol: xcol,
		},
		Content: String{
			Location: Location{
				Lin:  lin,
				Col:  col,
				XLin: lin,
				XCol: xcol,
			},
			Value: value,
		},
		Level: rs.level,
	})
}

// ValueNumber consumes value as number
func (rs *RawStorage) ValueNumber(lin, col, xcol int, value string) {
	vuint, err := strconv.ParseUint(value, 10, 64)
	if err == nil {
		rs.append(Unsigned{
			Location: Location{
				Lin:  lin,
				Col:  col,
				XLin: lin,
				XCol: xcol,
			},
			Value: value,
			Real:  vuint,
		})
		return
	}
	vint, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		rs.append(Integer{
			Location: Location{
				Lin:  lin,
				Col:  col,
				XLin: lin,
				XCol: xcol,
			},
			Value: value,
			Real:  vint,
		})
	}
	vfloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		rs.errors = append(rs.errors, fmt.Errorf("%d:%d: cannot convert `%s` into numeric type", lin+1, col+1, err))
	} else {
		rs.append(Float{
			Location: Location{
				Lin:  lin,
				Col:  col,
				XLin: lin,
				XCol: xcol,
			},
			Value: value,
			Real:  vfloat,
		})
	}
}

// ValueString consumes value as string
func (rs *RawStorage) ValueString(lin, col, xcol int, value string) {
	rs.append(String{
		Location: Location{
			Lin:  lin,
			Col:  col,
			XLin: lin,
			XCol: xcol,
		},
		Value: value,
	})
}

// Boolean consumes value as boolean
func (rs *RawStorage) Boolean(lin, col, xcol int, value string) {
	var val bool
	switch value {
	case "true":
		val = true
	case "false":
		val = false
	default:
		rs.errors = append(rs.errors, fmt.Errorf("%d:%d: not a boolean value", value))
	}
	rs.append(Boolean{
		Location: Location{
			Lin:  lin,
			Col:  col,
			XLin: lin,
			XCol: xcol,
		},
		Value: value,
		Real:  val,
	})
}

// Data returns collected data
func (rs *RawStorage) Data() []interface{} {
	return rs.items
}

// Err returns collected errors
func (rs *RawStorage) Err() []error {
	return rs.errors
}

type levelInfo struct {
	real    int
	nominal int
}

// FullTokenizer expands fenced blocks for `raw` syntax into the sequence of Header:Value items and 'normalizes' levels.
// Example of normalization, the following tree structure
//
// 1.
//   5.
//   2.
//   2.
//     4.
//
// will be translated into
//
// 1.
//   2.
//   2.
//   2.
//     3.
type FullTokenizer struct {
	levels    []levelInfo
	confirmed bool
	t7r       Tokenizer
	rawData   struct {
		items []interface{}
		index int
	}
	errors []error
}

// Next ...
func (f *FullTokenizer) Next() bool {
	if !f.confirmed {
		return true
	}
	if f.rawData.index < len(f.rawData.items)-1 {
		f.rawData.index++
		return true
	} else {
		f.rawData.items = nil
		f.rawData.index = 0
	}
	return f.t7r.Next()
}

func (f *FullTokenizer) curLevel() int {
	return len(f.levels)
}

// Token ...
func (f *FullTokenizer) Token() interface{} {
	if f.rawData.index < len(f.rawData.items) {
		return f.rawData.items[f.rawData.index]
	}
	res := f.t7r.Token()
	switch v := res.(type) {
	case Header:
		if len(f.levels) == 0 {
			f.levels = append(f.levels, levelInfo{
				real:    1,
				nominal: v.Level,
			})
		} else {
			var i int
			for i = len(f.levels) - 1; i >= 0; i-- {
				if f.levels[i].nominal < v.Level {
					break
				}
			}
			f.levels = f.levels[:i+1]
			f.levels = append(f.levels, levelInfo{
				real:    f.curLevel() + 1,
				nominal: v.Level,
			})
		}
		v.Level = f.curLevel()
		return v
	case Code:
		if v.Syntax.Value != "raw" {
			return res
		}
		// that is the `raw` code block, expanding it
		storage := NewRawStorage(f.curLevel() + 1)
		errors := rawparser.Parse(v.Content.Lin, v.Content.Col, v.Content.Value, storage)
		f.errors = append(f.errors, errors...)
		f.errors = append(f.errors, storage.errors...)
		f.rawData.items = storage.Data()
		f.rawData.index = 0
		f.t7r.Confirm()
		return f.rawData.items[0]
	}
	return res
}

// Confirm confirms confirmation in a confirmative way
func (f *FullTokenizer) Confirm() {
	if f.rawData.items == nil {
		f.t7r.Confirm()
	}
	f.confirmed = true
}

// NewFullTokenizer ...
func NewFullTokenizer(t7r Tokenizer) *FullTokenizer {
	return &FullTokenizer{
		t7r:       t7r,
		confirmed: true,
	}
}

// Err returns stack of errors
func (f *FullTokenizer) Err() []error {
	return f.errors
}
