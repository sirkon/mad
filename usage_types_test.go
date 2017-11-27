package mad

import (
	"testing"

	"bytes"

	"github.com/sirkon/mad/testdata"
	"github.com/stretchr/testify/require"
)

func TestCodeList_Decode(t *testing.T) {
	data, err := testdata.Asset("codelist.md")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewReader(data)
	decoder, err := NewDecoder(buf)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, decoder.tokens.Next())
	h := decoder.token().(header)
	require.Equal(t, "set trivial", h.Content.Value)
	decoder.tokens.Confirm()
	var dest1 CodeList
	if err := decoder.Decode(&dest1, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "python", dest1.dest[0].Syntax)
	require.Len(t, dest1.Codes(), 1)

	require.True(t, decoder.tokens.Next())
	h = decoder.token().(header)
	require.Equal(t, "set real", h.Content.Value)
	decoder.tokens.Confirm()
	var dest2 CodeList
	if err := decoder.Decode(&dest2, NewContext()); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "python", dest2.dest[0].Syntax)
	require.Equal(t, "sql", dest2.dest[1].Syntax)
	require.Equal(t, "SELECT * FROM table \n", dest2.dest[1].Code)
	require.Len(t, dest2.Codes(), 2)

	require.True(t, decoder.tokens.Next())
	h = decoder.token().(header)
	require.Equal(t, "set error", h.Content.Value)
	decoder.tokens.Confirm()
	var dest3 CodeList
	err = decoder.Decode(&dest3, NewContext())
	require.Error(t, err)
}
