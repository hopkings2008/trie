package storage

import (
	"fmt"
	"io"
	"path"

	"trie/lib/storage/driver"
)

type StorageMgr struct {
	Drivers map[string]driver.StorageDriver
	Root    string
	File    string
}

func CreateStorageMgr(root, file string) *StorageMgr {
	return &StorageMgr{
		Drivers: make(map[string]driver.StorageDriver),
		Root:    root,
		File:    file,
	}
}

func (s *StorageMgr) Init() error {
	set := "0123456789abcdef"
	for _, c := range set {
		for _, cc := range set {
			root := path.Join(s.Root, "sha256", fmt.Sprintf("%c%c", c, cc))
			if err := s.createDriver(root, "filesystem"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *StorageMgr) GetWriter(prefix string, bound int) (driver.StorageWriter, error) {
	root := s.getPrefix(prefix, bound)
	if driver, ok := s.Drivers[root]; ok {
		writer, err := driver.Writer(s.File, true)
		return writer, err
	}

	return nil, fmt.Errorf("GetWriter: Cannot find driver for prefix %s with bound %d", prefix, bound)
}

func (s *StorageMgr) GetReader(prefix string, bound int) (io.ReadCloser, error) {
	root := s.getPrefix(prefix, bound)
	if driver, ok := s.Drivers[root]; ok {
		reader, err := driver.Reader(s.File, int64(0))
		return reader, err
	}

	return nil, fmt.Errorf("GetReader: Cannot find driver for prefix %s with bound %d", prefix, bound)
}

func (s *StorageMgr) getPrefix(prefix string, bound int) string {
	id := ""
	for i, c := range prefix {
		if i >= bound {
			break
		}
		id = fmt.Sprintf("%s%c", id, c)
	}
	root := path.Join(s.Root, "sha256", id)
	return root
}

func (s *StorageMgr) createDriver(root, name string) error {
	opts := make(map[string]interface{})
	opts["root"] = root
	storageDriver, err := driver.Create(name, opts)
	if err != nil {
		return err
	}

	s.Drivers[root] = storageDriver
	return nil
}
