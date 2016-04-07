package suffix

var (
	sets         = "0123456789ABCDEF"
	num_children = 16
	idx          map[byte]uint8
)

type WalkFunc func(key string, value interface{}) error

type Trie struct {
	parent   *Trie
	children []*Trie
	value    interface{}
	childIdx uint8
	childNum uint8
}

func (t *Trie) Get(key string) interface{} {
	node := t
	lkey := len(key)
	for i := 0; i < lkey; i++ {
		node = node.children[idx[key[i]]]
		if node == nil {
			return nil
		}
	}

	return node.value
}

func (t *Trie) Put(key string, value interface{}) bool {
	node := t
	lkey := len(key)
	for i := 0; i < lkey; i++ {
		pos, _ := idx[key[i]]
		child := node.children[pos]
		if child == nil {
			child = NewTrie()
			child.parent = node
			child.childIdx = pos
			node.children[pos] = child
			node.childNum++
		}
		node = child
	}

	isNew := (node.value == nil)
	node.value = value
	return isNew
}

func (t *Trie) Delete(key string) bool {
	node := t
	lkey := len(key)
	for i := 0; i < lkey; i++ {
		node = node.children[idx[key[i]]]
		if node == nil {
			return false
		}
	}

	node.value = nil

	if node.isLeaf() {
		for i := lkey - 1; i >= 0; i-- {
			childIdx := node.childIdx
			node = node.parent
			node.children[childIdx] = nil
			node.childNum--
			if node.value != nil || !node.isLeaf() {
				break
			}
		}
	}

	return true
}

func (t *Trie) Walk(walker WalkFunc) error {
	return t.walk("", walker)
}

func (t *Trie) walk(key string, walker WalkFunc) error {
	if t.value != nil {
		walker(key, t.value)
	}
	for i, child := range t.children {
		err := child.walk(key+string(sets[i]), walker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Trie) isLeaf() bool {
	return t.childNum == 0
	/*for i := 0; i < num_children; i++ {
		if t.children[i] != nil {
			return false
		}
	}
	return true*/
}

func NewTrie() *Trie {
	return &Trie{
		parent:   nil,
		children: make([]*Trie, num_children),
		value:    nil,
		childIdx: uint8(0),
	}
}

func InitSets() {
	bs := []byte(sets)
	idx = make(map[byte]uint8)
	for i, b := range bs {
		idx[b] = uint8(i)
	}
}
