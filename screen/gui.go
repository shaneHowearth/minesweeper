package screen

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/shanehowearth/minesweeper/controller"
	"github.com/shanehowearth/minesweeper/state"
)

type FyneMinesweeperView struct {
	app         fyne.App
	win         fyne.Window
	board       *state.Board
	ctrl        *controller.GameController
	gridBox     *fyne.Container
	smileyBtn   *widget.Button
	mineCounter *canvas.Text
}

func NewFyneMinesweeperView(board *state.Board, ctrl *controller.GameController) *FyneMinesweeperView {
	a := app.New()
	w := a.NewWindow("Minesweeper")

	v := &FyneMinesweeperView{
		app:   a,
		win:   w,
		board: board,
		ctrl:  ctrl,
	}

	v.setupUI()
	return v
}

func (v *FyneMinesweeperView) setupUI() {
	// 1. Digital LED Counter Display (Mines Left)
	v.mineCounter = canvas.NewText("099", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	v.mineCounter.TextSize = 22
	v.mineCounter.TextStyle.Monospace = true
	counterBg := canvas.NewRectangle(WinBlack)
	counterBox := container.NewStack(counterBg, v.mineCounter)

	// 2. Central Smiley Face Button
	v.smileyBtn = widget.NewButton("🙂", func() {
		// Reset game
		*v.board = *state.NewBoard(v.board.Rows, v.board.Cols, v.board.MinesCount)
		v.Refresh()
	})

	header := container.NewBorder(nil, nil, counterBox, nil, container.NewCenter(v.smileyBtn))

	// 3. Grid of TileWidgets
	v.gridBox = container.NewGridWithColumns(v.board.Cols)
	v.buildGrid()

	// Outer sunken 3D frame wrapper
	mainLayout := container.NewVBox(
		header,
		widget.NewSeparator(),
		v.gridBox,
	)

	v.win.SetContent(container.NewPadded(mainLayout))
	v.win.SetFixedSize(true)
}

func (v *FyneMinesweeperView) buildGrid() {
	v.gridBox.Objects = nil
	for r := 0; r < v.board.Rows; r++ {
		for c := 0; c < v.board.Cols; c++ {
			// Pass 'v' (the FyneMinesweeperView pointer) as the 5th argument here:
			tileWidget := NewTileWidget(r, c, v.ctrl, v.board, v)
			v.gridBox.Add(tileWidget)
		}
	}
}

func (v *FyneMinesweeperView) Refresh() {
	// Update Smiley Face depending on GameState
	switch v.board.State {
	case state.StateReady, state.StateInProgress:
		v.smileyBtn.SetText("🙂")
	case state.StateWon:
		v.smileyBtn.SetText("😎")
	case state.StateLost:
		v.smileyBtn.SetText("😵")
	}

	// Update Mine Counter (Total Mines - Flagged Count)
	flagsUsed := 0
	for r := 0; r < v.board.Rows; r++ {
		for c := 0; c < v.board.Cols; c++ {
			if v.board.GetTile(r, c).IsFlagged {
				flagsUsed++
			}
		}
	}
	v.mineCounter.Text = fmt.Sprintf("%03d", v.board.MinesCount-flagsUsed)

	v.gridBox.Refresh()
}

func (v *FyneMinesweeperView) Run() {
	v.win.ShowAndRun()
}
