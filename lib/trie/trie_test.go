package trie

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/check.v1"
	"trie/lib/util"
)

func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&TrieSuites{})

type TrieSuites struct {
	trie *Trie
	rand *util.RandString
}

func (ts *TrieSuites) SetUpSuite(c *check.C) {
	ts.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
}

func (ts *TrieSuites) TearDownSuite(c *check.C) {
}

func (ts *TrieSuites) SetUpTest(c *check.C) {
	ts.trie = CreateTrie()
	c.Assert(ts.trie, check.NotNil)
}

func (ts *TrieSuites) TearDownTest(c *check.C) {
	ts.trie = nil
}

func (ts *TrieSuites) TestInsertDelete(c *check.C) {
	prefix := ts.rand.String()
	for i := 0; i < 128; i++ {
		err := ts.trie.Insert(prefix)
		c.Assert(err, check.IsNil)
	}
	for i := 0; i < 127; i++ {
		d, err := ts.trie.Delete(prefix)
		c.Assert(err, check.IsNil)
		c.Assert(d, check.Equals, false)
	}

	d, err := ts.trie.Delete(prefix)
	c.Assert(err, check.IsNil)
	c.Assert(d, check.Equals, true)
}

func (ts *TrieSuites) TestSaveLoad(c *check.C) {
	dir, err := ioutil.TempDir(".", "sha256")
	c.Assert(err, check.IsNil)
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "trie")
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0666)
	c.Assert(err, check.IsNil)
	c.Assert(fp, check.NotNil)
	prefixes := make([]string, 2048)

	for i := 0; i < 2048; i++ {
		prefix := ts.rand.String()
		prefixes[i] = prefix
		err := ts.trie.Insert(prefix)
		c.Assert(err, check.IsNil)
	}

	err = ts.trie.Save(fp)
	c.Assert(err, check.IsNil)
	fp.Close()

	fp, err = os.OpenFile(file, os.O_RDONLY, 0666)
	defer fp.Close()
	c.Assert(err, check.IsNil)
	c.Assert(fp, check.NotNil)
	ts.trie = nil
	ts.trie = CreateTrie()
	err = ts.trie.Load(fp)
	c.Assert(err, check.IsNil)
	for _, prefix := range prefixes {
		d, err := ts.trie.Delete(prefix)
		c.Assert(err, check.IsNil)
		c.Assert(d, check.Equals, true)
	}
}
