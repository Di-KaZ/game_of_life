package socket

import (
	"fmt"
	"gol_back/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ServerContext struct {
	Db     *gorm.DB
	Server *WebSocketServer
	Grid   *models.Grid
}

func newGen(ctx *ServerContext, Member *ChatMember, arg []string) {
	if len(arg) != 3 {
		ctx.Server.SendChat(Message{Name: "server", Content: "[ERROR] /newGen width height alive_count"})
		return
	}

	width, err_w := strconv.Atoi(arg[0])
	height, err_h := strconv.Atoi(arg[1])
	alive, err_a := strconv.Atoi(arg[2])

	if err_w != nil || err_a != nil || err_h != nil {
		ctx.Server.SendChat(Message{Name: "server", Content: "[ERROR] /newGen width height alive_count"})
		return
	}

	ctx.Grid.Paused = true

	ctx.Grid = models.Init(
		models.GridConfig{
			Width:      width,
			Height:     height,
			StartAlive: alive,
		},
	)

	ctx.Server.SendState(ctx.Grid)

	ctx.Server.SendChat(Message{Name: "server", Content: fmt.Sprintf("Requested new grid width:%dm height:%d, alive: %d", ctx.Grid.Width, ctx.Grid.Height, alive)})
}

func togglePlay(ctx *ServerContext, Member *ChatMember, arg []string) {
	ctx.Grid.TogglePause()
	ctx.Server.SendChat(Message{Name: "server", Content: "toggled play/pause"})

	if !ctx.Grid.Paused {
		go ctx.Grid.Loop(time.Second/10, func() {
			ctx.Server.SendState(ctx.Grid)
		})
	}
}

func ping(ctx *ServerContext, Member *ChatMember, arg []string) {
	ctx.Server.SendChat(Message{Name: "server", Content: "pong"})
}

var Commands = map[string]func(ctx *ServerContext, Member *ChatMember, arg []string){
	"/newGen":     newGen,
	"/ping":       ping,
	"/togglePlay": togglePlay,
}
