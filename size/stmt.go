package size

import (
	"fmt"
	"go/ast"
	"log"
)

func sizeOfStmt(stmt ast.Stmt) uint {
	var size uint

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		size = sizeOfAssignStmt(s)
		fmt.Println("Size of assign stmt:", size)
	case *ast.BlockStmt:
		size = sizeOfBlockStmt(s)
		fmt.Println("Size of block stmt:", size)
	case *ast.ReturnStmt:
		size = sizeOfReturnStmt(s)
		fmt.Println("Size of return stmt:", size)
	case *ast.RangeStmt:
		size = sizeOfRangeStmt(s)
		fmt.Println("Size of range stmt:", size)
	case *ast.IfStmt:
		size = sizeOfIfStmt(s)
		fmt.Println("Size of if stmt:", size)
	case *ast.ForStmt:
		size = sizeOfForStmt(s)
		fmt.Println("Size of for stmt:", size)
	case *ast.SwitchStmt:
		size = sizeOfSwitchStmt(s)
		fmt.Println("Size of switch stmt:", size)
	case *ast.TypeSwitchStmt:
		size = sizeOfTypeSwitchStmt(s)
		fmt.Println("Size of type switch stmt:", size)
	case *ast.CaseClause:
		size = sizeOfCaseClause(s)
		fmt.Println("Size of case clause stmt:", size)
	case *ast.SelectStmt:
		size = sizeOfSelectStmt(s)
		fmt.Println("Size of select stmt:", size)
	case *ast.CommClause:
		size = sizeOfCommClause(s)
		fmt.Println("Size of comm clause stmt:", size)
	case *ast.ExprStmt:
		size = sizeOfExprStmt(s)
		fmt.Println("Size of expr stmt:", size)
	case *ast.SendStmt:
		size = sizeOfSendStmt(s)
	case *ast.BranchStmt:
		size = sizeOfBranchStmt(s)
	case *ast.GoStmt:
		size = sizeOfGoStmt(s)
	case *ast.LabeledStmt:
		size = sizeOfLabeledStmt(s)
	case *ast.IncDecStmt:
		size = sizeOfIncDecStmt(s)
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
	if assign == nil {
		return 0
	}

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
	if ret == nil {
		return 0
	}

	var size uint

	for _, expr := range ret.Results {
		size += sizeOfExpr(expr)
	}
	return size
}

func sizeOfRangeStmt(rng *ast.RangeStmt) uint {
	if rng == nil {
		return 0
	}

	return sizeOfExpr(rng.Key) + sizeOfExpr(rng.Value) + sizeOfExpr(rng.X) +
		sizeOfBlockStmt(rng.Body)
}

func sizeOfIfStmt(ifs *ast.IfStmt) uint {
	if ifs == nil {
		return 0
	}

	return sizeOfStmt(ifs.Init) + sizeOfExpr(ifs.Cond) +
		sizeOfBlockStmt(ifs.Body) + sizeOfStmt(ifs.Else)
}

func sizeOfForStmt(fors *ast.ForStmt) uint {
	if fors == nil {
		return 0
	}

	return sizeOfStmt(fors.Init) + sizeOfExpr(fors.Cond) + sizeOfStmt(fors.Post) +
		sizeOfBlockStmt(fors.Body)
}

func sizeOfSwitchStmt(swtch *ast.SwitchStmt) uint {
	if swtch == nil {
		return 0
	}

	return sizeOfStmt(swtch.Init) + sizeOfExpr(swtch.Tag) +
		sizeOfBlockStmt(swtch.Body)
}

func sizeOfTypeSwitchStmt(typswitch *ast.TypeSwitchStmt) uint {
	if typswitch == nil {
		return 0
	}

	return sizeOfStmt(typswitch.Init) + sizeOfStmt(typswitch.Assign) +
		sizeOfBlockStmt(typswitch.Body)
}

func sizeOfCaseClause(clause *ast.CaseClause) uint {
	if clause == nil {
		return 0
	}

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
	if sel == nil {
		return 0
	}

	return 1 + sizeOfBlockStmt(sel.Body)
}

func sizeOfCommClause(clause *ast.CommClause) uint {
	if clause == nil {
		return 0
	}

	size := sizeOfStmt(clause.Comm)
	for _, stmt := range clause.Body {
		size += sizeOfStmt(stmt)
	}
	return size
}

func sizeOfDeclStmt(decl *ast.DeclStmt) uint {
	if decl == nil {
		return 0
	}

	return sizeOfDecl(decl.Decl)
}

func sizeOfBranchStmt(branch *ast.BranchStmt) uint {
	if branch == nil {
		return 0
	}

	return 1
}

func sizeOfLabeledStmt(label *ast.LabeledStmt) uint {
	if label == nil {
		return 0
	}

	return 1 + sizeOfStmt(label.Stmt) // labels add cognitive load
}

func sizeOfGoStmt(gos *ast.GoStmt) uint {
	if gos == nil {
		return 0
	}

	return sizeOfCallExpr(gos.Call)
}

func sizeOfSendStmt(send *ast.SendStmt) uint {
	if send == nil {
		return 0
	}

	return sizeOfExpr(send.Chan) + sizeOfExpr(send.Value)
}

func sizeOfDeferStmt(defe *ast.DeferStmt) uint {
	if defe == nil {
		return 0
	}

	return sizeOfCallExpr(defe.Call)
}

func sizeOfIncDecStmt(incdec *ast.IncDecStmt) uint {
	if incdec == nil {
		return 0
	}

	return sizeOfExpr(incdec.X)
}

func sizeOfExprStmt(expr *ast.ExprStmt) uint {
	if expr == nil {
		return 0
	}

	return sizeOfExpr(expr.X)
}
