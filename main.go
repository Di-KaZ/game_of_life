package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type cell struct {
	alive bool
	turns int
}

type grid struct {
	width      int
	height     int
	next_cells []cell
	cells      []cell
}

func init_grid() grid {
	var width = 50
	var height = 50

	var grid = grid{
		width:      width,
		height:     height,
		cells:      make([]cell, width*height),
		next_cells: make([]cell, width*height),
	}

	for i := 0; i < height*width; i++ {
		grid.cells[i] = cell{alive: rand.Intn(2) == 1}
		grid.next_cells[i] = cell{alive: false}
	}
	return grid
}

func (g grid) isValidCoordinateAndAlive(x int, y int) bool {
	return x > 0 &&
		x < g.width &&
		y > 0 &&
		y < g.height &&
		g.cells[y*g.width+x].alive
}

func (g grid) countNeighbors(y int, x int) int {
	var neighbors = 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if g.isValidCoordinateAndAlive(x+j, y+i) {
				neighbors++
			}
		}
	}
	return neighbors
}

func (g grid) _23Alive(pos int) {
	var y = pos / g.width
	var x = pos % g.width
	var neighbors = g.countNeighbors(y, x)
	if neighbors == 3 && !g.cells[pos].alive {
		g.next_cells[pos].alive = true
	} else if (neighbors == 3 || neighbors == 2) && g.cells[pos].alive {
		g.next_cells[pos].alive = true
	} else {
		g.next_cells[pos].alive = false
	}
	if g.cells[pos].alive && g.next_cells[pos].alive {
		g.next_cells[pos].turns += 1
	} else {
		g.next_cells[pos].turns = 0
	}
}

func (g grid) Tick() {
	for pos := 0; pos < g.width*g.height; pos++ {
		g._23Alive(pos)
	}
	for i := range g.cells {
		g.cells[i], g.next_cells[i] = g.next_cells[i], g.cells[i]
	}
}

func Update(grid grid, image *canvas.Raster) {
	time.Sleep(time.Second / 15)
	grid.Tick()
	image.Refresh()
	Update(grid, image)
}

func main() {
	app := app.New()
	w := app.NewWindow("Test")

	grid := init_grid()
	img := canvas.NewRaster(func(w, h int) image.Image {
		wRatio := w / grid.width
		hRatio := h / grid.height
		pixelSize := int(math.Min(float64(wRatio), float64(hRatio)))

		raster := image.NewRGBA(image.Rect(0, 0, w, h))

		for pos := range grid.cells {
			x := pos % grid.width
			y := pos / grid.height
			cell := raster.SubImage(image.Rect(x*pixelSize, y*pixelSize, (x+1)*pixelSize, (y+1)*pixelSize)).(*image.RGBA)

			if grid.cells[pos].alive {
				switch grid.cells[pos].turns {
				case 0:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{52, 102, 122, 255}}, image.Point{}, draw.Src)
				case 1:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{73, 121, 138, 255}}, image.Point{}, draw.Src)
				default:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{214, 151, 0, 255}}, image.Point{}, draw.Src)
				}
			} else if !grid.cells[pos].alive {
				draw.Draw(cell, cell.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
			}
		}
		return raster
	})

	go Update(grid, img)

	title := widget.NewLabel("Game of Life")
	title.Alignment = fyne.TextAlignCenter

	content := container.NewBorder(
		title,
		nil,
		nil,
		nil,
		img,
	)
	// change flyne app color background

	w.SetContent(content)
	w.ShowAndRun()
}
