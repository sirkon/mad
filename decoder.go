package mad

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

type level struct {
	level int // represents hierarchy level
	value int // represents header level set by user
}

// Decoder decodes sequence of tokens into destination object
type Decoder struct {
	levels  []level
	tokens  Tokenizer
	lastLin int
	lastCol int
}

// NewDecoder returns a new decoder that reads from r
func NewDecoder(r io.Reader) (*Decoder, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tzer := NewTokenizer(input)
	res := &Decoder{
		tokens: NewFullTokenizer(tzer),
	}
	return res, nil
}

func (d *Decoder) token() Locatable {
	res := d.tokens.Token()
	d.lastLin, d.lastCol = res.Finish()
	return res
}

// LocatedError points to a error position
type LocatedError struct {
	Lin int
	Col int
	Err error
}

// Error for error implementation
func (le LocatedError) Error() string {
	return fmt.Sprintf("%d:%d: %s", le.Lin+1, le.Col+1, le.Err)
}

func locerr(token Locatable, err error) error {
	lin, col := token.Start()
	return LocatedError{
		Lin: lin,
		Col: col,
		Err: err,
	}
}

func locerrf(token Locatable, format string, a ...interface{}) error {
	return locerr(token, fmt.Errorf(format, a...))
}

func (d *Decoder) noTokenErr(err error) error {
	return LocatedError{
		Lin: d.lastLin,
		Col: d.lastCol,
		Err: err,
	}
}

func (d *Decoder) noTokenErrf(format string, a ...interface{}) error {
	return d.noTokenErr(fmt.Errorf(format, a...))
}

// extracts comment from the underlying tokenizer
func (d *Decoder) extractComment(dest *Comment) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("comment required")
	}
	token := d.token()
	cmt, ok := token.(comment)
	if !ok {
		return locerrf(token, "comment expected, got %s instead", token)
	}
	*dest = Comment(cmt.Value)
	d.tokens.Confirm()
	return nil
}

// extracts comment from the underlying tokenizer if it is there
func (d *Decoder) extractOptionalComment(dest **Comment) error {
	if !d.tokens.Next() {
		*dest = nil
		return nil
	}
	token := d.token()
	cmt, ok := token.(comment)
	if !ok {
		*dest = nil
		return nil
	}
	res := Comment(cmt.Value)
	*dest = &res
	d.tokens.Confirm()
	return nil
}

var matchIdentifier *regexp.Regexp

func init() {
	matchIdentifier = regexp.MustCompile(`^[0-9a-zA-Z]+(?:-[0-9a-zA-Z]*)*$`)
}

func codeName(syntax string) (string, bool) {
	var syntaxName string
	ok := matchIdentifier.MatchString(syntax)
	if ok {
		syntaxName = syntax + " code"
	} else {
		syntaxName = "code"
	}
	return syntaxName, ok
}

func isBound(r rune) bool {
	return unicode.IsSpace(r) || r == ',' || r == ';'
}

// extracts code from the underlying tokenizer if it matches against the syntax filter
func (d *Decoder) extractCode(dest *Code) error {
	// if it is fixed syntax (not a list) then it is better be called as json syntax, sql syntax, etc
	syntaxName, expanded := codeName(dest.Syntax)
	if !d.tokens.Next() {
		return d.noTokenErrf("%s required", syntaxName)
	}

	token := d.token()
	cc, ok := token.(code)
	if !ok || (expanded && syntaxName != cc.String()) {
		if len(cc.Syntax.Value) == 0 {
			return locerrf(cc, "%s expected, got code block with unspecified syntax")
		}
		return locerrf(cc, "%s expected, got %s instead", syntaxName, cc)
	}

	// no syntax means any syntax is OK
	if len(dest.Syntax) == 0 {
		dest.Syntax = cc.Syntax.Value
		dest.Code = cc.Content.Value
		d.tokens.Confirm()
		return nil
	}
	if len(cc.Syntax.Value) == 0 {
		return locerrf(cc, "unspecified syntax, only these are allowed: %s", dest.Syntax)
	}

	// check if the input syntax is one of the allowed ones
	pos := strings.Index(dest.Syntax, cc.Syntax.Value)
	end := pos + len(cc.Syntax.Value)
	check := pos >= 0 &&
		// check left bound
		(pos == 0 || isBound(rune(dest.Syntax[pos]))) &&
		// check right bound
		(end == len(dest.Syntax) || isBound(rune(dest.Syntax[end])))
	if !check {
		return locerrf(cc, "unsupported syntax %s, only these are allowed: %s", dest.Syntax)
	}

	dest.Syntax = cc.Syntax.Value
	dest.Code = cc.Content.Value
	d.tokens.Confirm()
	return nil
}

// optionally extracts code from the underlying tokenizer if it matches against the syntax filter
func (d *Decoder) extractOptionalCode(dest **Code) error {
	// if it is fixed syntax (not a list) then it is better be called as json syntax, sql syntax, etc
	syntaxName, expanded := codeName((*dest).Syntax)
	if !d.tokens.Next() {
		*dest = nil
		return nil
	}

	token := d.token()
	cc, ok := token.(code)
	if !ok || (expanded && syntaxName != cc.String()) {
		*dest = nil
		return nil
	}

	// no syntax means any syntax is OK
	if len((*dest).Syntax) == 0 {
		(*dest).Syntax = cc.Syntax.Value
		(*dest).Code = cc.Content.Value
		d.tokens.Confirm()
		return nil
	}
	if len(cc.Syntax.Value) == 0 {
		*dest = nil
		return nil
	}

	// check if the input syntax is one of the allowed ones
	pos := strings.Index((*dest).Syntax, cc.Syntax.Value)
	end := pos + len(cc.Syntax.Value)
	check := pos >= 0 &&
		// check left bound
		(pos == 0 || isBound(rune((*dest).Syntax[pos]))) &&
		// check right bound
		(end == len((*dest).Syntax) || isBound(rune((*dest).Syntax[end])))
	if !check {
		*dest = nil
		return nil
	}

	(*dest).Syntax = cc.Syntax.Value
	(*dest).Code = cc.Content.Value
	d.tokens.Confirm()
	return nil
}

// extracts string from the underlying tokenizer
func (d *Decoder) extractString(dest *string) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("string required")
	}
	token := d.token()
	st, ok := token.(String)
	if !ok {
		return locerrf(token, "string expected, got %s instead", token)
	}
	*dest = st.Value
	d.tokens.Confirm()
	return nil
}

func overflowError(number interface{}, typeSample interface{}) error {
	return fmt.Errorf("overflow error, number %d is too large for %T", number, typeSample)
}

// extracts integer from the underlying tokenizer
func (d *Decoder) extractInt(dest interface{}) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("integer or unsigned required")
	}

	token := d.token()
	var value int64
	switch v := token.(type) {
	case integer:
		value = v.Real
	case unsigned:
		value = int64(v.Real)
		if value < 0 {
			return locerrf(token, "overflow error, number %d is too large for %s", v.Real, reflect.ValueOf(dest).Elem().Type())
		}
	default:
		return locerrf(token, "integer or unsigned required, got %s", token)
	}
	switch v := dest.(type) {
	case *int:
		if int64(int(value)) != value {
			return overflowError(value, *v)
		}
		*v = int(value)
	case *int8:
		if value%0xff != value {
			return overflowError(value, *v)
		}
		*v = int8(value)
	case *int16:
		if value%0xffff != value {
			return overflowError(value, *v)
		}
		*v = int16(value)
	case *int32:
		if value%0xffffffff != value {
			return overflowError(value, *v)
		}
		*v = int32(value)
	case *int64:
		*v = value
	}
	d.tokens.Confirm()
	return nil
}

// extracts uinteger from the underlying tokenizer
func (d *Decoder) extractUint(dest interface{}) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("unsigned required")
	}

	token := d.token()
	ut, ok := token.(unsigned)
	if !ok {
		return locerrf(token, "unsigned required, got %s", token)
	}
	value := ut.Real

	switch v := dest.(type) {
	case *uint:
		if uint64(uint(value)) != value {
			return overflowError(value, *v)
		}
		*v = uint(value)
	case *uint8:
		if value%0xff != value {
			return overflowError(value, *v)
		}
		*v = uint8(value)
	case *uint16:
		if value%0xffff != value {
			return overflowError(value, *v)
		}
		*v = uint16(value)
	case *uint32:
		if value%0xffffffff != value {
			return overflowError(value, *v)
		}
		*v = uint32(value)
	case *uint64:
		if value%0xffffffffffffffff != value {
			return overflowError(value, *v)
		}
		*v = uint64(value)
	}
	d.tokens.Confirm()
	return nil
}

// extracts float from the underlying tokenizer
func (d *Decoder) extractFloat(dest interface{}) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("float required")
	}

	token := d.token()
	var value float64
	switch v := token.(type) {
	case integer:
		value = float64(v.Real)
	case unsigned:
		value = float64(v.Real)
	case float:
		value = v.Real
	default:
		return locerrf(token, "integer or unsigned or float required, got %s", token)
	}
	switch v := dest.(type) {
	case *float32:
		if math.Abs(float64(float32(value))-value) >= 1e-6 {
			return overflowError(value, *v)
		}
		*v = float32(value)
	case *float64:
		*v = float64(value)
	}
	d.tokens.Confirm()
	return nil
}

// Decode decodes data from underlying tokenizer into the dest
// the dest must not be nil
func (d *Decoder) Decode(dest interface{}, context interface{}) error {
	// input must be pointer type
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		panic(fmt.Errorf("pointer type expected, got %T instead", dest))
	}

	// process atomic types
	switch v := dest.(type) {
	case *Comment:
		return d.extractComment(v)
	case **Comment:
		return d.extractOptionalComment(v)
	case *Code:
		return d.extractCode(v)
	case **Code:
		return d.extractOptionalCode(v)
	case *string:
		return d.extractString(v)
	case *int, *int8, *int16, *int32, *int64:
		return d.extractInt(dest)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		return d.extractUint(dest)
	case *float32, *float64:
		return d.extractFloat(dest)
	case []byte:
		panic("[]byte support doesn't make a sense – the idea is all about being as human readable as possible")
	}
	return nil
}
