package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type general struct {
	CpuRefreshDelay     int64  `json:"cpu_refresh_delay"`
	RamRefreshDelay     int64  `json:"ram_refresh_delay"`
	DiskRefreshDelay    int64  `json:"disk_refresh_delay"`
	ProcessRefreshDelay int64  `json:"process_refresh_delay"`
	UIRefreshDelay      int64  `json:"ui_refresh_delay"`
	DiskPath            string `json:"disk_path"`
}

type conf struct {
	General general `json:"general"`
}

var LoadedConfig conf

func LoadConfig() {
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)
	byteValue, _ := io.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &LoadedConfig); err != nil {
		return
	}
}
