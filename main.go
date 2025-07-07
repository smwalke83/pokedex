package main

import (
	"github.com/smwalke83/pokedex/internal/pokecache"
	"time"
)

func main() {
	interval := 5 * time.Second
	cache := pokecache.NewCache(interval)
	startRepl(cache)
}