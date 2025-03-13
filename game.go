package main

import (
	"math/rand"
	"sync"
	"time"
)

// Cell описывает состояние одной клетки на поле
type Cell struct {
	IsMine    bool
	IsOpen    bool
	IsFlagged bool
	Neighbors int
	X         int
	Y         int
}

// Game описывает состояние игры
type Game struct {
	Grid           [][]Cell
	Width          int
	Height         int
	MineCount      int
	FlagsRemaining int
	GameOver       bool
	Win            bool
}

// CellWithGame используется для передачи данных в шаблон
type CellWithGame struct {
	Cell *Cell
	Game *Game
}

// NewGame создаёт новую игру с заданными параметрами
func NewGame(width, height, mineCount int) *Game {
	rand.Seed(time.Now().UnixNano())

	grid := make([][]Cell, height)

	var wg sync.WaitGroup
	wg.Add(height)

	for y := range grid {
		go func(y int) {
			defer wg.Done()
			grid[y] = make([]Cell, width)
			for x := range grid[y] {
				grid[y][x] = Cell{
					X: x,
					Y: y,
				}
			}
		}(y)
	}

	wg.Wait()

	game := &Game{
		Grid:           grid,
		Width:          width,
		Height:         height,
		MineCount:      mineCount,
		FlagsRemaining: mineCount,
		GameOver:       false,
		Win:            false,
	}

	game.placeMines()
	game.calculateNeighbors()

	return game
}

// placeMines случайным образом расставляет мины на поле
func (g *Game) placeMines() {
	minesPlaced := 0
	for minesPlaced < g.MineCount {
		x := rand.Intn(g.Width)
		y := rand.Intn(g.Height)

		if g.Grid[y][x].IsMine {
			continue
		}

		g.Grid[y][x].IsMine = true
		minesPlaced++
	}
}

// calculateNeighbors вычисляет количество мин вокруг каждой клетки
func (g *Game) calculateNeighbors() {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Grid[y][x].IsMine {
				continue
			}

			count := 0
			for _, dir := range directions {
				nx, ny := x+dir.dx, y+dir.dy
				if nx >= 0 && nx < g.Width && ny >= 0 && ny < g.Height && g.Grid[ny][nx].IsMine {
					count++
				}
			}
			g.Grid[y][x].Neighbors = count
		}
	}
}

// OpenCell открывает клетку по координатам (x, y)
// OpenCell открывает клетку по координатам (x, y)
func (g *Game) OpenCell(x, y int) {
	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return
	}

	cell := &g.Grid[y][x]

	if cell.IsOpen || cell.IsFlagged {
		return
	}

	cell.IsOpen = true

	// Если это мина, игра заканчивается
	if cell.IsMine {
		g.GameOver = true
		g.RevealAllMines() // Раскрываем все мины
		return
	}

	// Если клетка пустая (Neighbors == 0), рекурсивно открываем соседей
	if cell.Neighbors == 0 {
		g.openNeighbors(x, y)
	}

	// Проверяем, выиграл ли игрок
	g.checkWin()
}

// openNeighbors открывает все соседние клетки
func (g *Game) openNeighbors(x, y int) {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir.dx, y+dir.dy
		if nx >= 0 && nx < g.Width && ny >= 0 && ny < g.Height {
			g.OpenCell(nx, ny)
		}
	}
}

// checkWin проверяет, выиграл ли игрок
func (g *Game) checkWin() {
	openedSafeCells := 0
	totalSafeCells := g.Width*g.Height - g.MineCount

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Grid[y][x].IsOpen && !g.Grid[y][x].IsMine {
				openedSafeCells++
			}
		}
	}

	if openedSafeCells == totalSafeCells {
		g.Win = true
		g.GameOver = true
	}
}

// RevealAllMines открывает все клетки с минами
func (g *Game) RevealAllMines() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			cell := &g.Grid[y][x]
			if cell.IsMine {
				cell.IsOpen = true // Открываем клетку с миной
			}
		}
	}
}
