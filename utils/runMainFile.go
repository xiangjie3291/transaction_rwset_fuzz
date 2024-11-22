package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// trans "[10 64]" (string type) to struct
func transToTxRWSetStruct(line string) *common.TxRWSet {
	line = strings.Trim(line, "[]")
	parts := strings.Fields(line)

	// 创建一个字节数组
	var byteArray []byte

	// 遍历每个部分，将其转换为byte类型并加入字节数组
	for _, part := range parts {
		// 将字符串转换为整数
		num, err := strconv.Atoi(part)
		if err != nil {
			fmt.Println("转换错误:", err)
			return nil
		}
		// 将整数转换为byte并加入数组
		byteArray = append(byteArray, byte(num))
	}

	rwSet := &common.TxRWSet{}
	err := rwSet.Unmarshal(byteArray)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return rwSet
}

// 计算该种子见读写集最大相似度
func calculateMaxSimilarity(baseTxRwSet, compareTxRwSet *common.TxRWSet) float64 {
	maxSimilarity := 0.00

	if len(baseTxRwSet.TxReads) != 0 {
		if len(compareTxRwSet.TxWrites) != 0 {
			for _, read := range baseTxRwSet.TxReads {
				for _, write := range compareTxRwSet.TxWrites {
					maxSimilarity = Max(maxSimilarity, CalculateSimilarity(string(read.Key), string(write.Key)))
				}
			}
		}
	}

	if len(baseTxRwSet.TxWrites) != 0 {
		if len(compareTxRwSet.TxReads) != 0 {
			for _, write := range baseTxRwSet.TxWrites {
				for _, read := range compareTxRwSet.TxReads {
					maxSimilarity = Max(maxSimilarity, CalculateSimilarity(string(write.Key), string(read.Key)))
				}
			}
		}
	}

	//debug:
	//fmt.Println("max: ", maxSimilarity)

	return maxSimilarity
}

// 获取交易读写集,并计算出最大相似度(两函数读集合与写集合之间)
// result：交易一的读写集合，交易二的读写集合
func RunMainFile() (*common.TxRWSet, *common.TxRWSet, float64) {
	// 创建一个管道
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
		return nil, nil, 0
	}

	defer reader.Close()
	defer writer.Close()

	// 执行子程序
	cmd := exec.Command("go", "run", "./fuzzengine/contractExecute/main.go")

	// 设置子程序的环境变量以传递管道的文件描述符
	cmd.ExtraFiles = []*os.File{writer}

	// 创建一个丢弃标准输出的文件
	nullFile, err := os.Open(os.DevNull)
	if err != nil {
		fmt.Println("Error opening /dev/null:", err)
		return nil, nil, 0
	}
	defer nullFile.Close()

	cmd.Stdout = nullFile // 重定向标准输出到 /dev/null

	// 启动子程序
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return nil, nil, 0
	}

	rwSetList := make([]*common.TxRWSet, 0)

	// 在一个goroutine中读取子程序的自定义日志输出
	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			rwSetList = append(rwSetList, transToTxRWSetStruct(line))

		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from pipe:", err)
		}
	}()

	// 等待子程序结束
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return nil, nil, 0
	}

	if len(rwSetList) != 2 {
		fmt.Println("run main file result error!")
		return nil, nil, 0
	}

	// 交易一的读写集合，交易二的读写集合
	return rwSetList[0], rwSetList[1], calculateMaxSimilarity(rwSetList[0], rwSetList[1])
}
