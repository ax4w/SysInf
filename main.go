package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"math"
	"runtime"
	"time"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "\\"
	}

	//Host
	HOSTInfoStat, _ := host.Info()
	HOSTInfo := widgets.NewParagraph()
	HOSTInfo.SetRect(70, 20, 36, 27)
	HOSTInfo.Title = "Host Info"
	HOSTInfo.Text = fmt.Sprintf("Hostname: %s\nOS: %s\nKernelArch: %s\nKernelVersion: %s\nProcesses: %d\n",
		HOSTInfoStat.Hostname,
		HOSTInfoStat.OS,
		HOSTInfoStat.KernelArch,
		HOSTInfoStat.KernelVersion,
		HOSTInfoStat.Procs)
	//CPU
	CPUInfoStat, _ := cpu.Info()
	CPUINfo := widgets.NewParagraph()
	CPUINfo.SetRect(0, 20, 36, 27)
	CPUINfo.Title = "CPU Stats"
	CPUINfo.Text = fmt.Sprintf("Model: %s\nCore: %d\nMHz: %.2f\n",
		CPUInfoStat[0].ModelName,
		CPUInfoStat[0].Cores,
		CPUInfoStat[0].Mhz)

	//RAM
	RAMData := []float64{1, 2}
	RAMPiChart := widgets.NewPieChart()
	RAMPiChart.Title = "RAM Usage %"
	RAMPiChart.SetRect(0, 0, 35, 20)
	RAMPiChart.AngleOffset = -.5 * math.Pi
	RAMPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f", v)
	}

	//Disk
	DiskData := []float64{1, 2}
	DISKPiChart := widgets.NewPieChart()
	DISKPiChart.Title = "Disk Usage"
	DISKPiChart.SetRect(70, 0, 36, 20)
	DISKPiChart.AngleOffset = -.5 * math.Pi
	DISKPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f GB", v)
	}

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		virtualMemInfo, _ := mem.VirtualMemory()
		diskInfo, _ := disk.Usage(diskPath)
		RamUsedInPercent := 10 + ((100 / float64(virtualMemInfo.Total)) * float64(virtualMemInfo.Used))
		DiskUsedInGB := diskInfo.Used / 1024 / 1024 / 1024
		RAMData[0] = RamUsedInPercent
		RAMData[1] = 100 - RamUsedInPercent
		DiskData[0] = float64(DiskUsedInGB)
		DiskData[1] = float64((diskInfo.Total / 1024 / 1024 / 1024) - DiskUsedInGB)
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return

			}
		case <-ticker:
			RAMPiChart.Data = RAMData
			DISKPiChart.Data = DiskData
			ui.Render(RAMPiChart, DISKPiChart, CPUINfo, HOSTInfo)
		}
	}
}
