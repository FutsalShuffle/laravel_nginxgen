package main

import (
	"encoding/base64"
	"fmt"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"laravel_nginxgen/php_deserialize"
	"os"
	"strings"
)

const (
	uriKey        = "'uri'"
	controllerKey = "'controller'"
	actionKey     = "'action'"
	methodsKey    = "'methods'"
)

type routeArrayTraverser struct {
	visitor.Null
	Routes    map[string]Route
	TempRoute Route
}

func NewRouteTraverser() *routeArrayTraverser {
	return &routeArrayTraverser{
		Routes:    map[string]Route{},
		TempRoute: Route{},
	}
}

func (rat *routeArrayTraverser) ExprArray(n *ast.ExprArray) {
	sm := rat.TempRoute
	for _, i := range n.Items {
		ai := i.(*ast.ExprArrayItem)
		key := ai.Key
		//TODO: Сделать вид что этого тут нет
		switch aa := key.(type) {
		case *ast.ScalarString:
			if string(aa.Value) == uriKey {
				v := ai.Val
				uri := string(v.(*ast.ScalarString).Value)
				sm.Uri = rat.sanitizeString(uri)
			}
			if string(aa.Value) == actionKey {
				v := ai.Val.(*ast.ExprArray)
				for _, ii := range v.Items {
					aii := ii.(*ast.ExprArrayItem)
					ikey := aii.Key
					switch iaa := ikey.(type) {
					case *ast.ScalarString:
						if string(iaa.Value) == controllerKey {
							iv := aii.Val
							c := string(iv.(*ast.ScalarString).Value)
							var a string
							sm.Controller, a = rat.processControllerString(c)
							sm.Action = append(sm.Action, a)
						}
					}
				}
			}
			if string(aa.Value) == methodsKey {
				v := ai.Val.(*ast.ExprArray)
				for _, ii := range v.Items {
					maii := ii.(*ast.ExprArrayItem).Val
					switch maiit := maii.(type) {
					case *ast.ScalarString:
						mval := maiit.Value
						sm.Methods = append(sm.Methods, rat.sanitizeString(string(mval)))
					}
				}
			}
		}
	}

	rat.TempRoute = sm
	rat.LeaveNode(n)
}

func (rat *routeArrayTraverser) ExprFunctionCall(n *ast.ExprFunctionCall) {
	fn := n.Function.(*ast.Name).Parts[0]
	afn := string(fn.(*ast.NamePart).Value)
	if afn == "base64_decode" {
		arg := rat.sanitizeString(string(n.Args[0].(*ast.Argument).Expr.(*ast.ScalarString).Value))
		rawDecodedText, err := base64.StdEncoding.DecodeString(arg)
		if err != nil {
			fmt.Println("couldn't parse base64 encoded laravel routes!")
			os.Exit(1)
		}

		tt := strings.TrimSpace(string(rawDecodedText))

		decoder := php_deserialize.NewUnSerializer(tt)
		val, err := decoder.Decode()
		if err != nil {
			panic(err)
		}

		var metv php_deserialize.PhpValue
		metv = string("methods")
		var actv php_deserialize.PhpValue
		actv = string("action")
		var ctv php_deserialize.PhpValue
		ctv = string("controller")
		var uriv php_deserialize.PhpValue
		uriv = string("uri")

		for k, v := range val.(*php_deserialize.PhpObject).GetMembers() {
			if strings.Contains(php_deserialize.PhpValueString(k), "routes") {
				for _, vv := range v.(php_deserialize.PhpArray) {
					for _, vvv := range vv.(php_deserialize.PhpArray) {
						if vvv == nil {
							rat.LeaveNode(n)
							continue
						}
						sm := rat.TempRoute
						obj := vvv.(*php_deserialize.PhpObject).GetMembers()
						methodsv := obj[metv].(php_deserialize.PhpArray)
						controllerv := php_deserialize.PhpValueString(obj[actv].(php_deserialize.PhpArray)[ctv])
						urival := php_deserialize.PhpValueString(obj[uriv])
						prc, pra := rat.processControllerString(controllerv)
						sm.Controller = prc
						sm.Action = append(sm.Action, pra)
						sm.Uri = urival
						for _, method := range methodsv {
							sm.Methods = append(sm.Methods, php_deserialize.PhpValueString(method))
						}
						rat.TempRoute = sm

						rat.LeaveNode(n)
					}
				}
			}
		}
	}
}

func (rat *routeArrayTraverser) LeaveNode(n ast.Vertex) {
	if rat.TempRoute.Uri != "" {
		i, exists := rat.Routes[rat.TempRoute.Uri]
		if exists {
			em := i.Methods
			for _, m := range rat.TempRoute.Methods {
				//Для web роутов - браузер не делает put/patch
				if m == "PUT" {
					em = append(em, "POST")
				}
				em = append(em, m)
			}
			i.Methods = em
			ea := i.Action
			for _, a := range rat.TempRoute.Action {
				ea = append(ea, a)
			}
			i.Action = ea
			rat.Routes[rat.TempRoute.Uri] = i
		} else {
			rat.Routes[rat.TempRoute.Uri] = rat.TempRoute
		}
	}

	rat.TempRoute = Route{}
}

func (rat *routeArrayTraverser) processControllerString(cs string) (string, string) {
	pcs := rat.sanitizeString(cs)
	ss := strings.Split(pcs, "@")
	if len(ss) < 2 {
		return "", ""
	}

	return ss[0], ss[1]
}

func (rat *routeArrayTraverser) sanitizeString(s string) string {
	return strings.Replace(s, "'", "", -1)
}

type Route struct {
	Uri        string
	Controller string
	Action     []string
	Params     *Params
	Methods    []string
}
