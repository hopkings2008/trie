package trie

import (
	"encoding/gob"
	"fmt"
	"io"
	"sync"

	trie "trie/lib/suffix"
)

type Trie struct {
	root  *trie.Trie
	mutex *sync.Mutex
}

type NodeInfo struct {
	ref int
}

type FileNode struct {
	Prefix string
	Ref    int
}

type Selector interface {
	Check(prefix string, node *NodeInfo) bool
	Get(prefix string, node *NodeInfo) error
}

func CreateTrie() *Trie {
	return &Trie{
		root:  trie.NewTrie(),
		mutex: &sync.Mutex{},
	}
}

func (tr *Trie) Insert(key string) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	item := tr.root.Get(key)
	if item != nil {
		node := tr.getNode(item)
		node.ref++
		//fmt.Printf("Got prefix %s, update ref to %d\n", key, node.ref)
		return nil
	}
	if ret := tr.root.Put(key, &NodeInfo{ref: 1}); ret {
		return nil
	}

	return fmt.Errorf("Failed to insert %s", key)
}

func (tr *Trie) GetRef(key string) (int, error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	if item := tr.root.Get(key); item != nil {
		node := tr.getNode(item)
		return node.ref, nil
	}
	return -1, fmt.Errorf("Not found %s", key)
}

func (tr *Trie) Update(key string, ref int) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if item := tr.root.Get(key); item != nil {
		node := tr.getNode(item)
		node.ref = ref
		return nil
	}

	if ret := tr.root.Put(key, &NodeInfo{ref: ref}); ret {
		return nil
	}
	return fmt.Errorf("Failed to update %s with ref %d", key, ref)
}

func (tr *Trie) Delete(key string) (bool, error) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if item := tr.root.Get(key); item != nil {
		node := tr.getNode(item)
		node.ref--
		if node.ref <= 0 {
			if d := tr.root.Delete(key); d {
				return true, nil
			}
			return false, fmt.Errorf("Failed to delete %s", key)
		}
	}
	return false, nil
}

func (tr *Trie) Select(selector Selector) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	vistor := func(prefix string, item interface{}) error {
		node := tr.getNode(item)
		if selector.Check(prefix, node) {
			if err := selector.Get(prefix, node); err != nil {
				return err
			}
		}
		return nil
	}

	return tr.root.Walk(vistor)
}

func (tr *Trie) Save(writer io.Writer) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	enc := gob.NewEncoder(writer)
	vistor := func(prefix string, item interface{}) error {
		node := tr.getNode(item)
		fileNode := FileNode{
			Prefix: prefix,
			Ref:    node.ref,
		}
		if err := enc.Encode(fileNode); err != nil {
			return fmt.Errorf("Failed to encode prefix %s with ref %d, err: %v", prefix, node.ref, err)
		}
		return nil
	}
	return tr.root.Walk(vistor)
}

func (tr *Trie) Load(reader io.Reader) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	decoder := gob.NewDecoder(reader)
	var fileNode FileNode
	for {
		err := decoder.Decode(&fileNode)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to decode %v, err: %v", reader, err)
		}
		tr.root.Put(fileNode.Prefix, &NodeInfo{ref: fileNode.Ref})
	}

	return nil
}

func (tr *Trie) Cleanup() {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()
	trie.FreeTrie(tr.root)
	tr.root = nil
}

func (tr *Trie) getNode(item interface{}) *NodeInfo {
	node, ok := item.(*NodeInfo)
	if !ok {
		return nil
	}
	return node
}
