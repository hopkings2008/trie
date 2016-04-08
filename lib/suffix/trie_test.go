package suffix

import (
	//"runtime"
	"testing"

	"gopkg.in/check.v1"
	"trie/lib/util"
)

func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&TrieSuites{})

type TrieSuites struct {
	rand     *util.RandString
	rand32   *util.RandString
	prefixes map[string]int
	tmap     map[string]int
}

func (ts *TrieSuites) SetUpSuite(c *check.C) {
	InitSets()
	ts.rand = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  64,
	}
	ts.rand32 = &util.RandString{
		Sets: "0123456789ABCDEF",
		Len:  32,
	}
	//ts.prefixes = ts.getStrings(500000, true)
	ts.prefixes = ts.getHalfSameStrings(500000)
	ts.tmap = make(map[string]int)
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

type testWalk struct {
	m map[string]interface{}
}

func (tw *testWalk) walk(key string, value interface{}) error {
	tw.m[key] = value
	return nil
}

func (ts *TrieSuites) TestWalkBasic(c *check.C) {
	walker := &testWalk{
		m: make(map[string]interface{}),
	}
	trie := NewTrie()
	ilen := 10
	prefixes := make(map[string]interface{})
	for i := 0; i < ilen; i++ {
		prefixes[ts.rand.String()] = i
	}

	for k, v := range prefixes {
		ret := trie.Put(k, v)
		c.Assert(ret, check.Equals, true)
	}

	err := trie.Walk(walker.walk)
	c.Assert(err, check.IsNil)

	for k, v := range prefixes {
		node := trie.Get(k)
		c.Assert(node, check.Equals, v)
		ret := trie.Delete(k)
		c.Assert(ret, check.Equals, true)
	}

	for k, _ := range walker.m {
		delete(walker.m, k)
	}
	c.Assert(len(walker.m), check.Equals, 0)

	err = trie.Walk(walker.walk)
	c.Assert(err, check.IsNil)
	c.Assert(len(walker.m), check.Equals, 0)
}

func (ts *TrieSuites) TestFreeBasic(c *check.C) {
	walker := &testWalk{
		m: make(map[string]interface{}),
	}
	trie := NewTrie()
	ilen := 10
	prefixes := make(map[string]interface{})
	for i := 0; i < ilen; i++ {
		prefixes[ts.rand.String()] = i
	}

	for k, v := range prefixes {
		ret := trie.Put(k, v)
		c.Assert(ret, check.Equals, true)
	}

	err := trie.Walk(walker.walk)
	c.Assert(err, check.IsNil)

	for k, v := range prefixes {
		node := trie.Get(k)
		c.Assert(node, check.Equals, v)
	}

	for k, _ := range walker.m {
		delete(walker.m, k)
	}
	c.Assert(len(walker.m), check.Equals, 0)

	FreeTrie(trie)

	err = trie.Walk(walker.walk)
	c.Assert(err, check.IsNil)
	c.Assert(len(walker.m), check.Equals, 0)
}

func (ts *TrieSuites) TestOpNum(c *check.C) {
	u := uint16(0xff00)
	for i := 0; i < 10; i++ {
		u = incNum(u)
		d := getNum(u)
		c.Assert(d, check.Equals, uint8(i+1))
	}

	for i := 10; i > 0; i-- {
		u = decNum(u)
		d := getNum(u)
		c.Assert(d, check.Equals, uint8(i-1))
	}

	for i := 0; i < 10; i++ {
		u = setNo(u, uint8(i+1))
		d := getNo(u)
		c.Assert(d, check.Equals, uint8(i+1))
	}
}

func (ts *TrieSuites) BenchmarkInsertDeleteMany(c *check.C) {
	//trie := NewTrie()
	for i := 0; i < c.N; i++ {
		trie := NewTrie()
		ts.insertDeleteMany(c, trie, ts.prefixes)
	}
}

func (ts *TrieSuites) BenchmarkMapInsertDeleteMany(c *check.C) {
	for i := 0; i < c.N; i++ {
		m := make(map[string]int)
		ts.mapInsertDeleteMany(c, m, ts.prefixes)
	}
}

func (ts *TrieSuites) insertDeleteMany(c *check.C, trie *Trie, prefixes map[string]int) {
	for k, _ := range prefixes {
		ret := trie.Put(k, 1)
		c.Assert(ret, check.Equals, true)
	}

	//runtime.GC()

	/*for k, _ := range prefixes {
		ret := trie.Delete(k)
		c.Assert(ret, check.Equals, true)
	}*/

	for k, _ := range prefixes {
		ref := trie.Get(k)
		c.Assert(ref, check.NotNil)
	}
}

func (ts *TrieSuites) mapInsertDeleteMany(c *check.C, m map[string]int, prefixes map[string]int) {
	for k, v := range prefixes {
		m[k] = v
	}
	/*for k, _ := range prefixes {
		delete(m, k)
	}*/
	for k, _ := range prefixes {
		_, ok := m[k]
		c.Assert(ok, check.Equals, true)
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

func (ts *TrieSuites) getHalfSameStrings(num int) map[string]int {
	prefixes := make(map[string]int)
	prefix := "0123456789ABCEDF0123456789ABCEDF"
	for i := 0; i < num; {
		str := prefix + ts.rand32.String()
		_, ok := prefixes[str]
		if ok {
			continue
		}
		i++
		prefixes[str] = 1
	}

	return prefixes
}

func incNum(u uint16) uint16 {
	t := uint8(u & 0x00ff)
	t++
	u &= 0xff00
	u ^= (0x00ff & uint16(t))
	return u
}

func decNum(u uint16) uint16 {
	t := uint8(u & 0x00ff)
	t--
	u &= 0xff00
	u ^= (0x00ff & uint16(t))
	return u
}

func getNum(u uint16) uint8 {
	return uint8(u & 0x00ff)
}

func setNo(u uint16, d uint8) uint16 {
	u &= 0x00ff
	u ^= ((uint16(d) << 8) & 0xff00)
	return u
}

func getNo(u uint16) uint8 {
	return uint8((u >> 8) & 0x00ff)
}
