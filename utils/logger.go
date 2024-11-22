package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogType string

const (
	FuncSeedsLog LogType = "funcseeds_log"
	ConflictLog  LogType = "conflict_test_logs"
	ExecutionLog LogType = "execution_info"
	FuzzLog      LogType = "fuzz_test_logs"
)

// Logger 结构体
type Logger struct {
	BaseDir string
}

// NewLogger 创建一个新的 Logger 实例
func NewLogger(contractName string) (*Logger, error) {
	timestamp := time.Now().Format("20060102_150405")
	baseDir := filepath.Join("./result", fmt.Sprintf("%s_%s", contractName, timestamp))

	// 创建目录
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	return &Logger{BaseDir: baseDir}, nil
}

// Log 将内容写入对应的日志文件
func (l *Logger) Log(logType LogType, message string) error {
	filePath := filepath.Join(l.BaseDir, string(logType)+".txt")

	// 打开文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 写入带时间戳的日志
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = file.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
	if err != nil {
		return fmt.Errorf("failed to write log: %v", err)
	}

	return nil
}
