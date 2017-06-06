package main

import (
	"testing"
)

func TestBoard(t *testing.T) {
	board := NewBoard(25, 25, nil)
	board.cells[0][0] = true
	board.cells[0][1] = true
	board.cells[0][2] = true

	board.Tick()

	if board.cells[0][0] != false {
		t.Error("Cells should be dead!")
	}

	if board.cells[1][1] != true {
		t.Error("Cells should be alive!")
	}

	if board.cells[4][1] != true {
		t.Error("Cells should be alive!")
	}
}

func TestBoardNorm(t *testing.T) {
	board := NewBoard(500, 500, nil)

	x := board.normX(101)

	if x != 1 {
		t.Error("Bad x: ", x)
	}

	y := board.normY(101)
	if y != 1 {
		t.Error("Bad y: ", y)
	}
}
