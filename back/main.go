package main

import (
	"gol_back/database"
	"gol_back/socket"
)

func main() {
	db := database.InitDb("game_of_life_users.db")

	server := &socket.WebSocketServer{}
	defer server.Close()

	socket.InitWebsocketServer(db, server)
}
