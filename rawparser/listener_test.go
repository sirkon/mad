package rawparser

import (
	"testing"

	"strings"

	"github.com/sirkon/mad/testdata"
	"github.com/stretchr/testify/require"
)

type item struct {
	name  string
	value string
	lin   int
	col   int
	xcol  int
}

type tokstor struct {
	data []item
}

func (t *tokstor) append(name string, lin, col, xcol int, value string) {
	t.data = append(t.data, item{
		name:  name,
		value: value,
		lin:   lin,
		col:   col,
		xcol:  xcol,
	})
}

func (t *tokstor) Header(lin, col, xcol int, value string) {
	t.append("header", lin, col, xcol, value)
}

func (t *tokstor) ValueNumber(lin, col, xcol int, value string) {
	t.append("number", lin, col, xcol, value)
}

func (t *tokstor) ValueString(lin, col, xcol int, value string) {
	t.append("string", lin, col, xcol, value)
}

func (t *tokstor) Boolean(lin, col, xcol int, value string) {
	t.append("bool", lin, col, xcol, value)
}

func TestListener(t *testing.T) {
	data, err := testdata.Asset("rawsection")
	if err != nil {
		t.Fatal(err)
	}
	ts := &tokstor{}
	errors := Parse(0, 0, string(data), ts)
	if len(errors) > 0 {
		res := []string{}
		for _, err := range errors {
			res = append(res, err.Error())
		}
		t.Errorf("\n%s", strings.Join(res, "\n"))
	}
	require.Equal(t,
		[]item{
			{
				name:  "header",
				value: "a",
				lin:   0,
				col:   0,
				xcol:  1,
			},
			{
				name:  "number",
				value: "1",
				lin:   0,
				col:   4,
				xcol:  5,
			},
			{
				name:  "header",
				value: "b",
				lin:   1,
				col:   0,
				xcol:  1,
			},
			{
				name:  "string",
				value: "be",
				lin:   1,
				col:   4,
				xcol:  8,
			},
			{
				name:  "header",
				value: "c",
				lin:   2,
				col:   1,
				xcol:  2,
			},
			{
				name:  "bool",
				value: "true",
				lin:   2,
				col:   5,
				xcol:  9,
			},
			{
				name:  "header",
				value: "d",
				lin:   3,
				col:   0,
				xcol:  1,
			},
			{
				name:  "string",
				value: "1kb",
				lin:   3,
				col:   6,
				xcol:  9,
			},
			{
				name:  "header",
				value: "e",
				lin:   4,
				col:   0,
				xcol:  1,
			},
			{
				name:  "number",
				value: "2.0",
				lin:   4,
				col:   2,
				xcol:  5,
			},
		}, ts.data)
}
