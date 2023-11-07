package main

import (
	"github.com/Tungsten-180/nasal-ls/internal/astparser"
)

func main() {
	ast := astparser.GetAST("/home/tungsten/Documents/flight/B-1B/Nasal/oso_management.nas")
	astparser.Run(ast)
	ast.DumpCalls()
}
