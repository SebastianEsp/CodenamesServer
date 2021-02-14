package main

import (
	"chaoticneutraltech/codenames/websocket"
	"log"
	"net/http"
)

func main() {
	r := websocket.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", r))
}
