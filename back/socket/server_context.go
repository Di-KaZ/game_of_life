package socket

import (
	"encoding/json"
	"gol_back/grid"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ServerContext struct {
	Db       *gorm.DB
	Grid     *grid.Grid
	StateSub []*websocket.Conn
	ChatSub  []*ChatMember
}

func (s *ServerContext) Close() {
	for i := range s.StateSub {
		s.StateSub[i].Close()
	}
	for i := range s.ChatSub {
		s.ChatSub[i].Sub.Close()
	}
}

func (s *ServerContext) SendState(grid *grid.Grid) {
	json, _ := json.Marshal(grid)
	print(json)
	for i := range s.StateSub {
		s.StateSub[i].WriteMessage(websocket.TextMessage, json)
	}
}

func (s *ServerContext) SendChat(message Message) {
	json, _ := json.Marshal(message)
	for i := range s.ChatSub {
		s.ChatSub[i].Sub.WriteMessage(websocket.TextMessage, json)
	}
}
