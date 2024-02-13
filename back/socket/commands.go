package socket

import (
	"fmt"
	"gol_back/constants"
	"gol_back/grid"
	"strconv"
	"time"
)

func gen(ctx *ServerContext, Member *ChatMember, arg []string) {
	var width, height, alive_count int
	var err error

	if len(arg) != 3 {
		goto result_error
	}

	if width, err = strconv.Atoi(arg[0]); err != nil {
		goto result_error
	}

	if height, err = strconv.Atoi(arg[1]); err != nil {
		goto result_error
	}

	if alive_count, err = strconv.Atoi(arg[2]); err != nil {
		goto result_error
	}

	// ctx.Grid.Paused = true

	ctx.Grid = grid.Init(
		grid.GridConfig{
			Width:      width,
			Height:     height,
			StartAlive: alive_count,
		},
	)

	ctx.SendState(ctx.Grid)

	ctx.SendChat(
		Message{
			Name: constants.SERVER_DISPLAY_NAME,
			Content: fmt.Sprintf(
				constants.OK_GEN_COMMAND,
				ctx.Grid.Width,
				ctx.Grid.Height,
				alive_count,
			),
		})

result_error:
	ctx.SendChat(
		Message{
			Name:    constants.SERVER_DISPLAY_NAME,
			Content: constants.ERROR_GEN_COMMAND,
		},
	)
}

func playpause(ctx *ServerContext, Member *ChatMember, arg []string) {
	ctx.Grid.TogglePlay()

	var Content string

	if ctx.Grid.Paused {
		Content = constants.PAUSED_PLAYPAUSE_COMMAND
		go ctx.Grid.Loop(time.Second/10, func() {
			ctx.SendState(ctx.Grid)
		})
	} else {
		Content = constants.RESUME_PLAYPAUSE_COMMAND
	}

	ctx.SendChat(Message{
		Name:    constants.SERVER_DISPLAY_NAME,
		Content: Content,
	})
}

func ping(ctx *ServerContext, Member *ChatMember, arg []string) {
	ctx.SendChat(Message{Name: "server", Content: "pong"})
}

var Commands = map[string]func(ctx *ServerContext, Member *ChatMember, arg []string){
	"/gen":       gen,
	"/ping":      ping,
	"/playpause": playpause,
}
