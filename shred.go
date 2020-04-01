// Package shred is a golang library to mimic the functionality of the linux shred command
package shred

import (
	"crypto/rand"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"sync"
)

// Shredder a shredder
type Shredder struct {
	wg sync.WaitGroup
}

// ShredderConf is a object containing all choices of the user
type ShredderConf struct {
	Shredder     *Shredder
	WriteOptions WriteOptions
	Times        int
	Delete       bool
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

// NewShredderConf create a new shredder
func NewShredderConf(shredder *Shredder, options WriteOptions, times int, delete bool) *ShredderConf {
	return &ShredderConf{
		Shredder:     shredder,
		WriteOptions: options,
		Times:        times,
		Delete:       delete,
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
	// Get fileinfo
	fileinfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	size := fileinfo.Size()

	if shredderConf.WriteOptions&WriteRand == WriteRand {
		// Write rand
		err = shredderConf.WriteRandom(path, size)
		if err != nil {
			return err
		}
	}

	if shredderConf.WriteOptions&WriteRandSecure == WriteRandSecure {
		// Write rand
		err = shredderConf.WriteRandomSecure(path, size)
		if err != nil {
			return err
		}
	}

	// Write zeros if desired
	if shredderConf.WriteOptions&WriteZeros == WriteZeros {
		err = doWriteZeros(path, size)
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
func (shredderConf *ShredderConf) WriteRandom(file string, size int64) error {
	// Do n times. Specified in conf
	for i := 0; i < shredderConf.Times; i++ {
		// Open file
		f, err := os.OpenFile(file, os.O_RDWR, 0)
		if err != nil {
			return err
		}

		// Seek to start
		offset, err := f.Seek(0, 0)
		if err != nil {
			return err
		}

		// Read random
		buff := make([]byte, size)
		_, err = mathrand.Read(buff)
		if err != nil {
			return err
		}

		// Write random
		_, err = f.WriteAt(buff, offset)
		if err != nil {
			return err
		}

		// Close file
		f.Close()
	}

	return nil
}

// WriteRandomSecure overwrites a File with random stuff.
func (shredderConf *ShredderConf) WriteRandomSecure(file string, size int64) error {
	// Do n times. Specified in conf
	for i := 0; i < shredderConf.Times; i++ {
		// Open file
		f, err := os.OpenFile(file, os.O_RDWR, 0)
		if err != nil {
			return err
		}

		// Seek to start
		offset, err := f.Seek(0, 0)
		if err != nil {
			return err
		}

		// Read random
		buff := make([]byte, size)
		_, err = rand.Read(buff)
		if err != nil {
			return err
		}

		// Write random
		_, err = f.WriteAt(buff, offset)
		if err != nil {
			return err
		}

		// Close file
		f.Close()
	}

	return nil
}

func doWriteZeros(path string, size int64) error {
	// Open file
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	// Goto 0
	offset, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Create zeroed buffer and write it
	buff := make([]byte, size)
	_, err = file.WriteAt(buff, offset)

	file.Close()
	return err
}
