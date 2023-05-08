package widgets

import (
	"SysInf/cpu"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	ControlsBox   *widgets.Paragraph
	ProcessList   *widgets.List
	RamPiChart    *widgets.PieChart
	DiskPiChart   *widgets.PieChart
	CpuCoresGraph *widgets.BarChart
)

func InitWidgets() {
	ControlsBox = widgets.NewParagraph()
	ProcessList = widgets.NewList()
	RamPiChart = widgets.NewPieChart()
	DiskPiChart = widgets.NewPieChart()
	CpuCoresGraph = widgets.NewBarChart()
}

func BuildWidgets() {
	w, h := ui.TerminalDimensions()
	//Quit Box
	ControlsBox.SetRect(0, h-3, w, h)
	ControlsBox.Title = "Controls"
	ControlsBox.Text = "General: q - quit, Processes: w - up, s - down, k - kill"
	ControlsBox.TextStyle = ui.Style{Fg: ui.ColorRed, Bg: ui.ColorClear}
	//CPU Graph
	CpuCoresGraph.SetRect(0, h/3, w, (h/2)+5)
	CpuCoresGraph.BarWidth = (w - (w)/int(cpu.Count())) / int(cpu.Count())
	CpuCoresGraph.NumFormatter = cpu.ChartFormat
	CpuCoresGraph.MaxVal = 100.0
	CpuCoresGraph.BarColors = []ui.Color{ui.ColorGreen, ui.ColorGreen}
	CpuCoresGraph.NumStyles = []ui.Style{{
		Fg: ui.ColorBlack,
	}, {
		Fg: ui.ColorBlack,
	}}
	CpuCoresGraph.LabelStyles = []ui.Style{{
		Fg: ui.ColorWhite,
		Bg: ui.ColorClear,
	}, {
		Fg: ui.ColorWhite,
		Bg: ui.ColorClear,
	}}
	//Processes
	ProcessList.Title = "Processes"
	ProcessList.WrapText = false
	ProcessList.SelectedRowStyle = ui.Style{Fg: ui.ColorGreen, Bg: ui.ColorClear}
	ProcessList.SetRect(0, (h/2)+5, w, h-3)

	//RAM
	RamPiChart.SetRect(0, 0, w/2, h/3)
	RamPiChart.Title = "RAM usage"
	RamPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f %s", v, "%")
	}

	//Disk
	DiskPiChart.Title = "Disk Space Used"
	DiskPiChart.SetRect(w/2, 0, w, h/3)
	DiskPiChart.LabelFormatter = func(i int, v float64) string {
		return fmt.Sprintf("%.02f GB", v)
	}
}
