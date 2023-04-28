package bittrie

import (
	"testing"
)

func enum(t *testing.T, tr Trie) {
	t.Helper()
	for _, v := range tr.Enumerate() {
		t.Logf("%v", v)
	}
}

func TestTrie(t *testing.T) {
	var tr Trie

	must := func(ok bool) {
		t.Helper()
		if !ok {
			t.Fatal("nope")
		}
	}
	must(tr.Insert(0b10, 2))
	must(tr.Insert(0b111, 3))
	must(tr.Insert(0, 2))
	must(tr.Insert(0b010, 3))
	enum(t, tr)

	t.Run("Search", func(t *testing.T) {
		for _, v := range []struct {
			key uint64
			len uint
			ok  bool
		}{
			{0b10, 2, true},
			{0b111, 3, true},
			{0, 2, true},
			{0b010, 3, true},
			{0b100, 3, false},
			{0b1000, 4, false},
		} {
			if tr.Search(v.key, v.len) != v.ok {
				t.Errorf("Search(%0*b, %[1]v) = %[3]v, want %v", v.len, v.key, !v.ok, v.ok)
			}
		}
	})

	for {
		i, ok := tr.Allocate(3)
		if !ok {
			break
		}
		t.Log(i)
	}
	enum(t, tr)
}

func TestAllocate(t *testing.T) {
	var tr Trie
	for {
		i, ok := tr.Allocate(3)
		if !ok {
			break
		}
		t.Log(i)
	}
	enum(t, tr)
}
