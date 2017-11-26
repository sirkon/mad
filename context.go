package mad

import "fmt"

// Context context provides generational generic storage
type Context struct {
	data []map[string]interface{}
}

// NewContext ...
func NewContext() Context {
	res := Context{}
	res.data = append(res.data, map[string]interface{}{})
	return res
}

// New creates new context generation
func (c Context) New() Context {
	data := c.data
	data = append(data, map[string]interface{}{})
	return Context{
		data: data,
	}
}

// Get gets a value of key from context. Returns value and ok = true in case if a key exists in one of generations and
// value = nil, ok = false otherwise
func (c Context) Get(key string) (value interface{}, ok bool) {
	for i := len(c.data) - 1; i >= 0; i-- {
		value, ok = c.data[i][key]
		if ok {
			return
		}
	}
	return
}

// MustGet gets a value from context using the Get. It replicates the Get except it raises the panic if no key has been
// found
func (c Context) MustGet(key string) interface{} {
	res, ok := c.Get(key)
	if !ok {
		panic(fmt.Errorf("key `%s` not found in a context", key))
	}
	return res
}

// Set sets a binds a value to a key in a last generation of context
func (c Context) Set(key string, value interface{}) {
	c.data[len(c.data)-1][key] = value
}

// GetString gets a string of a key if it exists. Raise assertion panic if key exists but not a string. Returns
// d(efault) when the key is absent
func (c Context) GetString(key string, d string) string {
	res, ok := c.Get(key)
	if !ok {
		return d
	}
	return res.(string)
}

// Steady sets the value into the first level
func (c Context) Steady(key string, value interface{}) {
	c.data[0][key] = value
}
