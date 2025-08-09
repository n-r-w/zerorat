package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mulInt64ByUint64ToInt64_AllBranches(t *testing.T) {
	// zero cases
	v, ok := mulInt64ByUint64ToInt64(0, 123)
	assert.True(t, ok)
	assert.Equal(t, int64(0), v)

	v, ok = mulInt64ByUint64ToInt64(123, 0)
	assert.True(t, ok)
	assert.Equal(t, int64(0), v)

	// overflow via hi!=0
	_, ok = mulInt64ByUint64ToInt64(math.MinInt64, 2)
	assert.False(t, ok)

	// negative exact MinInt64
	v, ok = mulInt64ByUint64ToInt64(-1, uint64(math.MaxInt64)+1)
	assert.True(t, ok)
	assert.Equal(t, int64(math.MinInt64), v)

	// negative lo>limit -> false
	_, ok = mulInt64ByUint64ToInt64(-3, 1<<62)
	assert.False(t, ok)

	// positive lo>MaxInt64 -> false
	_, ok = mulInt64ByUint64ToInt64(3, 1<<62)
	assert.False(t, ok)

	// positive within range
	v, ok = mulInt64ByUint64ToInt64(7, 9)
	assert.True(t, ok)
	assert.Equal(t, int64(63), v)
}

func TestRat_Div_MinInt64NegationOverflow(t *testing.T) {
	// r = MinInt64/1, other = -1/1 -> newNum=MinInt64, then negation would overflow, expect invalid
	r := New(math.MinInt64, 1)
	other := New(-1, 1)
	r.Div(other)
	assert.True(t, r.IsInvalid(), "negating MinInt64 during division should result in invalid state")
}
