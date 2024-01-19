package chunking

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/jotfs/fastcdc-go"
	resticchunker "github.com/restic/chunker"
	"github.com/skerkour/go-benchmarks/utils"
	tigerwill90fastcdc "github.com/tigerwill90/fastcdc"
)

type Chunker interface {
	Chunk(input []byte) (err error)
}

func BenchmarkChunking(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		16 * 1024,
		64 * 1024,
		1024 * 1024,
		10 * 1024 * 1024,
		100 * 1024 * 1024,
		1024 * 1024 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkChunker(size, "jotfs_fastcdc", jotfsFastCDCChunker{}, b)
		benchmarkChunker(size, "tigerwill90_fastcdc", tigerwill90FastCDCChunker{}, b)
		benchmarkChunker(size, "restic_chunker", resticChunker{}, b)
	}
}

func benchmarkChunker[C Chunker](size int64, algorithm string, chunker C, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			chunker.Chunk(buf)
		}
	})
}

type jotfsFastCDCChunker struct{}

func (jotfsFastCDCChunker) Chunk(input []byte) (err error) {
	data := bytes.NewReader(input)
	opts := fastcdc.Options{
		// MinSize:     32 * 1024,
		AverageSize: 64 * 1024,
		// MaxSize:     128 * 1024,
	}

	chunker, err := fastcdc.NewChunker(data, opts)
	if err != nil {
		return
	}

	for {
		_, err = chunker.Next()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}

type tigerwill90FastCDCChunker struct{}

func (tigerwill90FastCDCChunker) Chunk(input []byte) (err error) {
	data := bytes.NewReader(input)
	chunker, err := tigerwill90fastcdc.NewChunker(context.Background(), tigerwill90fastcdc.WithStreamMode(), tigerwill90fastcdc.With64kChunks())
	if err != nil {
		return
	}

	buf := make([]byte, 3*64*1024)
	for {
		var n int

		n, err = data.Read(buf)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		err = chunker.Split(bytes.NewReader(buf[:n]), func(offset, length uint, chunk []byte) error {
			return nil
		})
		if err != nil {
			return
		}
	}

	err = chunker.Finalize(func(offset, length uint, chunk []byte) error {
		return nil
	})

	return
}

type resticChunker struct{}

func (resticChunker) Chunk(input []byte) (err error) {
	data := bytes.NewReader(input)
	chnkr := resticchunker.New(data, resticchunker.Pol(0x3DA3358B4DC173))

	buf := make([]byte, 8*1024*1024)

	for {
		_, err = chnkr.Next(buf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}
