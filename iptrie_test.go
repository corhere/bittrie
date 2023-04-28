package bittrie

import (
	"net/netip"
	"testing"
)

func enumIPs(t *testing.T, tr IP4Trie) {
	t.Helper()
	for _, v := range tr.Enumerate() {
		t.Logf("%v", v)
	}
}

func TestIP4Trie(t *testing.T) {
	tr := NewIP4Trie(netip.MustParsePrefix("192.168.0.0/16"))

	must := func(ok bool) {
		t.Helper()
		if !ok {
			t.Fatal("nope")
		}
	}
	must(tr.Insert(netip.MustParsePrefix("192.168.42.0/24")))
	must(tr.Insert(netip.MustParsePrefix("192.168.43.0/24")))
	must(tr.Insert(netip.MustParsePrefix("192.168.1.128/25")))
	must(tr.Insert(netip.MustParsePrefix("192.168.99.240/28")))
	enumIPs(t, tr)

	for i := uint(20); i < 30; i++ {
		a, ok := tr.Allocate(i)
		t.Logf("Allocate(%v) = %v, %v", i, a, ok)
	}
	enumIPs(t, tr)
	t.Fail()
}
