package swap

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"
)

func BenchmarkParallelAtomicPointerRead(b *testing.B) {
	var ptr unsafe.Pointer
	data := make(map[int]int)
	data[1] = 1
	atomic.StorePointer(&ptr, unsafe.Pointer(&data))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localData := (*map[int]int)(atomic.LoadPointer(&ptr))
			_ = (*localData)[1]
		}
	})
}

func BenchmarkParallelAtomicValueRead(b *testing.B) {
	var ptr atomic.Value
	data := make(map[int]int)
	data[1] = 1
	ptr.Store(data)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			localData := ptr.Load().(map[int]int)
			_ = localData[1]
		}
	})
}

func BenchmarkParallelRWMutexRead(b *testing.B) {
	var mutex sync.RWMutex
	data := make(map[int]int)
	data[1] = 1

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mutex.RLock()
			_ = data[1]
			mutex.RUnlock()
		}
	})
}
