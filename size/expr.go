package size

import (
	"go/ast"
	"log"
	"reflect"
)

func sizeOfExpr(expr ast.Expr) uint {
	var size uint

	if isNilInterfaceOrPointer(expr) {
		return 0
	}

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
	case *ast.Ellipsis:
		size = sizeOfEllipsis(e)
	case nil:
		size = 0
	default:
		size = 1
		log.Printf("WARNING - Don't know size of unknown expr: %T", e)
	}
	return size
}

func sizeOfIdent(id *ast.Ident) uint {
	return 1
}

func sizeOfEllipsis(elli *ast.Ellipsis) uint {
	return sizeOfExpr(elli.Elt)
}

func sizeOfBasicLit(lit *ast.BasicLit) uint {
	return 1 + uint(len(lit.Value)/32)
}

func sizeOfCompositeLit(lit *ast.CompositeLit) uint {
	size := sizeOfExpr(lit.Type)

	for _, elt := range lit.Elts {
		size += sizeOfExpr(elt)
	}
	return size
}

func sizeOfStructType(typ *ast.StructType) uint {
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
	return sizeOfExpr(m.Key) + sizeOfExpr(m.Value)
}

func sizeOfKeyValueExpr(kv *ast.KeyValueExpr) uint {
	return sizeOfExpr(kv.Key) + sizeOfExpr(kv.Value)
}

func sizeOfArrayType(arr *ast.ArrayType) uint {
	return sizeOfExpr(arr.Len) + sizeOfExpr(arr.Elt)
}

func sizeOfFuncLit(fun *ast.FuncLit) uint {
	return sizeOfExpr(fun.Type) + sizeOfStmt(fun.Body)
}

func sizeOfSelectorExpr(sel *ast.SelectorExpr) uint {
	return sizeOfExpr(sel.X) + sizeOfExpr(sel.Sel)
}

func sizeOfIndexExpr(idx *ast.IndexExpr) uint {
	return sizeOfExpr(idx.X) + sizeOfExpr(idx.Index)
}

func sizeOfCallExpr(call *ast.CallExpr) uint {
	size := sizeOfExpr(call.Fun)

	for _, arg := range call.Args {
		size += sizeOfExpr(arg)
	}
	return size
}

func sizeOfBinaryExpr(bin *ast.BinaryExpr) uint {
	return sizeOfExpr(bin.X) + sizeOfExpr(bin.Y)
}

func sizeOfSliceExpr(slice *ast.SliceExpr) uint {
	return sizeOfExpr(slice.X) +
		sizeOfExpr(slice.Low) + sizeOfExpr(slice.High) + sizeOfExpr(slice.Max)
}

func sizeOfTypeAssertExpr(ass *ast.TypeAssertExpr) uint {
	return sizeOfExpr(ass.X) + sizeOfExpr(ass.Type)
}

func sizeOfChanType(ch *ast.ChanType) uint {
	return sizeOfExpr(ch.Value)
}

func sizeOfInterfaceType(iface *ast.InterfaceType) uint {
	return 1 + sizeOfFieldList(iface.Methods)
}

func isNilInterfaceOrPointer(v interface{}) bool {
	return v == nil ||
		(reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
