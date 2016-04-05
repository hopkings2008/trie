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
	/*for i := 0; i < 2048; i++ {
		c.Assert(prefixes[i], check.Equals, trash[i])
	}*/
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

func (dbms *InfoDbMgrSuites) TestSaveLoadMany(c *check.C) {
	prefixes := make([]string, 4096)
	for i := 0; i < 4096; i++ {
		prefixes[i] = dbms.rand.String()
		err := dbms.dbm.Add(prefixes[i])
		c.Assert(err, check.IsNil)
	}
	err := dbms.dbm.Save()
	c.Assert(err, check.IsNil)
	dbms.dbm = nil
	dbms.dbm, err = CreateInfoDbMgr(dbms.root, "db")
	c.Assert(err, check.IsNil)
	err = dbms.dbm.Load()
	c.Assert(err, check.IsNil)
	for i := 0; i < 4096; i++ {
		err := dbms.dbm.Delete(prefixes[i])
		c.Assert(err, check.IsNil)
	}
	trash := dbms.dbm.GetTrash()
	c.Assert(trash, check.NotNil)
	c.Assert(len(trash), check.Equals, 4096)
	/*for i := 0; i < 2048; i++ {
		c.Assert(prefixes[i] == trash[i], check.Equals, true)
	}*/
}
