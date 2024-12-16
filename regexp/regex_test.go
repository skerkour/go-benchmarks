package regex

import (
	"regexp"
	"strings"
	"testing"

	"github.com/skerkour/go-benchmarks/regexp/pcre"
)

var Sink bool

// func BenchmarkCompileRun(b *testing.B) {
//         for i := 0; i < b.N; i++ {
//                 rx := regexp.MustCompile(`[\w\.+-]+@[\w\.-]+\.[\w\.-]+`)
//                 Sink = rx.MatchString("123456789 foo@bar.etc")
//         }
// }

const URL = "https://www.example.com/something/16A1jfMtDjHsQNGmYHPucK82FXiPMpZZasZBi2HUBljRRHnZsk3CwuKPuERkNiKUa0Q3bQezPYXPyvEwx8danFfeGotki99ZaxQewP9VNveNwvaM60fz2W17JuxjhPM89edG3TyNvrw2RtKulDqw1vAuYYE1h6hCAHvGzfuApZg0pBfEK8op4N76GqV0hv4="

func BenchmarkRegexp(b *testing.B) {
	rx := regexp.MustCompile(`http(s?)://www\.example\.com/something/.*`)
	for i := 0; i < b.N; i++ {
		Sink = rx.MatchString(URL)
	}
}

func BenchmarkPCRE(b *testing.B) {
	rx := pcre.MustCompileJIT(`http(s?)://www\.example\.com/something/.*`, 0, pcre.CONFIG_JIT)
	for i := 0; i < b.N; i++ {
		Sink = rx.MatchStringWFlags(URL, 0)
	}
}

func BenchmarkHasPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sink = strings.HasPrefix(URL, "https://old.example.com/something/")
	}
}
