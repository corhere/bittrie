package bittrie

import (
	"encoding/binary"
	"net/netip"
)

type IP4Trie struct {
	base netip.Prefix
	// Only the variable portion of each prefix p
	// (i.e. the bits between base.Bits() and p.Bits())
	// is stored in the trie.
	tr Trie
}

// TODO: extend Trie to support 128-bit keys so IPv6 prefixes can be allocated.

func NewIP4Trie(base netip.Prefix) IP4Trie {
	if !base.Addr().Is4() {
		panic("ipv4 only")
	}
	return IP4Trie{
		base: base.Masked(),
	}
}

func (tr *IP4Trie) Insert(p netip.Prefix) bool {
	k, l, ok := tr.keyOf(p)
	if !ok {
		return false
	}
	return tr.tr.Insert(k, l)
}

func (tr IP4Trie) Search(p netip.Prefix) bool {
	k, l, ok := tr.keyOf(p)
	if !ok {
		return false
	}
	return tr.tr.Search(k, l)
}

func (tr IP4Trie) keyOf(p netip.Prefix) (uint64, uint, bool) {
	// Since the base prefix of all keys in the trie is the same (tr.base),
	// The trie key of p only needs to be the variable portion of the prefix.
	// Lopping off the base prefix ensures we subdivide base, not a

	if !p.Addr().Is4() {
		panic("ipv4 only")
	}

	p = p.Masked()
	if !p.IsValid() || p.Bits() < tr.base.Bits() {
		return 0, 0, false
	}
	if netip.PrefixFrom(p.Addr(), tr.base.Bits()).Masked() != tr.base {
		return 0, 0, false
	}

	bits := binary.BigEndian.Uint32(p.Addr().AsSlice()) >> (32 - uint(p.Bits()))

	return uint64(bits), uint(p.Bits() - tr.base.Bits()), true
}

func (tr IP4Trie) prefixFrom(v Item) netip.Prefix {
	bits := v.Len + uint(tr.base.Bits())
	a := binary.BigEndian.Uint32(tr.base.Addr().AsSlice())
	a |= uint32(v.V) << (32 - bits)
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], a)
	return netip.PrefixFrom(
		netip.AddrFrom4(buf),
		int(bits),
	)
}

func (tr IP4Trie) Enumerate() []netip.Prefix {
	var res []netip.Prefix
	for _, v := range tr.tr.Enumerate() {
		res = append(res, tr.prefixFrom(v))
	}
	return res
}

func (tr *IP4Trie) Allocate(bits uint) (netip.Prefix, bool) {
	v, ok := tr.tr.Allocate(bits - uint(tr.base.Bits()))
	if !ok {
		return netip.Prefix{}, false
	}
	return tr.prefixFrom(v), true
}
