package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const intOnlyPath = "([0-9]*)"
const stringOnlyPath = "(\\D+)"

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
		res = "    location / {\n        try_files $uri $uri/ /index.php?$query_string;\n    }"
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
		if len(route.Params.ParamIntOnly) > 0 {
			for _, i := range route.Params.ParamIntOnly {
				if i == key {
					path = intOnlyPath
				}
			}
		}

		if len(route.Params.ParamStringOnly) > 0 {
			for _, i := range route.Params.ParamStringOnly {
				if i == key {
					path = stringOnlyPath
				}
			}
		}

		if len(route.Params.ParamEnum) > 0 {
			i, exists := route.Params.ParamEnum[key]
			if exists {
				//Заменяем , на | для реги, убираем пробелы между
				values := strings.Replace(i, ",", "|", -1)
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
	if route.Params != nil {
		//TODO: ?
	}
	if formattedRoute != "/" {
		formattedRoute = "/" + formattedRoute
	}
	rp := "="
	if hdp {
		rp = "~"
	}

	return "location " + rp + " " + formattedRoute + " " +
		"{\n   " +
		"try_files $uri $uri/ /index.php?$query_string;" +
		"\n" +
		"}\n"
}

func (nf *NginxFormatter) sanitizeRgKey(key string) string {
	nkey := strings.Replace(key, "{", "", -1)
	nkey = strings.Replace(nkey, "}", "", -1)
	nkey = strings.Replace(nkey, "?", "", -1)

	return nkey
}
