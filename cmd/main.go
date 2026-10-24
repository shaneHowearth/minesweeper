package main

import (
	"github.com/shanehowearth/minesweeper/controller"
	"github.com/shanehowearth/minesweeper/screen"
	"github.com/shanehowearth/minesweeper/state"
)

func main() {
	// Standard Beginner Minesweeper dimensions: 9 rows, 9 cols, 10 mines
	rows := 9
	cols := 9
	minesCount := 10

	// 1. Initialize the Model (Pure data state from state package)
	board := state.NewBoard(rows, cols, minesCount)

	// 2. Initialize the Controller (Enforces game rules from control package)
	ctrl := controller.NewGameController(board)

	// 3. Initialize the View (Fyne graphical user interface from screen package)
	gameView := screen.NewFyneMinesweeperView(board, ctrl)

	// 4. Run the application event loop
	gameView.Run()
}
