package trie

import "fmt"

type trieNode struct {
	bit             byte
	avaliableBitmap [16]byte // bitmap of 128 chars
	data            []interface{}
	children        []*trieNode
}

func (tn *trieNode) put(data interface{}) {
	tn.data = append(tn.data, data)
}

func (tn *trieNode) append(char byte, data interface{}) error {
	if char == 0 {
		return fmt.Errorf("the given charecter is empty")
	}
	charVal := int(char)
	if charVal > 128 {
		return fmt.Errorf("invalid character for name server")
	}

	if tn.exists(charVal) {
		return fmt.Errorf("can't append - charecter already exists")
	}

	err := tn.setBitmap(charVal)
	if err != nil {
		return err
	}

	node := &trieNode{bit: char}
	node.data = append(node.data, data)
	tn.children = append(tn.children, node)

	return nil
}

// if char exists, return true
func (tn *trieNode) exists(charVal int) bool {
	bytePos := charVal / 8

	if bytePos < len(tn.avaliableBitmap) {
		var maskByte byte = 1 << (charVal % 8)

		resultBit := tn.avaliableBitmap[bytePos] & maskByte
		if resultBit == maskByte {
			return true
		}
	}
	return false
}

func (tn *trieNode) setBitmap(charVal int) error {
	if charVal == 0 {
		return fmt.Errorf("can't use zero as character value")
	}
	bytePos := charVal / 8

	if bytePos < len(tn.avaliableBitmap) {
		var maskByte byte = 1 << (charVal % 8)

		tn.avaliableBitmap[bytePos] = tn.avaliableBitmap[bytePos] | maskByte
		return nil
	}
	return fmt.Errorf("the given integer exceeds the character limit")
}

// gives the next node based on the character, returns
// the same node if not found with a false bool type
func iterateNode(tn *trieNode, char byte) (*trieNode, int, bool) {
	if tn.exists(int(char)) {
		for i, childNode := range tn.children {
			if childNode.bit == char {
				return childNode, i, true
			}
		}
	}
	return tn, -1, false
}
