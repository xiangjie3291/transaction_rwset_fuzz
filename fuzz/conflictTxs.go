package fuzz

import (
	nodecontrol "TransactionRwset/nodeControl"
	"TransactionRwset/picture"
	"TransactionRwset/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"golang.org/x/image/font/opentype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var maxWorkers int = 20

/*
 1. 确定最大并行限度 （以100笔交易为上限）
    只确定发送单类交易时，每100笔中最多能有多少交易
    直接计算比例发送。

2. 确定交易发送比例，发送比例分别为1:99 30:70 50:50 70:30 99:1

 3. 发送4w笔交易，15分钟后查询结果，主要查询：
    a. 交易丢失
    b. 交易执行结果为超时

4. 每秒发送两倍85笔交易，判断在当前压力下是否会出现交易池溢出
*/

type Tx struct {
	TxId           string
	TimeStamp      int64               // 交易发送时间戳 (int64 类型，表示 Unix 时间戳)
	InPool         bool                // 判断是否在交易池中
	OnChain        bool                // 判断是否在链中
	ExecuteResult  common.TxStatusCode // 判断交易执行情况
	ExecuteMessage string
	SendStatus     common.TxStatusCode // 判断交易是否成功发送
	SendMessage    string
}

type Txs struct {
	Mu  sync.Mutex
	Txs []*Tx
}

func (t *Txs) Append(tx *Tx) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	t.Txs = append(t.Txs, tx)
}

func (t *Txs) GetAllTxs() []*Tx {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return t.Txs
}

func (t *Txs) SortByTimeStamps() {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	sort.Slice(t.Txs, func(i, j int) bool {
		return t.Txs[i].TimeStamp < t.Txs[j].TimeStamp
	})
}

var (
	ratios = []struct {
		A int
		B int
	}{
		{1, 99},
		{30, 70},
		{50, 50},
		{70, 30},
		{99, 1},
	}
)

func generateAndTrackTransactions(aRatio, bRatio int, timeSleep int64, totalRound int, f *FuncPairSeed) []*Tx {
	txs := &Txs{}
	taskCh := make(chan func()) // 定义任务通道
	var wg sync.WaitGroup       // 主任务同步

	// 启动线程池，固定最大 Goroutine 数为 20
	const maxWorkers = 20
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for task := range taskCh { // 从通道中获取任务
				task() // 执行任务
			}
		}()
	}

	FuncOneInput := f.SeedOne.convertMapToKeyValuePair(f.SeedOne.FunctionInput)
	FuncOneName := f.SeedOne.FunctionName
	FuncTwoInput := f.SeedTwo.convertMapToKeyValuePair(f.SeedTwo.FunctionInput)
	FuncTwoName := f.SeedTwo.FunctionName

	for count := 0; count < totalRound; count++ {
		// 分发 aRatio 次 FuncOne 的任务
		wg.Add(aRatio + bRatio)

		go func() {
			for i := 0; i < aRatio; i++ {
				// wg.Add(1)
				taskCh <- func() {
					defer wg.Done()
					txid, message, status, _, _ := nodecontrol.ChainmakerController.UserContractInvoke(utils.GlobalContractInfo.ContractName, FuncOneName, FuncOneInput, false)
					txs.Append(&Tx{
						TxId:          txid,
						OnChain:       false,
						InPool:        false,
						TimeStamp:     time.Now().Unix(),
						ExecuteResult: -1,
						SendStatus:    status,
						SendMessage:   message,
					})
				}
			}
		}()

		// 分发 bRatio 次 FuncTwo 的任务
		go func() {
			for i := 0; i < bRatio; i++ {
				// wg.Add(1)
				taskCh <- func() {
					defer wg.Done()
					txid, message, status, _, _ := nodecontrol.ChainmakerController.UserContractInvoke(utils.GlobalContractInfo.ContractName, FuncTwoName, FuncTwoInput, false)
					txs.Append(&Tx{
						TxId:          txid,
						OnChain:       false,
						InPool:        false,
						TimeStamp:     time.Now().Unix(),
						ExecuteResult: -1,
						SendStatus:    status,
						SendMessage:   message,
					})
				}
			}
		}()

		wg.Wait() // 等待所有任务完成

		time.Sleep(time.Duration(timeSleep) * time.Millisecond)

	}

	close(taskCh) // 关闭任务通道

	// 根据时间戳信息排序
	txs.SortByTimeStamps()

	return txs.GetAllTxs()
}

func WaitForEmptyPool() {
	Log := utils.Log

	Log.Log(utils.ConflictLog, "Waiting for the transaction pool to be empty...")
	for {
		status, err := nodecontrol.ChainmakerController.Client.GetPoolStatus()
		if err != nil {
			fmt.Println("查询交易池失败：", err)
			return
		}
		if status.CommonTxNumInPending+status.CommonTxNumInQueue == 0 {
			Log.Log(utils.ConflictLog, "Transaction pool is empty.")
			return
		}
		Log.Log(utils.ConflictLog, fmt.Sprintf("Transaction pool is not empty, txs number: [%d]. Retrying in 30 seconds...", status.CommonTxNumInPending+status.CommonTxNumInQueue))
		time.Sleep(30 * time.Second)
	}
}

// 查询所有交易，获取交易时间戳信息
func GetTxsTimeStampInfo(txs []*Tx) {
	taskCh := make(chan *Tx) // 用于传递任务的通道
	var wg sync.WaitGroup

	// 启动线程池
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for tx := range taskCh { // 从通道中获取任务
				txInfo, err := nodecontrol.ChainmakerController.Client.GetTxByTxId(tx.TxId)
				if err != nil {
					fmt.Println("查询交易出现错误!", tx.TxId, err)
				}

				if txInfo != nil {
					tx.OnChain = true
					tx.ExecuteResult = txInfo.Transaction.Result.Code
					tx.ExecuteMessage = txInfo.Transaction.Result.ContractResult.Message
					tx.TimeStamp = txInfo.Transaction.Payload.Timestamp
				} else {
					tx.OnChain = false
					txInpool, _, _ := nodecontrol.ChainmakerController.Client.GetTxsInPoolByTxIds([]string{tx.TxId})
					if txInpool != nil {
						tx.InPool = true
						tx.TimeStamp = txInpool[0].Payload.Timestamp
					} else {
						tx.InPool = false
						// 跳过该时间的更新
					}
				}

				wg.Done() // 当前任务完成
			}
		}()
	}

	// 分发任务
	for _, tx := range txs {
		wg.Add(1)
		taskCh <- tx
	}

	close(taskCh) // 关闭通道，通知线程池中的 Goroutine 不再有新任务
	wg.Wait()     // 等待所有任务完成
}

func RecordTxsResultWithPool(txs []*Tx) {
	taskCh := make(chan *Tx) // 用于传递任务的通道
	var wg sync.WaitGroup

	// 启动线程池
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for tx := range taskCh { // 从通道中获取任务
				// 该交易已经上链，不需要再查询
				if tx.OnChain == true {
					continue
				}

				txInfo, err := nodecontrol.ChainmakerController.Client.GetTxByTxId(tx.TxId)
				if err != nil {
					fmt.Println("查询交易出现错误!", tx.TxId, err)
				}

				if txInfo != nil {
					tx.OnChain = true
					tx.ExecuteResult = txInfo.Transaction.Result.Code
				} else {
					tx.OnChain = false
					txInpool, _, _ := nodecontrol.ChainmakerController.Client.GetTxsInPoolByTxIds([]string{tx.TxId})
					if txInpool != nil {
						tx.InPool = true
					} else {
						tx.InPool = false
					}
				}

				wg.Done() // 当前任务完成
			}
		}()
	}

	// 分发任务
	for _, tx := range txs {
		wg.Add(1)
		taskCh <- tx
	}

	close(taskCh) // 关闭通道，通知线程池中的 Goroutine 不再有新任务
	wg.Wait()     // 等待所有任务完成
}

func ConflictPairSeedExperiment(f *FuncPairSeed) {
	Log := utils.Log

	// 生成时间戳
	timestamp := time.Now().Unix()

	// 创建目标目录
	targetDir := filepath.Join(Log.BaseDir, fmt.Sprintf("%s_%s_%d", f.SeedOne.FunctionName, f.SeedTwo.FunctionName, timestamp))
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		fmt.Printf("failed to create directory: %v\n", err)
		return
	}

	Log.Log(utils.ConflictLog, fmt.Sprintf("结果保存目录: [%s]", targetDir))

	for _, ratio := range ratios {
		txs := SendTxBatchAndQueryLater(ratio.A, ratio.B, f)
		// 准备保存目标文件
		err := SaveTxsToFile(txs, filepath.Join(targetDir, fmt.Sprintf("%d_%d.json", ratio.A, ratio.B)))
		if err != nil {
			fmt.Println("保存结果至文件失败!", err)
		}

		err = DrawTransactionLossStatus(txs, filepath.Join(targetDir, fmt.Sprintf("%d_%d_loss.png", ratio.A, ratio.B)))
		if err != nil {
			fmt.Println("生成交易上链情况图失败!", err)
		}

		err = DrawTransactionExecuteStatus(txs, filepath.Join(targetDir, fmt.Sprintf("%d_%d_execute.png", ratio.A, ratio.B)))
		if err != nil {
			fmt.Println("生成交易执行情况图失败!", err)
		}
	}

	txs := LongTermDDoSAttack(f)

	// 准备保存目标文件
	err = SaveTxsToFile(txs, filepath.Join(targetDir, "LongTermExperiment.json"))
	if err != nil {
		fmt.Println("保存结果至文件失败!", err)
	}

	err = PlotLostTransactions(txs, filepath.Join(targetDir, "LongTermExperiment.png"))
	if err != nil {
		fmt.Println("保存结果至文件失败!", err)
	}
}

// 短时间发送1w笔交易，直到交易池完全空后查询结果
func SendTxBatchAndQueryLater(ratioA, ratioB int, f *FuncPairSeed) []*Tx {
	Log := utils.Log

	Log.Log(utils.ConflictLog, fmt.Sprintf("Running experiment with ratio A:B = %d:%d", ratioA, ratioB))

	// 获取经过时间戳排序后的所有交易的TxId
	txs := generateAndTrackTransactions(ratioA, ratioB, 0, 100, f)
	Log.Log(utils.ConflictLog, fmt.Sprintf("Generated %d transactions", len(txs)))

	WaitForEmptyPool()

	RecordTxsResultWithPool(txs)

	return txs

}

// 长时间每秒发送，连续发送十分钟，主要观察交易发送状态，不需要获取交易的各种结果
func LongTermDDoSAttack(f *FuncPairSeed) []*Tx {
	Log := utils.Log
	Log.Log(utils.ConflictLog, "开始长时间交易发送测试")

	txs := generateAndTrackTransactions(85, 85, 1000, 600, f)

	Log.Log(utils.ConflictLog, fmt.Sprintf("长时间发送交易完成, Generated %d transactions", len(txs)))

	// 发送后等待交易池空
	WaitForEmptyPool()
	// RecordTxsResultWithPool(txs)
	return txs
}

// SaveTxsToFile 将 []*Tx 类型数据保存为可绘图的 JSON 文件
func SaveTxsToFile(txs []*Tx, filePath string) error {
	// 打开文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 将数据编码为 JSON 并写入文件
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 格式化输出
	if err := encoder.Encode(txs); err != nil {
		return fmt.Errorf("写入 JSON 数据失败: %w", err)
	}

	return nil
}

// 画交易丢失情况图
func DrawTransactionLossStatus(txs []*Tx, filePath string) error {
	txsLossStatus := make([]bool, 0)

	for i := 0; i < 1000; i++ {
		txsLossStatus = append(txsLossStatus, txs[i].InPool || txs[i].OnChain)
	}

	return picture.DrawTransactionGrid(txsLossStatus, filePath)
}

// 画交易执行情况图
func DrawTransactionExecuteStatus(txs []*Tx, filePath string) error {
	txsLossStatus := make([]bool, 0)

	txs = txs[0:1000]

	for _, tx := range txs {
		txsLossStatus = append(txsLossStatus, tx.ExecuteResult == 0)
	}

	return picture.DrawTransactionGrid(txsLossStatus, filePath)
}

// 画交易发送情况图
func PlotLostTransactions(txs []*Tx, savePath string) error {
	// 加载支持中文的字体
	fontBytes, err := ioutil.ReadFile("picture/simhei.ttf") // 请确保 simhei.ttf 字体文件在程序运行的目录下
	if err != nil {
		return err
	}

	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return err
	}

	const fontName = "SimHei"
	font.DefaultCache.Add([]font.Face{
		{
			Font: font.Font{
				Typeface: fontName,
			},
			Face: ttfFont,
		},
	})

	// 设置默认字体
	plot.DefaultFont = font.Font{Typeface: fontName}

	// 计算每秒钟丢失的交易数
	lostTxsPerSecond := make(map[int64]int)
	var minTime, maxTime int64
	for _, tx := range txs {
		timestamp := tx.TimeStamp
		if minTime == 0 || timestamp < minTime {
			minTime = timestamp
		}
		if maxTime == 0 || timestamp > maxTime {
			maxTime = timestamp
		}
		if tx.SendStatus != common.TxStatusCode_SUCCESS {
			lostTxsPerSecond[timestamp]++
		}
	}

	// 准备绘图数据
	var pts plotter.XYs
	for t := minTime; t <= maxTime; t++ {
		numLost := lostTxsPerSecond[t]
		pts = append(pts, struct{ X, Y float64 }{
			X: float64(t - minTime),
			Y: float64(numLost),
		})
	}

	// 创建绘图
	p := plot.New()

	p.Title.Text = "每秒钟丢失的交易数"
	p.Title.Padding = vg.Length(10)
	p.X.Label.Text = "时间（秒）"
	p.Y.Label.Text = "丢失的交易数"

	// 设置 Y 轴范围，至少展示到 170
	p.Y.Min = 0
	p.Y.Max = 200

	// 创建折线图
	line, err := plotter.NewLine(pts)
	if err != nil {
		return err
	}
	line.LineStyle.Width = vg.Points(1)
	p.Add(line)

	// 添加水平虚线，表示交易发送速率 (170 笔/秒)
	txPerSec := 170
	hLine := plotter.NewFunction(func(x float64) float64 { return float64(txPerSec) })
	hLine.LineStyle.Color = plotter.DefaultLineStyle.Color
	hLine.LineStyle.Width = vg.Points(1)
	hLine.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)} // 设置为虚线
	hLine.XMin = 0
	hLine.XMax = float64(maxTime - minTime)
	p.Add(hLine)

	// 添加图例
	p.Legend.Add("丢失的交易数", line)
	p.Legend.Add("交易发送速率 (170 笔/秒)", hLine)
	p.Legend.Top = true

	// 调整图例位置和边距
	p.Legend.Padding = vg.Points(5)
	p.Legend.YOffs = -vg.Length(15)

	// 增加顶部和底部边距
	p.Title.Padding = vg.Points(20)
	p.X.Padding = vg.Points(20)
	p.Y.Padding = vg.Points(20)

	// 保存为 PNG 文件
	if err := p.Save(10*vg.Inch, 6*vg.Inch, savePath); err != nil {
		return err
	}

	return nil
}
