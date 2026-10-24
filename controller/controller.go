package controller

import (
	"math/rand"
	"time"

	"github.com/shanehowearth/minesweeper/state"
)

type GameController struct {
	board *state.Board
	rng   *rand.Rand
}

func NewGameController(board *state.Board) *GameController {
	return &GameController{
		board: board,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// User Action: Primary Click
func (gc *GameController) HandleReveal(row, col int) {
	if !gc.board.IsValidCoord(row, col) || gc.board.State == state.StateWon || gc.board.State == state.StateLost {
		return
	}

	// Rule 1: First click generates mines safely
	if gc.board.State == state.StateReady {
		gc.generateMines(row, col)
		gc.calculateNeighbors()
		gc.board.State = state.StateInProgress
	}

	tile := gc.board.GetTile(row, col)
	if tile.IsRevealed || tile.IsFlagged {
		return
	}

	tile.IsRevealed = true

	// Rule 2: Loss Condition
	if tile.IsMine {
		gc.board.State = state.StateLost
		gc.revealAllMines()
		return
	}

	// Rule 3: Zero-Cascade Flood Fill
	if tile.AdjacentMines == 0 {
		gc.cascadeReveal(row, col)
	}

	// Rule 4: Win Condition Check
	gc.checkWinCondition()
}

// User Action: Secondary Click
func (gc *GameController) HandleToggleFlag(row, col int) {
	if !gc.board.IsValidCoord(row, col) {
		return
	}
	if gc.board.State != state.StateInProgress && gc.board.State != state.StateReady {
		return
	}

	tile := gc.board.GetTile(row, col)
	if !tile.IsRevealed {
		tile.IsFlagged = !tile.IsFlagged
	}
}

// Private Controller Rules & Algorithms

func (gc *GameController) generateMines(safeR, safeC int) {
	placed := 0
	for placed < gc.board.MinesCount {
		r := gc.rng.Intn(gc.board.Rows)
		col := gc.rng.Intn(gc.board.Cols)

		tile := gc.board.GetTile(r, col)
		if tile.IsMine || (r == safeR && col == safeC) {
			continue
		}

		tile.IsMine = true
		placed++
	}
}

func (c *GameController) calculateNeighbors() {
	for r := 0; r < c.board.Rows; r++ {
		for col := 0; col < c.board.Cols; col++ {
			tile := c.board.GetTile(r, col)
			if tile.IsMine {
				continue
			}

			count := 0
			c.forEachNeighbor(r, col, func(nr, nc int) {
				if c.board.GetTile(nr, nc).IsMine {
					count++
				}
			})
			tile.AdjacentMines = count
		}
	}
}

func (gc *GameController) cascadeReveal(row, col int) {
	gc.forEachNeighbor(row, col, func(nr, nc int) {
		neighbor := gc.board.GetTile(nr, nc)
		if !neighbor.IsRevealed && !neighbor.IsFlagged {
			gc.HandleReveal(nr, nc)
		}
	})
}

func (gc *GameController) forEachNeighbor(row, col int, fn func(nr, nc int)) {
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := row+dr, col+dc
			if gc.board.IsValidCoord(nr, nc) {
				fn(nr, nc)
			}
		}
	}
}

func (c *GameController) checkWinCondition() {
	unrevealed := 0
	for r := 0; r < c.board.Rows; r++ {
		for col := 0; col < c.board.Cols; col++ {
			if !c.board.GetTile(r, col).IsRevealed {
				unrevealed++
			}
		}
	}

	if unrevealed == c.board.MinesCount {
		c.board.State = state.StateWon

		// Authentic Polish: Automatically flag all remaining hidden mines upon winning
		for r := 0; r < c.board.Rows; r++ {
			for col := 0; col < c.board.Cols; col++ {
				tile := c.board.GetTile(r, col)
				if !tile.IsRevealed && tile.IsMine {
					tile.IsFlagged = true
				}
			}
		}
	}
}

func (gc *GameController) revealAllMines() {
	for r := 0; r < gc.board.Rows; r++ {
		for col := 0; col < gc.board.Cols; col++ {
			tile := gc.board.GetTile(r, col)
			if tile.IsMine {
				tile.IsRevealed = true
			}
		}
	}
}
