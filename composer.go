package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseComposerJson(d string) Composer {
	cv := Composer{}
	file, err := os.ReadFile(d + "/composer.json")
	if err != nil {
		fmt.Println("Can't open composer.json. Project path is probably invalid")
		os.Exit(1)
	}
	_ = json.Unmarshal([]byte(file), &cv)

	return cv
}

type Composer struct {
	Autoload struct {
		Psr map[string]string `json:"psr-4"`
	} `json:"autoload"`
	Require struct {
		Php string `json:"php"`
	}
}

func (c *Composer) GetPhpVersion() (int, int) {
	reg, _ := regexp.Compile(`(\d)`)
	match := reg.FindAllStringSubmatch(c.Require.Php, -1)
	if len(match) < 2 {
		return 8, 1
	}

	major, _ := strconv.Atoi(match[0][0])
	minor, _ := strconv.Atoi(match[1][0])

	return major, minor
}

func (c *Composer) PathFromPsrNs(ns string) string {
	if ns == "" {
		return ""
	}

	var ins string
	var nns string
	for i, v := range c.Autoload.Psr {
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
