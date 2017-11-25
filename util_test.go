package mad

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitter(t *testing.T) {
	type sample struct {
		name string
		tag  string
		res  []string
	}
	checks := []sample{
		{
			name: "easy",
			tag:  "tag",
			res:  []string{"tag"},
		},
		{
			name: "real",
			tag:  "aby,syntax=go",
			res:  []string{"aby", "syntax=go"},
		},
		{
			name: "hardest",
			tag:  `ab\,y,syntax=go`,
			res:  []string{"ab,y", "syntax=go"},
		},
	}
	for _, check := range checks {
		item := check
		t.Run(check.name, func(t *testing.T) {
			s := newSplitter(item.tag)
			res := []string{}
			for s.next() {
				res = append(res, s.text())
			}
			require.Equal(t, item.res, res)
		})
	}
}

func Test_keyVal(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantKey   string
		wantValue string
		wantOk    bool
	}{
		{
			name:   "degenerate case 1",
			input:  "",
			wantOk: false,
		},
		{
			name:      "degenerate case 2",
			input:     "=",
			wantKey:   "",
			wantValue: "",
			wantOk:    true,
		},
		{
			name:      "degenerate case 3",
			input:     "=12",
			wantKey:   "",
			wantValue: "12",
			wantOk:    true,
		},
		{
			name:      "degenerate case 4",
			input:     "12=",
			wantKey:   "12",
			wantValue: "",
			wantOk:    true,
		},
		{
			name:      "real case",
			input:     "a=bc",
			wantKey:   "a",
			wantValue: "bc",
			wantOk:    true,
		},
		{
			name:      "error case",
			input:     "abcdef",
			wantKey:   "",
			wantValue: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotValue, gotOk := keyVal(tt.input)
			if gotKey != tt.wantKey {
				t.Errorf("keyVal() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotValue != tt.wantValue {
				t.Errorf("keyVal() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("keyVal() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_extractMad(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "degenerate case 1",
		},
		{
			name: "degenerate case 2",
			tag:  `json:"value,omitempty"`,
			want: "",
		},
		{
			name: "easiest",
			tag:  `mad:"a"`,
			want: "a",
		},
		{
			name: "harder",
			tag:  `mad:"prepare,syntax=sql cql"`,
			want: "prepare,syntax=sql cql",
		},
		{
			name: "hardest",
			tag:  `mad:"status=\d+,syntax=go"`,
			want: `status=\d+,syntax=go`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMad(tt.tag); got != tt.want {
				t.Errorf("extractMad() = %v, want %v", got, tt.want)
			}
		})
	}
}
