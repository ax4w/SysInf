package ui

import (
	"SysInf/core/config"
	"SysInf/core/cpu"
	"SysInf/core/process"
	"SysInf/ui/widgets"
	"fmt"
	tui "github.com/gizak/termui/v3"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

var (
	alive = true
)

func resize(payload tui.Resize) {
	widgets.RamPiChart.SetRect(0, 0, payload.Width/2, payload.Height/3)
	widgets.DiskPiChart.SetRect(payload.Width/2, 0, payload.Width, payload.Height/3)
	widgets.ProcessList.SetRect(0, payload.Height/2, payload.Width, payload.Height)
	widgets.ControlsBox.SetRect(0, payload.Height-3, payload.Width, payload.Height)
	widgets.CpuCoresGraph.SetRect(0, payload.Height/3, payload.Width, payload.Height/2)
	widgets.CpuCoresGraph.BarWidth = payload.Width / int(cpu.Count()+(cpu.Count()/3))

}

func cpuRoutine() {
	for alive {
		widgets.CpuCoresGraph.Labels = cpu.Labels()
		widgets.CpuCoresGraph.Title = fmt.Sprintf("Total CPU usage by user %.2f %s", cpu.Usage()[0], "%")
		widgets.CpuCoresGraph.Data = cpu.CoresUsage()
		time.Sleep(time.Duration(config.LoadedConfig.General.CpuRefreshDelay) * time.Millisecond)
	}

}

func ramRoutine() {
	virtualMemInfo, err := mem.VirtualMemory()
	if err != nil {
		alive = false
		panic("error retrieving memory info")
	}
	for alive {
		RamUsedInPercent := 10 + ((100 / float64(virtualMemInfo.Total)) * float64(virtualMemInfo.Used))
		widgets.RamPiChart.Data = []float64{RamUsedInPercent, 100 - RamUsedInPercent}
		time.Sleep(time.Duration(config.LoadedConfig.General.RamRefreshDelay) * time.Millisecond)
	}
}

func diskRoutine() {
	diskInfo, err := disk.Usage(config.LoadedConfig.General.DiskPath)
	if err != nil {
		alive = false
		panic("error retrieving disk info")
	}
	for alive {
		DiskUsedInGB := process.ToGB(diskInfo.Used)
		widgets.DiskPiChart.Data = []float64{float64(DiskUsedInGB), float64((process.ToGB(diskInfo.Total)) - DiskUsedInGB)}
		time.Sleep(time.Duration(config.LoadedConfig.General.DiskRefreshDelay) * time.Millisecond)
	}
}

func processRoutine() {
	for alive {
		processes := process.Info()
		widgets.ProcessList.Title = fmt.Sprintf("Processes - %d running", len(process.SortedProcesses()))
		widgets.ProcessList.Rows = processes
		time.Sleep(time.Duration(config.LoadedConfig.General.ProcessRefreshDelay) * time.Millisecond)
	}
}

func StartThreads() {
	go ramRoutine()
	go diskRoutine()
	go cpuRoutine()
	go processRoutine()

}
