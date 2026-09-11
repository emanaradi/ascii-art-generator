package asciiartweb

import (
	"strings"
)

func GenerateASCIIArt(text string, banner Banner) string {
	if text == "" {
		return ""
	}

	text = strings.ReplaceAll(text, `\n`, "\n")
	text = strings.ReplaceAll(text, `\t`, "    ")

	var output strings.Builder

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if line == "" {
			if i == len(lines)-1 {
				continue
			}

			// output.WriteString("$")
			output.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, char := range line {

				index := int(char) - 32

				if index < 0 || index >= len(banner) {
					continue
				}

				if len(banner[index]) < 8 {
					continue
				}

				output.WriteString(banner[index][row])

			}
			// output.WriteString("$")
			output.WriteString("\n")
		}
	}
	return output.String()

}
