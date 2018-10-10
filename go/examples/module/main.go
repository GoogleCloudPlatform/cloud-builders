package main

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
)

func main() {
	var x uint64 = 1
	for i := 0; i < 15; i++ {
		fmt.Printf("%d bytes is %s\n", x, humanize.Bytes(uint64(x)))
		x = x * 10
	}
}
