package socket

import (
	"flag"
	"gol_back/grid"
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

func HandleBoardState(ctx *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		ctx.StateSub = append(ctx.StateSub, ws)
		ctx.SendState(ctx.Grid)
	}
}

func HandleChat(ctx *ServerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		// attribute the new connection a name
		member := &ChatMember{Sub: ws, Name: generator.Generate()}
		ctx.ChatSub = append(ctx.ChatSub, member)

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			// send back the message to everyone
			ctx.SendChat(Message{Name: member.Name, Content: string(message)})
			split := strings.Fields(string(message))

			// handle commands if any
			if command, ok := Commands[split[0]]; ok {
				command(ctx,
					member,
					split[1:],
				)
			} else if message[0] == '/' {
				ctx.SendChat(Message{Name: "server", Content: "command not found"})
			}
		}
	}
}

func Serve(db *gorm.DB, host string) {
	ctx := &ServerContext{
		Db:   db,
		Grid: &grid.Grid{Paused: true},
	}

	defer ctx.Close()

	addr := flag.String("addr", host, "http service address")

	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/state", HandleBoardState(ctx))

	http.HandleFunc("/chat", HandleChat(ctx))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
