package getContractStaticInfo

import (
	"go/ast"
	"go/token"
	"strings"
)

var GlobalNode ast.Node

func GetConstValue(findIdent *ast.Ident, fun ast.Node) string {
	result := ""
	// fmt.Println(findIdent.Name)
	ast.Inspect(fun, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {

			ast.Inspect(genDecl, func(spec ast.Node) bool {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {

					if len(valueSpec.Names) > 0 && valueSpec.Names[0].Name == findIdent.Name && len((valueSpec.Values)) > 0 {
						if basicLit, ok := valueSpec.Values[0].(*ast.BasicLit); ok {
							result = strings.Trim(basicLit.Value, "\"")
						}
						return false
					}
				}
				return true
			})
		}
		return true
	})
	return result
}

// 右侧赋值为GetArgs，获取左侧值Params
// 左侧赋值为Params，获取右侧字段名
// 如果发现有index，则左侧为字段名，如果没有，再去找
// 在assignstmt的rhs中存在sdk.Instance.GetArgs接口
func hasGetArgs(assignstmt ast.AssignStmt) bool {
	result := false
	if len(assignstmt.Rhs) == 1 && len(assignstmt.Lhs) == 1 {
		ast.Inspect(assignstmt.Rhs[0], func(n ast.Node) bool {
			if selectorExpr, ok := n.(*ast.SelectorExpr); ok {
				if x, ok := selectorExpr.X.(*ast.SelectorExpr); ok && selectorExpr.Sel.Name == "GetArgs" {
					if x.X.(*ast.Ident).Name == "sdk" && x.Sel.Name == "Instance" {
						result = true
						return false
					}
				}
			}
			return true
		})
	}
	return result
}

// 获取所有被GetArgs赋值的Ident(map类型)
// assignstmt -> rhs -> callexpr -> none args
// 获取lhs.ident.name
func getArgsNamesByGetArgs(fun ast.Node) string {
	argName := ""

	ast.Inspect(fun, func(n ast.Node) bool {
		//在hasGetArgs中已经判断了len(rhs)==1 len(lhs)==1
		if assignStmt, ok := n.(*ast.AssignStmt); ok && hasGetArgs(*assignStmt) {
			if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok && callExpr.Args == nil {
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
					argName = ident.Obj.Name
				}
			}
		}
		return true
	})

	return argName
}

// 获取所有被GetArgs赋值的Ident(map类型)
func getParamsNamesByArgsName(argsName string, fun ast.Node) []string {
	Params := make([]string, 0)

	ast.Inspect(fun, func(n ast.Node) bool {
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			ast.Inspect(assignStmt, func(n ast.Node) bool {
				if indexExpr, ok := n.(*ast.IndexExpr); ok {
					if ident, ok := indexExpr.X.(*ast.Ident); ok && ident.Name == argsName {
						if index, ok := indexExpr.Index.(*ast.BasicLit); ok {
							// ParamTmp := &utils.Param{
							// 	Name:         strings.Trim(index.Value, "\""),
							// 	ActualType:   reflect.Interface,
							// 	FirstObj:     assignStmt.Lhs[0].(*ast.Ident).Obj,
							// 	ReadRelated:  true,
							// 	WriteRelated: true,
							// }

							// ParamTmp = getActualType(ParamTmp, fun)
							Params = append(Params, strings.Trim(index.Value, "\""))

						} else if index, ok := indexExpr.Index.(*ast.Ident); ok && index.Obj.Kind == ast.Con {
							// ParamTmp := &utils.Param{
							// 	Name:         getConstValue(index, GlobalNode),
							// 	ActualType:   reflect.String,
							// 	FirstObj:     assignStmt.Lhs[0].(*ast.Ident).Obj,
							// 	ReadRelated:  true,
							// 	WriteRelated: true,
							// }

							// ParamTmp = getActualType(ParamTmp, fun)
							Params = append(Params, GetConstValue(index, GlobalNode))
						}
					}
				}
				return true
			})
		}

		return true
	})

	return Params
}

// 直接通过GetArgs获取字段值
// 对一整个funcDel进行分析
// lhs len == 1
// assignstmt -> rhs -> callexpr -> args -> indexexpr -> index -> value
// 在indexexpr下有sdk.Instance.GetArgs接口时，将Index.Value加入ParamsNames中
func getParamsNameByGetArgs(fun ast.Node) []string {
	Params := make([]string, 0)

	ast.Inspect(fun, func(n ast.Node) bool {
		//在hasGetArgs中已经判断了len(rhs)==1 len(lhs)==1
		if assignStmt, ok := n.(*ast.AssignStmt); ok && hasGetArgs(*assignStmt) {
			if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok && len(callExpr.Args) == 1 {
				if indexExpr, ok := callExpr.Args[0].(*ast.IndexExpr); ok {
					if basicLit, ok := indexExpr.Index.(*ast.BasicLit); ok {
						// ParamTmp := &utils.Param{
						// 	Name:         strings.Trim(basicLit.Value, "\""),
						// 	ActualType:   reflect.String,
						// 	FirstObj:     assignStmt.Lhs[0].(*ast.Ident).Obj,
						// 	ReadRelated:  true,
						// 	WriteRelated: true,
						// }

						// ParamTmp = getActualType(ParamTmp, fun)
						Params = append(Params, strings.Trim(basicLit.Value, "\""))

						// debug
						// fmt.Println(basicLit.Value)
						// fmt.Println(assignStmt.Lhs[0].(*ast.Ident).Obj)
					} else if index, ok := indexExpr.Index.(*ast.Ident); ok && index.Obj.Kind == ast.Con {
						// ParamTmp := &utils.Param{
						// 	Name:         getConstValue(index, GlobalNode),
						// 	ActualType:   reflect.String,
						// 	FirstObj:     assignStmt.Lhs[0].(*ast.Ident).Obj,
						// 	ReadRelated:  true,
						// 	WriteRelated: true,
						// }

						// ParamTmp = getActualType(ParamTmp, fun)
						Params = append(Params, GetConstValue(index, GlobalNode))
					}
				}
			}
		}
		return true
	})

	return Params
}

// 先通过GetArgs赋值给变量，再获取字段值
// 对一整个funcDel进行分析
func getParamsName(fun ast.Node) []string {
	Params := make([]string, 0)

	// 根据argsName找到所有的ParamsName(间接赋值)
	// for _, argsName := range getArgsNamesByGetArgs(fun) {
	argName := getArgsNamesByGetArgs(fun)

	Params = append(Params, getParamsNamesByArgsName(argName, fun)...)
	// }

	// 直接通过GetArgs接口获取的ParamsName(直接赋值)
	Params = append(Params, getParamsNameByGetArgs(fun)...)

	return Params
}

func getParamsNamesByInvokeFunc(name string, n ast.Node) []string {
	argName := ""
	funcName := ""
	params := make([]string, 0)

	funcName = name

	ast.Inspect(n, func(n ast.Node) bool {
		if fun, ok := n.(*ast.FuncDecl); ok && fun.Name.Name == "InvokeContract" {
			argName = getArgsNamesByGetArgs(fun)

			// 找到method相关switch case
			ast.Inspect(fun.Body, func(n ast.Node) bool {
				if switchStmt, ok := n.(*ast.SwitchStmt); ok {
					if ident, ok := switchStmt.Tag.(*ast.Ident); ok && ident.Name == "method" {
						ast.Inspect(switchStmt.Body, func(n ast.Node) bool {
							if caseStmt, ok := n.(*ast.CaseClause); ok && len(caseStmt.List) > 0 {
								if basicLit, ok := caseStmt.List[0].(*ast.BasicLit); ok && strings.Trim(basicLit.Value, "\"") == funcName {
									params = append(params, getParamsNamesByArgsName(argName, caseStmt)...)
								}
							}
							return true
						})
						return false
					}
				}
				return true
			})
			return false
		}
		return true
	})

	return params
}

func GetFuncParams(funcName string, node ast.Node) []string {
	params := make([]string, 0)

	// 在对应函数的实现下找到输入参数
	ast.Inspect(node, func(n ast.Node) bool {
		if fun, ok := n.(*ast.FuncDecl); ok && fun.Name.Name == funcName {
			params = append(params, getParamsName(fun)...)
		}
		return true
	})

	// 在switch语句中找到输入参数
	params = append(params, getParamsNamesByInvokeFunc(funcName, node)...)

	return params
}
