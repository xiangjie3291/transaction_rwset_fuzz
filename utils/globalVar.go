package utils

import (
	"fmt"
	"go/ast"
	"go/types"
)

var Log *Logger

var ChainMakerScriptsDir string = "./nodeControl/chainmaker/chainmaker-go/scripts"

var GlobalContractInfo *ContractInfo

type CandidateTypes struct {
	Types        map[types.Type]interface{}
	Confirm      bool
	ConfirmValue []interface{}
}

func (f *CandidateTypes) String() string {
	return fmt.Sprintf(`CandidateTypes{
		Len: %d,
		Types: %v,
		Confirm: %v,
		ConfirmValue: %v,
	}`,
		len(f.Types),
		f.Types,
		f.Confirm,
		f.ConfirmValue,
	)
}

type FuncAndParamsNameInfo struct {
	// 存储FuncName和InvokeName之间对应关系，每个FuncName对应一个InvokeName
	// 本工具全程使用FuncName,仅在调用过程中使用InvokeName
	InvokeName     string
	ParamsNameList []string
	Product        int
}

type ContractInfo struct {
	// 0. 合约AST Node
	Node ast.Node
	// 1. 合约路径
	ContractPath string
	// 2. 合约目录、合约名
	ContractDir string

	ContractName string
	// 3. 保存待测合约编译后二进制文件位置path, 路径与原始文件名有关
	ContractByteCodePath string
	// 4. 保存待测合约funcName和paramsName相关信息
	ContractFuncMap map[string]*FuncAndParamsNameInfo
	// 5. 保存合约中param与candidateTypes相关信息
	ParamAndCandidateTypes map[string]*CandidateTypes
}

func (info *ContractInfo) PrintContractInfo() {
	fmt.Println("Contract Path:", info.ContractPath)
	fmt.Println("Contract Dir :", info.ContractDir)
	fmt.Println("Contract Name:", info.ContractName)
	fmt.Println("Contract ByteCodePath:", info.ContractByteCodePath)
}
