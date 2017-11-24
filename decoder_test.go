package mad

import (
	"bytes"
	"testing"

	"fmt"

	"strings"

	"strconv"

	"github.com/sirkon/mad/testdata"
	"github.com/stretchr/testify/require"
)

func TestDecoderScalar(t *testing.T) {
	data, err := testdata.Asset("scalar_decoder.md")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoder(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	// pass first token to read unsigned
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destUint uint
	if err := d.Decode(&destUint, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, uint(1), destUint)

	// pass another token to read integer
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destInt int16
	if err := d.Decode(&destInt, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, int16(-1), destInt)

	// pass another token to read float
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destFloat float32
	if err := d.Decode(&destFloat, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, float32(1.12), destFloat)

	// pass another token to read string
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destString string
	if err := d.Decode(&destString, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "12", destString)

	// now read comment if possible
	var cmt *Comment
	if err := d.Decode(cmt, NewContext()); err == nil {
		t.Fatal("should be error")
	}
	if err := d.Decode(&cmt, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, (*Comment)(nil), cmt)

	ctx := NewContext()
	ctx = ctx.New()
	// now read sql code
	sqlCode := Code{}
	ctx.Set("syntax", "sql")
	if err := d.Decode(&sqlCode, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Code{
		loc:    Location{Lin: 8, Col: 0, XLin: 8, XCol: 19},
		Syntax: "sql",
		Code:   "SELECT * FROM table\n",
	}, sqlCode)

	// now read go or gohtml code
	goCode := &Code{}
	ctx.Set("syntax", "go, gohtml")
	if err := d.Decode(goCode, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &Code{
		loc:    Location{Lin: 11, Col: 0, XLin: 14, XCol: 1},
		Syntax: "go",
		Code:   "package main\nfunc main() {\n    panic(\"LOL\")\n}\n",
	}, goCode)

	// read comment
	tmp := Comment("")
	cmt = &tmp
	if err := d.Decode(&cmt, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Comment("this is just a random text\n"), *cmt)
	if err := d.Decode(&cmt, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, (*Comment)(nil), cmt)

	// read toml code
	tomlCode := &Code{}
	ctx.Set("syntax", "toml yaml json xml")
	if err := d.Decode(tomlCode, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &Code{
		loc:    Location{Lin: 18, Col: 0, XLin: 18, XCol: 9},
		Syntax: "toml",
		Code:   `a = "1kb"` + "\n",
	}, tomlCode)
}

func TestCodeComment(t *testing.T) {
	data, err := testdata.Asset("codecomment.md")
	if err != nil {
		t.Fatal(err)
	}
	var dest CodeComment
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	ctx := NewContext().New()
	ctx.Set("syntax", "sql")
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, CodeComment{
		Code: Code{
			loc:    Location{Lin: 1, Col: 0, XLin: 1, XCol: 19},
			Syntax: "sql",
			Code:   "SELECT * FROM table\n",
		},
		Comment: "This was just a request\n",
	}, dest)
}

func TestCommentCode(t *testing.T) {
	data, err := testdata.Asset("commentcode.md")
	if err != nil {
		t.Fatal(err)
	}
	var dest CommentCode
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	ctx := NewContext().New()
	ctx.Set("syntax", "sql")
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, CommentCode{
		Code: Code{
			loc:    Location{Lin: 2, Col: 0, XLin: 2, XCol: 19},
			Syntax: "sql",
			Code:   "SELECT * FROM table\n",
		},
		Comment: "This will be a request\n",
	}, dest)
}

func TestCodeSlice(t *testing.T) {
	data, err := testdata.Asset("codearray.md")
	if err != nil {
		t.Fatal(err)
	}
	var dest []Code
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	ctx := NewContext().New()
	ctx.Set("syntax", "sql")
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}

	require.Equal(t, []Code{
		{
			loc:    Location{Lin: 1, XLin: 1, XCol: 19},
			Syntax: "sql",
			Code:   "SELECT * FROM table\n",
		},
		{
			loc:    Location{Lin: 5, XLin: 5, XCol: 20},
			Syntax: "sql",
			Code:   "SELECT * FROM table2\n",
		},
		{
			loc:    Location{Lin: 9, XLin: 9, XCol: 20},
			Syntax: "sql",
			Code:   "SELECT * FROM table3\n",
		},
	}, dest)
}

func TestMap(t *testing.T) {
	data, err := testdata.Asset("maps.md")
	if err != nil {
		t.Fatal(err)
	}
	var dest map[string]Code
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext().New()
	ctx.Set("syntax", "yaml")
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, map[string]Code{
		"key1": {
			loc: Location{
				Lin:  2,
				XLin: 2,
				XCol: 4,
			},
			Syntax: "yaml",
			Code:   "a: 1\n",
		},
		"key2": {
			loc: Location{
				Lin:  7,
				XLin: 7,
				XCol: 4,
			},
			Syntax: "yaml",
			Code:   "a: 2\n",
		},
	}, dest)
}

type identity string

func (i *identity) Unmarshal(data string) (err error) {
	*i = identity(data)
	return nil
}

func TestUnmarshaler(t *testing.T) {
	data, err := testdata.Asset("rawunmarshaler.md")
	if err != nil {
		t.Fatal(err)
	}
	res := map[string]string{}
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	d, err = NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext().New()
	for d.tokens.Next() {
		token := d.tokens.Token()
		require.IsType(t, header{}, token)
		key := token.(header).Content.Value
		d.tokens.Confirm()

		dest := identity("")
		if err := d.Decode(&dest, ctx); err != nil {
			t.Fatal(err)
		}
		res[key] = string(dest)
	}
	require.Equal(t, map[string]string{
		"a": "1",
		"b": "ID",
		"c": `1`,
		"d": "true",
		"e": "128kb",
	}, res)
}

func TestRegression(t *testing.T) {
	data, err := testdata.Asset("rawunmarshaler.md")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	d.passComment()
	require.True(t, d.tokens.Next())
	require.IsType(t, header{}, d.token())
}

func TestStructDecoding(t *testing.T) {
	data, err := testdata.Asset("rawunmarshaler.md")
	if err != nil {
		t.Fatal(err)
	}
	type tmp struct {
		A Code    `mad:"a,syntax=go"`
		B Comment `mad:"b"`
	}
	var dest tmp
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	if err := d.Decode(&dest, ctx); err != nil {
		require.Error(t, err)
	}
}

func TestStructEasy(t *testing.T) {
	type tmp struct {
		A Code    `mad:"a,syntax=go"`
		B Comment `mad:"b"`
	}
	var dest tmp
	data, err := testdata.Asset("struct_easy.md")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tmp{
		A: Code{
			loc: Location{
				Lin:  2,
				XLin: 6,
				XCol: 1,
			},
			Syntax: "go",
			Code: `package main

func main() {
    panic("error")
}
`,
		},
		B: Comment("just a text\n"),
	}, dest)
}

func TestStructReal(t *testing.T) {
	type nested struct {
		Prepare Code `mad:"prepare,syntax=sql cql"`
	}
	type tmp struct {
		A     int     `mad:"a"`
		B     string  `mad:"b"`
		C     float64 `mad:"c"`
		Query nested  `mad:"query"`
	}
	var dest tmp
	data, err := testdata.Asset("struct_real.md")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContext()
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tmp{
		A: 1,
		B: "2",
		C: 3.5,
		Query: nested{
			Prepare: Code{
				loc: Location{
					Lin:  9,
					Col:  0,
					XLin: 9,
					XCol: 37,
				},
				Syntax: "sql",
				Code:   "CREATE TABLE a AS SELECT * FROM table\n",
			},
		},
	}, dest)
}

type resp map[int]Comment

func (r *resp) Decode(dest interface{}, header String, d *Decoder, ctx Context) (Sufficient, error) {
	dd, ok := dest.(*resp)
	if !ok {
		return nil, fmt.Errorf("dest must be %T, got %T", r, dest)
	}
	var tmp resp
	if dd == nil || *dd == nil {
		tmp = resp(map[int]Comment{})
	} else {
		tmp = *dd
	}
	chunk := strings.Split(header.Value, "=")
	statLit := chunk[1]
	status64, err := strconv.ParseInt(statLit, 10, 64)
	if err != nil {
		return nil, err
	}
	status := int(status64)
	var cmt Comment
	if err := d.Decode(&cmt, ctx); err != nil {
		return nil, err
	}
	tmp[status] = cmt
	return &tmp, nil
}

func (r *resp) Required() bool {
	return true
}

func TestStatuses(t *testing.T) {
	data, err := testdata.Asset("statuses.md")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewDecoder(bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	var dest struct {
		S *resp `mad:"status=\d+"`
	}
	ctx := NewContext()
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t,
		resp{
			200: Comment("к успеху пришёл\n"),
			404: Comment("нихуя не нашёл\n"),
			500: Comment("ебать пиздец\n"),
		},
		*dest.S,
	)
}
