package mad

import (
	"strings"
	"unicode"
)

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

type tagListener interface {
	name(name string)
	value(value string)
}

type tagScanner struct {
	rest     string
	listener tagListener
}

func newTagScanner(tag string, tl tagListener) *tagScanner {
	return &tagScanner{
		rest:     tag,
		listener: tl,
	}
}

func identChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-'
}

func (ts *tagScanner) next() bool {
	if len(ts.rest) == 0 {
		return false
	}

	// passing first spaces
	rs := []rune(ts.rest)
	pos := -1
	for i, r := range rs {
		if !unicode.IsSpace(r) {
			pos = i
			break
		}
	}
	if pos < 0 {
		pos = len(rs)
	}
	rs = rs[pos:]
	ts.rest = string(rs)

	//
	if len(rs) == 0 {
		return false
	}
	if rs[0] == '"' {
		rs := rs[1:]
		pos := -1
		for i, r := range rs {
			if r == '"' && (i == 0 || rs[i-1] != '\\') {
				ts.listener.value(strings.Replace(string(rs[:i]), `\"`, `""`, -1))
				pos = i + 1
				break
			}
		}
		if pos < 0 {
			pos = len(rs)
			ts.rest = ""
			return false
		}
		ts.rest = string(rs[pos:])
		return true
	} else if identChar(rs[0]) {
		pos := -1
		for i, r := range rs {
			if identChar(r) {
				continue
			}
			pos = i
			break
		}
		if pos < 0 {
			ts.rest = ""
			ts.listener.name(string(rs))
			return false
		}
		ts.listener.name(string(rs[:pos]))
		rs = rs[pos:]
		ts.rest = string(rs)
		ts.rest = strings.TrimLeftFunc(ts.rest, unicode.IsSpace)
		if len(ts.rest) == 0 {
			return false
		}
		if ts.rest[0] == ':' {
			ts.rest = ts.rest[1:]
			return true
		}
		return true
	}
	ts.rest = ""
	return false
}

type listener struct {
	n string
	v string
}

func (l *listener) name(name string) {
	l.n = name
}

func (l *listener) value(value string) {
	if l.n == "mad" {
		l.v = value
	} else {
		l.n = ""
	}
}

func (l *listener) madValue() string {
	return l.v
}

// extractMad extracs mad tag content from tag string, as the stdlib implementation refused to work with regular
// expressions inside
func extractMad(tag string) string {
	l := &listener{}
	ts := newTagScanner(tag, l)
	for ts.next() {
	}
	return l.madValue()
}
