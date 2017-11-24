package tagparser

import "testing"

func TestExtractMad(t *testing.T) {
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
			if got := ExtractMad(tt.tag); got != tt.want {
				t.Errorf("extractMad() = %v, want %v", got, tt.want)
			}
		})
	}
}
