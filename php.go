package main

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/conf"
	"github.com/VKCOM/php-parser/pkg/errors"
	"github.com/VKCOM/php-parser/pkg/parser"
	"github.com/VKCOM/php-parser/pkg/version"
	"github.com/VKCOM/php-parser/pkg/visitor/dumper"
	"github.com/VKCOM/php-parser/pkg/visitor/traverser"
	"log"
	"os"
)

type PhpParser struct {
}

func ParsePhp(code []byte, composer Composer) (ast.Vertex, error) {
	var parserErrors []*errors.Error
	errorHandler := func(e *errors.Error) {
		parserErrors = append(parserErrors, e)
	}

	majorVer, minorVer := composer.GetPhpVersion()

	rootNode, err := parser.Parse(code, conf.Config{
		Version:          &version.Version{Major: uint64(majorVer), Minor: uint64(minorVer)},
		ErrorHandlerFunc: errorHandler,
	})

	if err != nil {
		log.Fatal("Error:" + err.Error())
	}

	if len(parserErrors) > 0 {
		for _, e := range parserErrors {
			log.Println(e.String())
		}
		os.Exit(1)
	}

	return rootNode, nil
}

func DumpToStd(node ast.Vertex) {
	goDumper := dumper.NewDumper(os.Stdout)
	node.Accept(goDumper)
}

func Traverse(v ast.Visitor, code ast.Vertex) {
	traverser.NewTraverser(v).Traverse(code)
}
