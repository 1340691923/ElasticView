package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var args *config.CommandLineArgs

func init() {
	args = &config.CommandLineArgs{}
	flag.StringVar(&args.HomePath, "homePath", util.GetCurrentDirectory(), "ev程序所在文件夹")
	flag.StringVar(&args.CmdName, "cmdName", "build", "二进制名称")
	flag.StringVar(&args.ConfigFile, "configFile", "config/config.yml", "配置文件路径")
	flag.Parse()
}

func main() {
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

	//-H windowsgui
	ldFlags := fmt.Sprintf("-w -s%s%s ", " ", `-extldflags "-static"`)

	/*if cfg.OS == "windows" {
		ldFlags = fmt.Sprintf("-w -s%s%s -H windowsgui", " ", `-extldflags "-static"`)
	}*/

	outputPath := cfg.OutputPath

	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	args = append(args, "-ldflags", ldFlags)

	rootPackage := "./cmd/ev"

	args = append(args, rootPackage)

	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["CGO_ENABLED"] = "0"
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
	cmd.Dir = "resources/vue"

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
	os.RemoveAll(outputPath)
	fmt.Println("清理临时文件完毕")

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
	return
}
