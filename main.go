package main

import (
	"flag"
	"fmt"
	"laravel_nginxgen/common"
	"laravel_nginxgen/visitors/laravel"
	"os"
	"strings"
)

func main() {

	projectPath := flag.String("project", ".", "laravel project root path")
	outputPath := flag.String("output", "./locations.conf", "output path for nginx config")
	nginxHandle := flag.String("nginx-handle", common.NginxGateHandle, "Success route handle (default: try_files $uri $uri/ /index.php?$query_string;)")
	wms := flag.String("nginx-wms", "404", "Nginx response status for wrong method")
	wqps := flag.String("nginx-wqps", "404", "Nginx response status for wrong query params")
	appendPost := flag.String("add-post", "y", "Add POST method to routes when PATCH/PUT is present (for web routes)")

	flag.Parse()

	composer := NewComposerFromFile(*projectPath)
	nf := NewNginxFormatter("    "+*nginxHandle, *wms, *wqps, *appendPost == "y")

	router := getRoutes(*projectPath)
	src, _ := os.ReadFile(*projectPath + common.LaravelCacheRoutePath + router)

	parser := NewPhpParser(composer)
	visitor := laravel.NewRouteTraverser()
	code, _ := parser.Parse(src)
	parser.Traverse(visitor, code)

	for _, e := range visitor.Routes {
		//Get Controller method phpdoc tags
		path := composer.PathFromPsrNs(e.Controller)
		if path != "" {
			cf, _ := os.ReadFile(*projectPath + "/" + path + ".php")
			cc, _ := parser.Parse(cf)
			ct := laravel.NewControllerTraverser(e.Action)
			parser.Traverse(ct, cc)
			e.Params = &ct.Params
		}
		//Add to nginx rule block
		nf.AddToRules(e)
	}

	nf.WriteToConf(*outputPath)

	os.Exit(0)
}

func getRoutes(path string) string {
	//Route cache file search
	entries, err := os.ReadDir(path + common.LaravelCacheRoutePath)
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
