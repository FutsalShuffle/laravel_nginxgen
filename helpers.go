package main

import (
	"strings"
)

func PathFromPsrNs(ns string, composer Composer) string {
	if ns == "" {
		return ""
	}

	var ins string
	var nns string
	for i, v := range composer.Autoload.Psr {
		if strings.Contains(ns, i) {
			ins = i
			nns += v
			break
		}
	}
	if ins == "" {
		return ""
	}

	tns := strings.Replace(ns, ins+"\\", "", 1)
	tns = strings.Replace(tns, "\\", "/", -1)
	tns = strings.Replace(tns, "'", "", -1)
	tns = strings.Replace(tns, "//", "/", -1)

	return nns + tns
}
