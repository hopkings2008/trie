package infodb

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
)

type InfoDbMgr struct {
	Dbs    map[string]*InfoDb
	Root   string
	DbFile string
}

func (dbm *InfoDbMgr) Add(prefix string) error {
	db, err := dbm.getDb(prefix)
	if err != nil {
		return err
	}
	return db.Add(prefix)
}

func (dbm *InfoDbMgr) Delete(prefix string) error {
	db, err := dbm.getDb(prefix)
	if err != nil {
		return err
	}
	return db.Delete(prefix)
}

func (dbm *InfoDbMgr) GetTrash() []string {
	var trash []string
	for _, db := range dbm.Dbs {
		trash = append(trash, db.GetTrash()...)
	}

	return trash
}

func (dbm *InfoDbMgr) Save() error {
	for id, db := range dbm.Dbs {
		if err := db.Save(); err != nil {
			log.Errorf("Failed to save db %s, err: %v", id, err)
			return fmt.Errorf("Failed to save db %s, err: %v", id, err)
		}
	}

	log.Debugf("Succeed to save all dbs")

	return nil
}

func (dbm *InfoDbMgr) Load() error {
	for id, db := range dbm.Dbs {
		if err := db.Load(); err != nil {
			log.Errorf("Failed to load db %s, err: %v", id, err)
			return fmt.Errorf("Failed to load db %s, err: %v", id, err)
		}
	}
	log.Debugf("Succeed to load all dbs")
	return nil
}

func (dbm *InfoDbMgr) getDb(prefix string) (*InfoDb, error) {
	id := dbm.getPrefix(prefix, 2)
	if db, ok := dbm.Dbs[id]; ok {
		return db, nil
	}
	log.Errorf("Cannot find db for prefix %s, id %s", prefix, id)
	return nil, fmt.Errorf("Cannot find db for %s, id %s", prefix, id)
}

func (dbm *InfoDbMgr) getPrefix(prefix string, bound int) string {
	id := ""
	for i, c := range prefix {
		if i >= bound {
			break
		}
		id = fmt.Sprintf("%s%c", id, c)
	}

	return id
}

func CreateInfoDbMgr(root, file string) (*InfoDbMgr, error) {
	var err error
	set := "0123456789ABCDEF"
	dbm := &InfoDbMgr{
		Dbs:  make(map[string]*InfoDb),
		Root: root,
	}

	for _, c := range set {
		for _, cc := range set {
			id := fmt.Sprintf("%c%c", c, cc)
			bucket := path.Join(root, "sha256", id)
			if dbm.Dbs[id], err = CreateInfoDb(bucket, file); err != nil {
				log.Errorf("Failed to CreateInfoDb(%s, %s), err: %v", bucket, file, err)
				return nil, err
			}
		}
	}

	return dbm, nil
}
