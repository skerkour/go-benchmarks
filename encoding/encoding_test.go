package encoding

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	akamenskybase58 "github.com/akamensky/base58"
	mrtronbase58 "github.com/mr-tron/base58"
	"github.com/skerkour/go-benchmarks/utils"
	golibsbase32 "github.com/skerkour/golibs/base32"
)

type Encoder interface {
	Encode(data []byte)
	// Decode(str string)
}

func BenchmarkEncode(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		64 * 1024,
		100 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkEncode(size, "std_hex", stdHex{}, b)
		benchmarkEncode(size, "std_base64", stdBase64{}, b)
		benchmarkEncode(size, "std_base32", stdBase32{}, b)
		benchmarkEncode(size, "golibs_base32", golibsBase32{}, b)
		benchmarkEncode(size, "akamensky_base58", akamenskyBase58{}, b)
		benchmarkEncode(size, "mr-tron_base58", mrTronBase58{}, b)
	}
}

func benchmarkEncode[E Encoder](size int64, algorithm string, encoder E, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			encoder.Encode(buf)
		}
	})
}

type stdHex struct{}

func (stdHex) Encode(data []byte) {
	hex.EncodeToString(data)
}

type stdBase64 struct{}

func (stdBase64) Encode(data []byte) {
	base64.StdEncoding.EncodeToString(data)
}

type akamenskyBase58 struct{}

func (akamenskyBase58) Encode(data []byte) {
	akamenskybase58.Encode(data)
}

type mrTronBase58 struct{}

func (mrTronBase58) Encode(data []byte) {
	mrtronbase58.Encode(data)
}

type stdBase32 struct{}

func (stdBase32) Encode(data []byte) {
	base32.StdEncoding.EncodeToString(data)
}

type golibsBase32 struct{}

func (golibsBase32) Encode(data []byte) {
	golibsbase32.EncodeToString(data)
}
