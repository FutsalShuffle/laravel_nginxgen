package laravel

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"laravel_nginxgen/common"
)

const (
	uriKey        = "'uri'"
	controllerKey = "'controller'"
	actionKey     = "'action'"
	methodsKey    = "'methods'"
)

type ArrayProcessor struct{}

func (ap *ArrayProcessor) Process(i ast.Vertex, r *common.Route) {
	ai := i.(*ast.ExprArrayItem)
	key := ai.Key
	//TODO: Сделать вид что этого тут нет
	switch aa := key.(type) {
	case *ast.ScalarString:
		if string(aa.Value) == uriKey {
			v := ai.Val
			uri := string(v.(*ast.ScalarString).Value)
			r.Uri = common.SanitizeStringQuotes(uri)
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
						r.Controller, a = common.ProcessControllerStringLaravel(c)
						r.Action = append(r.Action, a)
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
					r.Methods = append(r.Methods, common.SanitizeStringQuotes(string(mval)))
				}
			}
		}
	}
}
