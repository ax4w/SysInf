package main

import (
	ui2 "SysInf/ui"
	"SysInf/widgets"
	ui "github.com/gizak/termui/v3"
	"log"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	widgets.InitWidgets()
	widgets.BuildWidgets()
	ui2.Run()
}
