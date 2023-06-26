package laravel

import (
	"encoding/base64"
	"fmt"
	"github.com/VKCOM/php-parser/pkg/ast"
	"laravel_nginxgen/common"
	"laravel_nginxgen/php_deserialize"
	"os"
	"strings"
)

type SerializedProcessor struct {
	Routes []common.Route
}

func NewSerializedProcessor() *SerializedProcessor {
	return &SerializedProcessor{
		Routes: []common.Route{},
	}
}

func (sp *SerializedProcessor) Process(i ast.Vertex, r *common.Route) {
	fno := i.(*ast.ExprFunctionCall)
	fn := fno.Function.(*ast.Name).Parts[0]
	afn := string(fn.(*ast.NamePart).Value)
	if afn == "base64_decode" {
		arg := common.SanitizeStringQuotes(string(fno.Args[0].(*ast.Argument).Expr.(*ast.ScalarString).Value))
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
							continue
						}

						sm := common.Route{}
						obj := vvv.(*php_deserialize.PhpObject).GetMembers()
						methodsv := obj[metv].(php_deserialize.PhpArray)
						controllerv := php_deserialize.PhpValueString(obj[actv].(php_deserialize.PhpArray)[ctv])
						urival := php_deserialize.PhpValueString(obj[uriv])
						prc, pra := common.ProcessControllerStringLaravel(controllerv)
						sm.Controller = prc
						sm.Action = append(sm.Action, pra)
						sm.Uri = urival
						for _, method := range methodsv {
							sm.Methods = append(sm.Methods, php_deserialize.PhpValueString(method))
						}
						sp.Routes = append(sp.Routes, sm)
					}
				}
			}
		}
	}
}
