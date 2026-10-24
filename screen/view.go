package screen

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	row, col   int
	ctrl       *controller.GameController
	board      *state.Board
	parentView *FyneMinesweeperView
}

func NewTileWidget(r, c int, ctrl *controller.GameController, board *state.Board, parentView *FyneMinesweeperView) *TileWidget {
	t := &TileWidget{
		row:        r,
		col:        c,
		ctrl:       ctrl,
		board:      board,
		parentView: parentView,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *TileWidget) Tapped(e *fyne.PointEvent) {
	t.ctrl.HandleReveal(t.row, t.col)
	t.parentView.Refresh() // Triggers global refresh, updating this tile
}

func (t *TileWidget) TappedSecondary(e *fyne.PointEvent) {
	t.ctrl.HandleToggleFlag(t.row, t.col)
	t.parentView.Refresh() // Triggers global refresh, updating this tile
}

func (t *TileWidget) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(WinGray)

	// 3D Bevel lines for unrevealed state
	top := canvas.NewRectangle(WinWhite)
	left := canvas.NewRectangle(WinWhite)
	bottom := canvas.NewRectangle(WinDarkGray)
	right := canvas.NewRectangle(WinDarkGray)

	label := canvas.NewText("", WinBlack)
	label.Alignment = fyne.TextAlignCenter
	label.TextSize = 14

	return &tileRenderer{
		tileWidget: t,
		bg:         bg,
		top:        top,
		left:       left,
		bottom:     bottom,
		right:      right,
		label:      label,
		objects:    []fyne.CanvasObject{bg, top, left, bottom, right, label},
	}
}

type tileRenderer struct {
	tileWidget               *TileWidget
	bg                       *canvas.Rectangle
	top, left, bottom, right *canvas.Rectangle
	label                    *canvas.Text
	objects                  []fyne.CanvasObject
}

func (r *tileRenderer) Destroy() {}

func (r *tileRenderer) MinSize() fyne.Size {
	return fyne.NewSize(24, 24)
}

func (r *tileRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)

	// Bevel lines positioning
	r.top.Resize(fyne.NewSize(size.Width, 2))
	r.top.Move(fyne.NewPos(0, 0))

	r.left.Resize(fyne.NewSize(2, size.Height))
	r.left.Move(fyne.NewPos(0, 0))

	r.bottom.Resize(fyne.NewSize(size.Width, 2))
	r.bottom.Move(fyne.NewPos(0, size.Height-2))

	r.right.Resize(fyne.NewSize(2, size.Height))
	r.right.Move(fyne.NewPos(size.Width-2, 0))

	r.label.Resize(size)
	r.label.Move(fyne.NewPos(0, 2)) // Center alignment tweak
}

func (r *tileRenderer) Refresh() {
	tile := r.tileWidget.board.GetTile(r.tileWidget.row, r.tileWidget.col)

	if !tile.IsRevealed {
		r.bg.FillColor = WinGray
		r.top.Show()
		r.left.Show()
		r.bottom.Show()
		r.right.Show()

		if tile.IsFlagged {
			r.label.Text = "🚩"
			r.label.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			r.label.Show()
		} else {
			r.label.Text = ""
			r.label.Hide()
		}
	} else {
		// Revealed state
		r.top.Hide()
		r.left.Hide()
		r.bottom.Hide()
		r.right.Hide()

		if tile.IsMine {
			r.bg.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			r.label.Text = "💣"
			r.label.Show()
		} else if tile.AdjacentMines > 0 {
			r.bg.FillColor = WinGray
			r.label.Text = fmt.Sprintf("%d", tile.AdjacentMines)
			r.label.Color = NumberColors[tile.AdjacentMines]
			r.label.TextStyle.Bold = true
			r.label.Show()
		} else {
			r.bg.FillColor = color.RGBA{R: 220, G: 220, B: 220, A: 255} // Sunken light gray
			r.label.Text = ""
			r.label.Hide()
		}
	}

	r.bg.Refresh()
	r.label.Refresh()
}

func (r *tileRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}
