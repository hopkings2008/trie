package storage

import (
	"encoding/gob"

	"trie/lib/storage/driver"
)

type Storage struct {
	Path          string
	storageDriver driver.StorageDriver
}

func CreateStorage(file string, driverName string) *Storage {
	return &Storage{
		Path:          file,
		storageDriver: nil,
	}
}
