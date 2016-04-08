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
	trie     *Trie
	rand     *util.RandString
	prefixes map[string]int
}

func (ts *TrieSuites) SetUpSuite(c *check.C) {
	ts.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
	ts.prefixes = ts.getStrings(500000, true)
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

func (ts *TrieSuites) insertDeleteMany(c *check.C, trie *Trie) {
	prefixes := ts.prefixes
	for k, _ := range prefixes {
		err := trie.Insert(k)
		c.Assert(err, check.IsNil)
	}

	for k, _ := range prefixes {
		_, err := trie.Delete(k)
		c.Assert(err, check.IsNil)
	}

	for k, _ := range prefixes {
		ref, err := trie.GetRef(k)
		c.Assert(err, check.NotNil)
		c.Assert(ref, check.Equals, -1)
	}
}

func (ts *TrieSuites) mapInsertDeleteMany(c *check.C, m map[string]int) {
	prefixes := ts.prefixes
	for k, v := range prefixes {
		m[k] = v
	}
	/*for k, _ := range prefixes {
		delete(m, k)
	}
	for k, _ := range prefixes {
		_, ok := m[k]
		c.Assert(ok, check.Equals, false)
	}*/
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
	num := 200000
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
	for i := 0; i < c.N; i++ {
		trie := CreateTrie()
		ts.insertDeleteMany(c, trie)
	}
}

func (ts *TrieSuites) BenchmarkMapInsertDeleteMany(c *check.C) {
	for i := 0; i < c.N; i++ {
		m := make(map[string]int)
		ts.mapInsertDeleteMany(c, m)
	}
}

func (ts *TrieSuites) getStrings(num int, diff bool) map[string]int {
	prefixes := make(map[string]int)
	for i := 0; i < num; {
		str := ts.rand.String()
		count, ok := prefixes[str]
		if ok {
			if !diff {
				count++
				i++
				prefixes[str] = count
			}
			continue
		}
		i++
		prefixes[str] = 1
	}

	return prefixes
}
