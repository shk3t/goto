package main

import (
	"fmt"
	"os"
)

func main() {
	cachehost := os.Getenv("CACHEHOST")
    if cachehost == "" {
        cachehost = "localhost"
    }

	fmt.Println(cachehost)
}