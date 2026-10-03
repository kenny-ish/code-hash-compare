package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

var rpcs = map[string]string{
	"eth":      "https://ethereum-rpc.publicnode.com",
	"base":     "https://mainnet.base.org",
	"arbitrum": "https://arb1.arbitrum.io/rpc",
	"optimism": "https://mainnet.optimism.io",
	"bsc":      "https://bsc-dataseed.bnbchain.org",
	"polygon":  "https://polygon-bor-rpc.publicnode.com",
}

func getCode(url, addr string) (string, error) {
	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": "eth_getCode", "params": []any{addr, "latest"}})
	c := http.Client{Timeout: 15 * time.Second}
	resp, err := c.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var r struct{ Result string }
	err = json.NewDecoder(resp.Body).Decode(&r)
	return r.Result, err
}

func main() {
	chains := flag.String("chains", "eth,base,arbitrum,optimism,bsc,polygon", "chains to compare")
	flag.Parse()
	addr := flag.Arg(0)
	if addr == "" {
		log.Fatal("usage: code-hash-compare [-chains a,b] ADDRESS")
	}
	codes := map[string]string{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, c := range strings.Split(*chains, ",") {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			code, err := getCode(rpcs[c], addr)
			if err != nil {
				log.Printf("%s: %v", c, err)
				return
			}
			mu.Lock()
			codes[c] = code
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	names := make([]string, 0, len(codes))
	for c := range codes {
		names = append(names, c)
	}
	sort.Strings(names)
	for _, c := range names {
		size := (len(codes[c]) - 2) / 2
		if size == 0 {
			fmt.Printf("%-9s no code\n", c)
			continue
		}
		fmt.Printf("%-9s %6d bytes  sha256 %s...\n", c, size, CodeHash(codes[c])[:16])
	}
	g := Groups(codes)
	switch len(g) {
	case 0:
		fmt.Println("\nno chain has code at this address")
	case 1:
		fmt.Println("\nall deployed copies are identical")
	default:
		fmt.Printf("\n%d different bytecodes at the same address\n", len(g))
	}
}
