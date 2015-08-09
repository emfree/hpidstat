package main

import (
	"flag"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getPids(comm string) []uint64 {
	var pids []uint64
	euid := strconv.Itoa(os.Geteuid())
	matches, err := exec.Command("pgrep", "-u", euid, comm).Output()
	if err != nil {
		// TODO(emfree) handle case of no results (exit status 1)
		panic(err)
	}
	fmt.Println(string(matches))
	for _, p := range strings.Split(string(matches), "\n") {
		if p != "" {
			pid, err := strconv.Atoi(p)
			if err != nil {
				panic(err)
			}
			pids = append(pids, uint64(pid))
		}
	}
	return pids
}

func main() {
	var comm string
	var interval float64
	flag.StringVar(&comm, "C", "",
		"Measure tasks whose command name includes the string comm")
	flag.Float64Var(&interval, "i", 1., "Plot every interval seconds")
	flag.Parse()
	pids := getPids(comm)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output216)
	drawLegend(getNames(pids), 60)

	x0 := 0
	y0 := 20
	w, h := termbox.Size()
	w *= 2
	h = h - y0
	canvas := Canvas{}
	canvas.Init(x0, y0, w, h)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

loop:
	for i := 0; ; i++ {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
				break loop
			}
		default:
			stats := getProcFractions(pids,
				time.Duration(int64(interval*1e3))*time.Millisecond)
			canvas.ScaleColumn(stats, i%w)
			var rot int
			if i < w {
				rot = 0
			} else {
				rot = (i + 1) % w
			}
			canvas.draw(rot)
		}
	}
}
