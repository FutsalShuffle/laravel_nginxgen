package main

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"regexp"
	"strings"
)

// Параметры phpdoc
const enumParamRegex = `@NGEnum (.+?( ))(.+)`
const intParamRegex = `@NGIntOnly (.+?( |\n))`
const stringParamRegex = `@NGStringOnly (.+?( |\n))`
const limitQueryParamRegex = `@NGQLimit (.+?( |\n))`

type controllerTraverser struct {
	visitor.Null
	Params  Params
	methods []string
}

func NewControllerTraverser(methods []string) *controllerTraverser {
	return &controllerTraverser{
		Params:  Params{},
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

func (ct *controllerTraverser) phpDocToParams(pd string) Params {
	p := Params{
		Enum:       map[string]string{},
		IntOnly:    []string{},
		StringOnly: []string{},
		LimitQuery: "",
	}

	enumRegex, _ := regexp.Compile(enumParamRegex)
	intRegex, _ := regexp.Compile(intParamRegex)
	stringRegex, _ := regexp.Compile(stringParamRegex)
	lqpRegex, _ := regexp.Compile(limitQueryParamRegex)

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

type Params struct {
	IntOnly    []string
	StringOnly []string
	Enum       map[string]string
	LimitQuery string
}
