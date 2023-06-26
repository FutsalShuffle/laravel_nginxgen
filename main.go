package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const cacheRoutePath = "/bootstrap/cache/"

func main() {
	nf := &NginxFormatter{
		Rules: map[string]string{},
	}
	projectPath := flag.String("project", ".", "laravel project root path")
	outputPath := flag.String("output", "./locations.conf", "output path for nginx config")
	flag.Parse()
	fmt.Println(*projectPath)
	composer := parseComposerJson(*projectPath)
	router := getRoutes(*projectPath)
	src, _ := os.ReadFile(*projectPath + cacheRoutePath + router)
	code, _ := ParsePhp(src, composer)
	//DumpToStd(code)
	//os.Exit(1)
	v := NewRouteTraverser()
	Traverse(v, code)

	for _, e := range v.Routes {
		path := composer.PathFromPsrNs(e.Controller)
		if path != "" {
			cf, _ := os.ReadFile(*projectPath + "/" + path + ".php")
			cc, _ := ParsePhp(cf, composer)
			ct := NewControllerTraverser(e.Action)
			Traverse(ct, cc)
			e.Params = &ct.Params
		}
		nf.AddToRules(e)
	}

	nf.WriteToConf(*outputPath)

	os.Exit(0)
}

func getRoutes(path string) string {
	entries, err := os.ReadDir(path + cacheRoutePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var router string

	for _, e := range entries {
		if strings.Contains(e.Name(), "routes") {
			router = e.Name()
		}
	}

	return router
}
