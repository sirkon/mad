package mad

//go:generate antlr4 -package tagparser -o tagparser -listener -Dlanguage=Go Tag.g4
//go:generate antlr4 -package rawparser -o rawparser -listener -Dlanguage=Go Raw.g4
