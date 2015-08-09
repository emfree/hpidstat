package main

import (
	"fmt"
	procinfo "github.com/c9s/goprocinfo/linux"
	"time"
)

func getProcTime(pid uint64) int64 {
	p, err := procinfo.ReadProcess(pid, "/proc")
	if err != nil {
		return 0
	}
	return int64(p.Stat.Utime) + p.Stat.Cutime
}

func getSysTime() int64 {
	sys, err := procinfo.ReadStat("/proc/stat")
	if err != nil {
		panic(err)
	}
	stat := sys.CPUStatAll
	return int64(stat.User + stat.Nice + stat.System + stat.Idle + stat.IOWait)
}

func getProcFractions(pids []uint64, freq time.Duration) []float64 {
	fractions := make([]float64, len(pids))
	startSysTime := getSysTime()
	startProcTimes := make([]int64, len(pids))
	for i, pid := range pids {
		startProcTimes[i] = getProcTime(pid)
	}
	time.Sleep(freq)
	endSysTime := getSysTime()
	diff := float64(endSysTime - startSysTime)
	for i, pid := range pids {
		fractions[i] = float64(getProcTime(pid)-startProcTimes[i]) / diff
	}
	return fractions
}

func getNames(pids []uint64) []string {
	var names []string
	for _, pid := range pids {
		path := fmt.Sprintf("/proc/%d/cmdline", pid)
		cmdline, err := procinfo.ReadProcessCmdline(path)
		if err != nil {
			cmdline = ""
		}
		name := fmt.Sprintf("%d %s", pid, cmdline)
		names = append(names, name)
	}
	return names
}
