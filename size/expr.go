package size

import (
	"fmt"
	"go/ast"
	"log"
)

func sizeOfExpr(expr ast.Expr) uint {
	var size uint

	switch e := expr.(type) {
	case *ast.BasicLit:
		size = sizeOfBasicLit(e)
	case *ast.CompositeLit:
		size = sizeOfCompositeLit(e)
	case *ast.Ident:
		size = sizeOfIdent(e)
	case *ast.SelectorExpr:
		size = sizeOfSelectorExpr(e)
	case *ast.UnaryExpr:
		size = sizeOfExpr(e.X)
	case *ast.CallExpr:
		size = sizeOfCallExpr(e)
	case *ast.KeyValueExpr:
		size = sizeOfKeyValueExpr(e)
	case *ast.StructType:
		size = sizeOfStructType(e)
	case *ast.ArrayType:
		size = sizeOfArrayType(e)
	case *ast.MapType:
		size = sizeOfMapType(e)
	case *ast.ChanType:
		size = sizeOfChanType(e)
	case *ast.InterfaceType:
		size = sizeOfInterfaceType(e)
	case *ast.FuncType:
		size = sizeOfFuncType(e)
	case *ast.FuncLit:
		size = sizeOfFuncLit(e)
	case *ast.TypeAssertExpr:
		size = sizeOfTypeAssertExpr(e)
	case *ast.StarExpr:
		size = sizeOfExpr(e.X)
	case *ast.SliceExpr:
		size = sizeOfSliceExpr(e)
	case *ast.IndexExpr:
		size = sizeOfIndexExpr(e)
	case *ast.BinaryExpr:
		size = sizeOfBinaryExpr(e)
	case *ast.ParenExpr:
		size = sizeOfExpr(e.X)
	case *ast.Ellipsis: // TODO: Test as function parameter type!!!
		size = sizeOfEllipsis(e)
		fmt.Println("Size of ellipsis expr:", size)
	case nil:
		size = 0
	default:
		size = 1
		log.Printf("WARNING - Don't know size of unknown expr: %T", e)
	}
	return size
}

func sizeOfIdent(id *ast.Ident) uint {
	if id == nil {
		return 0
	}

	return 1
}

func sizeOfEllipsis(elli *ast.Ellipsis) uint {
	if elli == nil {
		return 0
	}

	return sizeOfExpr(elli.Elt)
}

func sizeOfBasicLit(lit *ast.BasicLit) uint {
	if lit == nil {
		return 0
	}

	return 1 + uint(len(lit.Value)/32)
}

func sizeOfCompositeLit(lit *ast.CompositeLit) uint {
	if lit == nil {
		return 0
	}

	size := sizeOfExpr(lit.Type)

	for _, elt := range lit.Elts {
		size += sizeOfExpr(elt)
	}
	return size
}

func sizeOfStructType(typ *ast.StructType) uint {
	if typ == nil {
		return 0
	}
	return 1 + sizeOfFieldList(typ.Fields)
}

func sizeOfFieldList(list *ast.FieldList) uint {
	if list == nil {
		return 0
	}

	var size uint

	for _, field := range list.List {
		size += sizeOfField(field)
	}
	return size
}

func sizeOfField(field *ast.Field) uint {
	if field == nil {
		return 0
	}

	return sizeOfExpr(field.Type) + sizeOfExpr(field.Tag)
}

func sizeOfMapType(m *ast.MapType) uint {
	if m == nil {
		return 0
	}

	return sizeOfExpr(m.Key) + sizeOfExpr(m.Value)
}

func sizeOfKeyValueExpr(kv *ast.KeyValueExpr) uint {
	if kv == nil {
		return 0
	}

	return sizeOfExpr(kv.Key) + sizeOfExpr(kv.Value)
}

func sizeOfArrayType(arr *ast.ArrayType) uint {
	if arr == nil {
		return 0
	}

	return sizeOfExpr(arr.Len) + sizeOfExpr(arr.Elt)
}

func sizeOfFuncLit(fun *ast.FuncLit) uint {
	if fun == nil {
		return 0
	}

	return sizeOfFuncType(fun.Type) + sizeOfBlockStmt(fun.Body)
}

func sizeOfSelectorExpr(sel *ast.SelectorExpr) uint {
	if sel == nil {
		return 0
	}

	return sizeOfExpr(sel.X) + sizeOfIdent(sel.Sel)
}

func sizeOfIndexExpr(idx *ast.IndexExpr) uint {
	if idx == nil {
		return 0
	}

	return sizeOfExpr(idx.X) + sizeOfExpr(idx.Index)
}

func sizeOfCallExpr(call *ast.CallExpr) uint {
	if call == nil {
		return 0
	}

	size := sizeOfExpr(call.Fun)

	for _, arg := range call.Args {
		size += sizeOfExpr(arg)
	}
	return size
}

func sizeOfBinaryExpr(bin *ast.BinaryExpr) uint {
	if bin == nil {
		return 0
	}

	return sizeOfExpr(bin.X) + sizeOfExpr(bin.Y)
}

func sizeOfSliceExpr(slice *ast.SliceExpr) uint {
	if slice == nil {
		return 0
	}

	return sizeOfExpr(slice.X) +
		sizeOfExpr(slice.Low) + sizeOfExpr(slice.High) + sizeOfExpr(slice.Max)
}

func sizeOfTypeAssertExpr(ass *ast.TypeAssertExpr) uint {
	if ass == nil {
		return 0
	}

	return sizeOfExpr(ass.X) + sizeOfExpr(ass.Type)
}

func sizeOfChanType(ch *ast.ChanType) uint {
	if ch == nil {
		return 0
	}

	return sizeOfExpr(ch.Value)
}

func sizeOfInterfaceType(iface *ast.InterfaceType) uint {
	if iface == nil {
		return 0
	}

	return 1 + sizeOfFieldList(iface.Methods)
}
