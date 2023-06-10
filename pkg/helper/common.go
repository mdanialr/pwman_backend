package helper

import "strings"

// Pad give blank space / padding to each given string before concatenate them.
func Pad(msg ...string) string {
	var b strings.Builder
	for _, m := range msg {
		b.WriteString(m)
		b.WriteString(" ")
	}
	// trim last padding
	return strings.TrimSpace(b.String())
}
