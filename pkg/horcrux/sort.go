package horcrux

import (
	"bytes"
	"strconv"
	"unicode"
)

// NaturalSort implements sort.Interface based on natural numbers within the string
type NaturalSort []string

func (s NaturalSort) Len() int {
	return len(s)
}

// get a prefix number off of a string
func prefixInt(str string) (digits int, number int) {
	var b bytes.Buffer
	for i, v := range str {
		if unicode.IsDigit(v) {
			b.WriteRune(v)
		} else {
			parsedNumber, err := strconv.Atoi(b.String())
			if err != nil {
				panic(err)
			}
			return i + 1, parsedNumber
		}
	}
	return 0, 0
}

// is this string a valid integer
func isInt(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

// Less than the second arg
func (s NaturalSort) Less(i, j int) bool {
	first := s[i]
	second := s[j]

	index := -1
	for {
		index++
		if index >= len(first) || index >= len(second) {
			return first < second
		}

		firstChar := first[index : index+1]
		secondChar := second[index : index+1]

		if isInt(firstChar) && isInt(secondChar) {
			digits, first := prefixInt(first[index:])
			_, second := prefixInt(second[index:])
			if first != second {
				return first < second
			}
			index += digits
		}
	}
}

// Swap elements
func (s NaturalSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
