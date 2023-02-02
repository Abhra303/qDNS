package trie

import "fmt"

type trieNode struct {
	Bit             byte
	AvaliableBitmap [16]byte
	Data            interface{}
	Children        []*trieNode
}

type trie struct {
	Root  *trieNode
	Flags int
}

func (tn *trieNode) append(char byte, data interface{}) (bool, error) {
	if char == 0 || data == nil {
		return false, fmt.Errorf("either the character or the data is empty")
	}
	charVal := int(char)
	if charVal > 128 {
		return false, fmt.Errorf("invalid character for name server")
	}

	if tn.exists(charVal) {
		return false, fmt.Errorf("can't append - charecter already exists")
	}

	tn.Children = append(tn.Children, &trieNode{Bit: char, Data: data})
	tn.setBitmap(charVal)

	return true, nil
}

func (tn *trieNode) exists(charVal int) bool {
	bytePos := charVal / 8

	if bytePos < len(tn.AvaliableBitmap) {
		var maskByte byte = 1 << (charVal % 8)

		resultBit := tn.AvaliableBitmap[bytePos] & maskByte
		if resultBit == maskByte {
			return true
		}
	}
	return false
}

func (tn *trieNode) setBitmap(charVal int) bool {
	bytePos := charVal / 8

	if bytePos < len(tn.AvaliableBitmap) {
		var maskByte byte = 1 << (charVal % 8)

		tn.AvaliableBitmap[bytePos] = tn.AvaliableBitmap[bytePos] | maskByte
		return true
	}
	return false
}

func (t *trie) Insert(prefix string, data interface{}) (bool, error) {
	return true, nil
}

func (t *trie) Update(prefix string, data interface{}) (bool, error) {
	return true, nil
}

func (t *trie) Delete(prefix string) (bool, error) {
	return true, nil
}

func (t *trie) Search(prefix string) (interface{}, error) {
	return nil, nil
}

func NewTrie() *trie {
	return &trie{Root: &trieNode{}}
}
