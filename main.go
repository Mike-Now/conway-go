package main

import (
	"github.com/andlabs/ui"
	"math/rand"
	"time"
)

type Board struct {
	area     *ui.Area
	width    int
	height   int
	cells    [][]bool
	tmp      [][]bool
	cellSize int
	paused   bool
}

func (b *Board) Redraw() {
	ui.QueueMain(func() {
		b.area.QueueRedrawAll()
	})
}

func (b *Board) normY(y int) int {
	return Mod(y, b.height)
}

func (b *Board) normX(x int) int {
	return Mod(x, b.width)
}

func NewBoard(width int, height int, area *ui.Area) *Board {
	board := &Board{
		width:    width / 5,
		height:   height / 5,
		cells:    make([][]bool, height/5),
		tmp:      make([][]bool, height/5),
		cellSize: 5,
		area:     area,
	}
	cells1D := make([]bool, board.width*board.height)
	for i := range board.cells {
		board.cells[i] = cells1D[i*board.width : (i+1)*board.width]
	}
	tmp1D := make([]bool, board.width*board.height)
	for i := range board.tmp {
		board.tmp[i] = tmp1D[i*board.width : (i+1)*board.width]
	}
	return board
}

type CellOffset struct {
	y int
	x int
}

var Offsets = [...]CellOffset{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func Mod(dividend int, divisor int) int {
	result := dividend % divisor
	if result < 0 {
		return result + divisor
	}
	return result
}

func (b *Board) Tick() {
	for y, row := range b.tmp {
		for x := range row {
			b.tmp[y][x] = false
		}
	}
	for y, row := range b.cells {
		for x, alive := range row {
			aliveNeighbours := 0
			for _, offset := range Offsets {
				yn := b.normY(y + offset.y)
				xn := b.normX(x + offset.x)
				if b.cells[yn][xn] {
					aliveNeighbours += 1
				}
			}
			if alive && aliveNeighbours < 2 {
				b.tmp[y][x] = false
			} else if alive && aliveNeighbours < 4 {
				b.tmp[y][x] = true
			} else if alive {
				b.tmp[y][x] = false
			} else if aliveNeighbours == 3 {
				b.tmp[y][x] = true
			}
		}
	}

	for y, row := range b.cells {
		for x := range row {
			b.cells[y][x] = b.tmp[y][x]
		}
	}
}

func (b *Board) Draw(a *ui.Area, dp *ui.AreaDrawParams) {
	// Paint background
	p := ui.NewPath(ui.Winding)
	p.AddRectangle(dp.ClipX, dp.ClipY, dp.ClipWidth, dp.ClipHeight)
	p.End()

	dp.Context.Fill(p, &ui.Brush{
		Type: ui.Solid,
		R:    float64(255) / float64(255),
		G:    float64(240) / float64(255),
		B:    float64(204) / float64(255),
		A:    1,
	})
	p.Free()

	// Paint cells
	for y, row := range b.cells {
		for x, alive := range row {
			if alive {
				p = ui.NewPath(ui.Winding)
				p.AddRectangle(float64(x*b.cellSize), float64(y*b.cellSize), float64(b.cellSize), float64(b.cellSize))
				p.End()
				dp.Context.Fill(p, &ui.Brush{
					Type: ui.Solid,
					R:    float64(166) / float64(255),
					G:    float64(170) / float64(255),
					B:    float64(204) / float64(255),
					A:    1,
				})
				p.Free()
			}
		}
	}
}
func (b *Board) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
}
func (b *Board) MouseCrossed(a *ui.Area, left bool) {
}
func (b *Board) DragBroken(a *ui.Area)                                   {}
func (b *Board) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) { return true }

func (b *Board) addCell(y int, x int) {
	b.cells[b.normY(y)][b.normX(x)] = true
}
func (b *Board) addGlider(y int, x int) {
	b.addCell(y-1, x)
	b.addCell(y, x+1)
	b.addCell(y+1, x-1)
	b.addCell(y+1, x)
	b.addCell(y+1, x+1)
}

func main() {
	width := 1000
	height := 600
	err := ui.Main(func() {
		window := ui.NewWindow("automata", width, height, false)
		box := ui.NewVerticalBox()
		window.SetChild(box)
		board := NewBoard(width, height-50, nil)
		for i := 0; i < 6000; i += 1 {
			board.addCell(rand.Int(), rand.Int())
		}

		area := ui.NewScrollingArea(board, width, height-50)
		board.area = area
		button := ui.NewButton("Pause")
		stateChan := make(chan struct{})
		button.OnClicked(func(*ui.Button) {
			stateChan <- struct{}{}
		})

		box.Append(button, false)
		box.Append(area, true)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
		go tickTimer(board, stateChan)
	})
	if err != nil {
		panic(err)
	}
}

func tickTimer(b *Board, ch <-chan struct{}) {
	paused := false
	for {
		select {
		case <-ch:
			paused = !paused
		default:
		}
		if !paused {
			ui.QueueMain(func() {
				b.Tick()
			})
			b.Redraw()
		}
		time.Sleep(50 * time.Millisecond)
	}
}
