package mad

/*
All type items here represents allowed items for `raw` fenced code block
*/

///go:generate antlr4 -no-visitor -listener -o rawparser -Dlanguage=Go Raw.g4

// Integer integer number
type Integer struct {
	Location
	Value string
	Real  int64
}

// Unsigned unsigned integer number
type Unsigned struct {
	Location
	Value string
	Real  uint64
}

// Float represents floating point number
type Float struct {
	Location
	Value string
	Real  float64
}

// Boolean represents
type Boolean struct {
	Location
	Value string
	Real  bool
}
