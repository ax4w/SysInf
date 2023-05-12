package ui

import (
	"SysInf/core/config"
	"SysInf/core/process"
	"SysInf/ui/widgets"
	tui "github.com/gizak/termui/v3"
	"strconv"
	"strings"
	"time"
)

func Run() {
	uiEvents := tui.PollEvents()
	ticker := time.NewTicker(time.Duration(config.LoadedConfig.General.UIRefreshDelay) * time.Millisecond).C
	StartThreads()
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
				alive = false
				tui.Clear()
				return
			case "<Resize>":
				resize(e.Payload.(tui.Resize))

			}
		case <-ticker:
			//needs to be polled frequently
			tui.Clear()
			tui.Render(widgets.RamPiChart, widgets.DiskPiChart,
				widgets.CpuCoresGraph, widgets.ProcessList, widgets.ControlsBox)
		}
	}
}
