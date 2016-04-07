package suffix

import (
	"testing"

	"gopkg.in/check.v1"
	"trie/lib/util"
)

func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&TrieSuites{})

type TrieSuites struct {
	rand     *util.RandString
	prefixes map[string]int
}

func (ts *TrieSuites) SetUpSuite(c *check.C) {
	InitSets()
	ts.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
	ts.prefixes = ts.getStrings(500000, true)
}

func (ts *TrieSuites) TearDownSuite(c *check.C) {
}

func (ts *TrieSuites) SetUpTest(c *check.C) {
}

func (ts *TrieSuites) TearDownTest(c *check.C) {
}

func (ts *TrieSuites) TestInsertDelete(c *check.C) {
	prefix := ts.rand.String()
	trie := NewTrie()
	ret := trie.Put(prefix, 1)
	c.Assert(ret, check.Equals, true)
	val := trie.Get(prefix)
	c.Assert(val, check.NotNil)
	vali, _ := val.(int)
	c.Assert(vali, check.Equals, 1)
	ret = trie.Put(prefix, 2)
	c.Assert(ret, check.Equals, false)
	val = trie.Get(prefix)
	vali, _ = val.(int)
	c.Assert(vali, check.Equals, 2)
	ret = trie.Delete(prefix)
	c.Assert(ret, check.Equals, true)
	ret = trie.Delete(prefix)
	c.Assert(ret, check.Equals, false)
}

func (ts *TrieSuites) BenchmarkInsertDeleteMany(c *check.C) {
	trie := NewTrie()
	for i := 0; i < c.N; i++ {
		//trie := NewTrie()
		ts.insertDeleteMany(c, trie)
	}
}

func (ts *TrieSuites) BenchmarkMapInsertDeleteMany(c *check.C) {
	for i := 0; i < c.N; i++ {
		m := make(map[string]int)
		ts.mapInsertDeleteMany(c, m)
	}
}

func (ts *TrieSuites) insertDeleteMany(c *check.C, trie *Trie) {
	prefixes := ts.prefixes
	for k, _ := range prefixes {
		trie.Put(k, 1)
		//c.Assert(ret, check.Equals, true)
	}

	for k, _ := range prefixes {
		trie.Delete(k)
		//c.Assert(ret, check.Equals, true)
	}

	for k, _ := range prefixes {
		ref := trie.Get(k)
		c.Assert(ref, check.IsNil)
	}
}

func (ts *TrieSuites) mapInsertDeleteMany(c *check.C, m map[string]int) {
	prefixes := ts.prefixes
	for k, v := range prefixes {
		m[k] = v
	}
	for k, _ := range prefixes {
		delete(m, k)
	}
	for k, _ := range prefixes {
		_, ok := m[k]
		c.Assert(ok, check.Equals, false)
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
