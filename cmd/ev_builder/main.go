package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var args *config.CommandLineArgs
var execUpx bool

func init() {
	args = &config.CommandLineArgs{}
	flag.StringVar(&args.ConfigFile, "c", "config/config.yml", "配置文件路径")
	flag.BoolVar(&execUpx, "upx", false, "是否使用upx")
	flag.Parse()
}

func main() {
	currentDirectory, _ := os.Getwd()
	args.HomePath = currentDirectory
	args.CmdName = "ev_builder"

	cfg, err := config.InitConfig(args)
	if err != nil {
		fmt.Println("InitConfig err", err)
		panic(err)
	}

	err = BuildVue(cfg)
	if err != nil {
		fmt.Println("BuildVue err", err)
		panic(err)
	}

	outputZipPath := fmt.Sprintf("resources/dist/ev_%s", strings.ReplaceAll(cfg.Version, ".", "_"))

	err = util.WriteVersionGoFile(fmt.Sprintf(`package config

const Version = "v%s"

`, cfg.Version))

	if err != nil {
		fmt.Println("BuildVue err", err)
		panic(err)
	}

	fmt.Println("开始检测是否已有该版本打包", cfg.Version)
	if util.CheckFileIsExist(outputZipPath) {
		fmt.Println("检测到已经该版本打包，正在删除老包", cfg.Version)
		os.RemoveAll(outputZipPath)
	} else {
		fmt.Println("暂无该版本打包")
	}

	for _, osAndArch := range []string{
		"linux_arm64",
		"linux_amd64",
		"darwin_amd64",
		"darwin_arm64",
		"windows_amd64",
	} {
		osAndArchArr := strings.Split(osAndArch, "_")
		os := osAndArchArr[0]
		arch := osAndArchArr[1]
		err = BuildEvSvr(cfg, os, arch)
		if err != nil {
			fmt.Println("BuildEvSvr err", err)
			panic(err)
		}
	}

	fmt.Println("BuildEvSvr success")
}

type BuildConfig struct {
	OS         string // GOOS
	GOARCH     string
	Env        map[string]string
	OutputPath string
}

func getExecutableName(os, arch string) (string, error) {
	exname := "ev"

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

func buildBackend(cfg BuildConfig) error {

	exeName, err := getExecutableName(cfg.OS, cfg.GOARCH)
	if err != nil {
		return err
	}

	ldFlags := fmt.Sprintf("-w -s%s%s ", " ", `-extldflags "-static"`)

	if cfg.OS == "windows" {
		ldFlags = "-H windowsgui -w -s"
	}

	outputPath := cfg.OutputPath

	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	args = append(args, "-trimpath", "-ldflags", ldFlags)

	rootPackage := "./cmd/ev"

	args = append(args, rootPackage)

	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["GOARCH"] = cfg.GOARCH
	return RunGoBuild(cfg.Env, args...)
}

func newBuildConfig(os, arch, outputPath string) BuildConfig {
	return BuildConfig{
		OS:         os,
		GOARCH:     arch,
		OutputPath: outputPath,
		Env:        map[string]string{},
	}
}

func BuildVue(cfg *config.Config) (err error) {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = filepath.Join(filepath.Join(cfg.HomePath, "resources"), "vue")

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

func BuildEvSvr(cfg *config.Config, buildOs, buildArch string) (err error) {

	fmt.Println("开始编译ev二进制文件，版本：", cfg.Version, buildOs, buildArch)

	outputPath := fmt.Sprintf("resources/dist/ev_%s_%s_%s", cfg.Version, buildOs, buildArch)

	err = buildBackend(newBuildConfig(buildOs, buildArch, outputPath))
	if err != nil {
		return
	}

	os.MkdirAll(filepath.Join(outputPath, "config"), os.ModePerm)
	err = util.DirCopy("config", filepath.Join(outputPath, "config"))
	if err != nil {
		return
	}

	versionString := strings.ReplaceAll(cfg.Version, ".", "_")

	outputZipPath := fmt.Sprintf("resources/dist/ev_%s", versionString)

	fmt.Println("开始打包", cfg.Version, buildOs, buildArch)

	os.MkdirAll(outputZipPath, os.ModePerm)

	outputZip := filepath.Join(outputZipPath, fmt.Sprintf("ev_%s_%s_%s.zip", versionString, buildOs, buildArch))

	err = util.CompressPathToZip(outputPath, outputZip)
	if err != nil {
		return
	}
	fmt.Println("打包成功，开始清理临时文件", cfg.Version, buildOs, buildArch)
	err = os.RemoveAll(outputPath)
	if err != nil {
		fmt.Println("清理临时文件失敗", outputPath, err)
	} else {
		fmt.Println("清理临时文件完毕", outputPath)
	}

	return
}

func RunGoBuild(env map[string]string, args ...string) (err error) {
	if len(env) > 0 {
		envArr := []string{"env", "-w"}
		for k, v := range env {
			envArr = append(envArr, fmt.Sprintf("%s=%s", k, v))
		}

		fmt.Println(fmt.Sprintf("start build cmd: go %v", envArr))

		cmd := exec.Command("go", envArr...)

		err = cmd.Run()
		if err != nil {
			return
		}
	}

	fmt.Println(fmt.Sprintf("start build cmd: go %v", args))
	cmd := exec.Command("go", args...)
	err = cmd.Run()

	if err != nil {
		return
	}
	if execUpx {
		fmt.Println(fmt.Sprintf("start cmd %s %s %s", "upx", "--best", args[2]))
		cmd = exec.Command("upx", "--best", args[2])
		err = cmd.Run()
	}
	return
}
