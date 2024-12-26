package util

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os/exec"
	"strconv"
	"strings"
)

// KillProcessByPort 根据端口找到服务并杀死对应进程
func KillProcessByPort(port int) (string, error) {
	c1 := exec.Command("powershell.exe", "netstat", "-ano", "| findstr", strconv.Itoa(port), "| findstr", "LISTENING")
	output, err := c1.CombinedOutput()
	if err != nil {
		return "", err
	}
	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	strNetstat := string(decodeBytes)

	Pid := "0"

	for _, v := range strings.Split(strNetstat, "\r\n") {
		for {
			if strings.Contains(v, "  ") {
				v = strings.ReplaceAll(v, "  ", " ")
			} else {
				break
			}
		}
		strings.Trim(v, " ")

		arrTmp := strings.Split(v, " ")
		arr := make([]string, 0)
		for _, vv := range arrTmp {
			if vv != "" {
				arr = append(arr, vv)
			}
		}
		if len(arr) == 5 && strings.Contains(arr[1], fmt.Sprintf(":%v", port)) {
			Pid = arr[4]
			break
		}
	}

	if Pid == "0" {
		return "", errors.New(fmt.Sprintf("No service with port number %v found", port))
	}

	c2 := exec.Command("powershell.exe", "taskkill -PID ", Pid, "-F")
	output, err = c2.CombinedOutput()
	if err != nil {
		return "", err
	}
	decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	delStr := string(decodeBytes)
	return delStr, nil
}
