package infodb

import (
	"io/ioutil"
	"os"
	//"testing"

	"gopkg.in/check.v1"
	"trie/lib/util"
)

//func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&InfoDbMgrSuites{})

type InfoDbMgrSuites struct {
	dbm  *InfoDbMgr
	root string
	rand *util.RandString
}

func (dbms *InfoDbMgrSuites) SetUpSuite(c *check.C) {
	dbms.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
}

func (dbms *InfoDbMgrSuites) TearDownSuite(c *check.C) {
}

func (dbms *InfoDbMgrSuites) SetUpTest(c *check.C) {
	var err error
	dbms.root, err = ioutil.TempDir(".", "dbmgr")
	c.Assert(err, check.IsNil)
	dbms.dbm, err = CreateInfoDbMgr(dbms.root, "db")
	c.Assert(err, check.IsNil)
	c.Assert(dbms.dbm, check.NotNil)
}

func (dbms *InfoDbMgrSuites) TearDownTest(c *check.C) {
	os.RemoveAll(dbms.root)
	dbms.dbm = nil
}

func (dbms *InfoDbMgrSuites) TestAddDeleteSame(c *check.C) {
	prefix := dbms.rand.String()
	for i := 0; i < 2048; i++ {
		err := dbms.dbm.Add(prefix)
		c.Assert(err, check.IsNil)
	}
	for i := 0; i < 2047; i++ {
		err := dbms.dbm.Delete(prefix)
		c.Assert(err, check.IsNil)
	}
	c.Assert(len(dbms.dbm.GetTrash()), check.Equals, 0)
	err := dbms.dbm.Delete(prefix)
	c.Assert(err, check.IsNil)
	c.Assert(len(dbms.dbm.GetTrash()), check.Equals, 1)
	c.Assert(dbms.dbm.GetTrash()[0] == prefix, check.Equals, true)
}

func (dbms *InfoDbMgrSuites) TestAddDeleteDiff(c *check.C) {
	prefixes := make([]string, 2048)

	for i := 0; i < 2048; i++ {
		prefixes[i] = dbms.rand.String()
		err := dbms.dbm.Add(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	for i := 0; i < 2048; i++ {
		err := dbms.dbm.Delete(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	trash := dbms.dbm.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, 2048)
}

func (dbms *InfoDbMgrSuites) TestSaveLoadOne(c *check.C) {
	prefix := dbms.rand.String()
	err := dbms.dbm.Add(prefix)
	c.Assert(err, check.IsNil)
	err = dbms.dbm.Save()
	c.Assert(err, check.IsNil)
	dbms.dbm = nil
	dbms.dbm, err = CreateInfoDbMgr(dbms.root, "db")
	c.Assert(err, check.IsNil)
	err = dbms.dbm.Load()
	c.Assert(err, check.IsNil)
	err = dbms.dbm.Delete(prefix)
	c.Assert(err, check.IsNil)
	c.Assert(len(dbms.dbm.GetTrash()), check.Equals, 1)
	c.Assert(dbms.dbm.GetTrash()[0] == prefix, check.Equals, true)
}

func (dbms *InfoDbMgrSuites) TestDbMgrSaveLoadMany(c *check.C) {
	num := 2000000
	prefixes := make(map[string]int)
	for i := 0; i < num; i++ {
		prefix := dbms.rand.String()
		count, ok := prefixes[prefix]
		if !ok {
			prefixes[prefix] = 1
		} else {
			prefixes[prefix] = count + 1
		}
		err := dbms.dbm.Add(prefix)
		c.Assert(err, check.IsNil)
	}
	err := dbms.dbm.Save()
	c.Assert(err, check.IsNil)
	dbms.dbm = nil
	dbms.dbm, err = CreateInfoDbMgr(dbms.root, "db")
	c.Assert(err, check.IsNil)
	err = dbms.dbm.Load()
	c.Assert(err, check.IsNil)
	for prefix, count := range prefixes {
		for i := 0; i < count; i++ {
			err := dbms.dbm.Delete(prefix)
			c.Assert(err, check.IsNil)
		}
	}
	trash := dbms.dbm.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, len(prefixes))
	for _, prefix := range trash {
		delete(prefixes, prefix)
	}

	c.Assert(len(prefixes), check.Equals, 0)
}

func (dbms *InfoDbMgrSuites) BenchmarkDbMgrSaveLoadMany(c *check.C) {
	num := 2000000
	for i := 0; i < c.N; i++ {
		root, err := ioutil.TempDir(".", "dbmgr")
		c.Assert(err, check.IsNil)
		//defer os.RemoveAll(root)
		dbm, err := CreateInfoDbMgr(root, "db")
		c.Assert(err, check.IsNil)
		c.Assert(dbm, check.NotNil)
		prefixes := make(map[string]int)
		for i := 0; i < num; i++ {
			prefix := dbms.rand.String()
			count, ok := prefixes[prefix]
			if !ok {
				prefixes[prefix] = 1
			} else {
				prefixes[prefix] = count + 1
			}
			err := dbm.Add(prefix)
			c.Assert(err, check.IsNil)
		}
		err = dbm.Save()
		c.Assert(err, check.IsNil)
		dbm = nil
		dbm, err = CreateInfoDbMgr(root, "db")
		c.Assert(err, check.IsNil)
		err = dbm.Load()
		c.Assert(err, check.IsNil)
		for prefix, count := range prefixes {
			for i := 0; i < count; i++ {
				err := dbm.Delete(prefix)
				c.Assert(err, check.IsNil)
			}
		}
		trash := dbm.GetTrash()
		c.Assert(trash, check.NotNil)
		c.Assert(len(trash), check.Equals, len(prefixes))
		for _, prefix := range trash {
			delete(prefixes, prefix)
		}

		c.Assert(len(prefixes), check.Equals, 0)

	}
}

func (dbms *InfoDbMgrSuites) addDeleteMem(c *check.C, dbm *InfoDbMgr, num int) {
	prefixes := make(map[string]int)
	for i := 0; i < num; i++ {
		prefix := dbms.rand.String()
		count, ok := prefixes[prefix]
		if !ok {
			prefixes[prefix] = 1
		} else {
			prefixes[prefix] = count + 1
		}
		err := dbm.Add(prefix)
		c.Assert(err, check.IsNil)
	}
	for prefix, count := range prefixes {
		for i := 0; i < count; i++ {
			err := dbm.Delete(prefix)
			c.Assert(err, check.IsNil)
		}
	}
	trash := dbm.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, len(prefixes))
	for _, prefix := range trash {
		delete(prefixes, prefix)
	}

	c.Assert(len(prefixes), check.Equals, 0)
}

func (dbms *InfoDbMgrSuites) BenchmarkDbMgrAddDelete(c *check.C) {
	num := 2000000
	for i := 0; i < c.N; i++ {
		root, err := ioutil.TempDir(".", "dbmgr")
		c.Assert(err, check.IsNil)
		defer os.RemoveAll(root)
		dbm, err := CreateInfoDbMgr(root, "dbmgrfile")
		c.Assert(err, check.IsNil)
		c.Assert(dbm, check.NotNil)
		dbms.addDeleteMem(c, dbm, num)
	}
}