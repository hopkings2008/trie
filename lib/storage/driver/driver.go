package driver

import (
	"fmt"
	"io"
)

var driverFactory map[string]StorageDriverFactory = make(map[string]StorageDriverFactory)

type StorageDriverFactory interface {
	Create(parameters map[string]interface{}) (StorageDriver, error)
}

type StorageWriter interface {
	io.WriteCloser
	Commit() error
	Cancel() error
}

type StorageDriver interface {
	Writer(file string, append bool) (StorageWriter, error)
	Reader(file string, pos int64) (io.ReadCloser, error)
	Name() string
}

func Register(name string, factory StorageDriverFactory) {
	if _, ok := driverFactory[name]; ok {
		return
	}

	driverFactory[name] = factory
}

func Create(name string, parameters map[string]interface{}) (StorageDriver, error) {
	if factory, ok := driverFactory[name]; ok {
		return factory.Create(parameters)
	}
	return nil, fmt.Errorf("Unknown driver type %s", name)
}
