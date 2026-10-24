package state

type GameState int

const (
	StateReady GameState = iota
	StateInProgress
	StateWon
	StateLost
)

type Tile struct {
	IsMine        bool
	IsRevealed    bool
	IsFlagged     bool
	AdjacentMines int
}

type Board struct {
	Rows       int
	Cols       int
	MinesCount int
	State      GameState
	Grid       [][]Tile
}

func NewBoard(rows, cols, minesCount int) *Board {
	grid := make([][]Tile, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]Tile, cols)
	}
	return &Board{
		Rows:       rows,
		Cols:       cols,
		MinesCount: minesCount,
		State:      StateReady,
		Grid:       grid,
	}
}

func (b *Board) IsValidCoord(r, c int) bool {
	return r >= 0 && r < b.Rows && c >= 0 && c < b.Cols
}

func (b *Board) GetTile(r, c int) *Tile {
	if !b.IsValidCoord(r, c) {
		return nil
	}
	return &b.Grid[r][c]
}
