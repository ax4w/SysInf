package ui

import (
	"SysInf/cpu"
	"SysInf/process"
	"SysInf/widgets"
	"fmt"
	tui "github.com/gizak/termui/v3"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"strconv"
	"strings"
	"time"
)

const refreshDelay = 250

var startTimeStamp = time.Now().UTC().UnixMilli()

var diskPath string

func resize(payload tui.Resize) {
	widgets.RamPiChart.SetRect(0, 0, payload.Width/2, payload.Height/3)
	widgets.DiskPiChart.SetRect(payload.Width/2, 0, payload.Width, payload.Height/3)
	widgets.ProcessList.SetRect(0, payload.Height/2, payload.Width, payload.Height)
	widgets.ControlsBox.SetRect(0, payload.Height-3, payload.Width, payload.Height)
	widgets.CpuCoresGraph.SetRect(0, payload.Height/3, payload.Width, payload.Height/2)
	widgets.CpuCoresGraph.BarWidth = (payload.Width - 8) / int(cpu.Count())

}

func update() {
	virtualMemInfo, err := mem.VirtualMemory()
	diskInfo, err := disk.Usage(diskPath)
	if err != nil {
		log.Fatalf("Could not retrieve host info")
	}
	//Calc used RAM in % and set RAN PiChart values
	RamUsedInPercent := 10 + ((100 / float64(virtualMemInfo.Total)) * float64(virtualMemInfo.Used))
	DiskUsedInGB := process.ToGB(diskInfo.Used)
	//Update Processes
	processes := process.Info()
	//CPU stats don't need to be updated every 250 milliseconds
	if time.Now().UTC().UnixMilli()-startTimeStamp > 1000 {
		startTimeStamp = time.Now().UTC().UnixMilli()
		widgets.CpuCoresGraph.Labels = cpu.Labels()
		widgets.CpuCoresGraph.Title = fmt.Sprintf("Total CPU usage by user %.2f %s", cpu.Usage()[0], "%")
		widgets.CpuCoresGraph.Data = cpu.CoresUsage()
	}

	widgets.ProcessList.Title = fmt.Sprintf("Processes - %d running", len(process.SortedProcesses()))
	//Update Values
	widgets.RamPiChart.Data = []float64{RamUsedInPercent, 100 - RamUsedInPercent}
	widgets.DiskPiChart.Data = []float64{float64(DiskUsedInGB), float64((process.ToGB(diskInfo.Total)) - DiskUsedInGB)}
	widgets.ProcessList.Rows = processes
}

func Run() {
	//overwrite shading blocks for the pi charts
	tui.SHADED_BLOCKS = [...]rune{'▒', '█', ' ', ' ', ' '}
	uiEvents := tui.PollEvents()
	ticker := time.NewTicker(refreshDelay * time.Millisecond).C

	diskPath = "/"
	if widgets.IsWindows() {
		diskPath = "\\"
	}
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "r":
				if !widgets.IsWindows() {
					t := strings.TrimSpace(strings.
						Split(widgets.ProcessList.Rows[widgets.ProcessList.SelectedRow], "|")[0])
					parsed, _ := strconv.ParseInt(t, 10, 32)
					process.ResumeProcess(int32(parsed))
				}
			case "p":
				if !widgets.IsWindows() {
					t := strings.TrimSpace(strings.
						Split(widgets.ProcessList.Rows[widgets.ProcessList.SelectedRow], "|")[0])
					parsed, _ := strconv.ParseInt(t, 10, 32)
					process.SuspendProcess(int32(parsed))
				}
			case "k":
				t := strings.TrimSpace(strings.
					Split(widgets.ProcessList.Rows[widgets.ProcessList.SelectedRow], "|")[0])
				parsed, _ := strconv.ParseInt(t, 10, 32)
				process.KillProcess(int32(parsed))

			case "w":
				if widgets.ProcessList.SelectedRow > 0 {
					widgets.ProcessList.SelectedRow--
				}
			case "s":
				if widgets.ProcessList.SelectedRow < len(widgets.ProcessList.Rows) {
					widgets.ProcessList.SelectedRow++
				}
			case "q", "<C-c>":
				tui.Clear()
				return
			case "<Resize>":
				resize(e.Payload.(tui.Resize))

			}
		case <-ticker:
			//needs to be polled frequently
			update()
			tui.Clear()
			tui.Render(widgets.RamPiChart, widgets.DiskPiChart,
				widgets.CpuCoresGraph, widgets.ProcessList, widgets.ControlsBox)
		}
	}
}
