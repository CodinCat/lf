package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func isRoot(name string) bool { return path.Dir(name) == name }

// This function converts a size in bytes to a human readable form. For this
// purpose metric suffixes are used (e.g. 1K = 1000). For values less than 10
// the first significant digit is shown, otherwise it is hidden. Numbers are
// always rounded down. For these reasons this function always show somewhat
// smaller values but it should be fine for most human beings.
func humanize(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%d", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / 1000
	for _, s := range suffix {
		if curr < 10 {
			return fmt.Sprintf("%.1f%s", curr-0.0499, s)
		} else if curr < 1000 {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= 1000
	}

	return ""
}

// This function extracts numbers from a string and returns with the rest.
// It is used for numeric sorting of files when the file name consists of
// both digits and letters.
//
// For instance if your input is 'foo123bar456' you get a slice of number
// consisting of elements '123' and '456' and rest of the string as a slice
// consisting of elements 'foo' and 'bar'. The last return argument denotes
// whether or not the first partition is a number.
func extractNums(s string) (nums []int, rest []string, numFirst bool) {
	var buf []rune

	r, _ := utf8.DecodeRuneInString(s)
	digit := unicode.IsDigit(r)
	numFirst = digit

	for i, c := range s {
		if unicode.IsDigit(c) == digit {
			buf = append(buf, c)
			if i != len(s)-1 {
				continue
			}
		}

		if digit {
			i, err := strconv.Atoi(string(buf))
			if err != nil {
				// TODO: handle error
				log.Printf("extracting numbers: %s", err)
			}
			nums = append(nums, i)
		} else {
			rest = append(rest, string(buf))
		}

		buf = nil
		buf = append(buf, c)
		digit = !digit
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// We don't need no generic code
// We don't need no thought control
// No dark templates in compiler
// Haskell leave them kids alone
// Hey Bjarne leave them kids alone
// All in all it's just another brick in the code
// All in all you're just another brick in the code
//
// -- Pink Trolled --
