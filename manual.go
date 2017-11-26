package mad

// Manual types command their processing themselves
type Manual interface {
	Decode(dest interface{}, header String, d *Decoder, ctx Context) (Manual, error)

	// Required must be static function, i.e. must be able to execute on zero instance of a type implementing Manual
	Required() bool
}
