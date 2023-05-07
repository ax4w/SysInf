package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ProcInf struct {
	Name  string
	Usage float64
	pID   int32
}

type ByUsage []ProcInf

func (b ByUsage) Len() int {
	return len(b)
}
func (b ByUsage) Swap(c, d int) {
	b[c], b[d] = b[d], b[c]
}

func (b ByUsage) Less(c, d int) bool {
	return b[c].Usage > b[d].Usage
}

func toGB(val uint64) uint64 {
	return val / 1024 / 1024 / 1024
}

func killProcessByID(id int32) {
	PROCESSInfoStat, err := process.Processes()
	if err != nil {
		log.Fatalf("Could not retrieve host info")
	}
	for _, p := range PROCESSInfoStat {
		if p.Pid == id {
			err := p.Kill()
			if err != nil {
				return
			}
		}
	}
}

func getSortedProcesses(coresCount int32) ([]ProcInf, float64) {
	PROCESSInfoStat, err := process.Processes()
	if err != nil {
		log.Fatalf("Could not retrieve host info")
	}
	var totalUsage float64
	var procInfos []ProcInf
	for _, p := range PROCESSInfoStat {
		pName, _ := p.Name()
		pUsage, _ := p.CPUPercent()
		pPID := p.Pid
		procInfos = append(procInfos, ProcInf{pName, pUsage, pPID})
		totalUsage += pUsage
	}
	totalUsage /= float64(coresCount)
	sort.Sort(ByUsage(procInfos))
	return procInfos, totalUsage
}

func getProcInfos(coresCount int32) ([]string, float64) {

	var topProcs []string
	procInfos, totalUsage := getSortedProcesses(coresCount)
	for _, p := range procInfos {
		topProcs = append(topProcs, fmt.Sprintf("%-6d| %-55s - [Usage: %.2f %s] ", p.pID, p.Name, p.Usage, "%"))
	}
	return topProcs, totalUsage

}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	w, h := ui.TerminalDimensions()

	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "\\"
	}
	//Widget Inits
	CPUInfo := widgets.NewParagraph()
	QUITBox := widgets.NewParagraph()
	PROCESSList := widgets.NewList()
	HOSTInfo := widgets.NewParagraph()
	RAMPiChart := widgets.NewPieChart()
	DISKPiChart := widgets.NewPieChart()
	//Stats
	HOSTInfoStat, err := host.Info()
	CPUInfoStat, err := cpu.Info()
	if err != nil {
		log.Fatalf("Could creating widgets")
	}

	//Quit Box
	QUITBox.SetRect(0, h-3, w, h)
	QUITBox.Title = "Controls"
	QUITBox.Text = "General: q - quit, Processes: w - up, s - down, k - kill"
	QUITBox.TextStyle = ui.Style{Fg: ui.ColorRed, Bg: ui.ColorClear}

	//Processes
	PROCESSList.WrapText = false
	PROCESSList.SelectedRowStyle = ui.Style{Fg: ui.ColorGreen, Bg: ui.ColorClear}
	PROCESSList.SetRect(0, h/2, w, h-3)

	//Host
	HOSTInfo.SetRect(w/2, h/3, w, h/2)
	HOSTInfo.Title = "Host Info"
	HOSTInfo.Text = fmt.Sprintf("Hostname: %s\nOS: %s\nKernelArch: %s\nKernelVersion: %s\nProcesses: %d\n",
		HOSTInfoStat.Hostname,
		HOSTInfoStat.OS,
		HOSTInfoStat.KernelArch,
		HOSTInfoStat.KernelVersion,
		HOSTInfoStat.Procs)
	//CPU
	CPUInfo.SetRect(0, h/3, w/2, h/2)
	CPUInfo.Title = "CPU Stats"
	CPUInfo.Text = fmt.Sprintf("Model: %s\nCore: %d\nMHz: %.2f\n",
		CPUInfoStat[0].ModelName,
		CPUInfoStat[0].Cores,
		CPUInfoStat[0].Mhz)

	//RAM
	RAMPiChart.SetRect(0, 0, w/2, h/3)
	RAMPiChart.Title = "RAM usage"
	RAMPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f %s", v, "%")
	}

	//Disk
	DISKPiChart.Title = "Disk Space Used"
	DISKPiChart.SetRect(w/2, 0, w, h/3)
	DISKPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f GB", v)
	}

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(250 * time.Millisecond).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "k":

				t := strings.TrimSpace(strings.Split(PROCESSList.Rows[PROCESSList.SelectedRow], "|")[0])
				parsed, _ := strconv.ParseInt(t, 10, 32)
				killProcessByID(int32(parsed))

			case "w":
				if PROCESSList.SelectedRow > 0 {
					PROCESSList.SelectedRow--
				}
			case "s":
				if PROCESSList.SelectedRow < len(PROCESSList.Rows) {
					PROCESSList.SelectedRow++
				}
			case "q", "<C-c>":
				ui.Clear()
				return
			case "<Resize>":
				//Resize widgets
				payload := e.Payload.(ui.Resize)
				RAMPiChart.SetRect(0, 0, payload.Width/2, payload.Height/3)
				DISKPiChart.SetRect(payload.Width/2, 0, payload.Width, payload.Height/3)
				CPUInfo.SetRect(0, payload.Height/3, payload.Width/2, payload.Height/2)
				HOSTInfo.SetRect(payload.Width/2, payload.Height/3, payload.Width, payload.Height/2)
				PROCESSList.SetRect(0, payload.Height/2, payload.Width, payload.Height)
				QUITBox.SetRect(0, payload.Height-3, payload.Width, payload.Height)
			}
		case <-ticker:
			//needs to be polled frequently
			virtualMemInfo, err := mem.VirtualMemory()
			diskInfo, err := disk.Usage(diskPath)
			if err != nil {
				log.Fatalf("Could not retrieve host info")
			}
			//Calc used RAM in % and set RAN PiChart values
			RamUsedInPercent := 10 + ((100 / float64(virtualMemInfo.Total)) * float64(virtualMemInfo.Used))
			DiskUsedInGB := toGB(diskInfo.Used)
			//Update Processes
			processes, usage := getProcInfos(CPUInfoStat[0].Cores)
			PROCESSList.Title = fmt.Sprintf("Total CPU usage %.2f %s", usage, "%")
			//Update Values
			RAMPiChart.Data = []float64{RamUsedInPercent, 100 - RamUsedInPercent}
			DISKPiChart.Data = []float64{float64(DiskUsedInGB), float64((toGB(diskInfo.Total)) - DiskUsedInGB)}
			PROCESSList.Rows = processes
			ui.Clear()
			ui.Render(RAMPiChart, DISKPiChart, CPUInfo, HOSTInfo, PROCESSList, QUITBox)
		}
	}
}
