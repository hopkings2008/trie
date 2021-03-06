package infodb

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/check.v1"
	"trie/lib/util"
)

func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&InfoDbSuites{})

type InfoDbSuites struct {
	db   *InfoDb
	root string
	rand *util.RandString
}

func (dbs *InfoDbSuites) SetUpSuite(c *check.C) {
	dbs.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
}

func (dbs *InfoDbSuites) TearDownSuite(c *check.C) {
}

func (dbs *InfoDbSuites) SetUpTest(c *check.C) {
	var err error
	dbs.root, err = ioutil.TempDir(".", "sha256")
	c.Assert(err, check.IsNil)
	dbs.db, err = CreateInfoDb(dbs.root, "db")
	c.Assert(err, check.IsNil)
	c.Assert(dbs.db, check.NotNil)
}

func (dbs *InfoDbSuites) TearDownTest(c *check.C) {
	os.RemoveAll(dbs.root)
	dbs.db = nil
}

func (dbs *InfoDbSuites) TestAddDeleteSame(c *check.C) {
	prefix := dbs.rand.String()
	for i := 0; i < 2048; i++ {
		err := dbs.db.Add(prefix)
		c.Assert(err, check.IsNil)
	}
	for i := 0; i < 2047; i++ {
		err := dbs.db.Delete(prefix)
		c.Assert(err, check.IsNil)
	}
	c.Assert(len(dbs.db.GetTrash()), check.Equals, 0)
	err := dbs.db.Delete(prefix)
	c.Assert(err, check.IsNil)
	c.Assert(len(dbs.db.GetTrash()), check.Equals, 1)
	c.Assert(dbs.db.GetTrash()[0] == prefix, check.Equals, true)
}

func (dbs *InfoDbSuites) TestAddDeleteDiff(c *check.C) {
	prefixes := make([]string, 2048)

	for i := 0; i < 2048; i++ {
		prefixes[i] = dbs.rand.String()
		err := dbs.db.Add(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	for i := 0; i < 2048; i++ {
		err := dbs.db.Delete(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	trash := dbs.db.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, 2048)
	for i := 0; i < 2048; i++ {
		c.Assert(prefixes[i] == trash[i], check.Equals, true)
	}
}

func (dbs *InfoDbSuites) TestSaveLoadOne(c *check.C) {
	prefix := dbs.rand.String()
	err := dbs.db.Add(prefix)
	c.Assert(err, check.IsNil)
	err = dbs.db.Save()
	c.Assert(err, check.IsNil)
	dbs.db = nil
	dbs.db, err = CreateInfoDb(dbs.root, "db")
	c.Assert(err, check.IsNil)
	err = dbs.db.Load()
	c.Assert(err, check.IsNil)
	err = dbs.db.Delete(prefix)
	c.Assert(err, check.IsNil)
	c.Assert(len(dbs.db.GetTrash()), check.Equals, 1)
	c.Assert(dbs.db.GetTrash()[0] == prefix, check.Equals, true)
}

func (dbs *InfoDbSuites) TestSaveLoadMany(c *check.C) {
	prefixes := make([]string, 2048)
	for i := 0; i < 2048; i++ {
		prefixes[i] = dbs.rand.String()
		err := dbs.db.Add(prefixes[i])
		c.Assert(err, check.IsNil)
	}
	err := dbs.db.Save()
	c.Assert(err, check.IsNil)
	dbs.db = nil
	dbs.db, err = CreateInfoDb(dbs.root, "db")
	c.Assert(err, check.IsNil)
	err = dbs.db.Load()
	c.Assert(err, check.IsNil)
	for i := 0; i < 2048; i++ {
		err := dbs.db.Delete(prefixes[i])
		c.Assert(err, check.IsNil)
	}
	trash := dbs.db.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, 2048)
	for i := 0; i < 2048; i++ {
		c.Assert(prefixes[i] == trash[i], check.Equals, true)
	}
}

func (dbs *InfoDbSuites) addDeleteMem(c *check.C, db *InfoDb, num int) {
	prefixes := make(map[string]int)

	for i := 0; i < num; i++ {
		prefix := dbs.rand.String()
		count, ok := prefixes[prefix]
		if !ok {
			prefixes[prefix] = 1
		} else {
			prefixes[prefix] = count + 1
		}
		err := db.Add(prefix)
		c.Assert(err, check.IsNil)
	}

	for prefix, count := range prefixes {
		for i := 0; i < count; i++ {
			err := db.Delete(prefix)
			c.Assert(err, check.IsNil)
		}
	}

	trash := db.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, len(prefixes))
	for _, prefix := range trash {
		delete(prefixes, prefix)
	}
	c.Assert(len(prefixes), check.Equals, 0)
}

func (dbs *InfoDbSuites) TestInfoDbAddDeleteMem(c *check.C) {
	num := 200000

	root, err := ioutil.TempDir(".", "sha256")
	c.Assert(err, check.IsNil)
	defer os.RemoveAll(root)
	db, err := CreateInfoDb(root, "dbbench")
	c.Assert(err, check.IsNil)
	c.Assert(db, check.NotNil)
	dbs.addDeleteMem(c, db, num)
}
func (dbs *InfoDbSuites) BenchmarkInfoDbAddDeleteMem(c *check.C) {
	num := 2000000

	for i := 0; i < c.N; i++ {
		root, err := ioutil.TempDir(".", "sha256")
		c.Assert(err, check.IsNil)
		defer os.RemoveAll(root)
		db, err := CreateInfoDb(root, "dbbench")
		c.Assert(err, check.IsNil)
		c.Assert(db, check.NotNil)
		dbs.addDeleteMem(c, db, num)
	}
}
