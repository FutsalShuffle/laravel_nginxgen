package laravel

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"laravel_nginxgen/common"
)

type DecodeInterface interface {
	Process(i ast.Vertex, r *common.Route)
}
