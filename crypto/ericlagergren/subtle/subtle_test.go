package subtle

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"math/rand/v2"
	"testing"
	"time"
)

func TestConstantTimeByteGreater(t *testing.T) {
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			x := byte(i)
			y := byte(j)
			if (ConstantTimeByteGreater(x, y) == 1) != (x > y) {
				t.Fatalf("(%d, %d): expected %t", x, y, x > y)
			}
		}
	}
}

func TestConstantTimeByteLessOrEq(t *testing.T) {
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			x := byte(i)
			y := byte(j)
			if (ConstantTimeByteLessOrEq(x, y) == 1) != (x <= y) {
				t.Fatalf("(%d, %d): expected %t", x, y, x <= y)
			}
		}
	}
}

func TestConstantTimeBigEndianLessOrEq(t *testing.T) {
	d := 2 * time.Second
	if testing.Short() {
		d = 100 * time.Millisecond
	}
	tm := time.NewTimer(d)

	var seed [32]byte
	binary.LittleEndian.PutUint64(seed[0:8], uint64(time.Now().UnixNano()))
	t.Logf("seed: %s", hex.EncodeToString(seed[:]))
	rng := rand.NewChaCha8(seed)

	x := make([]byte, 32)
	y := make([]byte, 32)
	var bx, by big.Int
	for i := 0; ; i++ {
		select {
		case <-tm.C:
			t.Logf("iter: %d", i)
			return
		default:
		}

		rng.Read(x)
		rng.Read(y)

		bx.SetBytes(x)
		by.SetBytes(y)
		want := bx.Cmp(&by) <= 0
		got := ConstantTimeBigEndianLessOrEq(x, y) == 1
		if got != want {
			t.Fatalf("#%d: ConstantTimeBigEndianLessOrEq(%x, %x) != %t", i, x, y, want)
		}
	}
}
