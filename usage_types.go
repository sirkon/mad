package mad

// Decodable is a type that can decode itself using source decoder d, optionally using context ctx
// And nullify is a handle for setting the object into nil which is the case for optionals
type Decodable interface {
	Decode(d *Decoder, ctx Context) error
}

// Comment is for comment decoding
type Comment string

// Decode decodes comment from source decoder
func (c *Comment) Decode(d *Decoder, ctx Context) error {
	return d.extractComment(c)
}

// Code is for code decoding.
type Code struct {
	Syntax string
	Code   string
}

// Decode ...
func (c *Code) Decode(d *Decoder, ctx Context) error {
	return d.extractCode(c, ctx)
}

// CommentCode is for code preceded with comment
type CommentCode struct {
	Comment Comment
	Code    Code
}

// Decode ...
func (c *CommentCode) Decode(d *Decoder, ctx Context) error {
	if err := c.Comment.Decode(d, ctx); err != nil {
		return err
	}
	if err := c.Code.Decode(d, ctx); err != nil {
		return err
	}
	return nil
}

// CodeComment for code prolonged by comment
type CodeComment struct {
	Code    Code
	Comment Comment
}

// Decode ...
func (c *CodeComment) Decode(d *Decoder, ctx Context) error {
	if err := c.Code.Decode(d, ctx); err != nil {
		return err
	}
	if err := c.Comment.Decode(d, ctx); err != nil {
		return err
	}
	return nil
}
