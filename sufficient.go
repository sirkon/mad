package mad

// Sufficient types command their processing themselves
type Sufficient interface {
	Decode(dest interface{}, header String, d *Decoder, ctx Context) (Sufficient, error)

	// Required must be static function, i.e. must be able to execute on zero instance of a type implementing Sufficient
	Required() bool
}
