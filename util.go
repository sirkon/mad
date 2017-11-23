package mad

import "strings"

type splitter struct {
	rest     string
	cur      string
	notFirst bool
}

func newSplitter(data string) *splitter {
	return &splitter{
		rest: data,
	}
}

func (s *splitter) next() bool {
	if len(s.rest) == 0 && s.notFirst {
		return false
	}
	s.notFirst = true
	pos := -1
	for i := 0; i < len(s.rest); i++ {
		if s.rest[i] == ',' && (i == 0 || s.rest[i-1] != '\\') {
			pos = i
			break
		}
	}
	if pos >= 0 {
		s.cur = s.rest[:pos]
		s.rest = s.rest[pos+1:]
	} else {
		s.cur = s.rest
		s.rest = ""
	}
	return true
}

func (s *splitter) text() string {
	return strings.Replace(s.cur, `\,`, `,`, -1)
}

// keyVal
func keyVal(input string) (key string, value string, ok bool) {
	pos := strings.IndexByte(input, '=')
	if pos < 0 {
		return "", "", false
	}
	return input[:pos], input[pos+1:], true
}
