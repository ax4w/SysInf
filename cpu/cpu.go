package cpu

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
)

func ChartFormat(f float64) string {
	return fmt.Sprintf("%.f%s", f, "%")
}

func Count() int32 {
	CPUInfoStat, err := cpu.Info()
	if err != nil {
		log.Fatalf("Error while getting cpu Info")
	}
	return CPUInfoStat[0].Cores
}

func CoresUsage() []float64 {
	CPUInfoStat, err := cpu.Percent(0, true)
	if err != nil {
		log.Fatalf("Error while getting cpu Info")
	}
	return CPUInfoStat
}

func Usage() []float64 {
	CPUInfoStat, err := cpu.Percent(0, false)
	if err != nil {
		log.Fatalf("Error while getting cpu Info")
	}
	return CPUInfoStat
}

func Labels() []string {
	var result []string
	for i := 1; i <= int(Count()); i++ {
		result = append(result, fmt.Sprintf("Core %d", i))
	}
	return result
}
