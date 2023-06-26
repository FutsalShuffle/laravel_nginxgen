package main

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
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
	Routes    []Route
	TempRoute Route
}

func NewRouteTraverser() *routeArrayTraverser {
	return &routeArrayTraverser{
		Routes:    []Route{},
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
							sm.Controller, sm.Action = rat.processControllerString(c)
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

func (rat *routeArrayTraverser) LeaveNode(n ast.Vertex) {
	if rat.TempRoute.Uri != "" {
		rat.Routes = append(rat.Routes, rat.TempRoute)
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
	Action     string
	Params     *Params
	Methods    []string
}
