package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	st := time.Now()

	for i := 0; i < 100000000; i++ {
		testString("helloths", "ths")
	}

	fmt.Println(time.Since(st))
}

func fb(n int) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return fb(n-1) + fb(n-2)
	}
}

func testString(a, b string) bool {
	return strings.Contains(a, b)
}
