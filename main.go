package main

import (
	"chaoticneutraltech/codenames/codenames"
	"chaoticneutraltech/codenames/websocket"
	"fmt"
	"log"
	"net/http"
)

func main() {

	board := codenames.GenerateBoard(codenames.TeamRed)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			fmt.Println(board.Codenames[i][j].Word)
			fmt.Println(board.Codenames[i][j].Agent)
			fmt.Println(board.Codenames[i][j].Position)
		}
	}

	//codenames.PublishGameState()

	r := websocket.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
