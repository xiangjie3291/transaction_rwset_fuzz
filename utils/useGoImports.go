package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// 检查 goimports 是否已经安装
func isGoimportsInstalled() bool {
	_, err := exec.LookPath("goimports")
	return err == nil
}

// 安装 goimports
func installGoimports() {
	cmd := exec.Command("go", "install", "golang.org/x/tools/cmd/goimports@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to install goimports: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("goimports installed successfully.")
}

// 运行 goimports 格式化文件
func runGoimports(filePath string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Printf("Failed to get absolute path of %s: %s\n", filePath, err)
		os.Exit(1)
	}

	cmd := exec.Command("goimports", "-w", absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to run goimports: %s\n", err)
		os.Exit(1)
	}

	//debug:
	//fmt.Printf("File %s formatted successfully by goimports.\n", filePath)
}

func UseGoimports(filePath string) {
	if !isGoimportsInstalled() {
		fmt.Println("goimports not found. Installing...")
		installGoimports()
	} else {
		//debug:
		//fmt.Println("goimports is already installed.")
	}

	runGoimports(filePath)
}
