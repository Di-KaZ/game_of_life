package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gol_front/chat"
	"gol_front/grid"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func init_ws(path string) *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	return c
}

func messageHandler(ws *websocket.Conn, w *chat.ChatWidget) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		msg := &chat.ChatMessage{}
		json.Unmarshal(message, msg)
		w.AddMessage(*msg)
	}
}

func gridWsUpdateHanlder(ws *websocket.Conn, grid *grid.GridWidget) {

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		json.Unmarshal(message, grid.Grid)
		grid.Refresh()
	}
}

func main() {
	app := app.New()
	w := app.NewWindow("Test")

	wsGridState := init_ws("/state")
	defer wsGridState.Close()

	wsChatState := init_ws("/chat")
	defer wsChatState.Close()

	chat := chat.Init(
		chat.ChatWidgetProps{
			OnSubmit: func(text string) {
				wsChatState.WriteMessage(websocket.TextMessage, []byte(text))
			},
		},
	)

	grid := grid.Init()

	go gridWsUpdateHanlder(wsGridState, grid)

	go messageHandler(wsChatState, chat)

	grid_layout := container.NewHSplit(grid.Widget, chat.Widget)

	title := widget.NewLabel("Game of Life")
	title.Alignment = fyne.TextAlignCenter

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
