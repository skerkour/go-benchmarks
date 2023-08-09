package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/skerkour/golibs/cpuinfo"
	"golang.org/x/sys/cpu"
)

func main() {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("-- SYSTEM INFO")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Print("\n")

	fmt.Println("Date:", time.Now().UTC().Format("2006-01-02"))
	fmt.Print("\n")

	fmt.Println("Go version:", runtime.Version())
	fmt.Print("\n")

	fmt.Println("CPU:")
	fmt.Println("- arch:", runtime.GOARCH)
	fmt.Println("- physical cores:", cpuinfo.CPU.PhysicalCores)
	fmt.Println("- logical cores:", cpuinfo.CPU.LogicalCores)
	fmt.Print("\n")

	fmt.Println("CPU features:")
	// more details here: https://pkg.go.dev/internal/cpu
	// and here: https://github.com/golang/go/blob/master/src/go/build/syslist.go

	if runtime.GOARCH == "amd64" {
		fmt.Println("- AVX:", cpu.X86.HasAVX)
		fmt.Println("- AVX2:", cpu.X86.HasAVX2)
		fmt.Println("- AVX512:", cpu.X86.HasAVX512)
		fmt.Println("- SSE:", cpuinfo.CPU.Supports(cpuinfo.SSE))
		fmt.Println("- SSE2:", cpuinfo.CPU.Supports(cpuinfo.SSE2))
		fmt.Println("- AES:", cpuinfo.CPU.Supports(cpuinfo.AESNI))
		fmt.Println("- SHA1:", cpuinfo.CPU.Supports(cpuinfo.SHA1))
		fmt.Println("- SHA2:", cpuinfo.CPU.Supports(cpuinfo.SHA2))
		fmt.Println("- SHA512:", cpuinfo.CPU.Supports(cpuinfo.SHA512))
		fmt.Println("- CRC32:", cpuinfo.CPU.Supports(cpuinfo.CRC32))
		fmt.Println("- ATOMICS:", cpuinfo.CPU.Supports(cpuinfo.ATOMICS))
	} else if runtime.GOARCH == "arm64" {
		// fmt.Println("- NEON:", cpu.ARM.HasNEON)
		fmt.Println("- SVE:", cpu.ARM64.HasSVE)
		fmt.Println("- AES:", cpu.ARM64.HasAES)
		fmt.Println("- SHA1:", cpu.ARM64.HasSHA1)
		fmt.Println("- SHA2:", cpu.ARM64.HasSHA2)
		fmt.Println("- SHA512:", cpu.ARM64.HasSHA512)
		fmt.Println("- SHA3:", cpu.ARM64.HasSHA3)
		fmt.Println("- CRC32:", cpuinfo.CPU.Supports(cpuinfo.CRC32))
		fmt.Println("- ATOMICS:", cpuinfo.CPU.Supports(cpuinfo.ATOMICS))
	}
	fmt.Print("\n")

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Print("\n")
}
