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

func (ts *TrieSuites) insertDeleteMany(c *check.C, trie *Trie, num int) {
	prefixes := make([]string, num)
	for i := 0; i < num; i++ {
		prefixes[i] = ts.rand.String()
		err := trie.Insert(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	for i := 0; i < num; i++ {
		_, err := trie.Delete(prefixes[i])
		c.Assert(err, check.IsNil)
	}

	for i := 0; i < num; i++ {
		ref, err := trie.GetRef(prefixes[i])
		c.Assert(err, check.NotNil)
		c.Assert(ref, check.Equals, -1)
	}
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

type testSelector struct {
	ref   int
	trash []string
}

func (tsl *testSelector) Check(prefix string, node *NodeInfo) bool {
	return node.ref == node.ref
}

func (tsl *testSelector) Get(prefix string, node *NodeInfo) error {
	tsl.trash = append(tsl.trash, prefix)
	return nil
}

func (ts *TrieSuites) TestZeroSelector(c *check.C) {
	tsl := &testSelector{ref: 0}
	num := 2000000
	prefixes := make(map[string]int)

	for count := 0; count < num; {
		var ok bool
		prefix := ts.rand.String()
		if _, ok = prefixes[prefix]; ok {
			continue
		}
		count++
		prefixes[prefix] = count
		err := ts.trie.Update(prefix, 0)
		c.Assert(err, check.IsNil)
	}

	err := ts.trie.Select(tsl)
	c.Assert(err, check.IsNil)
	c.Assert(len(tsl.trash), check.Equals, num)
	for _, prefix := range tsl.trash {
		delete(prefixes, prefix)
	}

	c.Assert(len(prefixes), check.Equals, 0)
}

func (ts *TrieSuites) BenchmarkInsertDeleteMany(c *check.C) {
	num := 2000000
	for i := 0; i < c.N; i++ {
		trie := CreateTrie()
		ts.insertDeleteMany(c, trie, num)
	}
}
