package trie

import (
	"fmt"
	"io"
	"sync"

	trie "github.com/tchap/go-patricia/patricia"
)

type Trie struct {
	root  *trie.Trie
	mutex *sync.Mutex
}

type NodeInfo struct {
	ref int
}

func (tr *Trie) Insert(key string) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	item := tr.root.Get(trie.Prefix(key))
	if item != nil {
		node := tr.getNode(item)
		node.ref++
		return nil
	}
	if ret := tr.root.Insert(trie.Prefix(key), &NodeInfo{ref: 1}); ret {
		return nil
	}

	return fmt.Errorf("Failed to insert %s", key)
}

func (tr *Trie) Delete(key string) (string, error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if item := tr.root.Get(trie.Prefix(key)); item != nil {
		node := tr.getNode(item)
		if node.ref--; node.ref <= 0 {
			if d := tr.root.Delete(trie.Prefix(key)); d {
				return key, nil
			}
			return "", fmt.Errorf("Failed to delete %s", key)
		}
	}
	return "", nil
}

func (tr *Trie) Dump(writer io.Writer) error {
	/*vistor := func(prefix trie.Prefix, item trie.Item) error {
	}()
	return nil*/
}

func (tr *Trie) getNode(item trie.Item) *NodeInfo {
	node, ok := item.(*NodeInfo)
	if !ok {
		return nil
	}
	return node
}
