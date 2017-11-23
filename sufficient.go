package mad

// Sufficient types command their processing themselves
type Sufficient interface {
	New(dest interface{}, header String, d *Decoder, ctx interface{}) Sufficient

	// Required must be static function, i.e. must be able to execute on zero instance of a type implementing Sufficient
	Required() bool
}
