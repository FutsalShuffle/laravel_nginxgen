package common

import "strings"

func ProcessControllerStringLaravel(cs string) (string, string) {
	pcs := SanitizeStringQuotes(cs)
	ss := strings.Split(pcs, "@")
	if len(ss) < 2 {
		return "", ""
	}

	return ss[0], ss[1]
}

func SanitizeStringQuotes(s string) string {
	return strings.Replace(s, "'", "", -1)
}
