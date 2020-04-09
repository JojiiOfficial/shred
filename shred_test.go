package shred

import (
	"testing"
)

func TestShredder(t *testing.T) {
	shredder := Shredder{}
	shredderConfg := NewShredderConf(&shredder, WriteZeros, 1, false)
	err := shredderConfg.ShredFile("./toShredder")
	if err != nil {
		t.Errorf("shredder error: %s", err)
	}
}

func BenchmarkShredderSecure(b *testing.B) {
	b.ReportAllocs()
	shredderConfig := initShredder()

	for i := 0; i < b.N; i++ {
		err := shredderConfig.ShredFile("./toShredder")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkShredder(b *testing.B) {
	b.ReportAllocs()
	shredderConfig := initShredder()

	for i := 0; i < b.N; i++ {
		err := shredderConfig.ShredFile("./toShredder")
		if err != nil {
			b.Error(err)
		}
	}
}

// Big
func BenchmarkShredderBigSecure(b *testing.B) {
	b.ReportAllocs()
	shredderConfig := initShredder()

	for i := 0; i < b.N; i++ {
		err := shredderConfig.ShredFile("./toShredderBig")
		if err != nil {
			b.Error(err)
		}
	}
}

// Big
func BenchmarkShredderBig(b *testing.B) {
	b.ReportAllocs()
	shredderConfig := initShredder()

	for i := 0; i < b.N; i++ {
		err := shredderConfig.ShredFile("./toShredderBig")
		if err != nil {
			b.Error(err)
		}
	}
}

func initShredder() *ShredderConf {
	shredder := Shredder{}
	shredderConfg := NewShredderConf(&shredder, WriteRand, 1, false)
	shredderConfg.WriteRandBufferSize = 4 * 1024
	return shredderConfg
}
