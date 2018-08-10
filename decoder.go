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

	"os"

	"bytes"

	"github.com/sirkon/message"
)

// Unmarshaler is implemented by types that can unmarshal
// text, identifier, unsinged, integer, numeric, inline_string
// or boolean.
type Unmarshaler interface {
	Unmarshal(data string) (err error)
}

// Decoder decodes sequence of tokens into destination object
type Decoder struct {
	levels  []int
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

// TokenError returns LocatedError using token to find Lin and Col position
func TokenError(token Locatable, err error) LocatedError {
	lin, col := token.Start()
	return LocatedError{
		Lin: lin,
		Col: col,
		Err: err,
	}
}

// TokenErrorf returns formatted LocatedError using TokenError
func TokenErrorf(token Locatable, format string, a ...interface{}) LocatedError {
	return TokenError(token, fmt.Errorf(format, a...))
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

func (d *Decoder) passComment() bool {
	for d.tokens.Next() {
		token := d.token()
		if _, ok := token.(comment); ok {
			d.tokens.Confirm()
		} else {
			return true
		}
	}
	return false
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
func (d *Decoder) extractCode(dest *Code, ctx Context) error {
	d.passComment()
	syntax := ctx.GetString("syntax", "")
	// if it is fixed syntax (not a list) then it is better be called as json syntax, sql syntax, etc
	syntaxName, expanded := codeName(syntax)
	if !d.tokens.Next() {
		return d.noTokenErrf("%s required", syntaxName)
	}

	token := d.token()
	cc, ok := token.(code)
	if !ok {
		var format string
		if len(syntaxName) > 0 {
			format = syntaxName + " was expected, got %s"
		} else {
			format = "code block was expected, got %s"
		}
		return locerrf(cc, format, token)
	}
	if expanded && syntaxName != cc.String() {
		if len(cc.Syntax.Value) == 0 {
			return locerrf(cc, "%s expected, got code block with unspecified syntax")
		}
		return locerrf(cc, "%s expected, got %s instead", syntaxName, cc)
	}

	// no syntax means any syntax is OK
	if len(syntax) == 0 {
		dest.Syntax = cc.Syntax.Value
		dest.Code = cc.Content.Value
		d.tokens.Confirm()
		return nil
	}
	if len(cc.Syntax.Value) == 0 {
		return locerrf(cc, "unspecified syntax, only these are allowed: %s", dest.Syntax)
	}

	// check if the input syntax is one of the allowed ones
	pos := strings.Index(syntax, cc.Syntax.Value)
	end := pos + len(cc.Syntax.Value)
	check := pos >= 0 &&
		// check left bound
		(pos == 0 || isBound(rune(syntax[pos-1]))) &&
		// check right bound
		(end == len(syntax) || isBound(rune(syntax[end])))
	if !check {
		return locerrf(cc, "unsupported syntax %s, only these are allowed: %s", cc.Syntax.Value, syntax)
	}

	dest.loc = cc.Content.Location
	dest.Syntax = cc.Syntax.Value
	dest.Code = cc.Content.Value
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

// extracts data from the underlying tokenizer using type own capabilities
func (d *Decoder) extractUnmarshaler(dest Unmarshaler) error {
	if !d.tokens.Next() {
		return d.noTokenErrf("token required")
	}

	token := d.token()
	var value string
	switch v := token.(type) {
	case integer:
		value = v.Value
	case unsigned:
		value = v.Value
	case float:
		value = v.Value
	case String:
		value = v.Value
	case boolean:
		value = v.Value
	default:
		return locerrf(token, "only items of raw fenced code blocks can be unmarshalled, got %s", token)
	}
	if err := dest.Unmarshal(value); err != nil {
		return locerr(token, err)
	}

	d.tokens.Confirm()
	return nil
}

// fill input slice
func (d *Decoder) extractSlice(tmp reflect.Value, dest interface{}, ctx Context) error {
	for {
		value := reflect.New(tmp.Type().Elem())
		if err := d.Decode(value.Interface(), ctx); err != nil {
			reflect.ValueOf(dest).Elem().Set(tmp)
			break
		}
		tmp = reflect.Append(tmp, value.Elem())
	}
	return nil
}

func (d *Decoder) level() int {
	return len(d.levels)
}

func (d *Decoder) extractMap(dest interface{}, ctx Context) error {
	for d.tokens.Next() {
		token := d.token()
		h, ok := token.(header)
		if !ok {
			return nil
		}
		if h.Level <= d.level() {
			return nil
		} else if h.Level > d.level()+1 {
			return locerrf(h, "unexpected header level %d, may be cut some #s?", h.Level)
		}
		vv := reflect.ValueOf(dest).Elem().MapIndex(reflect.ValueOf(h.Content.Value))
		if vv.IsValid() {
			return locerrf(h, "duplicate key `%s`", h.Content.Value)
		}
		d.levels = append(d.levels, h.Level)
		d.tokens.Confirm()

		// create map value and decode it then fill it
		vdest := reflect.New(reflect.ValueOf(dest).Elem().Type().Elem()).Interface()
		err := d.Decode(vdest, ctx)
		d.levels = d.levels[:len(d.levels)-1]
		if err != nil {
			return err
		}
		reflect.ValueOf(dest).Elem().SetMapIndex(reflect.ValueOf(h.Content.Value), reflect.ValueOf(vdest).Elem())
	}
	return nil
}

type fieldDescription struct {
	regex    *regexp.Regexp
	index    []int
	labels   map[string]string
	required bool
	checked  bool
}

func extractFieldInfo(tag string, errf func(format string, a ...interface{}) error) (f fieldDescription, err error) {
	f.labels = map[string]string{}
	first := true
	s := newSplitter(tag)
	for s.next() {
		retext := s.text()
		if first {
			first = false
			if !strings.HasPrefix(retext, "^") {
				retext = "^" + retext
			}
			if !strings.HasSuffix(retext, "$") {
				retext += "$"
			}
			f.regex, err = regexp.Compile(retext)
			if err != nil {
				return f, errf("incorrect regular expression `%s`", retext)
			}
			continue
		}
		key, value, ok := keyVal(retext)
		if !ok {
			return f, errf("incorrect tag fragment `%s`", retext)
		}
		if _, ok = f.labels[key]; ok {
			return f, errf("duplicate lable name `%s`", key)
		}
		f.labels[key] = value
	}
	return
}

func extractFieldsMetainfo(dest interface{}) (fields []fieldDescription, err error) {
	tmp := reflect.ValueOf(dest).Elem()
	limit := tmp.NumField()
	for i := 0; i < limit; i++ {
		fieldType := tmp.Type().Field(i)
		if fieldType.Type.Kind() == reflect.Struct && fieldType.Anonymous {
			embeds, err := extractFieldsMetainfo(tmp.Field(i).Addr().Interface())
			if err != nil {
				return nil, err
			}
			for _, embed := range embeds {
				embed.index = append([]int{i}, embed.index...)
				fields = append(fields, embed)
			}
		}
		rawMad := extractMad(string(fieldType.Tag))
		if len(rawMad) == 0 {
			continue
		}
		f, err := extractFieldInfo(rawMad, func(format string, a ...interface{}) error {
			return fmt.Errorf(fmt.Sprintf(format, a...) +
				fmt.Sprintf(" for field %s of type %T", fieldType.Name, tmp.Interface()))
		})
		if err != nil {
			return fields, err
		}
		f.index = []int{i}
		switch fieldType.Type.Kind() {
		case reflect.Map, reflect.Slice, reflect.Ptr:
			f.required = false
		default:
			f.required = true
		}
		if v, ok := dest.(Manual); ok {
			// This is sufficient type
			f.required = v.Required()
		}
		fields = append(fields, f)
	}
	return
}

// getFieldData extracts reflect.Value and reflect.Type of field.
func getFieldData(dest interface{}, indices []int) (field reflect.Value, fieldType reflect.StructField) {
	field = reflect.ValueOf(dest).Elem()
	for i, fieldIndex := range indices {
		if i == len(indices)-1 {
			fieldType = field.Type().Field(fieldIndex)
		}
		field = field.Field(fieldIndex)
	}
	return
}

func (d *Decoder) extractStruct(dest interface{}, ctx Context) (err error) {
	fields, err := extractFieldsMetainfo(dest)
	if err != nil {
		return
	}
	tmp := reflect.TypeOf(dest).Elem()
	value := reflect.ValueOf(dest).Elem()
	realType := value.Interface()
	taken := map[string]Locatable{}
	for d.tokens.Next() {
		if !d.passComment() {
			break
		}
		token := d.token()

		h, ok := token.(header)
		if !ok {
			break
		}
		if h.Level <= d.level() {
			return nil
		} else if h.Level > d.level()+1 {
			return locerrf(h, "unexpected header level %d, may be cut some #s?", h.Level)
		}
		hname := h.Content.Value
		if prev, ok := taken[hname]; ok {
			plin, pcol := prev.Start()
			return locerrf(token, "key `%s` has been taken already at (%d, %d)", hname, plin, pcol)
		}

		// checks on general header validity done
		taken[hname] = token
		d.levels = append(d.levels, h.Level)
		d.tokens.Confirm()

		// now look for the right field
		indices := []int{}
		for i, field := range fields {
			if field.regex.MatchString(hname) {
				indices = append(indices, i)
			}
		}
		if len(indices) == 0 {
			return locerrf(token, "no match for header `%s` in type %T", hname, realType)
		}
		if len(indices) > 1 {
			fnames := []string{}
			for _, i := range indices {
				_, fieldType := getFieldData(dest, fields[i].index)
				fnames = append(fnames, fieldType.Name)
			}
			return locerrf(
				token,
				"ambigious mapping of header `%s`, several fields were matched in type %T: %s",
				hname,
				strings.Join(fnames, ", "),
				realType,
			)
		}

		// process the field
		index := indices[0]
		fieldMeta := fields[index]
		newCtx := ctx.New()
		for k, v := range fieldMeta.labels {
			ctx.Set(k, v)
		}
		fieldValue, _ := getFieldData(dest, fieldMeta.index)
		if v, ok := fieldValue.Interface().(Manual); ok {
			n, err := v.Decode(v, h.Content, d, ctx)
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(n))
		} else {
			fv := fieldValue.Addr().Interface()
			if err = d.Decode(fv, newCtx); err != nil {
				return err
			}
		}
		d.levels = d.levels[:len(d.levels)-1]
		fieldMeta.checked = true
		fields[index] = fieldMeta
	}

	// and now the last check, all required fields must be checked
	missed := []string{}
	for _, field := range fields {
		if field.required && !field.checked {
			missed = append(missed, tmp.Field(field.index[0]).Name)
		}
	}
	if len(missed) > 0 {
		return fmt.Errorf(
			"%d:%d: fields %s of type %T are required but were missed",
			d.lastLin+1, d.lastCol+1,
			strings.Join(missed, ", "),
			realType,
		)
	}

	return nil
}

// Decode decodes data from underlying tokenizer into the dest
// the dest must not be nil
func (d *Decoder) Decode(dest interface{}, ctx Context) error {
	// input must be pointer type
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		panic(fmt.Errorf("pointer type expected, got %T instead", dest))
	}

	// process atomic types
	switch v := dest.(type) {
	case *string:
		return d.extractString(v)
	case *int, *int8, *int16, *int32, *int64:
		return d.extractInt(dest)
	case *uint, *uint8, *uint16, *uint32, *uint64:
		return d.extractUint(dest)
	case *float32, *float64:
		return d.extractFloat(dest)
	case Unmarshaler:
		return d.extractUnmarshaler(v)
	case []byte:
		panic("[]byte support doesn't make a sense – the idea is all about being as human readable as possible")
	case Decodable:
		return v.Decode(d, ctx)
	}

	// may be a pointer to decodable?
	tmp := reflect.ValueOf(dest)
	decodable := reflect.TypeOf((*Decodable)(nil)).Elem()
	if tmp.Elem().Type().Implements(decodable) {
		v := tmp.Elem().Interface().(Decodable)
		if err := v.Decode(d, ctx); err != nil {
			// setting up
			tmp.Elem().Set(reflect.Zero(tmp.Elem().Type()))
			return nil
		}
		return nil
	}

	// may be an slice of something that should be decodable
	tmp = reflect.ValueOf(dest).Elem()
	switch tmp.Kind() {
	case reflect.Slice:
		return d.extractSlice(tmp, dest, ctx)

		// may be map of string → something
	case reflect.Map:
		if tmp.Type().Key().Kind() != reflect.String {
			panic(fmt.Errorf("only map[string]T are allowed, got %T", tmp.Interface()))
		}
		if tmp.IsNil() {
			ddest := reflect.MakeMap(tmp.Type())
			reflect.ValueOf(dest).Elem().Set(ddest)
		}
		return d.extractMap(dest, ctx)

	case reflect.Struct:
		return d.extractStruct(dest, ctx)

	case reflect.Ptr:
		if tmp.Elem().Kind() == reflect.Struct {
			realDest := reflect.New(tmp.Elem().Elem().Type())
			ggg := realDest.Addr().Interface()
			err := d.extractStruct(ggg, ctx)
			if err == nil {
				reflect.ValueOf(dest).Set(realDest)
			} else {
				reflect.ValueOf(dest).Set(reflect.Zero(reflect.TypeOf(dest)))
			}
			return err
		}
		fallthrough

	default:
		panic(fmt.Errorf("type %T cannot be a target for decoding", dest))
	}
	return nil
}

func decode(r io.Reader, dest interface{}, ctx Context) error {
	var tmp struct {
		Dest Manual `mad:".*"`
	}
	if v, ok := dest.(Manual); ok {
		tmp.Dest = v
		err := decode(r, &tmp, ctx)
		v = tmp.Dest
		return err
	}

	decoder, err := NewDecoder(r)
	if err != nil {
		return err
	}
	if err := decoder.Decode(dest, ctx); err != nil {
		return err
	}
	if decoder.tokens.Next() {
		token := decoder.token()
		lin, col := token.Start()
		return LocatedError{
			Lin: lin,
			Col: col,
			Err: fmt.Errorf("syntax error"),
		}
	}
	return nil
}

// Unmarshal unmarshals data from input with context ctx provided into a dest object. Return error != nil on error.
// err may be of type mad.LocatedError what provides positional information in (line, column) pair as well as error
// message.
func Unmarshal(input []byte, dest interface{}, ctx Context) error {
	return decode(bytes.NewReader(input), dest, ctx)
}

// UnmarshalFile unmarshals data right from the input file and returns error message in the form of <file>:<lin>:<col>: <msg>
func UnmarshalFile(fileName string, dest interface{}, ctx Context) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			message.Error(err)
		}
	}()
	err = decode(file, dest, ctx)
	if err == nil {
		return nil
	} else if _, ok := err.(LocatedError); ok {
		return fmt.Errorf("%s:%s", fileName, err)
	}
	return err
}
