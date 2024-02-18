// from https://github.com/learning-go-book-2e/ch06/blob/main/sample_code/pointer_perf/perf_test.go
package pointers

import "testing"

var x byte

//go:noinline
func BenchmarkPointer_10_In(b *testing.B) {
	var v v10
	for i := 0; i < b.N; i++ {
		readMiddle10p(&v)
	}
}

//go:noinline
func BenchmarkValue_10_In(b *testing.B) {
	var v v10
	for i := 0; i < b.N; i++ {
		readMiddle10(v)
	}
}

//go:noinline
func BenchmarkPointer_10_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get10p()
		x = v.data[5]
	}
}

//go:noinline
func BenchmarkValue_10_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get10v()
		x = v.data[5]
	}
}

type v10 struct {
	data [10]byte
}

//go:noinline
func readMiddle10(v v10) {
	x = v.data[5]
}

//go:noinline
func readMiddle10p(v *v10) {
	x = v.data[5]
}

//go:noinline
func get10p() *v10 {
	var v v10
	return &v
}

//go:noinline
func get10v() v10 {
	var v v10
	return v
}

//go:noinline
func BenchmarkPointer_100_In(b *testing.B) {
	var v v100
	for i := 0; i < b.N; i++ {
		readMiddle100p(&v)
	}
}

//go:noinline
func BenchmarkValue_100_In(b *testing.B) {
	var v v100
	for i := 0; i < b.N; i++ {
		readMiddle100(v)
	}
}

//go:noinline
func BenchmarkPointer_100_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get100p()
		x = v.data[50]
	}
}

//go:noinline
func BenchmarkValue_100_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get100v()
		x = v.data[50]
	}
}

type v100 struct {
	data [100]byte
}

//go:noinline
func readMiddle100(v v100) {
	x = v.data[50]
}

//go:noinline
func readMiddle100p(v *v100) {
	x = v.data[50]
}

//go:noinline
func get100p() *v100 {
	var v v100
	return &v
}

//go:noinline
func get100v() v100 {
	var v v100
	return v
}

//go:noinline
func BenchmarkPointer_1_000_In(b *testing.B) {
	var v v1_000
	for i := 0; i < b.N; i++ {
		readMiddle1_000p(&v)
	}
}

//go:noinline
func BenchmarkValue_1_000_In(b *testing.B) {
	var v v1_000
	for i := 0; i < b.N; i++ {
		readMiddle1_000(v)
	}
}

//go:noinline
func BenchmarkPointer_1_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get1_000p()
		x = v.data[500]
	}
}

//go:noinline
func BenchmarkValue_1_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get1_000v()
		x = v.data[500]
	}
}

type v1_000 struct {
	data [1_000]byte
}

//go:noinline
func readMiddle1_000(v v1_000) {
	x = v.data[500]
}

//go:noinline
func readMiddle1_000p(v *v1_000) {
	x = v.data[500]
}

//go:noinline
func get1_000p() *v1_000 {
	var v v1_000
	return &v
}

//go:noinline
func get1_000v() v1_000 {
	var v v1_000
	return v
}

//go:noinline
func BenchmarkPointer_100_000_In(b *testing.B) {
	var v v100_000
	for i := 0; i < b.N; i++ {
		readMiddle100_000p(&v)
	}
}

//go:noinline
func BenchmarkValue_100_000_In(b *testing.B) {
	var v v100_000
	for i := 0; i < b.N; i++ {
		readMiddle100_000(v)
	}
}

//go:noinline
func BenchmarkPointer_100_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get100_000p()
		x = v.data[50_000]
	}
}

//go:noinline
func BenchmarkValue_100_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get100_000v()
		x = v.data[50_000]
	}
}

type v100_000 struct {
	data [100_000]byte
}

//go:noinline
func readMiddle100_000(v v100_000) {
	x = v.data[50_000]
}

//go:noinline
func readMiddle100_000p(v *v100_000) {
	x = v.data[50_000]
}

//go:noinline
func get100_000p() *v100_000 {
	var v v100_000
	return &v
}

//go:noinline
func get100_000v() v100_000 {
	var v v100_000
	return v
}

//go:noinline
func BenchmarkPointer_1_000_000_In(b *testing.B) {
	var v v1_000_000
	for i := 0; i < b.N; i++ {
		readMiddle1_000_000p(&v)
	}
}

//go:noinline
func BenchmarkValue_1_000_000_In(b *testing.B) {
	var v v1_000_000
	for i := 0; i < b.N; i++ {
		readMiddle1_000_000(v)
	}
}

//go:noinline
func BenchmarkPointer_1_000_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get1_000_000p()
		x = v.data[500_000]
	}
}

//go:noinline
func BenchmarkValue_1_000_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get1_000_000v()
		x = v.data[500_000]
	}
}

type v1_000_000 struct {
	data [1_000_000]byte
}

//go:noinline
func readMiddle1_000_000(v v1_000_000) {
	x = v.data[500_000]
}

//go:noinline
func readMiddle1_000_000p(v *v1_000_000) {
	x = v.data[500_000]
}

//go:noinline
func get1_000_000p() *v1_000_000 {
	var v v1_000_000
	return &v
}

//go:noinline
func get1_000_000v() v1_000_000 {
	var v v1_000_000
	return v
}

//go:noinline
func BenchmarkPointer_10_000_000_In(b *testing.B) {
	var v v10_000_000
	for i := 0; i < b.N; i++ {
		readMiddle10_000_000p(&v)
	}
}

//go:noinline
func BenchmarkValue_10_000_000_In(b *testing.B) {
	var v v10_000_000
	for i := 0; i < b.N; i++ {
		readMiddle10_000_000(v)
	}
}

//go:noinline
func BenchmarkPointer_10_000_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get10_000_000p()
		x = v.data[5_000_000]
	}
}

//go:noinline
func BenchmarkValue_10_000_000_Out(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v := get10_000_000v()
		x = v.data[5_000_000]
	}
}

type v10_000_000 struct {
	data [10_000_000]byte
}

//go:noinline
func readMiddle10_000_000(v v10_000_000) {
	x = v.data[5_000_000]
}

//go:noinline
func readMiddle10_000_000p(v *v10_000_000) {
	x = v.data[5_000_000]
}

//go:noinline
func get10_000_000p() *v10_000_000 {
	var v v10_000_000
	return &v
}

//go:noinline
func get10_000_000v() v10_000_000 {
	var v v10_000_000
	return v
}
