package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
)

var pluginJsonFile string
var execUpx bool

func init() {
	flag.StringVar(&pluginJsonFile, "pj", "plugin.json", "插件配置文件")
	flag.BoolVar(&execUpx, "upx", false, "是否使用upx")
	flag.Parse()
}

func main() {

	BuildVue()

	err := build.BuildPluginSvr(pluginJsonFile, execUpx)

	if err != nil {
		fmt.Println("BuildPluginSvr err", err)
	} else {
		fmt.Println("BuildPluginSvr success")
	}

}

func BuildVue() (err error) {
	cmd := exec.Command("npm", "run", "build")

	sourceDir, _ := os.Getwd()

	cmd.Dir = filepath.Join(sourceDir, "frontend")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithStack(err)
	}
	if err = cmd.Start(); err != nil {
		return errors.WithStack(err)
	}
	scanner := bufio.NewScanner(stdout)
	fmt.Println("=================build vue================")

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err = cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return
}
