package getContractStaticInfo

import (
	"fmt"
	"go/ast"
	"strings"
)

// 获取合约中对应的函数名
// 及从客户端发送交易时的调用名

// 调用名称
var invokeNames []string

// 函数名称
var txsNames []string

func txsNameInspectSelectorExpr(n ast.Node) bool {
	if selectorExpr, ok := n.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name != "Error" {
		txsNames = append(txsNames, selectorExpr.Sel.Name)
		return false
	}
	return true
}

func txsNameInspectReturnExpr(n ast.Node) bool {
	if returnStmt, ok := n.(*ast.ReturnStmt); ok {
		ast.Inspect(returnStmt, txsNameInspectSelectorExpr)
		return false
	}
	return true
}

func txsNameInspectCaseClause(n ast.Node) bool {
	if caseStmt, ok := n.(*ast.CaseClause); ok && len(caseStmt.List) == 1 {
		if basicLit, ok := caseStmt.List[0].(*ast.BasicLit); ok {
			invokeNames = append(invokeNames, strings.Trim(basicLit.Value, "\""))

		} else if ident, ok := caseStmt.List[0].(*ast.Ident); ok {
			invokeNames = append(invokeNames, GetConstValue(ident, GlobalNode))

		}
		for _, caseStmtBody := range caseStmt.Body {
			ast.Inspect(caseStmtBody, txsNameInspectReturnExpr)
		}
	}
	return true
}

func txsNameInspectSwitchStmt(n ast.Node) bool {
	if switchStmt, ok := n.(*ast.SwitchStmt); ok {
		if ident, ok := switchStmt.Tag.(*ast.Ident); ok && ident.Name == "method" {
			ast.Inspect(switchStmt.Body, txsNameInspectCaseClause)
			return false
		}
	}
	return true
}

func txsNameInspectFuncDecl(n ast.Node) bool {
	if fun, ok := n.(*ast.FuncDecl); ok && fun.Name.Name == "InvokeContract" {
		ast.Inspect(fun.Body, txsNameInspectSwitchStmt)
		return false
	}
	return true
}

func GetTxsName(node ast.Node) ([]string, []string) {
	GlobalNode = node
	ast.Inspect(node, txsNameInspectFuncDecl)

	if len(txsNames) != len(invokeNames) {
		fmt.Println("Get Txs Name Error!")
		return nil, nil
	}

	return txsNames, invokeNames
}
