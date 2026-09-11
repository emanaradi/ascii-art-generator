package asciiartweb

import (
	"fmt"
	"os"
	"strings"
)

type Banner [][]string

// parsing and loading the text files into a map for easy access to the ascii art representations of the characters
/**
* @return map[rune][]string - a map where the key is the character and the value is its ascii art representation in the form of a slice of strings, where each string represents a line of the ascii art
 */
func Loader(fileName string) (Banner, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	var banner Banner

	for i := 1; i+7 < len(lines); i += 9 {
		char := make([]string, 8)

		for j := 0; j < 8; j++ {
			char[j] = lines[i+j]
		}
		banner = append(banner, char)
	}

	return banner, nil
}
