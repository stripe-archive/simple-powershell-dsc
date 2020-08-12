package fsmock

import (
	"io"
	"os"
)

// Mock type for a filesystem.
type FileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
}

// Mock return value for a file
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	Stat() (os.FileInfo, error)
}

// OSFS implements FileSystem using the local disk.
type OSFS struct{}

func (OSFS) Open(name string) (File, error)               { return os.Open(name) }
func (OSFS) Create(name string) (File, error)             { return os.Create(name) }
func (OSFS) Stat(name string) (os.FileInfo, error)        { return os.Stat(name) }
func (OSFS) MkdirAll(path string, perm os.FileMode) error { return os.Mkdir(path, perm) }
