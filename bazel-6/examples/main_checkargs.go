package main

import (
	"fmt"
	"os"
	"reflect"
)

var (
	ExpectedArgs = []string{"a", "b", "c"}
)

func main() {
	argv := os.Args[1:]
	if !reflect.DeepEqual(argv, ExpectedArgs) {
		fmt.Printf("Got unexpected args: %v\nExpected: %v\n", argv, ExpectedArgs)
		os.Exit(1)
	}
	fmt.Printf("Hello, %+v!\n", argv)
}
