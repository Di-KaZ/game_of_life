package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/gorilla/websocket"
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

func init_grid_canva(grid *Grid) *canvas.Raster {

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

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func init_ws(path string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	return c
}

type Message struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func messageHandler(ws *websocket.Conn, messages *[]Message, list *widget.List) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		msg := &Message{}
		json.Unmarshal(message, msg)
		*messages = append(*messages, *msg)
		list.Refresh()
	}
}

func gridWsUpdateHanlder(ws *websocket.Conn, grid *Grid, raster *canvas.Raster) {

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		json.Unmarshal(message, grid)
		raster.Refresh()
	}
}

func main() {
	messages := []Message{}

	grid := &Grid{}

	app := app.New()
	w := app.NewWindow("Test")

	wsGridState := init_ws("/state")
	defer wsGridState.Close()

	// time.Sleep(time.Second * 4)

	grid_canva := init_grid_canva(grid)
	go gridWsUpdateHanlder(wsGridState, grid, grid_canva)

	board := container.NewBorder(
		container.NewHBox(),
		nil,
		nil,
		nil,
		grid_canva,
	)

	title := widget.NewLabel("Game of Life")
	title.Alignment = fyne.TextAlignCenter

	messageInput := widget.NewMultiLineEntry()

	chatMessageList := widget.NewList(
		func() int {
			return len(messages)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Test")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			msg := messages[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %s", msg.Name, msg.Content))
		},
	)

	wsChatState := init_ws("/chat")
	defer wsChatState.Close()

	go messageHandler(wsChatState, &messages, chatMessageList)

	messageForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "New Message", Widget: messageInput},
		},
		OnSubmit: func() {
			wsChatState.WriteMessage(websocket.TextMessage, []byte(messageInput.Text))
			messageInput.Text = ""
			messageInput.Refresh()
			chatMessageList.Refresh()
		},
	}

	chat := container.NewBorder(
		widget.NewLabel("Chat"),
		messageForm,
		nil,
		nil,
		chatMessageList,
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
