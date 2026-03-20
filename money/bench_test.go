package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
)

var (
	benchmarkMoney Money
	errBenchmark   error
)

// BenchmarkNewMoneyFloat measures exact float-to-money construction.
func BenchmarkNewMoneyFloat(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkMoney, errBenchmark = NewMoneyFloat("USD", 0.85)
	}
}

// BenchmarkMoneyMulRat measures scalar multiplication through Rat.
func BenchmarkMoneyMulRat(b *testing.B) {
	value := NewMoneyInt("USD", 100)
	scalar := zerorat.New(5, 2)

	b.ReportAllocs()
	for b.Loop() {
		m := value
		errBenchmark = m.MulRat(scalar)
		benchmarkMoney = m
	}
}

// BenchmarkMoneyDivRat measures scalar division through Rat.
func BenchmarkMoneyDivRat(b *testing.B) {
	value := NewMoneyInt("USD", 100)
	scalar := zerorat.New(5, 2)

	b.ReportAllocs()
	for b.Loop() {
		m := value
		errBenchmark = m.DivRat(scalar)
		benchmarkMoney = m
	}
}
