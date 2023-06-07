//go:build ignore
// +build ignore

package main

import (
	"context"
	"flag"
	"fmt"
	build_mageimport "github.com/1340691923/ElasticView/pkg/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	// Use local types and functions in order to avoid name conflicts with additional magefiles.
	type arguments struct {
		Verbose bool          // print out log statements
		List    bool          // print out a list of targets
		Help    bool          // print out help for a specific target
		Timeout time.Duration // set a timeout to running the targets
		Args    []string      // args contain the non-flag command-line arguments
	}

	parseBool := func(env string) bool {
		val := os.Getenv(env)
		if val == "" {
			return false
		}
		b, err := strconv.ParseBool(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid bool value: %v", env, val)
			return false
		}
		return b
	}

	parseDuration := func(env string) time.Duration {
		val := os.Getenv(env)
		if val == "" {
			return 0
		}
		d, err := time.ParseDuration(val)
		if err != nil {
			log.Printf("warning: environment variable %s is not a valid duration value: %v", env, val)
			return 0
		}
		return d
	}
	args := arguments{}
	fs := flag.FlagSet{}
	fs.SetOutput(os.Stdout)

	// default flag set with ExitOnError and auto generated PrintDefaults should be sufficient
	fs.BoolVar(&args.Verbose, "v", parseBool("MAGEFILE_VERBOSE"), "show verbose output when running targets")
	fs.BoolVar(&args.List, "l", parseBool("MAGEFILE_LIST"), "list targets for this binary")
	fs.BoolVar(&args.Help, "h", parseBool("MAGEFILE_HELP"), "print out help for a specific target")
	fs.DurationVar(&args.Timeout, "t", parseDuration("MAGEFILE_TIMEOUT"), "timeout in duration parsable format (e.g. 5m30s)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stdout, `
%s [options] [target]

Commands:
  -l    list targets in this binary
  -h    show this help

Options:
  -h    show description of a target
  -t <string>
        timeout in duration parsable format (e.g. 5m30s)
  -v    show verbose output when running targets
 `[1:], filepath.Base(os.Args[0]))
	}
	if err := fs.Parse(os.Args[1:]); err != nil {
		// flag will have printed out an error already.
		return
	}
	args.Args = fs.Args()
	if args.Help && len(args.Args) == 0 {
		fs.Usage()
		return
	}

	// color is ANSI color type
	type color int

	// If you add/change/remove any items in this constant,
	// you will need to run "stringer -type=color" in this directory again.
	// NOTE: Please keep the list in an alphabetical order.
	const (
		black color = iota
		red
		green
		yellow
		blue
		magenta
		cyan
		white
		brightblack
		brightred
		brightgreen
		brightyellow
		brightblue
		brightmagenta
		brightcyan
		brightwhite
	)

	// AnsiColor are ANSI color codes for supported terminal colors.
	var ansiColor = map[color]string{
		black:         "\u001b[30m",
		red:           "\u001b[31m",
		green:         "\u001b[32m",
		yellow:        "\u001b[33m",
		blue:          "\u001b[34m",
		magenta:       "\u001b[35m",
		cyan:          "\u001b[36m",
		white:         "\u001b[37m",
		brightblack:   "\u001b[30;1m",
		brightred:     "\u001b[31;1m",
		brightgreen:   "\u001b[32;1m",
		brightyellow:  "\u001b[33;1m",
		brightblue:    "\u001b[34;1m",
		brightmagenta: "\u001b[35;1m",
		brightcyan:    "\u001b[36;1m",
		brightwhite:   "\u001b[37;1m",
	}

	const _color_name = "blackredgreenyellowbluemagentacyanwhitebrightblackbrightredbrightgreenbrightyellowbrightbluebrightmagentabrightcyanbrightwhite"

	var _color_index = [...]uint8{0, 5, 8, 13, 19, 23, 30, 34, 39, 50, 59, 70, 82, 92, 105, 115, 126}

	colorToLowerString := func(i color) string {
		if i < 0 || i >= color(len(_color_index)-1) {
			return "color(" + strconv.FormatInt(int64(i), 10) + ")"
		}
		return _color_name[_color_index[i]:_color_index[i+1]]
	}

	// ansiColorReset is an ANSI color code to reset the terminal color.
	const ansiColorReset = "\033[0m"

	// defaultTargetAnsiColor is a default ANSI color for colorizing targets.
	// It is set to Cyan as an arbitrary color, because it has a neutral meaning
	var defaultTargetAnsiColor = ansiColor[cyan]

	getAnsiColor := func(color string) (string, bool) {
		colorLower := strings.ToLower(color)
		for k, v := range ansiColor {
			colorConstLower := colorToLowerString(k)
			if colorConstLower == colorLower {
				return v, true
			}
		}
		return "", false
	}

	// Terminals which  don't support color:
	// 	TERM=vt100
	// 	TERM=cygwin
	// 	TERM=xterm-mono
	var noColorTerms = map[string]bool{
		"vt100":      false,
		"cygwin":     false,
		"xterm-mono": false,
	}

	// terminalSupportsColor checks if the current console supports color output
	//
	// Supported:
	// 	linux, mac, or windows's ConEmu, Cmder, putty, git-bash.exe, pwsh.exe
	// Not supported:
	// 	windows cmd.exe, powerShell.exe
	terminalSupportsColor := func() bool {
		envTerm := os.Getenv("TERM")
		if _, ok := noColorTerms[envTerm]; ok {
			return false
		}
		return true
	}

	// enableColor reports whether the user has requested to enable a color output.
	enableColor := func() bool {
		b, _ := strconv.ParseBool(os.Getenv("MAGEFILE_ENABLE_COLOR"))
		return b
	}

	// targetColor returns the ANSI color which should be used to colorize targets.
	targetColor := func() string {
		s, exists := os.LookupEnv("MAGEFILE_TARGET_COLOR")
		if exists == true {
			if c, ok := getAnsiColor(s); ok == true {
				return c
			}
		}
		return defaultTargetAnsiColor
	}

	// store the color terminal variables, so that the detection isn't repeated for each target
	var enableColorValue = enableColor() && terminalSupportsColor()
	var targetColorValue = targetColor()

	printName := func(str string) string {
		if enableColorValue {
			return fmt.Sprintf("%s%s%s", targetColorValue, str, ansiColorReset)
		} else {
			return str
		}
	}

	list := func() error {

		targets := map[string]string{
			"build:backend":     "build a production build for the current platform",
			"build:darwin":      "builds the back-end plugin for OSX.",
			"build:darwinARM64": "builds the back-end plugin for OSX on ARM (M1).",
			"build:linux":       "builds the back-end plugin for Linux.",
			"build:linuxARM":    "builds the back-end plugin for Linux on ARM.",
			"build:linuxARM64":  "builds the back-end plugin for Linux on ARM64.",
			"build:windows":     "builds the back-end plugin for Windows.",
			"buildAll*":         "builds production executables for all supported platforms.",
		}

		keys := make([]string, 0, len(targets))
		for name := range targets {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		fmt.Println("Targets:")
		w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
		for _, name := range keys {
			fmt.Fprintf(w, "  %v\t%v\n", printName(name), targets[name])
		}
		err := w.Flush()
		if err == nil {
			fmt.Println("\n* default target")
		}
		return err
	}

	var ctx context.Context
	var ctxCancel func()

	getContext := func() (context.Context, func()) {
		if ctx != nil {
			return ctx, ctxCancel
		}

		if args.Timeout != 0 {
			ctx, ctxCancel = context.WithTimeout(context.Background(), args.Timeout)
		} else {
			ctx = context.Background()
			ctxCancel = func() {}
		}
		return ctx, ctxCancel
	}

	runTarget := func(fn func(context.Context) error) interface{} {
		var err interface{}
		ctx, cancel := getContext()
		d := make(chan interface{})
		go func() {
			defer func() {
				err := recover()
				d <- err
			}()
			err := fn(ctx)
			d <- err
		}()
		select {
		case <-ctx.Done():
			cancel()
			e := ctx.Err()
			fmt.Printf("ctx err: %v\n", e)
			return e
		case err = <-d:
			cancel()
			return err
		}
	}
	// This is necessary in case there aren't any targets, to avoid an unused
	// variable error.
	_ = runTarget

	handleError := func(logger *log.Logger, err interface{}) {
		if err != nil {
			logger.Printf("Error: %+v\n", err)
			type code interface {
				ExitStatus() int
			}
			if c, ok := err.(code); ok {
				os.Exit(c.ExitStatus())
			}
			os.Exit(1)
		}
	}
	_ = handleError

	// Set MAGEFILE_VERBOSE so mg.Verbose() reflects the flag value.
	if args.Verbose {
		os.Setenv("MAGEFILE_VERBOSE", "1")
	} else {
		os.Setenv("MAGEFILE_VERBOSE", "0")
	}

	log.SetFlags(0)
	if !args.Verbose {
		log.SetOutput(ioutil.Discard)
	}
	logger := log.New(os.Stderr, "", 0)
	if args.List {
		if err := list(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}

	if args.Help {
		if len(args.Args) < 1 {
			logger.Println("no target specified")
			os.Exit(2)
		}
		switch strings.ToLower(args.Args[0]) {
		case "build:backend":
			fmt.Println("Backend build a production build for the current platform")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:backend\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:darwin":
			fmt.Println("Darwin builds the back-end plugin for OSX.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:darwin\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:darwinarm64":
			fmt.Println("DarwinARM64 builds the back-end plugin for OSX on ARM (M1).")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:darwinarm64\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:linux":
			fmt.Println("Linux builds the back-end plugin for Linux.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:linux\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:linuxarm":
			fmt.Println("LinuxARM builds the back-end plugin for Linux on ARM.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:linuxarm\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:linuxarm64":
			fmt.Println("LinuxARM64 builds the back-end plugin for Linux on ARM64.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:linuxarm64\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "build:windows":
			fmt.Println("Windows builds the back-end plugin for Windows.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage build:windows\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		case "buildall":
			fmt.Println("BuildAll builds production executables for all supported platforms.")
			fmt.Println()

			fmt.Print("Usage:\n\n\tmage buildall\n\n")
			var aliases []string
			if len(aliases) > 0 {
				fmt.Printf("Aliases: %s\n\n", strings.Join(aliases, ", "))
			}
			return
		default:
			logger.Printf("Unknown target: %q\n", args.Args[0])
			os.Exit(2)
		}
	}
	if len(args.Args) < 1 {
		ignoreDefault, _ := strconv.ParseBool(os.Getenv("MAGEFILE_IGNOREDEFAULT"))
		if ignoreDefault {
			if err := list(); err != nil {
				logger.Println("Error:", err)
				os.Exit(1)
			}
			return
		}

		wrapFn := func(ctx context.Context) error {
			build_mageimport.BuildAll()
			return nil
		}
		ret := runTarget(wrapFn)
		handleError(logger, ret)
		return
	}
	for x := 0; x < len(args.Args); {
		target := args.Args[x]
		x++

		// resolve aliases
		switch strings.ToLower(target) {

		}

		switch strings.ToLower(target) {

		case "build:backend":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:Backend\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:Backend")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.Backend()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:darwin":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:Darwin\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:Darwin")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.Darwin()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:darwinarm64":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:DarwinARM64\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:DarwinARM64")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.DarwinARM64()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:linux":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:Linux\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:Linux")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.Linux()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:linuxarm":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:LinuxARM\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:LinuxARM")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.LinuxARM()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:linuxarm64":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:LinuxARM64\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:LinuxARM64")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.LinuxARM64()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "build:windows":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"Build:Windows\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "Build:Windows")
			}

			wrapFn := func(ctx context.Context) error {
				return build_mageimport.Build{}.Windows()
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		case "buildall":
			expected := x + 0
			if expected > len(args.Args) {
				// note that expected and args at this point include the arg for the target itself
				// so we subtract 1 here to show the number of args without the target.
				logger.Printf("not enough arguments for target \"BuildAll\", expected %v, got %v\n", expected-1, len(args.Args)-1)
				os.Exit(2)
			}
			if args.Verbose {
				logger.Println("Running target:", "BuildAll")
			}

			wrapFn := func(ctx context.Context) error {
				build_mageimport.BuildAll()
				return nil
			}
			ret := runTarget(wrapFn)
			handleError(logger, ret)
		default:
			logger.Printf("Unknown target specified: %q\n", target)
			os.Exit(2)
		}
	}
}
