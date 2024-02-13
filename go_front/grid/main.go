package grid

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Cell struct {
	Alive bool `json:"alive"`
	Turns int  `json:"turns"`
}

type Grid struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Cells  []Cell `json:"cells"`
	Paused bool
}

type GridWidget struct {
	Grid   *Grid
	raster *canvas.Raster
	Widget *fyne.Container
}

func initCanva(grid *Grid) *canvas.Raster {

	img := canvas.NewRaster(func(w, h int) image.Image {
		if grid.Height == 0 || grid.Width == 0 {
			return image.Black
		}
		wRatio := w / grid.Width
		hRatio := h / grid.Height
		pixelSize := int(math.Min(float64(wRatio), float64(hRatio)))

		raster := image.NewRGBA(image.Rect(0, 0, w, h))

		for pos := range grid.Cells {
			x := pos % grid.Width
			y := pos / grid.Height
			cell := raster.SubImage(image.Rect(x*pixelSize, y*pixelSize, (x+1)*pixelSize, (y+1)*pixelSize)).(*image.RGBA)

			if grid.Cells[pos].Alive {
				switch grid.Cells[pos].Turns {
				case 0:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{52, 102, 122, 255}}, image.Point{}, draw.Src)
				case 1:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{73, 121, 138, 255}}, image.Point{}, draw.Src)
				default:
					draw.Draw(cell, cell.Bounds(), &image.Uniform{color.RGBA{214, 151, 0, 255}}, image.Point{}, draw.Src)
				}
			} else if !grid.Cells[pos].Alive {
				draw.Draw(cell, cell.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
			}
		}
		return raster
	})
	return img
}

func Init() *GridWidget {
	gridWidget := &GridWidget{}
	gridWidget.Grid = &Grid{}

	gridWidget.raster = initCanva(gridWidget.Grid)

	gridWidget.Widget = container.NewBorder(
		container.NewHBox(),
		nil,
		nil,
		nil,
		gridWidget.raster,
	)
	return gridWidget
}

func (g *GridWidget) Refresh() {
	g.raster.Refresh()
}
