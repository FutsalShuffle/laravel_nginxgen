package main

import (
	"fmt"
	"laravel_nginxgen/common"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

type NginxFormatter struct {
	Rules                 map[string]string
	gateHandle            string
	wrongMethodStatus     string
	wrongQueryParamStatus string
	appendPost            bool
}

func NewNginxFormatter(gh string, wms string, wqps string, ap bool) *NginxFormatter {
	return &NginxFormatter{
		Rules:                 map[string]string{},
		gateHandle:            gh,
		wrongMethodStatus:     wms,
		wrongQueryParamStatus: wqps,
		appendPost:            ap,
	}
}

func (nf *NginxFormatter) AddToRules(route common.Route) {
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
		res = "location / {\n    try_files $uri $uri/ /index.php?$query_string;\n}"
	}
	for _, i := range nf.Rules {
		res += i
	}
	_ = os.WriteFile(fp, []byte(res), 0644)
}

func (nf *NginxFormatter) formatRouteToRegex(route common.Route) (string, bool) {
	regexUri, _ := regexp.Compile(`({.+?(}))`)
	regexUriEmpty, _ := regexp.Compile(`{}`)
	path := "(?P<item>.+)"
	match := regexUri.FindAllStringSubmatch(route.Uri, -1)
	hdp := false
	res := route.Uri
	for _, m := range match {
		key := nf.sanitizeRgKey(m[0])
		path = "(?P<" + key + nf.randString(4) + ">.+)"

		if route.Params != nil {
			if len(route.Params.IntOnly) > 0 {
				for _, i := range route.Params.IntOnly {
					if i == key {
						path = common.NginxIntOnlyPath
					}
				}
			}

			if len(route.Params.StringOnly) > 0 {
				for _, i := range route.Params.StringOnly {
					if i == key {
						path = common.NginxStringOnlyPath
					}
				}
			}

			if len(route.Params.Enum) > 0 {
				i, exists := route.Params.Enum[key]
				if exists {
					//Заменяем , на | для реги, убираем пробелы между
					values := strings.Replace(i, common.PDocParamSeparator, "|", -1)
					values = strings.Replace(values, " ", "", -1)
					path = "(" + values + ")"
				}
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
		b[i] = common.RandomLetterBytes[rand.Int63()%int64(len(common.RandomLetterBytes))]
	}
	return string(b)
}

func (nf *NginxFormatter) makeLocationBlock(route common.Route, formattedRoute string, hdp bool) string {
	body := ""

	if len(route.Methods) > 0 {
		rm := nf.getRouteMethods(route.Methods)
		body += "    if ($request_method !~ ^(" + rm + ")$) {\n        return " + nf.wrongMethodStatus + ";\n    }\n"
	}

	if route.Params != nil {
		if route.Params.LimitQuery != "" {
			sp := strings.Split(route.Params.LimitQuery, common.PDocParamSeparator)
			if len(sp) > 0 {
				qlr := "^$"
				for _, e := range sp {
					qlr += "|(^|&)" + e + "="
				}

				body += "    if ($args !~* " + qlr + ") {\n        return " + nf.wrongQueryParamStatus + ";\n    }\n"
			}
		}
	}

	body += nf.gateHandle
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

func (nf *NginxFormatter) getRouteMethods(m []string) string {
	if nf.appendPost == false {
		return strings.Join(m, "|")
	}
	for _, i := range m {
		if i == "PUT" || i == "PATCH" {
			m = append(m, "POST")
			break
		}
	}

	return strings.Join(m, "|")
}
