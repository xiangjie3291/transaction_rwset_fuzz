package getContractStaticInfo

import (
	"TransactionRwset/utils"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// 指针指向元素
func pointerToElem(t types.Type) types.Type {
	if t, ok := t.(*types.Pointer); ok {
		return pointerToElem(t.Elem())
	}
	return t
}

func getLookupIndexName(lookupInstr *ssa.Lookup) string {
	var paramName string
	re := regexp.MustCompile(`"([^"]*?)":`)
	match := re.FindStringSubmatch(lookupInstr.Index.String())
	if len(match) >= 1 {
		paramName = match[1]
	}

	return paramName
}

func dfs(val ssa.Value, visited map[ssa.Value]bool, path []ssa.Value, candidateTypes *utils.CandidateTypes) {
	if visited[val] {
		return
	}

	visited[val] = true
	path = append(path, val)

	hasChildren := false

	if referres := val.Referrers(); referres != nil {
		// hasChildren = true
		for _, instr := range *referres {
			switch instr := instr.(type) {

			case ssa.CallInstruction:
				// 处理函数调用
				callee := instr.Common().StaticCallee()
				if callee == nil {
					break
				}

				if fmt.Sprintf("%s", callee) == "encoding/json.Unmarshal" {
					hasChildren = true
					dfs(instr.Common().Args[1].(*ssa.MakeInterface).X, visited, path, candidateTypes)
					break
				}

				// 获取函数参数的实际值，并递归追踪
				for i, arg := range instr.Common().Args {
					if arg == val && len(instr.Common().Args) == len(callee.Params) {
						param := callee.Params[i]
						hasChildren = true
						dfs(param, visited, path, candidateTypes)
					}
				}

				// 对返回值进行递归追踪
				if retVal := instr.Value(); retVal != nil {
					hasChildren = true
					dfs(retVal, visited, path, candidateTypes)
				}

			case ssa.Value:
				hasChildren = true
				dfs(instr, visited, path, candidateTypes)
			}

		}
	}

	if !hasChildren {
		// fmt.Println("Path:")
		for _, instr := range path {
			// fmt.Println(instr, pointerToElem(instr.Type().Underlying()))
			// fmt.Println(ParseType(instr.Type()))
			candidateTypes.Types[pointerToElem(instr.Type().Underlying())] = utils.ParseType(instr.Type())
		}
		// fmt.Println("End of Path\n")
	}

	path = path[:len(path)-1]
}

// 以getArgs接口为起始，获取输入可能的数据类型
func GetParamCandidateTypeBySSA(contractPath string) map[string]*utils.CandidateTypes {
	path := filepath.Dir(contractPath)

	// 构建 go.mod 文件的路径
	modFilePath := filepath.Join(path, "go.mod")

	// 读取 go.mod 文件内容
	modFileBytes, err := os.ReadFile(modFilePath)
	if err != nil {
		fmt.Println("读取 go.mod 文件时出错:", err)
		return nil
	}

	// 解析 go.mod 文件
	modFile, err := modfile.Parse("go.mod", modFileBytes, nil)
	if err != nil {
		fmt.Println("解析 go.mod 文件时出错:", err)
		return nil
	}

	// 获取模块名称
	moduleName := modFile.Module.Mod.Path

	loadMode :=
		packages.NeedName |
			packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedModule |
			packages.NeedTypes |
			packages.NeedImports |
			packages.NeedSyntax |
			packages.NeedTypesInfo

	parseMode := parser.SkipObjectResolution

	// patterns := []string{dir}
	// patterns := []string{"./...", "chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"}
	// patterns := []string{"./..."}
	patterns := []string{"all"}

	pkgs, err := packages.Load(&packages.Config{
		Mode:    loadMode,
		Context: context.Background(),
		Env:     os.Environ(),
		Dir:     path,
		Tests:   false,
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			return parser.ParseFile(fset, filename, src, parseMode)
		},
	}, patterns...)
	if err != nil {
		fmt.Println(err)
	}

	ssaBuildMode := ssa.InstantiateGenerics // ssa.SanityCheckFunctions | ssa.GlobalDebug

	// Analyze the package.
	ssaProg, ssaPkgs := ssautil.Packages(pkgs, ssaBuildMode)

	ssaProg.Build()

	// fmt.Println(moduleName)

	var txSeqPkg *ssa.Package
	for _, pkg := range ssaPkgs {
		if pkg.Pkg.Path() == moduleName {
			txSeqPkg = pkg
			break
		}
	}

	txSeqPkg.Build()

	// var mainFn *ssa.Function
	var srcFns []*ssa.Function

	scope := txSeqPkg.Pkg.Scope()
	// fmt.Print("scope start:", scope)
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		if typeName, ok := obj.(*types.TypeName); ok {
			namedType := typeName.Type().(*types.Named)
			if _, ok := namedType.Underlying().(*types.Struct); ok {
				// 获取和分析 struct 的所有方法
				for i := 0; i < namedType.NumMethods(); i++ {
					meth := namedType.Method(i)

					if ssaFunc := txSeqPkg.Prog.FuncValue(meth); ssaFunc != nil {
						var addAnons func(f *ssa.Function)
						addAnons = func(f *ssa.Function) {
							srcFns = append(srcFns, f)
							for _, anon := range f.AnonFuncs {
								addAnons(anon)
							}
						}
						addAnons(ssaFunc)
					}
				}
			}
		}
	}

	paramCandidateTypes := make(map[string]*utils.CandidateTypes)

	for _, fun := range srcFns {
		for _, block := range fun.Blocks {
			for _, instr := range block.Instrs {
				switch v := instr.(type) {
				case *ssa.Call:
					call := v.Call
					// if strings.Contains(fmt.Sprintf("%s", call.Method), "InvokeContract") {
					if strings.Contains(fmt.Sprintf("%s", call.Method), "(chainmaker.org/chainmaker/contract-sdk-go/v2/sdk.SDKInterface).GetArgs()") {
						if referrers := v.Referrers(); referrers != nil {
							for _, instr := range *referrers {
								if _, ok := instr.(*ssa.Lookup); ok {
									// 获取变量name
									paramName := getLookupIndexName(instr.(*ssa.Lookup))
									//根据变量name获取对应的typeMap
									candidateTypes, ok := paramCandidateTypes[paramName]
									if !ok {
										candidateTypes = &utils.CandidateTypes{
											Types:        make(map[types.Type]interface{}),
											Confirm:      false,
											ConfirmValue: nil,
										}
									}

									visited := make(map[ssa.Value]bool)
									dfs(instr.(ssa.Value), visited, []ssa.Value{}, candidateTypes)

									paramCandidateTypes[paramName] = candidateTypes
								}
							}
						}
					}
				}
			}
		}
	}

	return paramCandidateTypes
}
