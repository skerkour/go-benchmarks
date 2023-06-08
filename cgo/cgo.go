package cgobench

//#include <unistd.h>
//void foo() { }
//void fooSleep() { sleep(100); }
import "C"

//go:noinline
func foo() {}

func CallCgo(n int) {
	for i := 0; i < n; i++ {
		C.foo()
	}
}

func CallGo(n int) {
	for i := 0; i < n; i++ {
		foo()
	}
}
