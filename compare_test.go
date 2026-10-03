package main

import "testing"

func TestGroups(t *testing.T) {
	g := Groups(map[string]string{"eth": "0x6080", "base": "0x6080", "bsc": "0x6081", "op": "0x"})
	if len(g) != 2 {
		t.Fatalf("groups %v", g)
	}
	if len(g[CodeHash("0x6080")]) != 2 {
		t.Fatal("eth and base should match")
	}
}
