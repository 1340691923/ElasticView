package build

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"path/filepath"
	"runtime"
)

var Version = "1.0.0"

// Config holds the setup variables required for a build
type Config struct {
	OS               string // GOOS
	Arch             string // GOOS
	Env              map[string]string
	OutputBinaryPath string
}

func getExecutableName(os string, arch string) (string, error) {
	exname := "elastic_view"

	exeName := fmt.Sprintf("%s_%s_%s", exname, os, arch)
	if os == "windows" {
		exeName = fmt.Sprintf("%s.exe", exeName)
	}
	return exeName, nil
}

func buildBackend(cfg Config) error {

	exeName, err := getExecutableName(cfg.OS, cfg.Arch)
	if err != nil {
		return err
	}

	ldFlags := `-extldflags "-static"`

	// Add linker flags to drop debug information
	prefix := ""
	if ldFlags != "" {
		prefix = " "
	}
	ldFlags = fmt.Sprintf("-w -s%s%s", prefix, ldFlags)

	outputPath := cfg.OutputBinaryPath
	if outputPath == "" {
		outputPath = fmt.Sprintf("resources/dist/%s", Version)
	}
	args := []string{
		"build", "-o", filepath.Join(outputPath, exeName),
	}

	args = append(args, "-ldflags", ldFlags)

	rootPackage := "./cmd/ev"

	args = append(args, rootPackage)

	cfg.Env["GOARCH"] = cfg.Arch
	cfg.Env["GOOS"] = cfg.OS
	cfg.Env["CGO_ENABLED"] = "0"

	// TODO: Change to sh.RunWithV once available.
	return sh.RunWith(cfg.Env, "go", args...)
}

func newBuildConfig(os string, arch string) Config {
	return Config{
		OS:   os,
		Arch: arch,
		Env:  map[string]string{},
	}
}

// Build is a namespace.
type Build mg.Namespace

// Linux builds the back-end plugin for Linux.
func (Build) Linux() error {
	return buildBackend(newBuildConfig("linux", "amd64"))
}

// LinuxARM builds the back-end plugin for Linux on ARM.
func (Build) LinuxARM() error {
	return buildBackend(newBuildConfig("linux", "arm"))
}

// LinuxARM64 builds the back-end plugin for Linux on ARM64.
func (Build) LinuxARM64() error {
	return buildBackend(newBuildConfig("linux", "arm64"))
}

// Windows builds the back-end plugin for Windows.
func (Build) Windows() error {
	return buildBackend(newBuildConfig("windows", "amd64"))
}

// Darwin builds the back-end plugin for OSX.
func (Build) Darwin() error {
	return buildBackend(newBuildConfig("darwin", "amd64"))
}

// DarwinARM64 builds the back-end plugin for OSX on ARM (M1).
func (Build) DarwinARM64() error {
	return buildBackend(newBuildConfig("darwin", "arm64"))
}

// Backend build a production build for the current platform
func (Build) Backend() error {
	cfg := newBuildConfig(runtime.GOOS, runtime.GOARCH)
	return buildBackend(cfg)
}

// BuildAll builds production executables for all supported platforms.
func BuildAll() { //revive:disable-line
	b := Build{}
	mg.Deps(b.Linux, b.Windows, b.Darwin, b.DarwinARM64, b.LinuxARM64, b.LinuxARM)
	//util.DirCopy("./resources/config",fmt.Sprintf("./resources/dist/%s/config",Version))
}
