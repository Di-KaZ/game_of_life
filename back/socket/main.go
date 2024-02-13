package socket

import (
	"encoding/json"
	"flag"
	"gol_back/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var seed = time.Now().UTC().UnixNano()
var generator = namegenerator.NewNameGenerator(seed)

var upgrader = websocket.Upgrader{}

type ChatMember struct {
	Sub  *websocket.Conn
	Name string
}

type WebSocketServer struct {
	StateSub []*websocket.Conn
	ChatSub  []*ChatMember
}

func HandleBoardState(ctx *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		ctx.Server.StateSub = append(ctx.Server.StateSub, ws)
		ctx.Server.SendState(ctx.Grid)
	}
}

func HandleChat(ctx *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		member := &ChatMember{Sub: ws, Name: generator.Generate()}
		ctx.Server.ChatSub = append(ctx.Server.ChatSub, member)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			// handle commands if any
			ctx.Server.SendChat(Message{Name: member.Name, Content: string(message)})
			split := strings.Fields(string(message))
			if command, ok := Commands[split[0]]; ok {
				command(ctx,
					member,
					split[1:],
				)
			} else if message[0] == '/' {
				ctx.Server.SendChat(Message{Name: "server", Content: "command not found"})
			}
		}
	}
}

func (s *WebSocketServer) Close() {
	for i := range s.StateSub {
		s.StateSub[i].Close()
	}
	for i := range s.ChatSub {
		s.ChatSub[i].Sub.Close()
	}
}

func (s *WebSocketServer) SendState(grid *models.Grid) {
	json, _ := json.Marshal(grid)
	print(json)
	for i := range s.StateSub {
		s.StateSub[i].WriteMessage(websocket.TextMessage, json)
	}
}

func (s *WebSocketServer) SendChat(message Message) {
	json, _ := json.Marshal(message)
	for i := range s.ChatSub {
		s.ChatSub[i].Sub.WriteMessage(websocket.TextMessage, json)
	}
}

func InitWebsocketServer(db *gorm.DB, server *WebSocketServer) {
	ctx := &ServerContext{
		Db:     db,
		Grid:   &models.Grid{Paused: true},
		Server: server,
	}

	addr := flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/state", HandleBoardState(ctx))

	http.HandleFunc("/chat", HandleChat(ctx))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
