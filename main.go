package main

import (
	"TransactionRwset/engine"
	"TransactionRwset/fuzz"
	nodecontrol "TransactionRwset/nodeControl"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var contractPath []string = []string{
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/encdata/enc_data.go",  // 部分读写集使用到哈希计算，超过限制时间
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/erc721/erc721.go",     // 无法触发冲突
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/erc1155/erc1155.go",   // 1
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/exchange/exchange.go", // 可能涉及到跨合约调用
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/fact/fact.go",         // success
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/itinerary/main.go",    // 无冲突
	"/data/wxj/transaction_rwset_fuzz/contract/contracts-go/raffle/raffle.go",       // 存在不可触及的分支
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/standard-evidence/evidence.go",
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/standard-identity/identity.go",
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/standard-nfa/nfa.go",
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/trace/trace.go", // success
	"/home/ChainMaker/transaction_rwset/contract/contracts-go/vote/vote.go",   // 输入参数为结构体，无法自动构造
}

func main() {
	// 创建一个通道来接收信号
	signalChan := make(chan os.Signal, 1)
	// 监听 SIGINT（Ctrl+C）和 SIGTERM（终止）信号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 使用 goroutine 来等待信号
	go func() {
		// 等待信号的到来
		sig := <-signalChan
		fmt.Printf("Received signal: %s\n", sig)

		// 执行清理操作
		engine.Stop()

		// 退出程序
		os.Exit(0)
	}()

	pairSeedsfilePath := flag.String("load", "", "Path to JSON file to load the pair seeds pool")
	flag.Parse()

	var pool *fuzz.FuncPairSeedsPool
	var err error
	engine.Start("/data/wxj/transaction_rwset_fuzz/contract/contracts-go/raffle/raffle.go")

	if *pairSeedsfilePath != "" {
		// 从 JSON 文件加载种子池
		fmt.Printf("Loading seed pool from file: %s\n", *pairSeedsfilePath)
		pool, err = fuzz.LoadPairSeedPoolFromFile(*pairSeedsfilePath)
		if err != nil {
			fmt.Printf("Error loading seed pool: %v\n", err)
			os.Exit(1)
		}
	} else {
		pool = engine.GenerateSeeds()
	}

	engine.HandleFuncPairSeedsPool(pool)
	nodecontrol.ChainmakerController.StopChainmaker()

}
