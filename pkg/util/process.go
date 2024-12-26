package util

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v3/process"
)

func GetProcCPUPercent(pid int32) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Fatalf("Failed to create process object: %v", err)
	}

	// 获取 CPU 使用率
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		log.Printf("Failed to get CPU usage: %v", err)
	} else {
		fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent)
	}

}

func GetProcMemoryInfo(pid int32) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Fatalf("Failed to create process object: %v", err)
	}
	// 获取内存信息
	memInfo, err := proc.MemoryInfo()
	if err != nil {
		log.Printf("Failed to get memory usage: %v", err)
	} else {
		fmt.Printf("Memory Usage: %v mb\n", memInfo.RSS/1024/1024)
	}
}
