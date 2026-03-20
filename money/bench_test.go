package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
)

var (
	benchmarkMoney Money
	benchmarkErr   error
)

// BenchmarkNewMoneyFloat measures exact float-to-money construction.
func BenchmarkNewMoneyFloat(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkMoney, benchmarkErr = NewMoneyFloat("USD", 0.85)
	}
}

// BenchmarkMoneyMulRat measures scalar multiplication through Rat.
func BenchmarkMoneyMulRat(b *testing.B) {
	value := NewMoneyInt("USD", 100)
	scalar := zerorat.New(5, 2)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := value
		benchmarkErr = m.MulRat(scalar)
		benchmarkMoney = m
	}
}

// BenchmarkMoneyDivRat measures scalar division through Rat.
func BenchmarkMoneyDivRat(b *testing.B) {
	value := NewMoneyInt("USD", 100)
	scalar := zerorat.New(5, 2)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m := value
		benchmarkErr = m.DivRat(scalar)
		benchmarkMoney = m
	}
}
