package mad

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	ctx := NewContext()

	ctx = ctx.New()
	ctx.Set("syntax", "go")
	val, ok := ctx.Get("syntax")
	require.True(t, ok)
	require.Equal(t, "go", val)

	ctx = ctx.New()
	ctx1 := ctx
	ctx.Set("value", true)
	val, ok = ctx.Get("value")
	require.True(t, ok)
	require.Equal(t, val, true)
	val, ok = ctx.Get("syntax")
	require.True(t, ok)
	require.Equal(t, "go", val)

	ctx = ctx.New()
	ctx.Set("syntax", "sql")
	val, ok = ctx.Get("syntax")
	require.True(t, ok)
	require.Equal(t, "sql", val)
	val, ok = ctx.Get("no such key")
	require.False(t, ok)

	val, ok = ctx1.Get("syntax")
	require.True(t, ok)
	require.Equal(t, "go", val)

}
