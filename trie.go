package bittrie

import "fmt"

// Trie should not be used in real code as a radix tree would be significantly
// more space and time efficient.
type Trie struct {
	Root node
}

type node struct {
	C0, C1 *node
	// Is this flag needed? Does leaf node imply IsTerminal?
	IsTerminal bool
}

type Item struct {
	V   uint64
	Len uint
}

func (i Item) String() string {
	return fmt.Sprintf("%0*b", i.Len, i.V)
}

// Insert inserts an item comprised of keyLen least significant bits of key
// into the true. Returns true if the item was inserted.
func (t *Trie) Insert(key uint64, keyLen uint) bool {
	n := &t.Root
	key <<= 64 - keyLen
	for i := uint(0); i < keyLen; i++ {
		if key&(1<<(63-i)) == 0 {
			if n.C0 == nil {
				n.C0 = &node{}
			}
			n = n.C0
		} else {
			if n.C1 == nil {
				n.C1 = &node{}
			}
			n = n.C1
		}

		if n.IsTerminal {
			return false
		}
	}
	n.IsTerminal = true
	return true
}

func (t Trie) Search(key uint64, keyLen uint) bool {
	n := &t.Root
	key <<= 64 - keyLen
	for i := uint(0); i < keyLen; i++ {
		if key&(1<<(63-i)) == 0 {
			if n.C0 == nil {
				return false
			}
			n = n.C0
		} else {
			if n.C1 == nil {
				return false
			}
			n = n.C1
		}
	}
	return n.IsTerminal
}

func (t *Trie) Allocate(keyLen uint) (i Item, ok bool) {
	return t.allocate(0, 0, keyLen, &t.Root)
}

func (t *Trie) allocate(key uint64, keyLen uint, keyLenMax uint, n *node) (i Item, ok bool) {
	if n.IsTerminal || keyLen == keyLenMax {
		return Item{}, false
	}

	if n.C0 == nil || n.C1 == nil {
		// It's free real estate.
		if n.C0 == nil {
			n.C0 = &node{}
			n = n.C0
			key <<= 1
		} else {
			n.C1 = &node{}
			n = n.C1
			key = key<<1 | 1
		}
		for i := keyLen + 1; i < keyLenMax; i++ {
			n.C0 = &node{}
			n = n.C0
			key <<= 1
		}
		n.IsTerminal = true
		return Item{V: key, Len: keyLenMax}, true
	}

	// Appending both 0 or 1 to key will yield a prefix of an existing key.
	// We must go deeper.
	if i, ok = t.allocate(key<<1, keyLen+1, keyLenMax, n.C0); ok {
		return i, ok
	}
	if i, ok = t.allocate(key<<1|1, keyLen+1, keyLenMax, n.C1); ok {
		return i, ok
	}
	return Item{}, false
}

func (t Trie) Enumerate() []Item {
	return t.enumerate(0, 0, &t.Root)
}

func (t Trie) enumerate(key uint64, keyLen uint, n *node) []Item {
	if n == nil {
		return nil
	}

	if n.IsTerminal {
		return []Item{{V: key, Len: keyLen}}
	}

	var keys []Item
	keys = append(keys, t.enumerate(key<<1, keyLen+1, n.C0)...)
	keys = append(keys, t.enumerate(key<<1|1, keyLen+1, n.C1)...)
	return keys
}
