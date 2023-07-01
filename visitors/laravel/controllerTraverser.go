package laravel

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"laravel_nginxgen/common"
	"regexp"
	"strings"
)

type controllerTraverser struct {
	visitor.Null
	Params  common.Params
	methods []string
}

func NewControllerTraverser(methods []string) *controllerTraverser {
	return &controllerTraverser{
		Params:  common.Params{},
		methods: methods,
	}
}

func (ct *controllerTraverser) StmtClassMethod(n *ast.StmtClassMethod) {
	p := false
	for _, m := range ct.methods {
		if string(n.Name.(*ast.Identifier).Value) == m {
			p = true
			break
		}
	}

	if !p {
		return
	}

	phpdocs := n.Modifiers
	phpdoc := ""
	for _, v := range phpdocs {
		for _, iv := range v.(*ast.Identifier).IdentifierTkn.FreeFloating {
			phpdoc += string(iv.Value)
		}
	}

	ct.Params = ct.phpDocToParams(phpdoc)
}

func (ct *controllerTraverser) phpDocToParams(pd string) common.Params {
	p := common.Params{
		Enum:       map[string]string{},
		IntOnly:    []string{},
		StringOnly: []string{},
		LimitQuery: "",
	}

	enumRegex, _ := regexp.Compile(common.EnumParamRegex)
	intRegex, _ := regexp.Compile(common.IntParamRegex)
	stringRegex, _ := regexp.Compile(common.StringParamRegex)
	lqpRegex, _ := regexp.Compile(common.LimitQueryParamRegex)

	match := enumRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.Enum[key] = match[0][3]
	}

	match = intRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.IntOnly = append(p.IntOnly, key)
	}

	match = stringRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.StringOnly = append(p.StringOnly, key)
	}

	match = lqpRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		if p.LimitQuery != "" {
			key = "," + key
		}
		p.LimitQuery += key
	}

	return p
}
