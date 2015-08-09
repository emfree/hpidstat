package main

import (
	"github.com/nsf/termbox-go"
	"math"
)

type Canvas struct {
	X0, Y0, Width, Height int
	Cells                 [][]int
}

func (canvas *Canvas) Init(x0, y0, width, height int) {
	canvas.X0 = x0
	canvas.Y0 = y0
	canvas.Width = width
	canvas.Height = height
	canvas.Cells = make([][]int, canvas.Width)
	for i := 0; i < canvas.Width; i++ {
		canvas.Cells[i] = make([]int, canvas.Height)
	}
}

func getColor(val int) int {
	return (109*val)%216 + 1
}

func (canvas *Canvas) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x+canvas.X0, y+canvas.Y0, ch, fg, bg)
}

func (canvas *Canvas) draw(rot int) {
	for col := 0; col < canvas.Width; col += 2 {
		for row := 0; row < canvas.Height; row++ {
			lval := canvas.Cells[(col+rot)%canvas.Width][row]
			rval := canvas.Cells[(col+rot+1)%canvas.Width][row]
			if lval != 0 {
				canvas.SetCell(col/2, row, '\u258C', termbox.Attribute(lval),
					termbox.Attribute(rval))
			} else if rval != 0 {
				canvas.SetCell(col/2, row, '\u2590', termbox.Attribute(rval),
					termbox.Attribute(lval))
			} else {
				canvas.SetCell(col/2, row, ' ', termbox.ColorDefault,
					termbox.ColorDefault)
			}
		}
	}
	termbox.Flush()
}

func scale(fracs []float64, dest []int) {
	offset := len(dest) - 1
	for val, frac := range fracs {
		band := int(math.Ceil(float64(len(dest)) * frac))
		for i := 0; i < band && offset >= 0; i++ {
			dest[offset] = getColor(val)
			offset--
		}
	}
}

func (canvas *Canvas) ScaleColumn(fracs []float64, x int) {
	scale(fracs, canvas.Cells[x])
}

func tbprint(x int, y int, fg termbox.Attribute, bg termbox.Attribute,
	msg string) {
	for _, char := range msg {
		termbox.SetCell(x, y, char, fg, bg)
		x++
	}
}

func drawLegend(names []string, maxlen int) {
	for ind, name := range names {
		col := getColor(ind)
		termbox.SetCell(0, ind, ' ', termbox.ColorDefault,
			termbox.Attribute(col))
		if maxlen > len(name) {
			maxlen = len(name)
		}
		tbprint(2, ind, termbox.ColorDefault, termbox.ColorDefault, name[:maxlen])
	}
}
