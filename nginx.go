package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

const gateHandle = "    try_files $uri $uri/ /index.php?$query_string;"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const intOnlyPath = "([0-9]*)"
const stringOnlyPath = "(\\D+)"
const paramSeparator = ","

type NginxFormatter struct {
	Rules map[string]string
}

func (nf *NginxFormatter) AddToRules(route Route) {
	formattedRoute, hdp := nf.formatRouteToRegex(route)
	_, ok := nf.Rules[formattedRoute]
	if ok {
		fmt.Println("Duplicated route", formattedRoute)
		return
	}
	nf.Rules[formattedRoute] = nf.makeLocationBlock(route, formattedRoute, hdp)
}

func (nf *NginxFormatter) WriteToConf(fp string) {
	res := ""
	if len(nf.Rules) == 0 {
		res = "location / {\n        try_files $uri $uri/ /index.php?$query_string;\n}"
	}
	for _, i := range nf.Rules {
		res += i
	}
	_ = os.WriteFile(fp, []byte(res), 0644)
}

func (nf *NginxFormatter) formatRouteToRegex(route Route) (string, bool) {
	regexUri, _ := regexp.Compile(`({.+?(}))`)
	regexUriEmpty, _ := regexp.Compile(`{}`)
	path := "(?P<item>.+)"
	match := regexUri.FindAllStringSubmatch(route.Uri, -1)
	hdp := false
	res := route.Uri
	for _, m := range match {
		key := nf.sanitizeRgKey(m[0])
		path = "(?P<" + key + nf.randString(5) + ">.+)"
		if len(route.Params.IntOnly) > 0 {
			for _, i := range route.Params.IntOnly {
				if i == key {
					path = intOnlyPath
				}
			}
		}

		if len(route.Params.StringOnly) > 0 {
			for _, i := range route.Params.StringOnly {
				if i == key {
					path = stringOnlyPath
				}
			}
		}

		if len(route.Params.Enum) > 0 {
			i, exists := route.Params.Enum[key]
			if exists {
				//Заменяем , на | для реги, убираем пробелы между
				values := strings.Replace(i, paramSeparator, "|", -1)
				values = strings.Replace(values, " ", "", -1)
				path = "(" + values + ")"
			}
		}

		hdp = true
		res = strings.Replace(res, m[0], path, 1)
	}

	res = regexUriEmpty.ReplaceAllString(res, path)
	return res, hdp
}

func (nf *NginxFormatter) randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (nf *NginxFormatter) makeLocationBlock(route Route, formattedRoute string, hdp bool) string {
	body := ""
	if len(route.Methods) > 0 {
		rm := strings.Join(route.Methods, "|")
		body += "    if ($request_method !~ ^(" + rm + ")$) {\n        return 404;\n    }\n"
	}

	if route.Params.LimitQuery != "" {
		sp := strings.Split(route.Params.LimitQuery, paramSeparator)
		if len(sp) > 0 {
			qlr := "^$"
			for _, e := range sp {
				qlr += "|(^|&)" + e + "="
			}

			body += "    if ($args !~* " + qlr + ") {\n        return 404;\n    }\n"
		}
	}

	body += gateHandle
	if formattedRoute != "/" {
		formattedRoute = "/" + formattedRoute
	}
	rp := "="
	if hdp {
		rp = "~"
	}

	return "location " + rp + " " + formattedRoute + " " +
		"{\n" +
		body +
		"\n" +
		"}\n"
}

func (nf *NginxFormatter) sanitizeRgKey(key string) string {
	nkey := strings.Replace(key, "{", "", -1)
	nkey = strings.Replace(nkey, "}", "", -1)
	nkey = strings.Replace(nkey, "?", "", -1)

	return nkey
}
