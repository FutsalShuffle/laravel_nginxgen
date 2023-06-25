package main

import (
	"encoding/json"
	"fmt"
	"os"
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
