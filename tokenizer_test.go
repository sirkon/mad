package mdl

import (
	"reflect"
	"testing"

	"strings"

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
	tokenSimplest = func(lin int, offset int) Header {
		return Header{
			Location: Location{
				Lin:  lin,
				XLin: lin,
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

	tokenHarder = func(lin int) Header {
		return Header{
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

	tokenHardest = func(lin int) Header {
		return Header{
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
			token: tokenSimplest(0, 0),
		},
		{
			name:  "harder case",
			input: "## хеадер",
			scan:  true,
			token: tokenHarder(0),
		},
		{
			name:  "harder case",
			input: "#  хе д р    ",
			scan:  true,
			token: tokenHardest(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tz := NewTokenizer([]byte(tt.input))
			if tt.scan != tz.Next() {
				t.Errorf("tokenizer expected to scan on %s but it didn't", tt.input)
				return
			}
			require.Equal(t, tt.token, tz.Token())
			require.False(t, tz.Next())
			require.Empty(t, tz.Warnings())
		})
	}
}

func TestTokenizerSeveralLines(t *testing.T) {
	inputList := []string{"#header", " #header", "## хеадер", "#  хе д р    "}
	tz := NewTokenizer([]byte(strings.Join(inputList, "\n")))
	samples := []Header{
		tokenSimplest(0, 0),
		tokenSimplest(0, 1),
		tokenHarder(1),
		tokenHardest(2),
	}
	var i = 0
	for tz.Next() {
		require.Equal(t, samples[i], tz.Token())
		i++
	}
	require.Len(t, tz.Warnings(), 1)
}
