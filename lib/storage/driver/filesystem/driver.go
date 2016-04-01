package filesystem

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"trie/lib/storage/driver"
)

const (
	DriverName = "filesystem"
)

type FileWriter struct {
	file   *os.File
	offset int64
	bw     *bufio.Writer
}

type FileSystem struct {
	root string
}

func (fs *FileSystem) Name() string {
	return DriverName
}

func (fs *FileSystem) Reader(path string, pos int64) (io.ReadCloser, error) {
	file, err := os.OpenFile(fs.fullpath(path), os.O_RDONLY, 0644)
	if err != nil {
		log.Errorf("Failed to open %s, err: %v", fs.fullpath(path), err)
		return nil, err
	}
	seekPos, err := file.Seek(int64(pos), os.SEEK_SET)
	if err != nil {
		file.Close()
		log.Errorf("Failed to seek %d in file %s, err: %v", pos, fs.fullpath(path), err)
		return nil, err
	}
	if seekPos < pos {
		file.Close()
		log.Errorf("seekpos %d < pos %d", seekPos, pos)
		return nil, fmt.Errorf("seekpos %d < pos %d", seekPos, pos)
	}
	return file, nil
}

func (fs *FileSystem) Writer(subpath string, append bool) (driver.StorageWriter, error) {
	fullpath := fs.fullpath(subpath)
	if err := os.MkdirAll(path.Dir(fullpath), 0755); err != nil {
		log.Errorf("Failed to create folder %s, err: %v", path.Dir(fullpath), err)
		return nil, err
	}
	fp, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("Failed to open %s to write, err: %v", fullpath, err)
		return nil, err
	}
	offset := int64(0)
	if append {
		nn, err := fp.Seek(0, os.SEEK_END)
		if err != nil {
			fp.Close()
			log.Errorf("Failed to seek to end in file %s, err: %v", fullpath, err)
			return nil, fmt.Errorf("Failed to seek to end in file %s, err: %v", fullpath, err)
		}
		offset = int64(nn)
	} else {
		if err := fp.Truncate(0); err != nil {
			log.Errorf("Failed to truncate the file %s, err: %v", fullpath, err)
			return nil, err
		}
	}

	return fs.newFileWriter(fp, offset), nil
}

func (fs *FileSystem) fullpath(subpath string) string {
	return path.Join(fs.root, subpath)
}

func (fs *FileSystem) newFileWriter(fp *os.File, offset int64) *FileWriter {
	return &FileWriter{
		file:   fp,
		offset: offset,
		bw:     bufio.NewWriter(fp),
	}
}

// FileWriter member functions
func (fw *FileWriter) Commit() error {
	if err := fw.bw.Flush(); err != nil {
		return err
	}
	if err := fw.file.Sync(); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) Write(p []byte) (int, error) {
	n, err := fw.bw.Write(p)
	fw.offset += int64(n)
	return n, err
}

func (fw *FileWriter) Close() error {
	if err := fw.bw.Flush(); err != nil {
		return err
	}
	if err := fw.file.Sync(); err != nil {
		return err
	}
	if err := fw.file.Close(); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) Cancel() error {
	if err := fw.file.Close(); err != nil {
		return err
	}
	return os.Remove(fw.file.Name())
}

func New(parameters map[string]interface{}) (driver.StorageDriver, error) {
	rootv, ok := parameters["root"]
	if !ok {
		log.Errorf("root is not specified.")
		return nil, fmt.Errorf("root is not specified.")
	}

	root, _ := rootv.(string)
	return &FileSystem{
		root: root,
	}, nil
}

type FileSystemFactory struct {
}

func (fsf *FileSystemFactory) Create(parameters map[string]interface{}) (driver.StorageDriver, error) {
	return New(parameters)
}

func init() {
	driver.Register(DriverName, &FileSystemFactory{})
}
