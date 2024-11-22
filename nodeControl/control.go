package nodecontrol

import (
	"TransactionRwset/utils"
	"fmt"
	"os"
	"os/exec"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

var ChainmakerController *NodeController = NewNodeController()

// 路径以main执行目录路径为基准
const (
	sdkConfPath  = "./sdk_config.yml"
	claimVersion = "1.0.0"
)

type NodeController struct {
	Client *sdk.ChainClient
}

func NewNodeController() *NodeController {
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)

	if err != nil {
		return nil
	}

	controller := &NodeController{
		Client: client,
	}
	return controller
}

func (n *NodeController) StartChainmaker() {
	scriptDir := utils.ChainMakerScriptsDir

	// 执行启动脚本
	cmdStart := exec.Command("sudo", "./cluster_quick_start.sh", "normal")
	cmdStart.Dir = scriptDir // 设置启动脚本的工作目录
	cmdStart.Stdout = os.Stdout
	cmdStart.Stderr = os.Stderr
	err := cmdStart.Run()
	if err != nil {
		fmt.Println("启动集群时出错：", err)
		return
	}
}

func (n *NodeController) StopChainmaker() {
	scriptDir := utils.ChainMakerScriptsDir

	// 执行停止脚本
	cmdStop := exec.Command("sudo", "./cluster_quick_stop.sh", "clean")
	cmdStop.Dir = scriptDir // 设置停止脚本的工作目录
	cmdStop.Stdout = os.Stdout
	cmdStop.Stderr = os.Stderr
	err := cmdStop.Run()
	if err != nil {
		fmt.Println("停止集群时出错：", err)
		return
	}
}

func (n *NodeController) GetChainClient() (*sdk.ChainClient, error) {
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)

	if err != nil {
		return nil, err
	}

	return client, nil

}

// 部署合约
func (n *NodeController) UserContractClaimCreate(claimContractName string, claimByteCodePath string,
	withSyncResult bool, isIgnoreSameContract bool) (string, error) {
	client := ChainmakerController.Client
	usernames := []string{"org1admin1", "org2admin1",
		"org3admin1", "org4admin1"}

	resp, err := n.createUserContract(client, claimContractName, claimVersion, claimByteCodePath,
		common.RuntimeType_DOCKER_GO, []*common.KeyValuePair{}, withSyncResult, usernames...)
	if err != nil {
		if !isIgnoreSameContract {
			return "", err
		} else {
			return "", fmt.Errorf("CREATE claim contract failed, err: %s, resp: %+v\n", err, resp)
		}
	}

	if resp != nil {
		return resp.TxId, nil
	}

	return "", nil
}

func (n *NodeController) createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string, runtime common.RuntimeType,
	kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	payload = client.AttachGasLimit(payload, &common.Limit{
		GasLimit: 60000000,
	})

	//endorsers, err := examples.GetEndorsers(payload, usernames...)
	endorsers, err := GetEndorsersWithAuthType(client.GetHashType(),
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	var createContractTimeout int64 = 5

	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return resp, err
	}

	err = CheckProposalRequestResp(resp, true)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// 调用合约
func (n *NodeController) UserContractInvoke(contractName, method string, kvs []*common.KeyValuePair, withSyncResult bool) (string, string, common.TxStatusCode, bool, error) {
	// 将FuncName转化为InvokeName
	client := ChainmakerController.Client
	method = utils.GlobalContractInfo.ContractFuncMap[method].InvokeName

	txId, message, code, success, err := n.invokeUserContract(client, contractName, method, "", kvs, withSyncResult, &common.Limit{GasLimit: 200000})
	if err != nil {
		return txId, message, code, success, err
	}

	return txId, message, code, success, nil
}

func (n *NodeController) invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair,
	withSyncResult bool, limit *common.Limit) (string, string, common.TxStatusCode, bool, error) {

	resp, err := client.InvokeContractWithLimit(contractName, method, txId, kvs, -1, withSyncResult, limit)

	if err != nil {
		return resp.GetTxId(), resp.Message, resp.Code, false, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		// fmt.Println(resp.TxId)
		return resp.GetTxId(), resp.Message, resp.Code, false, fmt.Errorf("invoke contract failed, %s[code:%d]/[method:%s]/[height:%d]/[message:%s]\n", resp.GetTxId(), resp.Code, method, resp.TxBlockHeight, resp.GetMessage())
	}

	return resp.GetTxId(), resp.Message, resp.Code, true, nil
}

// used to debug:
// 判断合约是否成功部署上链
func TestContractGetTxByTxId(txId string) *common.TransactionInfo {
	client := ChainmakerController.Client

	transactionInfo, err := client.GetTxByTxId(txId)
	if err != nil {
		fmt.Println(err)
	}
	return transactionInfo
}

func ShowPath() {
	cmdStart := exec.Command("pwd")
	cmdStart.Stdout = os.Stdout
	cmdStart.Stderr = os.Stderr
	cmdStart.Run()
}
