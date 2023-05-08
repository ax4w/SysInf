package process

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"sort"
)

type ProcInf struct {
	Name  string
	Usage float64
	pID   int32
}

type asProcessList []ProcInf

func (b asProcessList) Len() int {
	return len(b)
}
func (b asProcessList) Swap(c, d int) {
	b[c], b[d] = b[d], b[c]
}

func (b asProcessList) Less(c, d int) bool {
	return b[c].Usage > b[d].Usage
}

func ToGB(val uint64) uint64 {
	return val / 1024 / 1024 / 1024
}

func processByID(id int32) *process.Process {
	PROCESSInfoStat, err := process.Processes()
	if err != nil {
		log.Fatalf("Could not retrieve host info")
	}
	for _, p := range PROCESSInfoStat {
		if p.Pid == id {
			return p
		}
	}
	return nil
}

func ResumeProcess(id int32) {
	err := processByID(id).Resume()
	if err != nil {
		return
	}
}

func SuspendProcess(id int32) {
	err := processByID(id).Suspend()
	if err != nil {
		return
	}
}

func KillProcess(id int32) {
	err := processByID(id).Terminate()
	if err != nil {
		return
	}
}

func SortedProcesses() []ProcInf {
	PROCESSInfoStat, err := process.Processes()
	if err != nil {
		log.Fatalf("Could not retrieve host info")
	}
	var procInfos []ProcInf
	for _, p := range PROCESSInfoStat {
		pName, _ := p.Name()
		pUsage, _ := p.CPUPercent()
		pPID := p.Pid
		procInfos = append(procInfos, ProcInf{pName, pUsage, pPID})
	}
	sort.Sort(asProcessList(procInfos))
	return procInfos
}

func Info() []string {
	var processes []string
	procInfos := SortedProcesses()
	for _, p := range procInfos {
		processes = append(processes, fmt.Sprintf("%-6d| %-55s - [Usage: %.2f %s] ", p.pID, p.Name, p.Usage, "%"))
	}
	return processes

}
