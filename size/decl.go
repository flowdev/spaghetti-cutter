package size

import (
	"go/ast"
	"log"
)

func sizeOfDecl(decl ast.Decl) uint {
	var size uint

	if isNilInterfaceOrPointer(decl) {
		return 0
	}

	switch d := decl.(type) {
	case *ast.FuncDecl:
		size += sizeOfFuncDecl(d)
	case *ast.GenDecl:
		size += sizeOfGenDecl(d)
	default:
		size = 1
		log.Printf("WARNING - Don't know size of unknown decl: %T", d)
	}
	return size
}

func sizeOfGenDecl(decl *ast.GenDecl) uint {
	var size uint

	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			size += sizeOfExpr(s.Type)
		case *ast.ValueSpec:
			size += uint(len(s.Names))
			for _, v := range s.Values {
				size += sizeOfExpr(v)
			}
		}
	}
	return size
}

func sizeOfFuncDecl(fun *ast.FuncDecl) uint {
	size := sizeOfFieldList(fun.Recv)
	size += sizeOfFuncType(fun.Type)
	size += sizeOfStmt(fun.Body)
	return size
}

func sizeOfFuncType(fun *ast.FuncType) uint {
	return 1 + sizeOfFieldList(fun.Params) + sizeOfFieldList(fun.Results)
}
