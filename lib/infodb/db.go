package infodb

import (
	log "github.com/Sirupsen/logrus"
	"trie/lib/storage/driver"
	_ "trie/lib/storage/driver/filesystem"
	"trie/lib/trie"
)

type InfoDb struct {
	Driver driver.StorageDriver
	MemDb  *trie.Trie
	Trash  []string
	DbFile string
}

func (db *InfoDb) Add(prefix string) error {
	return db.MemDb.Insert(prefix)
}

func (db *InfoDb) Delete(prefix string) error {
	var deleted bool = false
	var err error
	if deleted, err = db.MemDb.Delete(prefix); err == nil {
		if deleted {
			db.Trash = append(db.Trash, prefix)
		}
		return nil
	}

	log.Errorf("Failed to delete %s, err: %v", prefix, err)
	return err
}

func (db *InfoDb) GetTrash() []string {
	return db.Trash
}

func (db *InfoDb) Save() error {
	writer, err := db.Driver.Writer(db.DbFile, true)
	if err != nil {
		log.Errorf("Cannot get writer for %s, err: %v", db.DbFile, err)
		return err
	}
	defer writer.Close()
	if err = db.MemDb.Save(writer); err != nil {
		log.Errorf("Failed to save %s, err: %v", db.DbFile, err)
		return err
	}
	log.Debugf("Succeed to save %s", db.DbFile)
	return nil
}

func (db *InfoDb) Load() error {
	reader, err := db.Driver.Reader(db.DbFile, int64(0))
	if err != nil {
		log.Errorf("Cannot get reader for %s, err: %v", db.DbFile, err)
		return err
	}
	defer reader.Close()
	if err = db.MemDb.Load(reader); err != nil {
		log.Errorf("Failed to load %s, err: %v", db.DbFile, err)
		return err
	}
	log.Debugf("Succeed to load %s", db.DbFile)
	return nil
}

func CreateInfoDb(root, file string) (*InfoDb, error) {
	opts := make(map[string]interface{})
	opts["root"] = root
	driver, err := driver.Create("filesystem", opts)
	if err != nil {
		log.Errorf("Failed to create driver for root %s, err: %v", root, err)
		return nil, err
	}

	db := &InfoDb{
		Driver: driver,
		MemDb:  trie.CreateTrie(),
		DbFile: file,
	}

	return db, nil
}
