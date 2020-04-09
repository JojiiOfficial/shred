// Package shred is a golang library to mimic the functionality of the linux shred command
package shred

import (
	crand "crypto/rand"
	"io"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Shredder a shredder
type Shredder struct {
	wg sync.WaitGroup
}

// ShredderConf is a object containing all choices of the user
type ShredderConf struct {
	Shredder            *Shredder
	WriteOptions        WriteOptions
	Times               int
	Delete              bool
	WriteRandBufferSize int
}

// WriteOptions options how to shred
type WriteOptions int

// Available write options
const (
	NoWrite WriteOptions = 1 << iota
	WriteZeros
	WriteRand
	WriteRandSecure
)

// DefaultBufferSize the default buffersize used for writing
// operations
const DefaultBufferSize = 10 * 1024

// NewShredderConf create a new shredder
func NewShredderConf(shredder *Shredder, options WriteOptions, times int, delete bool) *ShredderConf {
	return &ShredderConf{
		Shredder:            shredder,
		WriteOptions:        options,
		Times:               times,
		Delete:              delete,
		WriteRandBufferSize: DefaultBufferSize,
	}
}

// ShredPath shreds all files in the location of path
// recursively. If remove is set to true files will be deleted
// after shredding. When a file is shredded its content
// is NOT recoverable so !!USE WITH CAUTION!!
func (shredderConf *ShredderConf) ShredPath(path string) error {
	stats, err := os.Stat(path)
	if err != nil {
		return err
	} else if stats.IsDir() {
		err := shredderConf.ShredDir(path)
		if err != nil {
			return err
		}
	} else {
		err := shredderConf.ShredFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

// ShredDir overwrites every File in the location of path and everything in its subdirectories
func (shredderConf *ShredderConf) ShredDir(path string) error {
	// For each file
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		stats, _ := os.Stat(path)

		if !stats.IsDir() {
			shredderConf.Shredder.wg.Add(1)
			go (func() {
				shredderConf.ShredFile(path)
				shredderConf.Shredder.wg.Done()
			})()
			shredderConf.Shredder.wg.Wait()
		}
		return nil
	})
	return err
}

// ShredFile overwrites a given ShredFile in the location of path
func (shredderConf *ShredderConf) ShredFile(path string) error {
	// Write rand
	if shredderConf.WriteOptions&WriteRand == WriteRand {
		err := shredderConf.WriteRandom(path, false)
		if err != nil {
			return err
		}
	}

	// Write rand secure (using crypto/rand)
	if shredderConf.WriteOptions&WriteRandSecure == WriteRandSecure {
		err := shredderConf.WriteRandom(path, true)
		if err != nil {
			return err
		}
	}

	// Write zeros if desired
	if shredderConf.WriteOptions&WriteZeros == WriteZeros {
		err := shredderConf.DoWriteZeros(path)
		if err != nil {
			return err
		}
	}

	// Delete file if desired
	if shredderConf.Delete {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteRandom overwrites a File with random stuff.
// If secure is true, crypto/rand is used, otherwise math/rand
func (shredderConf *ShredderConf) WriteRandom(file string, secure bool) error {
	// Do n times. Specified in conf
	buff := make([]byte, shredderConf.WriteRandBufferSize)

	var r io.Reader
	if secure {
		r = crand.Reader
	} else {
		r = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	}

	var err error
	var f *os.File

	// Overwrite file shredderConf.Times times
	for i := 0; i < shredderConf.Times; i++ {
		f, err = shredderConf.OverwriteFile(file, r, buff)
		f.Close()
		if err != nil {
			break
		}
	}

	return err
}

// OverwriteFile opens and overwrites a files ONCE by reading from r
// returned file must be closed
func (shredderConf *ShredderConf) OverwriteFile(file string, r io.Reader, buff []byte) (*os.File, error) {
	var wrCounter int
	var b bool
	var n int
	var err error
	var f *os.File

	// Opens a file Write mode only
	// reading is not necessary
	f, err = os.OpenFile(file, os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}

	// Retrieve fileinfo to get the filesize
	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}

	readLen := int64(shredderConf.WriteRandBufferSize)

	for {
		// Current offset to start write call from
		offs := int64(wrCounter * shredderConf.WriteRandBufferSize)
		if offs > fs.Size() {
			break
		}

		// If end of write overflows the filesize
		// Use the end of file as end to read
		if offs+int64(shredderConf.WriteRandBufferSize) > fs.Size() {
			readLen = fs.Size() - offs
			b = true
		}

		n, err = r.Read(buff[:readLen])
		if err != nil {
			return nil, err
		}

		// Write at offset n bytes
		_, err = f.WriteAt(buff[:n], offs)
		if err != nil {
			return nil, err
		}

		if b {
			break
		}

		wrCounter++
	}

	return f, nil
}

// DoWriteZeros overwrite file with zeros
func (shredderConf *ShredderConf) DoWriteZeros(file string) error {
	buff := make([]byte, shredderConf.WriteRandBufferSize)
	r := ZeroReader{}
	f, err := shredderConf.OverwriteFile(file, r, buff)
	defer f.Close()

	return err
}

// ZeroReader struct implementing the io.Reader interface
// used to overwrite files with zeros
type ZeroReader struct{}

func (z ZeroReader) Read(b []byte) (int, error) {
	memset(b, 0)
	return len(b), nil
}

// Sets all indexes of b to val
func memset(b []byte, val byte) {
	for i := range b {
		b[i] = val
	}
}
