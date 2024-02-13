package main

import (
	"gol_back/database"
	"gol_back/socket"
)

func main() {
	db := database.InitDb("game_of_life_users.db")

	socket.Serve(db, "localhost:8080")
}
