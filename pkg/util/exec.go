package util

import (
	"log"
	"os/exec"
	"runtime"
	"syscall"
)

func OpenWinBrowser(uri string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command(`cmd`, `/c`, `start`, uri)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		err := cmd.Start()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
