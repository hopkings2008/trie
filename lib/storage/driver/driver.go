package driver

import (
	"fmt"
	"io"
)

var StorageDrivers map[string]StorageDriver

type StorageDriver interface {
	Writer(file string, pos int64) io.Writer
	Reader(file string, pos int64) io.Reader
}
