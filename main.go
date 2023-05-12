package main

import (
	"SysInf/core/config"
	ui "SysInf/ui"
	"SysInf/ui/widgets"
	tui "github.com/gizak/termui/v3"
	"log"
)

func main() {
	config.LoadConfig()
	if err := tui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer tui.Close()

	widgets.InitWidgets()
	widgets.BuildWidgets()
	ui.Run()
}
