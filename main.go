package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
	Width    int
	Height   int
	Alive    int
}

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

func init_grid(width int, height int, coord_count int) grid {

	var grid = grid{
		width:      width,
		height:     height,
		cells:      make([]cell, width*height),
		next_cells: make([]cell, width*height),
	}

	for i := 0; i < height*width; i++ {
		grid.cells[i] = cell{alive: false}
		grid.next_cells[i] = cell{alive: false}
	}

	current_count_alive := 0
	retry_count := 0
	for current_count_alive < coord_count {
		c := rand.Intn(width * height)
		if grid.cells[c].alive {
			retry_count++
			continue
		}
		grid.cells[c].alive = true
		current_count_alive++
		if retry_count > 100 {
			break
		}
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

func init_db() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gol.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	db.Create(&User{ID: 1, Password: "123456", Name: "23Alive", Width: 50, Height: 50, Alive: 20})
	db.Create(&User{ID: 2, Password: "123456", Name: "23Alive_big", Width: 100, Height: 100, Alive: 600})
	return db
}

func init_grid_canva(grid grid) *canvas.Raster {
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
	return img
}

func main() {
	db := init_db()

	app := app.New()
	w := app.NewWindow("Test")

	var user User

	db.Where("name = ?", "23Alive_big").First(&user)

	grid := init_grid(user.Width, user.Height, user.Alive)

	grid_canva := init_grid_canva(grid)

	board := container.NewBorder(
		container.NewHBox(
			widget.NewButton("Test", func() {}),
			widget.NewButton("TOTO", func() {}),
		),
		nil,
		nil,
		nil,
		grid_canva,
	)

	go Update(grid, grid_canva)

	title := widget.NewLabel("Game of Life")
	title.Alignment = fyne.TextAlignCenter

	messageInput := widget.NewMultiLineEntry()

	messageForm := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "New Message", Widget: messageInput}},
		OnSubmit: func() { // optional, handle form submission
			fmt.Println("multiline:", messageInput.Text)
		},
	}

	chat := container.NewBorder(
		widget.NewLabel("Chat"),
		messageForm,
		nil,
		nil,
		widget.NewList(
			func() int {
				return 100
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Test")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(fmt.Sprintf("%d", i))
			},
		),
	)

	grid_layout := container.NewHSplit(board, chat)

	w.SetContent(
		container.NewBorder(
			title,
			nil,
			nil,
			nil,
			grid_layout,
		),
	)
	w.ShowAndRun()
}
