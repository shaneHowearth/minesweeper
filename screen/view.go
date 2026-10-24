package view

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/shanehowearth/minesweeper/controller"
	"github.com/shanehowearth/minesweeper/state"
)

// Classic Windows 95 Color Palette
var (
	WinGray      = color.RGBA{R: 192, G: 192, B: 192, A: 255}
	WinWhite     = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	WinDarkGray  = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	WinBlack     = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	NumberColors = []color.Color{
		nil,
		color.RGBA{R: 0, G: 0, B: 255, A: 255},     // 1: Blue
		color.RGBA{R: 0, G: 128, B: 0, A: 255},     // 2: Green
		color.RGBA{R: 255, G: 0, B: 0, A: 255},     // 3: Red
		color.RGBA{R: 0, G: 0, B: 128, A: 255},     // 4: Navy
		color.RGBA{R: 128, G: 0, B: 0, A: 255},     // 5: Maroon
		color.RGBA{R: 0, G: 128, B: 128, A: 255},   // 6: Teal
		color.RGBA{R: 0, G: 0, B: 0, A: 255},       // 7: Black
		color.RGBA{R: 128, G: 128, B: 128, A: 255}, // 8: Gray
	}
)

// Custom Tile Widget representing a single Minesweeper square
type TileWidget struct {
	widget.BaseWidget
	row, col  int
	ctrl      *controller.GameController
	board     *state.Board
	isHovered bool
}

func NewTileWidget(r, c int, ctrl *controller.GameController, board *state.Board) *TileWidget {
	t := &TileWidget{row: r, col: c, ctrl: ctrl, board: board}
	t.ExtendBaseWidget(t)
	return t
}

// Mouse event handling (Primary/Left Click & Secondary/Right Click)
func (t *TileWidget) Tapped(e *fyne.PointEvent) {
	t.ctrl.HandleReveal(t.row, t.col)
}

func (t *TileWidget) TappedSecondary(e *fyne.PointEvent) {
	t.ctrl.HandleToggleFlag(t.row, t.col)
}

// Implement Desktop canvas mouse listener for classic face reaction during click-down
func (t *TileWidget) MouseDown(e *desktop.MouseEvent) {
	if e.Button == desktop.MouseButtonPrimary {
		// Trigger "O" face on smiley during click down
		// (Signal your view manager to update face state if desired)
	}
}

func (t *TileWidget) MouseUp(e *desktop.MouseEvent) {}

func (t *TileWidget) CreateRenderer() fyne.WidgetRenderer {
	return &tileRenderer{tileWidget: t}
}

// Custom Renderer to draw classic 3D bevels and digits
type tileRenderer struct {
	tileWidget *TileWidget
}

func (r *tileRenderer) Destroy()           {}
func (r *tileRenderer) Layout(s fyne.Size) {}
func (r *tileRenderer) MinSize() fyne.Size {
	return fyne.NewSize(24, 24) // Classic 24x24 pixel tile size
}

func (r *tileRenderer) Refresh() {
	canvas.Refresh(r.tileWidget)
}

func (r *tileRenderer) Objects() []fyne.CanvasObject {
	t := r.tileWidget
	tile := t.board.GetTile(t.row, t.col)
	size := fyne.NewSize(24, 24)

	var objects []fyne.CanvasObject

	bg := canvas.NewRectangle(WinGray)
	bg.Resize(size)
	objects = append(objects, bg)

	if !tile.IsRevealed {
		// Raised 3D Bevel for Unrevealed Tiles
		top := canvas.NewLine(WinWhite, fyne.NewPos(0, 0), fyne.NewPos(24, 0))
		left := canvas.NewLine(WinWhite, fyne.NewPos(0, 0), fyne.NewPos(0, 24))
		bottom := canvas.NewLine(WinDarkGray, fyne.NewPos(0, 23), fyne.NewPos(24, 23))
		right := canvas.NewLine(WinDarkGray, fyne.NewPos(23, 0), fyne.NewPos(23, 24))

		top.StrokeWidth = 2
		left.StrokeWidth = 2
		bottom.StrokeWidth = 2
		right.StrokeWidth = 2

		objects = append(objects, top, left, bottom, right)

		if tile.IsFlagged {
			flagText := canvas.NewText("🚩", color.RGBA{R: 255, G: 0, B: 0, A: 255})
			flagText.TextSize = 14
			flagText.Alignment = fyne.TextAlignCenter
			flagText.Resize(size)
			objects = append(objects, flagText)
		}
	} else {
		// Sunken 3D Border for Revealed Tiles
		border := canvas.NewRectangle(color.Transparent)
		border.StrokeColor = WinDarkGray
		border.StrokeWidth = 1
		border.Resize(size)
		objects = append(objects, border)

		if tile.IsMine {
			mineText := canvas.NewText("💣", WinBlack)
			if t.board.State == state.StateLost && tile.IsRevealed {
				// Mine that was stepped on
				bg.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			}
			mineText.TextSize = 14
			mineText.Alignment = fyne.TextAlignCenter
			mineText.Resize(size)
			objects = append(objects, mineText)
		} else if tile.AdjacentMines > 0 {
			numText := canvas.NewText(
				fmt.Sprintf("%d", tile.AdjacentMines),
				NumberColors[tile.AdjacentMines],
			)
			numText.TextSize = 16
			numText.Style.Bold = true
			numText.Alignment = fyne.TextAlignCenter
			numText.Resize(size)
			objects = append(objects, numText)
		}
	}

	return objects
}
