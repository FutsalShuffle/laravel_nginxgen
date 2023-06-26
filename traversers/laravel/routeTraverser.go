package laravel

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"laravel_nginxgen/common"
	"laravel_nginxgen/decoders/laravel"
)

type routeArrayTraverser struct {
	visitor.Null
	Routes    map[string]common.Route
	TempRoute common.Route
}

func NewRouteTraverser() *routeArrayTraverser {
	return &routeArrayTraverser{
		Routes:    map[string]common.Route{},
		TempRoute: common.Route{},
	}
}

func (rat *routeArrayTraverser) ExprArray(n *ast.ExprArray) {
	proc := laravel.ArrayProcessor{}
	for _, i := range n.Items {
		proc.Process(i, &rat.TempRoute)
	}

	rat.LeaveNode(n)
}

func (rat *routeArrayTraverser) ExprFunctionCall(n *ast.ExprFunctionCall) {
	proc := laravel.NewSerializedProcessor()

	for _, i := range proc.Routes {
		rat.TempRoute = i
		rat.LeaveNode(n)
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

	rat.TempRoute = common.Route{}
}
