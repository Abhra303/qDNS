package trie

import (
	"fmt"
)

type TrieContext struct {
	KeyLimit uint16
}

type Trie interface {
	Put(key string, data interface{}) error
	Update(key string, data interface{}) error
	Delete(key string) (interface{}, error)
	Search(key string) ([]interface{}, error)
	IsEmpty() bool
}

type trie struct {
	root   *trieNode
	config trieConfig
}

func (t *trie) IsExceedingKeyLimit(key string) bool {
	return t.config.isExceedingKeyLimit(key)
}

func (t *trie) Put(key string, data interface{}) error {
	var err error
	if key == "" {
		t.root.put(data)
		return nil
	}
	if t.IsExceedingKeyLimit(key) {
		return fmt.Errorf("error: the given key exceeds the key length limit")
	}

	var node *trieNode = t.root
	var exist bool
	for i, char := range key {
		charByte := byte(char)
		if i == len(key)-1 {
			node, _, exist = iterateNode(node, charByte)
			if exist {
				node.put(data)
			} else {
				err = node.append(charByte, data)
				if err != nil {
					return err
				}
			}
			return nil
		}
		node, _, exist = iterateNode(node, charByte)
		if !exist {
			err = node.append(charByte, nil)
			if err != nil {
				return err
			}
			node, _, _ = iterateNode(node, charByte)
		}
	}
	return nil
}

func (t *trie) Update(key string, data interface{}) error {
	tn, err := t.search(key)
	if err != nil {
		return err
	}
	tn.put(data)
	return nil
}

func (t *trie) Delete(key string) (interface{}, error) {
	var node *trieNode = t.root
	var ci int = -1
	var exist bool
	for i, char := range key {
		charByte := byte(char)
		prevNode := node
		node, ci, exist = iterateNode(node, charByte)
		if !exist {
			return nil, fmt.Errorf("the given key doesn't exist")
		}

		if i == len(key)-1 {
			copy(prevNode.children[ci:], prevNode.children[ci+1:])
			prevNode.children[len(prevNode.children)-1] = nil // or the zero value of T
			prevNode.children = prevNode.children[:len(prevNode.children)-1]
		}
	}
	return node.data, nil
}

func (t *trie) Search(key string) ([]interface{}, error) {
	tn, err := t.search(key)
	if err != nil {
		return nil, err
	}

	return tn.data, nil
}

func (t *trie) IsEmpty() bool {
	return t.root.children == nil
}

func (t *trie) search(key string) (*trieNode, error) {
	var node *trieNode = t.root
	if key == "" {
		return node, nil
	}
	var exist bool
	for _, char := range key {
		charByte := byte(char)
		node, _, exist = iterateNode(node, charByte)
		if !exist {
			return nil, fmt.Errorf("the given key doesn't exist")
		}
	}
	return node, nil
}

func NewTrie(context *TrieContext) Trie {
	config := CreateNewTrieConfig(context)
	return &trie{root: &trieNode{}, config: config}
}
