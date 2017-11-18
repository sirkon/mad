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
	sqlCode := Code{
		Syntax: "sql",
	}
	if err := d.Decode(&sqlCode, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Code{
		Syntax: "sql",
		Code:   "SELECT * FROM table\n",
	}, sqlCode)

	// now read go or gohtml code
	goCode := &Code{
		Syntax: "go, gohtml",
	}
	if err := d.Decode(&sqlCode, nil); err == nil {
		t.Fatal("should be error")
	}
	if err := d.Decode(&goCode, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &Code{
		Syntax: "go",
		Code:   "package main\nfunc main() {\n    panic(\"LOL\")\n}\n",
	}, goCode)

	// read comment
	tmp := Comment("")
	cmt = &tmp
	if err := d.Decode(cmt, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, Comment("this is just a random text\n"), *cmt)

	// read toml code
	tomlCode := &Code{
		Syntax: "toml yaml json xml",
	}
	if err := d.Decode(tomlCode, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, &Code{
		Syntax: "toml",
		Code:   `a = "1kb"` + "\n",
	}, tomlCode)
}
