package grid

import (
	"math/rand"
	"time"
)

type Grid struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	nextCells []Cell `json:"-"`
	Cells     []Cell `json:"cells"`
	Paused    bool
}

func (g *Grid) Reset(config GridConfig) {
	cell_count := config.Width * config.Height
	g.Paused = true
	g.Width = config.Width
	g.Height = config.Height
	g.Cells = make([]Cell, cell_count)
	g.nextCells = make([]Cell, cell_count)

	for i := 0; i < cell_count; i++ {
		g.Cells[i] = Cell{Alive: false}
		g.nextCells[i] = Cell{Alive: false}
	}

	current_count_alive := 0
	retry_count := 0
	for current_count_alive < config.StartAlive {
		c := rand.Intn(cell_count)
		if g.Cells[c].Alive {
			retry_count++
			continue
		}
		g.Cells[c].Alive = true
		current_count_alive++
		if retry_count > 100 {
			break
		}
	}

}

func Init(config GridConfig) *Grid {

	grid := &Grid{}

	grid.Reset(config)

	return grid
}

func (g *Grid) _isValidCoordinateAndAlive(x int, y int) bool {
	return x > 0 &&
		x < g.Width &&
		y > 0 &&
		y < g.Height &&
		g.Cells[y*g.Width+x].Alive
}

func (g *Grid) _countNeighbors(y int, x int) int {
	var neighbors = 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if g._isValidCoordinateAndAlive(x+j, y+i) {
				neighbors++
			}
		}
	}
	return neighbors
}

func (g *Grid) _23Alive(pos int) {
	var y = pos / g.Width
	var x = pos % g.Width
	var neighbors = g._countNeighbors(y, x)
	if neighbors == 3 && !g.Cells[pos].Alive {
		g.nextCells[pos].Alive = true
	} else if (neighbors == 3 || neighbors == 2) && g.Cells[pos].Alive {
		g.nextCells[pos].Alive = true
	} else {
		g.nextCells[pos].Alive = false
	}
	if g.Cells[pos].Alive && g.nextCells[pos].Alive {
		g.nextCells[pos].Turns += 1
	} else {
		g.nextCells[pos].Turns = 0
	}
}

func (g *Grid) TogglePlay() {
	g.Paused = !g.Paused
}

func (g *Grid) Loop(rate time.Duration, callback func()) {
	for pos := range g.Cells {
		g._23Alive(pos)
	}

	for i := range g.Cells {
		g.Cells[i], g.nextCells[i] = g.nextCells[i], g.Cells[i]
	}

	time.Sleep(rate)

	callback()

	if !g.Paused {
		g.Loop(rate, callback)
	}
}
