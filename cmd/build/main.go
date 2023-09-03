package main

import (
	"flag"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"log"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var args *config.CommandLineArgs

func init() {
	args = &config.CommandLineArgs{}
	flag.StringVar(&args.HomePath, "homePath", util.GetCurrentDirectory(), "ev程序所在文件夹")
	flag.StringVar(&args.CmdName, "cmdName", "build", "二进制名称")
	flag.StringVar(&args.ConfigFile, "configFile", "D:\\eve\\new_ev\\ElasticView\\config\\config.yml", "配置文件路径")
	flag.Parse()
}

func main() {
	cfg, err := config.InitConfig(args)
	if err != nil {
		log.Println("InitConfig err", err)
		panic(err)
	}
	err = BuildEvSvr(cfg)
	if err != nil {
		log.Println("BuildEvSvr err", err)
		panic(err)
	}

	log.Println("BuildEvSvr success")
}

func buildBackground() {

}

type Build struct {
	EvVersion string
}

type BuildConfig struct {
	OS         string // GOOS
	Arch       string // GOOS
	Env        map[string]string
	EvVersion  string
	OutputPath string
}

func getExecutableName(os string, arch string) (string, error) {
	exname := "elastic_view"

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

func buildBackend(cfg BuildConfig) error {

	exeName, err := getExecutableName(cfg.OS, cfg.Arch)
	if err != nil {
		return err
	}

	ldFlags := fmt.Sprintf("-w -s%s%s -H windowsgui", " ", `-extldflags "-static"`)

	outputPath := cfg.OutputPath

	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	args = append(args, "-ldflags", ldFlags)

	rootPackage := "./cmd/ev"

	args = append(args, rootPackage)

	cfg.Env["GOARCH"] = cfg.Arch
	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["CGO_ENABLED"] = "0"

	return RunGoBuild(cfg.Env, args...)
}

func newBuildConfig(os, arch, evVersion string) BuildConfig {
	return BuildConfig{
		OS:         os,
		Arch:       arch,
		EvVersion:  evVersion,
		OutputPath: fmt.Sprintf("resources/dist/%s", evVersion),
		Env:        map[string]string{},
	}
}

func (this *Build) Linux() error {
	return buildBackend(newBuildConfig("linux", "amd64", this.EvVersion))
}

func (this *Build) LinuxARM() error {
	return buildBackend(newBuildConfig("linux", "arm", this.EvVersion))
}

func (this *Build) LinuxARM64() error {
	return buildBackend(newBuildConfig("linux", "arm64", this.EvVersion))
}

func (this *Build) Windows() error {
	return buildBackend(newBuildConfig("windows", "amd64", this.EvVersion))
}

func (this *Build) Darwin() error {
	return buildBackend(newBuildConfig("darwin", "amd64", this.EvVersion))
}

func (this *Build) DarwinARM64() error {
	return buildBackend(newBuildConfig("darwin", "arm64", this.EvVersion))
}

func (this *Build) Backend() error {
	cfg := newBuildConfig(runtime.GOOS, runtime.GOARCH, this.EvVersion)
	return buildBackend(cfg)
}

func BuildEvSvr(cfg *config.Config) (err error) {

	log.Println("开始编译ev二进制文件")

	b := Build{
		EvVersion: cfg.Version,
	}

	runErrFn := func(fnArr ...func() error) (err error) {
		for _, fn := range fnArr {
			if err = fn(); err != nil {
				log.Println("err", err)
				return
			}
		}
		return
	}

	err = runErrFn(
		b.Linux,
		b.Darwin,
		b.DarwinARM64,
		b.LinuxARM64,
		b.LinuxARM,
		b.Windows,
	)
	if err != nil {
		return
	}

	outputPath := fmt.Sprintf("resources/dist/%s", cfg.Version)
	os.Mkdir(filepath.Join(outputPath, "config"), 0755)
	err = util.DirCopy("config", filepath.Join(outputPath, "config"))
	if err != nil {
		return
	}

	outputZipPath := fmt.Sprintf("resources/dist/%s", strings.ReplaceAll(cfg.Version, ".", "_"))
	log.Println("开始检测是否已有该版本打包")
	if util.CheckFileIsExist(outputZipPath) {
		log.Println("检测到已经该版本打包，正在删除老包")
		os.RemoveAll(outputZipPath)
	} else {
		log.Println("暂无该版本打包")
	}
	log.Println("开始打包")

	os.Mkdir(outputZipPath, 0755)

	outputZip := filepath.Join(outputZipPath, "ev.zip")

	err = util.CompressPathToZip(outputPath, outputZip)
	if err != nil {
		return
	}
	log.Println("打包成功，开始清理临时文件")
	os.RemoveAll(outputPath)
	log.Println("清理临时文件完毕")

	return
}

func RunGoBuild(env map[string]string, args ...string) (err error) {
	if len(env) > 0 {
		envArr := []string{"env", "-w"}
		for k, v := range env {
			envArr = append(envArr, fmt.Sprintf("%s=%s", k, v))
		}

		log.Println(fmt.Sprintf("start build cmd: go %v", envArr))

		cmd := exec.Command("go", envArr...)

		err = cmd.Run()
		if err != nil {
			return
		}
	}

	log.Println(fmt.Sprintf("start build cmd: go %v", args))
	cmd := exec.Command("go", args...)
	err = cmd.Run()
	return
}

/*func BuildVue() (err error) {
	cmd := exec.Command("cd resources/vue", "&&", "npm run build:prod")
	err = cmd.Run()
	return
}*/
