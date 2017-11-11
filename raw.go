package mad

/*
All type items here represents allowed items for `raw` fenced code block
*/

///go:generate antlr4 -no-visitor -listener -o rawparser -Dlanguage=Go Raw.g4

// Integer integer number
type Integer struct {
	Location
	Value int64
}

// Unsigned unsigned integer number
type Unsigned struct {
	Location
	Value uint64
}

// Float represents floating point number
type Float struct {
	Location
	Value float64
}

// Boolean represents
type Boolean struct {
	Location
	Value bool
}
