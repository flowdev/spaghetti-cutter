package size

import (
	"fmt"
	"go/ast"
)

func sizeOfDecl(decl ast.Decl) uint {
	var size uint

	switch d := decl.(type) {
	case *ast.FuncDecl:
		size += sizeOfFuncDecl(d)
	case *ast.GenDecl:
		size += sizeOfGenDecl(d)
	default:
		size = 1
		fmt.Printf("Size of unknown decl: %T 1\n", d)
	}
	return size
}

func sizeOfGenDecl(decl *ast.GenDecl) uint {
	var size uint
	names := ""

	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			size += sizeOfExpr(s.Type)
			names += s.Name.Name + " "
		case *ast.ValueSpec:
			size += uint(len(s.Names))
			for _, v := range s.Values {
				size += sizeOfExpr(v)
			}
			names += fmt.Sprint(s.Names) + " "
		}
	}
	fmt.Println("Size of gen decl:", names, size)
	return size
}

func sizeOfFuncDecl(fun *ast.FuncDecl) uint {
	size := sizeOfFieldList(fun.Recv)
	size += sizeOfFuncType(fun.Type)
	size += sizeOfBlockStmt(fun.Body)

	fmt.Println("Size of func:", fun.Name.Name, size)
	return size
}

func sizeOfFuncType(fun *ast.FuncType) uint {
	return sizeOfFieldList(fun.Params) + sizeOfFieldList(fun.Results)
}
