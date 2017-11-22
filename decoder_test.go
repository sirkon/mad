package mad

import (
	"testing"

	"bytes"

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
			Syntax: "sql",
			Code:   "SELECT * FROM table\n",
		},
		{
			Syntax: "sql",
			Code:   "SELECT * FROM table2\n",
		},
		{
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
	ctx.Set("sytnax", "yaml")
	if err := d.Decode(&dest, ctx); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, map[string]Code{
		"key1": {
			Syntax: "yaml",
			Code:   "a: 1\n",
		},
		"key2": {
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
