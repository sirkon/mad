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
	if err := d.Decode(&destUint, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, uint(1), destUint)

	// pass another token to read integer
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destInt int16
	if err := d.Decode(&destInt, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, int16(-1), destInt)

	// pass another token to read float
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destFloat float32
	if err := d.Decode(&destFloat, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, float32(1.12), destFloat)

	// pass another token to read string
	require.True(t, d.tokens.Next())
	d.tokens.Confirm()
	var destString string
	if err := d.Decode(&destString, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "12", destString)

	// now read comment if possible
	var cmt *Comment
	if err := d.Decode(cmt, nil); err == nil {
		t.Fatal("should be error")
	}
	if err := d.Decode(&cmt, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, (*Comment)(nil), cmt)

	// now read sql code
	sqlCode := Code{}
	if err := d.Decode(&sqlCode, "sql"); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Code{
		Syntax: "sql",
		Code:   "SELECT * FROM table\n",
	}, sqlCode)

	// now read go or gohtml code
	goCode := &Code{}
	if err := d.Decode(goCode, "go, gohtml"); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &Code{
		Syntax: "go",
		Code:   "package main\nfunc main() {\n    panic(\"LOL\")\n}\n",
	}, goCode)

	// read comment
	tmp := Comment("")
	cmt = &tmp
	if err := d.Decode(&cmt, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Comment("this is just a random text\n"), *cmt)
	if err := d.Decode(&cmt, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, (*Comment)(nil), cmt)

	// read toml code
	tomlCode := &Code{}
	if err := d.Decode(tomlCode, "toml yaml json xml"); err != nil {
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

	if err := d.Decode(&dest, "sql"); err != nil {
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

	if err := d.Decode(&dest, "sql"); err != nil {
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

	if err := d.Decode(&dest, "sql"); err != nil {
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
	if err := d.Decode(&dest, "yaml"); err != nil {
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
