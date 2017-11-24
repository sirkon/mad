package mad

import (
	"reflect"
	"testing"

	"strings"

	"github.com/sirkon/mad/testdata"
	"github.com/stretchr/testify/require"
)

func Test_throwTrailingSpaces(t *testing.T) {
	type args struct {
		src []rune
	}
	tests := []struct {
		name    string
		args    args
		wantRes []rune
	}{
		{
			name: "degenerate case (nil)",
			args: args{
				src: nil,
			},
			wantRes: nil,
		},
		{
			name: "degenerate case (empty slice)",
			args: args{
				src: []rune{},
			},
			wantRes: []rune{},
		},
		{
			name: "degenerate case (all spaces)",
			args: args{
				src: []rune("    "),
			},
			wantRes: []rune{},
		},
		{
			name: "simplest case",
			args: args{
				src: []rune("12345"),
			},
			wantRes: []rune("12345"),
		},
		{
			name: "real case",
			args: args{
				src: []rune("12345  "),
			},
			wantRes: []rune("12345"),
		},
		{
			name: "regression case",
			args: args{
				src: []rune("12 12345  "),
			},
			wantRes: []rune("12 12345"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := throwTrailingSpaces(tt.args.src); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("throwTrailingSpaces() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_passHeadSpaces(t *testing.T) {
	type args struct {
		src []rune
	}
	tests := []struct {
		name       string
		args       args
		wantSpaces []rune
		wantRest   []rune
	}{
		{
			name: "degenerate case (nil)",
			args: args{
				src: nil,
			},
			wantSpaces: nil,
			wantRest:   nil,
		},
		{
			name: "degenerate case (empty)",
			args: args{
				src: []rune{},
			},
			wantSpaces: []rune{},
			wantRest:   []rune{},
		},
		{
			name: "degenerate case (all spaces)",
			args: args{
				src: []rune("    "),
			},
			wantSpaces: []rune("    "),
			wantRest:   []rune{},
		},
		{
			name: "simplest case",
			args: args{
				src: []rune("12345"),
			},
			wantSpaces: []rune{},
			wantRest:   []rune("12345"),
		},
		{
			name: "real world case",
			args: args{
				src: []rune("   12345"),
			},
			wantSpaces: []rune("   "),
			wantRest:   []rune("12345"),
		},
		{
			name: "probable regression case",
			args: args{
				src: []rune("   12345 12"),
			},
			wantSpaces: []rune("   "),
			wantRest:   []rune("12345 12"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpaces, gotRest := passHeadSpaces(tt.args.src)
			if !reflect.DeepEqual(gotSpaces, tt.wantSpaces) {
				t.Errorf("passHeadSpaces() gotSpaces = %v, want %v", gotSpaces, tt.wantSpaces)
			}
			if !reflect.DeepEqual(gotRest, tt.wantRest) {
				t.Errorf("passHeadSpaces() gotRest = %v, want %v", gotRest, tt.wantRest)
			}
		})
	}
}

var (
	tokenHeaderSimplest = func(lin int, offset int) header {
		return header{
			Location: Location{
				Lin:  lin,
				XLin: lin,
				Col:  offset,
				XCol: 7 + offset,
			},
			Content: String{
				Location: Location{
					Lin:  lin,
					XLin: lin,
					Col:  1 + offset,
					XCol: 7 + offset,
				},
				Value: "header",
			},
			Level: 1,
		}
	}

	tokenHeaderHarder = func(lin int) header {
		return header{
			Location: Location{
				Lin:  lin,
				XLin: lin,
				XCol: 9,
			},
			Content: String{
				Location: Location{
					Lin:  lin,
					XLin: lin,
					Col:  3,
					XCol: 9,
				},
				Value: "хеадер",
			},
			Level: 2,
		}
	}

	tokenHeaderHardest = func(lin int) header {
		return header{
			Location: Location{
				Lin:  lin,
				XLin: lin,
				XCol: 9,
			},
			Content: String{
				Location: Location{
					Lin:  lin,
					XLin: lin,
					Col:  3,
					XCol: 9,
				},
				Value: "хе д р",
			},
			Level: 1,
		}
	}
)

func TestTokeniserLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		scan  bool
		token interface{}
	}{
		{
			name:  "degenerate case (empty)",
			input: "",
			token: nil,
		},
		{
			name:  "simplest case",
			input: "#header",
			scan:  true,
			token: tokenHeaderSimplest(0, 0),
		},
		{
			name:  "harder case",
			input: "## хеадер",
			scan:  true,
			token: tokenHeaderHarder(0),
		},
		{
			name:  "harder case",
			input: "#  хе д р    ",
			scan:  true,
			token: tokenHeaderHardest(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz := newTokenizer([]byte(tt.input))
			require.Equal(t, tt.scan, tz.next())
			require.Equal(t, tt.token, tz.getToken())
			require.False(t, tz.next())
			require.Empty(t, tz.Warnings())
		})
	}
}

func TestTokenizerSeveralLines(t *testing.T) {
	inputList := []string{"#header", "", " #header", "## хеадер", "#  хе д р    "}
	tz := newTokenizer([]byte(strings.Join(inputList, "\n")))
	samples := []header{
		tokenHeaderSimplest(0, 0),
		tokenHeaderSimplest(2, 1),
		tokenHeaderHarder(3),
		tokenHeaderHardest(4),
	}
	var i = 0
	for tz.next() {
		token := tz.getToken()
		require.Equal(t, samples[i], token)
		i++
	}
	require.Len(t, tz.Warnings(), 1)
}

func TestTokenizerCodeBlockRealWorld(t *testing.T) {
	input := strings.Join([]string{
		"",
		"```sql",
		"SELECT 1, 2, 3 FROM a",
		"WHERE date > '2017-06-01'",
		"```",
	}, "\n")

	tz := newTokenizer([]byte(input))
	require.True(t, tz.next())
	require.Empty(t, tz.Warnings())
	require.Equal(t,
		code{
			Location: Location{
				Lin:  1,
				Col:  0,
				XLin: 4,
				XCol: 3,
			},
			Syntax: String{
				Location: Location{
					Lin:  1,
					Col:  3,
					XLin: 1,
					XCol: 6,
				},
				Value: "sql",
			},
			Content: String{
				Location: Location{
					Lin:  2,
					Col:  0,
					XLin: 3,
					XCol: 25,
				},
				Value: "SELECT 1, 2, 3 FROM a\nWHERE date > '2017-06-01'\n",
			},
		},
		tz.getToken())
}

func TestTokenizerRealWorld(t *testing.T) {
	sample := []string{
		"bugaga",
		"",
		"lol",
		"#header",
		"lol again",
		"again",
		"```sql",
		"SELECT 1, 2, 3 FROM a",
		"WHERE date > '2017-06-01'",
		"```",
	}
	input := strings.Join(sample, "\n")
	tz := newTokenizer([]byte(input))
	var tokens []interface{}
	for tz.next() {
		tokens = append(tokens, tz.getToken())
	}
	require.Equal(t, []interface{}{
		comment{
			Location: Location{
				Lin:  0,
				Col:  0,
				XLin: 2,
				XCol: 3,
			},
			Value: "bugaga\n\nlol\n",
		},
		tokenHeaderSimplest(3, 0),
		comment{
			Location: Location{
				Lin:  4,
				Col:  0,
				XLin: 5,
				XCol: 5,
			},
			Value: "lol again\nagain\n",
		},
		code{
			Location: Location{
				Lin:  6,
				Col:  0,
				XLin: 9,
				XCol: 3,
			},
			Syntax: String{
				Location: Location{
					Lin:  6,
					Col:  3,
					XLin: 6,
					XCol: 6,
				},
				Value: "sql",
			},
			Content: String{
				Location: Location{
					Lin:  7,
					Col:  0,
					XLin: 8,
					XCol: 25,
				},
				Value: "SELECT 1, 2, 3 FROM a\nWHERE date > '2017-06-01'\n",
			},
		},
	}, tokens)
}

func TestTabProcessing(t *testing.T) {
	tz := newTokenizer([]byte("#\theader"))
	require.True(t, tz.next())
	token := tz.getToken()
	require.False(t, tz.next())
	require.Equal(t, header{
		Location: Location{
			Col:  0,
			XCol: 11,
		},
		Level: 1,
		Content: String{
			Location: Location{
				Col:  5,
				XCol: 11,
			},
			Value: "header",
		},
	}, token)
}

func TestTokenStream(t *testing.T) {
	input := []string{
		"# header",
		"comment",
	}
	ttt := NewTokenizer([]byte(strings.Join(input, "\n")))
	var tokens []interface{}
	for i := 0; i < 3; i++ {
		require.True(t, ttt.Next())
		token := ttt.Token()
		tokens = append(tokens, token)
		if i == 0 {
			require.IsType(t, header{}, token)
		} else {
			require.Equal(t, tokens[0], token)
		}
	}
	for ttt.Next() {
		tokens = append(tokens, ttt.Token())
		ttt.Confirm()
	}
	require.Len(t, tokens, 5)
	require.Equal(t, tokens[0], tokens[3])
	require.IsType(t, comment{}, tokens[4])
}

func TestFullTokenSteam(t *testing.T) {
	data, err := testdata.Asset("test.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	var tokens []interface{}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	require.Equal(t, []interface{}{
		header{
			Location: Location{
				Lin:  1,
				Col:  0,
				XLin: 1,
				XCol: 1,
			},
			Level: 1,
			Content: String{
				Location: Location{
					Lin:  1,
					Col:  0,
					XLin: 1,
					XCol: 1,
				},
				Value: "a",
			},
		},
		String{
			Location: Location{
				Lin:  1,
				Col:  4,
				XLin: 1,
				XCol: 7,
			},
			Value: "1",
		},
		header{
			Location: Location{
				Lin:  4,
				Col:  0,
				XLin: 4,
				XCol: 8,
			},
			Level: 1,
			Content: String{
				Location: Location{
					Lin:  4,
					Col:  1,
					XLin: 4,
					XCol: 8,
				},
				Value: "request",
			},
		},
		header{
			Location: Location{
				Lin:  6,
				Col:  0,
				XLin: 6,
				XCol: 1,
			},
			Level: 2,
			Content: String{
				Location: Location{
					Lin:  6,
					Col:  0,
					XLin: 6,
					XCol: 1,
				},
				Value: "b",
			},
		},
		unsigned{
			Location: Location{
				Lin:  6,
				Col:  4,
				XLin: 6,
				XCol: 5,
			},
			Value: "1",
			Real:  uint64(1),
		},
	}, tokens)
}

func TestFullTokenRegressionSteam(t *testing.T) {
	data, err := testdata.Asset("test.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	require.True(t, t2.Next())
	t2.Confirm()
	require.True(t, t2.Next())
	require.Equal(t, String{
		Location: Location{
			Lin:  1,
			Col:  4,
			XLin: 1,
			XCol: 7,
		},
		Value: "1",
	}, t2.Token())
	t2.Confirm()
}

func TestFullTokenRegression2(t *testing.T) {
	data, err := testdata.Asset("regression.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	tokens := []interface{}{}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	for _, err := range t2.Err() {
		t.Error(err)
	}
	require.Equal(t, []interface{}{
		header{
			Location: Location{
				Lin:  1,
				Col:  0,
				XLin: 1,
				XCol: 2,
			},
			Content: String{
				Location: Location{
					Lin:  1,
					Col:  0,
					XLin: 1,
					XCol: 2,
				},
				Value: "na",
			},
			Level: 1,
		},
		integer{
			Location: Location{
				Lin:  1,
				Col:  5,
				XLin: 1,
				XCol: 7,
			},
			Value: "-1",
			Real:  -1,
		},
		header{
			Location: Location{
				Lin:  2,
				Col:  0,
				XLin: 2,
				XCol: 2,
			},
			Content: String{
				Location: Location{
					Lin:  2,
					Col:  0,
					XLin: 2,
					XCol: 2,
				},
				Value: "nb",
			},
			Level: 1,
		},
		float{
			Location: Location{
				Lin:  2,
				Col:  5,
				XLin: 2,
				XCol: 9,
			},
			Value: "12.0",
			Real:  12.0,
		},
	}, tokens)
}

func TestLevelNormalization(t *testing.T) {
	data, err := testdata.Asset("level_normalization.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	var tokens []interface{}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	require.Equal(t,
		[]interface{}{
			header{
				Location: Location{
					XCol: 6,
				},
				Level: 1,
				Content: String{
					Location: Location{
						Col:  4,
						XCol: 6,
					},
					Value: "h1",
				},
			},

			header{
				Location: Location{
					Lin:  1,
					Col:  0,
					XLin: 1,
					XCol: 8,
				},
				Level: 2,
				Content: String{
					Location: Location{
						Lin:  1,
						Col:  6,
						XLin: 1,
						XCol: 8,
					},
					Value: "h2",
				},
			},

			header{
				Location: Location{
					Lin:  2,
					Col:  0,
					XLin: 2,
					XCol: 7,
				},
				Level: 2,
				Content: String{
					Location: Location{
						Lin:  2,
						Col:  5,
						XLin: 2,
						XCol: 7,
					},
					Value: "h3",
				},
			},
		},
		tokens)
}

func TestRegression2(t *testing.T) {
	data, err := testdata.Asset("regression2.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	var tokens = []interface{}{}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	for _, err := range t2.Err() {
		t.Error(err)
	}
	require.Equal(t,
		[]interface{}{
			header{
				Location: Location{XCol: 7},
				Content:  String{Location: Location{Lin: 0, Col: 1, XLin: 0, XCol: 7}, Value: "field1"},
				Level:    1,
			},
			header{
				Location: Location{Lin: 1, Col: 0, XLin: 1, XCol: 9},
				Content:  String{Location: Location{Lin: 1, Col: 3, XLin: 1, XCol: 9}, Value: "field2"},
				Level:    2},
			code{
				Location: Location{Lin: 2, Col: 0, XLin: 4, XCol: 3},
				Syntax:   String{Location: Location{Lin: 2, Col: 3, XLin: 2, XCol: 6}, Value: "sql"},
				Content:  String{Location: Location{Lin: 3, Col: 0, XLin: 3, XCol: 6}, Value: "SELECT\n"}},
		}, tokens)
}

func TestRegression3(t *testing.T) {
	data, err := testdata.Asset("codecomment.md")
	if err != nil {
		t.Fatal(err)
	}
	t1 := NewTokenizer(data)
	t2 := NewFullTokenizer(t1)
	var tokens = []interface{}{}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	for t2.Next() {
		tokens = append(tokens, t2.Token())
		t2.Confirm()
	}
	require.Len(t, tokens, 2)
}
