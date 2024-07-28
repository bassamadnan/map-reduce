package internal

import (
	"fmt"
	"strconv"
)

// user input functions

// count the number of characters
func Map(_, value string) []KeyValue {
	charCount := make(map[rune]int)
	for _, char := range value {
		charCount[char]++
	}

	var result []KeyValue
	for char, count := range charCount {
		result = append(result, KeyValue{
			Key:   string(char),
			Value: count, // storing as int, not string
		})
	}
	return result
}

// aggregate data
func Reduce(key string, values []string) string {
	count := 0
	for _, v := range values {
		n, err := strconv.Atoi(v)
		if err != nil {
			print(err)
			continue
		}
		count += n
	}
	return fmt.Sprintf("%s\t%d", key, count)
}
