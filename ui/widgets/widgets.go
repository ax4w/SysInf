package widgets

import (
	"SysInf/core/config"
	"SysInf/core/cpu"
	"fmt"
	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"runtime"
)

var (
	ControlsBox   *widgets.Paragraph
	ProcessList   *widgets.List
	RamPiChart    *widgets.PieChart
	DiskPiChart   *widgets.PieChart
	CpuCoresGraph *widgets.BarChart
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func InitWidgets() {
	ControlsBox = widgets.NewParagraph()
	ProcessList = widgets.NewList()
	RamPiChart = widgets.NewPieChart()
	DiskPiChart = widgets.NewPieChart()
	CpuCoresGraph = widgets.NewBarChart()
}

func BuildWidgets() {
	w, h := tui.TerminalDimensions()
	//Controls Box
	ControlsBox.SetRect(0, h-3, w, h)
	ControlsBox.Title = "Controls"
	text := "General: q - quit, Processes: w - up, s - down, k - kill"
	if !IsWindows() {
		text += ", p - pause, r resume"
	}
	ControlsBox.Text = text
	ControlsBox.TextStyle = tui.Style{Fg: tui.ColorGreen, Bg: tui.ColorClear}
	//CPU Graph
	CpuCoresGraph.SetRect(0, h/3, w, (h/2)+5)
	CpuCoresGraph.BarWidth = w / int(cpu.Count()+(cpu.Count()/3))
	CpuCoresGraph.BarGap = 2
	CpuCoresGraph.NumFormatter = cpu.ChartFormat
	CpuCoresGraph.MaxVal = 100.0
	CpuCoresGraph.BarColors = []tui.Color{tui.ColorGreen, tui.ColorGreen}
	CpuCoresGraph.NumStyles = []tui.Style{{
		Fg: tui.ColorBlack,
	}, {
		Fg: tui.ColorBlack,
	}}
	CpuCoresGraph.LabelStyles = []tui.Style{{
		Fg: tui.ColorWhite,
		Bg: tui.ColorClear,
	}, {
		Fg: tui.ColorWhite,
		Bg: tui.ColorClear,
	}}
	//Processes
	ProcessList.Title = "Processes"
	ProcessList.WrapText = false
	ProcessList.SelectedRowStyle = tui.Style{Fg: tui.ColorGreen, Bg: tui.ColorClear}
	ProcessList.SetRect(0, (h/2)+5, w, h-3)

	//RAM
	RamPiChart.SetRect(0, 0, w/2, h/3)
	RamPiChart.Title = "RAM usage"
	RamPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f%s", v, "%")
	}

	//Disk
	DiskPiChart.Title = "Disk Space Used for disk: " + config.LoadedConfig.General.DiskPath
	DiskPiChart.SetRect(w/2, 0, w, h/3)
	DiskPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f GB", v)
	}
}
