package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func CodeHash(code string) string {
	b, _ := hex.DecodeString(strings.TrimPrefix(code, "0x"))
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// Groups maps each distinct hash to the chains that have it (empty code excluded).
func Groups(codes map[string]string) map[string][]string {
	out := map[string][]string{}
	for chain, code := range codes {
		if code == "" || code == "0x" {
			continue
		}
		h := CodeHash(code)
		out[h] = append(out[h], chain)
	}
	return out
}
