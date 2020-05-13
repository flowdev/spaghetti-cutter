package size

import (
	"go/ast"
	"log"
)

func sizeOfStmt(stmt ast.Stmt) uint {
	var size uint

	if isNilInterfaceOrPointer(stmt) {
		return 0
	}

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		size = sizeOfAssignStmt(s)
	case *ast.IncDecStmt:
		size = sizeOfIncDecStmt(s)
	case *ast.ReturnStmt:
		size = sizeOfReturnStmt(s)
	case *ast.ExprStmt:
		size = sizeOfExprStmt(s)
	case *ast.IfStmt:
		size = sizeOfIfStmt(s)
	case *ast.ForStmt:
		size = sizeOfForStmt(s)
	case *ast.RangeStmt:
		size = sizeOfRangeStmt(s)
	case *ast.BlockStmt:
		size = sizeOfBlockStmt(s)
	case *ast.SwitchStmt:
		size = sizeOfSwitchStmt(s)
	case *ast.TypeSwitchStmt:
		size = sizeOfTypeSwitchStmt(s)
	case *ast.CaseClause:
		size = sizeOfCaseClause(s)
	case *ast.SelectStmt:
		size = sizeOfSelectStmt(s)
	case *ast.CommClause:
		size = sizeOfCommClause(s)
	case *ast.SendStmt:
		size = sizeOfSendStmt(s)
	case *ast.BranchStmt:
		size = sizeOfBranchStmt(s)
	case *ast.GoStmt:
		size = sizeOfGoStmt(s)
	case *ast.LabeledStmt:
		size = sizeOfLabeledStmt(s)
	case *ast.DeferStmt:
		size = sizeOfDeferStmt(s)
	case *ast.DeclStmt:
		size = sizeOfDeclStmt(s)
	case *ast.EmptyStmt:
		size = 0
	case nil:
		size = 0
	default:
		size = 1
		log.Printf("WARNING - Don't know size of unknown stmt: %T", s)
	}
	return size
}

func sizeOfBlockStmt(block *ast.BlockStmt) uint {
	var size uint

	for _, stmt := range block.List {
		size += sizeOfStmt(stmt)
	}
	return size
}

func sizeOfAssignStmt(assign *ast.AssignStmt) uint {
	var size uint

	for _, expr := range assign.Lhs {
		size += sizeOfExpr(expr)
	}
	for _, expr := range assign.Rhs {
		size += sizeOfExpr(expr)
	}
	return size
}

func sizeOfReturnStmt(ret *ast.ReturnStmt) uint {
	var size uint

	for _, expr := range ret.Results {
		size += sizeOfExpr(expr)
	}
	return size
}

func sizeOfRangeStmt(rng *ast.RangeStmt) uint {
	return sizeOfExpr(rng.Key) + sizeOfExpr(rng.Value) + sizeOfExpr(rng.X) +
		sizeOfStmt(rng.Body)
}

func sizeOfIfStmt(ifs *ast.IfStmt) uint {
	return sizeOfStmt(ifs.Init) + sizeOfExpr(ifs.Cond) +
		sizeOfStmt(ifs.Body) + sizeOfStmt(ifs.Else)
}

func sizeOfForStmt(fors *ast.ForStmt) uint {
	return sizeOfStmt(fors.Init) + sizeOfExpr(fors.Cond) + sizeOfStmt(fors.Post) +
		sizeOfStmt(fors.Body)
}

func sizeOfSwitchStmt(swtch *ast.SwitchStmt) uint {
	return sizeOfStmt(swtch.Init) + sizeOfExpr(swtch.Tag) +
		sizeOfStmt(swtch.Body)
}

func sizeOfTypeSwitchStmt(typswitch *ast.TypeSwitchStmt) uint {
	return sizeOfStmt(typswitch.Init) + sizeOfStmt(typswitch.Assign) +
		sizeOfStmt(typswitch.Body)
}

func sizeOfCaseClause(clause *ast.CaseClause) uint {
	var size uint

	for _, expr := range clause.List {
		size += sizeOfExpr(expr)
	}
	for _, stmt := range clause.Body {
		size += sizeOfStmt(stmt)
	}
	return size
}

func sizeOfSelectStmt(sel *ast.SelectStmt) uint {
	return 1 + sizeOfStmt(sel.Body)
}

func sizeOfCommClause(clause *ast.CommClause) uint {
	size := sizeOfStmt(clause.Comm)
	for _, stmt := range clause.Body {
		size += sizeOfStmt(stmt)
	}
	return size
}

func sizeOfDeclStmt(decl *ast.DeclStmt) uint {
	return sizeOfDecl(decl.Decl)
}

func sizeOfBranchStmt(branch *ast.BranchStmt) uint {
	return 1
}

func sizeOfLabeledStmt(label *ast.LabeledStmt) uint {
	return 1 + sizeOfStmt(label.Stmt) // labels add cognitive load
}

func sizeOfGoStmt(gos *ast.GoStmt) uint {
	return 1 + sizeOfExpr(gos.Call)
}

func sizeOfSendStmt(send *ast.SendStmt) uint {
	return sizeOfExpr(send.Chan) + sizeOfExpr(send.Value)
}

func sizeOfDeferStmt(defe *ast.DeferStmt) uint {
	return 1 + sizeOfExpr(defe.Call)
}

func sizeOfIncDecStmt(incdec *ast.IncDecStmt) uint {
	return sizeOfExpr(incdec.X)
}

func sizeOfExprStmt(expr *ast.ExprStmt) uint {
	return sizeOfExpr(expr.X)
}
