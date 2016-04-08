package suffix

var (
	sets         = "0123456789ABCDEF"
	num_children = 16
)

type WalkFunc func(key string, value interface{}) error

type Trie struct {
	parent   *Trie
	children []*Trie
	value    interface{}
	childIdx uint16
	childNum uint16
}

func (t *Trie) Get(key string) interface{} {
	node := t
	lkey := len(key)
	for i := 0; i < lkey; i++ {
		node = node.children[getIdx(key[i])]
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
		pos := getIdx(key[i])
		child := node.children[pos]
		if child == nil {
			child = NewTrie()
			child.parent = node
			child.childIdx = uint16(pos)
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
		node = node.children[getIdx(key[i])]
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
		if child == nil {
			continue
		}
		if err := child.walk(key+string(sets[i]), walker); err != nil {
			return err
		}
	}
	return nil
}

func (t *Trie) isLeaf() bool {
	return t.childNum == 0
}

func getIdx(b byte) uint8 {
	switch b {
	case 48:
		return 0
	case 49:
		return 1
	case 50:
		return 2
	case 51:
		return 3
	case 52:
		return 4
	case 53:
		return 5
	case 54:
		return 6
	case 55:
		return 7
	case 56:
		return 8
	case 57:
		return 9
	case 65:
		return 10
	case 66:
		return 11
	case 67:
		return 12
	case 68:
		return 13
	case 69:
		return 14
	case 70:
		return 15
	default:
		return 16
	}
}

func NewTrie() *Trie {
	return &Trie{
		parent:   nil,
		children: make([]*Trie, num_children),
		value:    nil,
		childIdx: uint16(0),
		childNum: uint16(0),
	}
}

func InitSets() {
	/*bs := []byte(sets)
	idx = make(map[byte]uint8)
	for i, b := range bs {
		idx[b] = uint8(i)
	}*/
}
