package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	f, err := os.Create("path1_large.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := 1; i <= 1500; i++ {
		path := fmt.Sprintf("path%d", rand.Intn(500))
		size := rand.Intn(1_000_000) + 100
		fmt.Fprintf(f, "%s\t%d\n", path, size)
	}
}

