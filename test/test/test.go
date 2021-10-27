package main

import (
	"fmt"
	"unicode/utf8"
)

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func main() {
	s := "Bob,Mark,"
	fmt.Println(s)
	s = trimLastChar(s)
	fmt.Println(s)
}