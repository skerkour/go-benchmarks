package subtle

import "runtime"

// Wipe sets every byte in x to zero.
//
//go:noinline
func Wipe(x []byte) {
	// You don't have to twist the Go compiler's arm to keep it
	// from optimizing a piece of code. But, for insurance
	// reasons we mark Wipe as "noinline" so that the compiler
	// (hopefully) won't peer inside it and notice that x can be
	// DCEd.
	for i := range x {
		x[i] = 0
	}
	// Additionally, KeepAlive should (hopefully) nudge the
	// compiler away from DCEing the for-loop.
	runtime.KeepAlive(x)
}
