package picture

import (
	"image/color"
	"math"
	"strconv"

	"github.com/fogleman/gg"
)

const (
	gridSize   = 20                  // 每个矩形（小方格）的边长
	numCols    = 100                 // 每行展示 100 个交易
	scaleStepX = 10                  // 上侧刻度以 10 为单位
	scaleStepY = 100                 // 左侧刻度以 100 为单位
	fontSize   = 20                  // 字体大小
	fontPath   = "picture/Arial.ttf" // 字体路径
	margin     = 70                  // 图片四周的边缘空间
)

// DrawTransactionGrid 生成交易网格图片
func DrawTransactionGrid(transactionStatus []bool, outputPath string) error {
	totalTx := len(transactionStatus) // 获取交易总数

	// 计算图片宽高
	numRows := int(math.Ceil(float64(totalTx) / float64(numCols))) // 行数
	imageWidth := numCols*gridSize + margin*2                      // 图片宽度（网格宽度 + 两侧边距）
	imageHeight := numRows*gridSize + margin*2                     // 图片高度（网格高度 + 两侧边距）

	// 创建画布
	dc := gg.NewContext(imageWidth, imageHeight)
	dc.SetColor(color.White) // 设置背景颜色为白色
	dc.Clear()

	// 加载字体
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return err
	}

	// 绘制刻度
	drawScales(dc, numCols, numRows)

	// 绘制每笔交易的矩形
	for i := 0; i < totalTx; i++ {
		x := (i%numCols)*gridSize + margin // x 坐标偏移，留出边距
		y := (i/numCols)*gridSize + margin // y 坐标偏移，留出边距

		// 绘制矩形填充颜色
		if transactionStatus[i] {
			dc.SetColor(color.White) // 成功（白色）
		} else {
			dc.SetColor(color.RGBA{255, 0, 0, 255}) // 失败（红色）
		}
		dc.DrawRectangle(float64(x), float64(y), gridSize, gridSize)
		dc.Fill()

		// 绘制边框
		dc.SetColor(color.Black)
		dc.SetLineWidth(1)
		dc.DrawRectangle(float64(x), float64(y), gridSize, gridSize)
		dc.Stroke()
	}

	// 保存图片
	return dc.SavePNG(outputPath)
}

// 绘制左侧和上侧刻度
func drawScales(dc *gg.Context, numCols, numRows int) {
	dc.SetColor(color.Black)

	// 绘制上侧刻度
	for j := 0; j <= numCols; j += scaleStepX {
		x := float64(j*gridSize + margin) // 刻度位于矩形左边缘
		label := strconv.Itoa(j)
		// 绘制刻度值
		dc.DrawStringAnchored(label, x, float64(margin-30), 0.5, 1)
		// 绘制横线
		dc.DrawLine(x, float64(margin-10), x, float64(margin))
		dc.Stroke()
	}

	// 绘制左侧刻度
	for i := 0; i <= numRows; i += scaleStepY / numCols {
		y := float64(i*gridSize + margin) // 刻度位于矩形顶部边缘
		label := strconv.Itoa(i * scaleStepY)
		// 绘制刻度值
		dc.DrawStringAnchored(label, float64(margin-15), y, 1, 0.5)
		// 绘制竖线
		dc.DrawLine(float64(margin-10), y, float64(margin), y)
		dc.Stroke()
	}
}
