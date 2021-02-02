package websocket

import (
	"chaoticneutraltech/codenames/codenames"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type chat struct {
	MsgType string `json:"MsgType"`
	Data    string `json:"Data"`
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader_test = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func reader(conn *websocket.Conn) {
	for {
		message := &chat{}

		// read in a message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		log.Println(message)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader_test.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	game := codenames.CreateGame(codenames.GenerateBoard(codenames.TeamRed), codenames.TeamRed)

	g, err := json.Marshal(game)

	err = ws.WriteMessage(1, []byte(g))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader_test.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	board := codenames.GenerateBoard(codenames.TeamRed)

	game := codenames.CreateGame(board, codenames.TeamRed)

	b, err := json.Marshal(game.Board)

	g, err := json.Marshal(game)

	fmt.Print(g)

	err = ws.WriteMessage(1, []byte(b))
}

func SetupRoutes() *mux.Router {
	hub := newHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws", wsEndpoint)
	r.HandleFunc("/ws/{game}", wsEndpoint)
	r.HandleFunc("/ws2", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	r.HandleFunc("/ws2/{game}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	return r
}
