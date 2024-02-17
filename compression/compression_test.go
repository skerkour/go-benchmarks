package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	zstddd "github.com/DataDog/zstd"
	snappygolang "github.com/golang/snappy"
	"github.com/klauspost/compress/s2"
	snappykp "github.com/klauspost/compress/snappy"
	zstdkp "github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
)

var (
	FILES = []string{
		"alices_adventures_in_wonderland.txt.gz",
		"illiad.txt.gz",
		"country_asn.csv.gz",
		"country_asn.json.gz",
		"country_asn.mmdb.gz",
	}
)

type Compresser interface {
	Compress(destination io.Writer, source io.Reader) (err error)
	Decompress(destination io.Writer, source io.Reader) (err error)
}

// func BenchmarkCompress(b *testing.B) {
// 	for _, file := range FILES {
// 		benchmarkCompress(file, "klausp_s2_default", newKlauspostS2Compresser(1), b)
// 		benchmarkCompress(file, "klausp_s2_better_compression", newKlauspostS2Compresser(2), b)
// 		benchmarkCompress(file, "klausp_s2_best_compression", newKlauspostS2Compresser(3), b)
// 		benchmarkCompress(file, "golang_snappy", golangSnappyCompresser{}, b)
// 		benchmarkCompress(file, "klausp_snappy", klauspostSnappyCompresser{}, b)
// 		benchmarkCompress(file, "pierrec_lz4", pierrecLz4Compresser{}, b)
// 		benchmarkCompress(file, "klausp_zstd_1", newklauspostZstdCompresser(zstdkp.SpeedFastest), b)
// 		benchmarkCompress(file, "klausp_zstd_3", newklauspostZstdCompresser(zstdkp.SpeedDefault), b)
// 		benchmarkCompress(file, "klausp_zstd_better_compression", newklauspostZstdCompresser(zstdkp.SpeedBetterCompression), b)
// 		benchmarkCompress(file, "klausp_zstd_best_compression", newklauspostZstdCompresser(zstdkp.SpeedBestCompression), b)
// 		benchmarkCompress(file, "datadog_zstd_1", newdatadogZstdCompresser(zstddd.BestSpeed), b)
// 		benchmarkCompress(file, "datadog_zstd_3", newdatadogZstdCompresser(3), b)
// 		benchmarkCompress(file, "datadog_zstd_5", newdatadogZstdCompresser(zstddd.DefaultCompression), b)
// 		benchmarkCompress(file, "datadog_zstd_7", newdatadogZstdCompresser(7), b)
// 		benchmarkCompress(file, "datadog_zstd_20", newdatadogZstdCompresser(zstddd.BestCompression), b)
// 		benchmarkCompress(file, "golang_gzip_fastest", newGolangGzipCompresser(gzip.BestSpeed), b)
// 		benchmarkCompress(file, "golang_gzip_default", newGolangGzipCompresser(gzip.DefaultCompression), b)
// 		benchmarkCompress(file, "golang_gzip_best_compression", newGolangGzipCompresser(gzip.BestCompression), b)
// 	}
// }

func benchmarkCompress[C Compresser](file, algorithm string, compresser C, b *testing.B) {
	originaleFilename := strings.TrimSuffix(file, ".gz")
	b.Run(fmt.Sprintf("%s-%s", originaleFilename, algorithm), func(b *testing.B) {
		originalData, err := readGzippedFile(filepath.Join("..", "testdata", file))
		if err != nil {
			b.Error(err)
		}

		originalDataReader := bytes.NewReader(originalData)
		destinationBuffer := bytes.NewBuffer(make([]byte, 0, len(originalData)*2))

		runtime.GC()
		b.ReportAllocs()
		b.SetBytes(int64(len(originalData)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			originalDataReader.Seek(0, io.SeekStart)
			destinationBuffer.Reset()
			err = compresser.Compress(destinationBuffer, originalDataReader)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkDecompress(b *testing.B) {
	for _, file := range FILES {
		benchmarkDecompress(file, "klausp_s2_default", newKlauspostS2Compresser(1), b)
		benchmarkDecompress(file, "klausp_s2_better_compression", newKlauspostS2Compresser(2), b)
		benchmarkDecompress(file, "klausp_s2_best_compression", newKlauspostS2Compresser(3), b)
		benchmarkDecompress(file, "golang_snappy", golangSnappyCompresser{}, b)
		benchmarkDecompress(file, "klausp_snappy", klauspostSnappyCompresser{}, b)
		benchmarkDecompress(file, "pierrec_lz4", pierrecLz4Compresser{}, b)
		benchmarkDecompress(file, "klausp_zstd_1", newklauspostZstdCompresser(zstdkp.SpeedFastest), b)
		benchmarkDecompress(file, "klausp_zstd_3", newklauspostZstdCompresser(zstdkp.SpeedDefault), b)
		benchmarkDecompress(file, "klausp_zstd_better_compression", newklauspostZstdCompresser(zstdkp.SpeedBetterCompression), b)
		benchmarkDecompress(file, "klausp_zstd_best_compression", newklauspostZstdCompresser(zstdkp.SpeedBestCompression), b)
		benchmarkDecompress(file, "datadog_zstd_1", newdatadogZstdCompresser(zstddd.BestSpeed), b)
		benchmarkDecompress(file, "datadog_zstd_3", newdatadogZstdCompresser(3), b)
		benchmarkDecompress(file, "datadog_zstd_5", newdatadogZstdCompresser(zstddd.DefaultCompression), b)
		benchmarkDecompress(file, "datadog_zstd_7", newdatadogZstdCompresser(7), b)
		benchmarkDecompress(file, "datadog_zstd_best_20", newdatadogZstdCompresser(zstddd.BestCompression), b)
		benchmarkDecompress(file, "golang_gzip_fastest", newGolangGzipCompresser(gzip.BestSpeed), b)
		benchmarkDecompress(file, "golang_gzip_default", newGolangGzipCompresser(gzip.DefaultCompression), b)
		benchmarkDecompress(file, "golang_gzip_best_compression", newGolangGzipCompresser(gzip.BestCompression), b)
	}
}

func benchmarkDecompress[C Compresser](file, algorithm string, compresser C, b *testing.B) {
	originaleFilename := strings.TrimSuffix(file, ".gz")
	b.Run(fmt.Sprintf("%s-%s", originaleFilename, algorithm), func(b *testing.B) {
		originalData, err := readGzippedFile(filepath.Join("..", "testdata", file))
		if err != nil {
			b.Error(err)
		}

		originalDataBuffer := bytes.NewBuffer(originalData)
		compressedDataBuffer := bytes.NewBuffer(make([]byte, 0, len(originalData)*2))
		err = compresser.Compress(compressedDataBuffer, originalDataBuffer)
		if err != nil {
			b.Error(err)
		}

		compressedDataReader := bytes.NewReader(compressedDataBuffer.Bytes())
		destinationBuffer := bytes.NewBuffer(make([]byte, 0, len(originalData)*2))

		runtime.GC()
		b.ReportAllocs()
		b.SetBytes(int64(len(originalData)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			destinationBuffer.Reset()
			compressedDataReader.Seek(0, io.SeekStart)
			err = compresser.Decompress(destinationBuffer, compressedDataReader)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func readGzippedFile(filePath string) (data []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("openning file %s: %w", filePath, err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("getting info for file %s: %w", filePath, err)
		return
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		err = fmt.Errorf("creating gzip reader %s: %w", filePath, err)
		return
	}

	dataBuffer := bytes.NewBuffer(make([]byte, 0, fileInfo.Size()*3))

	_, err = io.Copy(dataBuffer, gzipReader)
	if err != nil {
		err = fmt.Errorf("reading file %s: %w", filePath, err)
		return
	}

	data = dataBuffer.Bytes()

	return
}

type klauspostS2Compresser struct {
	level int
}

func newKlauspostS2Compresser(level int) klauspostS2Compresser {
	return klauspostS2Compresser{
		level: level,
	}
}

func (compresser klauspostS2Compresser) Compress(destination io.Writer, source io.Reader) (err error) {
	var encoder *s2.Writer

	switch compresser.level {
	case 2:
		encoder = s2.NewWriter(destination, s2.WriterBetterCompression(), s2.WriterConcurrency(1))
	case 3:
		encoder = s2.NewWriter(destination, s2.WriterBestCompression(), s2.WriterConcurrency(1))
	default:
		encoder = s2.NewWriter(destination, s2.WriterConcurrency(1))
	}

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	// Blocks until compression is done.
	return encoder.Close()
}

func (klauspostS2Compresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder := s2.NewReader(source)
	_, err = io.Copy(destination, decoder)
	return err
}

type klauspostZstdCompresser struct {
	level zstdkp.EncoderLevel
}

func newklauspostZstdCompresser(level zstdkp.EncoderLevel) klauspostZstdCompresser {
	return klauspostZstdCompresser{
		level: level,
	}
}

func (compresser klauspostZstdCompresser) Compress(destination io.Writer, source io.Reader) (err error) {
	var encoder *zstdkp.Encoder

	encoder, err = zstdkp.NewWriter(destination, zstdkp.WithEncoderLevel(compresser.level), zstdkp.WithEncoderConcurrency(1))
	if err != nil {
		return
	}

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	// Blocks until compression is done.
	return encoder.Close()
}

func (klauspostZstdCompresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder, err := zstdkp.NewReader(source, zstdkp.WithDecoderConcurrency(1))
	if err != nil {
		return
	}
	_, err = io.Copy(destination, decoder)
	return err
}

type pierrecLz4Compresser struct {
}

func (pierrecLz4Compresser) Compress(destination io.Writer, source io.Reader) (err error) {
	encoder := lz4.NewWriter(destination)
	if err != nil {
		return
	}

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	return encoder.Close()
}

func (pierrecLz4Compresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder := lz4.NewReader(source)

	_, err = io.Copy(destination, decoder)
	return err
}

type datadogZstdCompresser struct {
	level int
}

func newdatadogZstdCompresser(level int) datadogZstdCompresser {
	return datadogZstdCompresser{
		level: level,
	}
}

func (compresser datadogZstdCompresser) Compress(destination io.Writer, source io.Reader) (err error) {
	encoder := zstddd.NewWriterLevel(destination, compresser.level)

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}

	// Blocks until compression is done.
	return encoder.Close()
}

func (datadogZstdCompresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder := zstddd.NewReader(source)

	_, err = io.Copy(destination, decoder)
	return
}

type golangSnappyCompresser struct {
}

func (golangSnappyCompresser) Compress(destination io.Writer, source io.Reader) (err error) {
	encoder := snappygolang.NewBufferedWriter(destination)

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	return encoder.Close()
}

func (golangSnappyCompresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder := snappygolang.NewReader(source)

	_, err = io.Copy(destination, decoder)
	return err
}

type klauspostSnappyCompresser struct {
}

func (klauspostSnappyCompresser) Compress(destination io.Writer, source io.Reader) (err error) {
	// here we use s2 directly as it's what NewBufferedWriter does under the hood
	// but we also need to set s2.WriterConcurrency to 1 to be fair with other implementations
	encoder := s2.NewWriter(destination, s2.WriterSnappyCompat(), s2.WriterBetterCompression(), s2.WriterConcurrency(1))

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	return encoder.Close()
}

func (klauspostSnappyCompresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder := snappykp.NewReader(source)

	_, err = io.Copy(destination, decoder)
	return err
}

type golangGzipCompresser struct {
	level int
}

func newGolangGzipCompresser(level int) golangGzipCompresser {
	return golangGzipCompresser{
		level: level,
	}
}

func (compresser golangGzipCompresser) Compress(destination io.Writer, source io.Reader) (err error) {
	encoder, err := gzip.NewWriterLevel(destination, compresser.level)
	if err != nil {
		return
	}

	_, err = io.Copy(encoder, source)
	if err != nil {
		encoder.Close()
		return err
	}
	return encoder.Close()
}

func (golangGzipCompresser) Decompress(destination io.Writer, source io.Reader) (err error) {
	decoder, err := gzip.NewReader(source)
	if err != nil {
		return
	}

	_, err = io.Copy(destination, decoder)
	return err
}
