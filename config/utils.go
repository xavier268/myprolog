package config

import "strings"

// CutString is a utility to intelligently cut the string in pieces with a max width.
func CutString(s string, width int) (lines []string) {

	words := strings.Split(s, " ")
	line := ""
	for _, w := range words {
		w = strings.TrimSpace(w)
		if len(w)+len(line) <= width {
			line = line + " " + w
		} else {
			lines = append(lines, strings.TrimSpace(line))
			line = w
		}
	}
	line = strings.TrimSpace(line)
	if len(line) != 0 {
		lines = append(lines, line)
	}

	return lines
}
