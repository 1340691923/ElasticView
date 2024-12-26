package main

import (
	"flag"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"os"
	"runtime"
)

func main() {
	// 解析命令行参数

	destZip := flag.String("o", "dist/source_code.zip", "输出的zip文件路径")
	excludeDir := flag.String("e", "node_modules,dist,.idea,.vscode,.git", "要排除的文件夹路径")
	flag.Parse()

	sourceDir, _ := os.Getwd()
	execFileName := "plugin"
	if runtime.GOOS == "windows" {
		execFileName = execFileName + ".exe"
	}

	os.MkdirAll("dist", os.ModePerm)

	err := build.CompressPathToZip(sourceDir, *excludeDir, execFileName, *destZip)

	if err != nil {
		fmt.Println("压缩过程中出错:", err)
		os.Exit(1)
	}

	fmt.Println("压缩完成，输出文件:", *destZip)
}
