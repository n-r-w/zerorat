package zerorat

import (
	"math/big"
	"testing"
)

// BenchmarkZeroRat_Construction benchmarks ZeroRat construction
func BenchmarkZeroRat_Construction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(int64(i%1000), uint64((i%999)+1))
	}
}

// BenchmarkBigRat_Construction benchmarks big.Rat construction
func BenchmarkBigRat_Construction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := big.NewRat(int64(i%1000), int64((i%999)+1))
		_ = r
	}
}

// BenchmarkZeroRat_Add benchmarks ZeroRat addition
func BenchmarkZeroRat_Add(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(1, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := r1
		r.Add(r2)
		_ = r
	}
}

// BenchmarkBigRat_Add benchmarks big.Rat addition
func BenchmarkBigRat_Add(b *testing.B) {
	r1 := big.NewRat(3, 4)
	r2 := big.NewRat(1, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := new(big.Rat).Set(r1)
		r.Add(r, r2)
		_ = r
	}
}

// BenchmarkZeroRat_AddImmutable benchmarks ZeroRat immutable addition
func BenchmarkZeroRat_AddImmutable(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(1, 3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := r1.Added(r2)
		_ = r
	}
}

// BenchmarkZeroRat_Mul benchmarks ZeroRat multiplication
func BenchmarkZeroRat_Mul(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(5, 7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := r1
		r.Mul(r2)
		_ = r
	}
}

// BenchmarkBigRat_Mul benchmarks big.Rat multiplication
func BenchmarkBigRat_Mul(b *testing.B) {
	r1 := big.NewRat(3, 4)
	r2 := big.NewRat(5, 7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := new(big.Rat).Set(r1)
		r.Mul(r, r2)
		_ = r
	}
}

// BenchmarkZeroRat_Div benchmarks ZeroRat division
func BenchmarkZeroRat_Div(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(5, 7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := r1
		r.Div(r2)
		_ = r
	}
}

// BenchmarkBigRat_Div benchmarks big.Rat division
func BenchmarkBigRat_Div(b *testing.B) {
	r1 := big.NewRat(3, 4)
	r2 := big.NewRat(5, 7)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := new(big.Rat).Set(r1)
		r.Quo(r, r2)
		_ = r
	}
}

// BenchmarkZeroRat_Equal benchmarks ZeroRat equality comparison
func BenchmarkZeroRat_Equal(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(6, 8) // equivalent fraction
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := r1.Equal(r2)
		_ = result
	}
}

// BenchmarkBigRat_Equal benchmarks big.Rat equality comparison
func BenchmarkBigRat_Equal(b *testing.B) {
	r1 := big.NewRat(3, 4)
	r2 := big.NewRat(6, 8) // equivalent fraction
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := r1.Cmp(r2) == 0
		_ = result
	}
}

// BenchmarkZeroRat_String benchmarks ZeroRat string conversion
func BenchmarkZeroRat_String(b *testing.B) {
	r := New(355, 113) // approximation of π
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := r.String()
		_ = s
	}
}

// BenchmarkBigRat_String benchmarks big.Rat string conversion
func BenchmarkBigRat_String(b *testing.B) {
	r := big.NewRat(355, 113) // approximation of π
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := r.String()
		_ = s
	}
}

// BenchmarkZeroRat_ComplexExpression benchmarks complex arithmetic expression
func BenchmarkZeroRat_ComplexExpression(b *testing.B) {
	// Compute: (3/4 + 1/3) * (5/7 - 2/9) / (11/13)
	a := New(3, 4)
	b1 := New(1, 3)
	c := New(5, 7)
	d := New(2, 9)
	e := New(11, 13)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Using immutable operations to avoid modifying original values
		temp1 := a.Added(b1)             // 3/4 + 1/3
		temp2 := c.Subtracted(d)         // 5/7 - 2/9
		temp3 := temp1.Multiplied(temp2) // (3/4 + 1/3) * (5/7 - 2/9)
		result := temp3.Divided(e)       // ... / (11/13)
		_ = result
	}
}

// BenchmarkBigRat_ComplexExpression benchmarks complex arithmetic expression with big.Rat
func BenchmarkBigRat_ComplexExpression(b *testing.B) {
	// Compute: (3/4 + 1/3) * (5/7 - 2/9) / (11/13)
	a := big.NewRat(3, 4)
	b1 := big.NewRat(1, 3)
	c := big.NewRat(5, 7)
	d := big.NewRat(2, 9)
	e := big.NewRat(11, 13)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp1 := new(big.Rat).Add(a, b1)        // 3/4 + 1/3
		temp2 := new(big.Rat).Sub(c, d)         // 5/7 - 2/9
		temp3 := new(big.Rat).Mul(temp1, temp2) // (3/4 + 1/3) * (5/7 - 2/9)
		result := new(big.Rat).Quo(temp3, e)    // ... / (11/13)
		_ = result
	}
}

// BenchmarkZeroRat_ArrayOperations benchmarks operations on arrays of rationals
func BenchmarkZeroRat_ArrayOperations(b *testing.B) {
	// Create array of rational numbers
	rationals := make([]Rat, 100)
	for i := range rationals {
		rationals[i] = New(int64(i+1), uint64(i+2))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := Zero()
		for j := range rationals {
			sum.Add(rationals[j])
		}
		_ = sum
	}
}

// BenchmarkBigRat_ArrayOperations benchmarks operations on arrays of big.Rat
func BenchmarkBigRat_ArrayOperations(b *testing.B) {
	// Create array of rational numbers
	rationals := make([]*big.Rat, 100)
	for i := range rationals {
		rationals[i] = big.NewRat(int64(i+1), int64(i+2))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum := big.NewRat(0, 1)
		for j := range rationals {
			sum.Add(sum, rationals[j])
		}
		_ = sum
	}
}

// Memory allocation benchmarks to verify zero-allocation claims

// BenchmarkZeroRat_MemoryAllocation tests memory allocation for ZeroRat operations
func BenchmarkZeroRat_MemoryAllocation(b *testing.B) {
	r1 := New(3, 4)
	r2 := New(5, 7)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test various operations that should not allocate
		r := r1
		r.Add(r2)
		r.Mul(r2)
		r.Sub(r2)
		r.Div(r2)

		// Test immutable operations (should allocate only for return value)
		result1 := r1.Added(r2)
		result2 := r1.Multiplied(r2)
		result3 := r1.Subtracted(r2)
		result4 := r1.Divided(r2)

		// Test comparisons (should not allocate)
		equal := r1.Equal(r2)
		less := r1.Less(r2)

		// Test utilities (should not allocate except String())
		sign := r1.Sign()
		isZero := r1.IsZero()

		// Prevent compiler optimization
		_ = r
		_ = result1
		_ = result2
		_ = result3
		_ = result4
		_ = equal
		_ = less
		_ = sign
		_ = isZero
	}
}

// BenchmarkBigRat_MemoryAllocation tests memory allocation for big.Rat operations
func BenchmarkBigRat_MemoryAllocation(b *testing.B) {
	r1 := big.NewRat(3, 4)
	r2 := big.NewRat(5, 7)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test various operations
		r := new(big.Rat).Set(r1)
		r.Add(r, r2)
		r.Mul(r, r2)
		r.Sub(r, r2)
		r.Quo(r, r2)

		// Test immutable-style operations
		result1 := new(big.Rat).Add(r1, r2)
		result2 := new(big.Rat).Mul(r1, r2)
		result3 := new(big.Rat).Sub(r1, r2)
		result4 := new(big.Rat).Quo(r1, r2)

		// Test comparisons
		equal := r1.Cmp(r2) == 0
		less := r1.Cmp(r2) < 0

		// Test utilities
		sign := r1.Sign()
		isZero := r1.Sign() == 0

		// Prevent compiler optimization
		_ = r
		_ = result1
		_ = result2
		_ = result3
		_ = result4
		_ = equal
		_ = less
		_ = sign
		_ = isZero
	}
}
