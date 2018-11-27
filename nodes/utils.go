package nodes

import (
	"strings"
)

const Indentation = "    "

func Indentate(text string, level int) string {
	result := ""
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i != 0 {
			result += "\n"
		}
		result += strings.Repeat(Indentation, level)
		result += strings.Trim(line, " \n\t\r")
	}
	return result
}
