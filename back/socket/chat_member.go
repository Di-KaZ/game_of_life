package socket

import "github.com/gorilla/websocket"

type ChatMember struct {
	Sub  *websocket.Conn
	Name string
}
