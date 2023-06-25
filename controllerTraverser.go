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

type controllerTraverser struct {
	visitor.Null
	Params Params
	method string
}

func NewControllerTraverser(method string) *controllerTraverser {
	return &controllerTraverser{
		Params: Params{},
		method: method,
	}
}

func (ct *controllerTraverser) StmtClassMethod(n *ast.StmtClassMethod) {
	if string(n.Name.(*ast.Identifier).Value) != ct.method {
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
		ParamEnum:       map[string]string{},
		ParamIntOnly:    []string{},
		ParamStringOnly: []string{},
	}

	enumRegex, _ := regexp.Compile(enumParamRegex)
	intRegex, _ := regexp.Compile(intParamRegex)
	stringRegex, _ := regexp.Compile(stringParamRegex)

	match := enumRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.ParamEnum[key] = match[0][3]
	}

	match = intRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.ParamIntOnly = append(p.ParamIntOnly, key)
	}

	match = stringRegex.FindAllStringSubmatch(pd, -1)
	if len(match) > 0 {
		key := strings.TrimSpace(match[0][1])
		p.ParamStringOnly = append(p.ParamStringOnly, key)
	}

	return p
}

type Params struct {
	ParamIntOnly    []string
	ParamStringOnly []string
	ParamEnum       map[string]string
}
