package codenames

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

type Team int

const (
	TeamRed Team = iota
	TeamBlue
)

type Agent int

const (
	None Agent = iota
	Red
	Blue
	DoubleAgent
	Bystander
	Assassin
)

type Board struct {
	Codenames [5][5]Codename
}

type Codename struct {
	Word     string
	Agent    string
	Position [2]int
}

type Game struct {
	Board        *Board
	StartingTeam Team
}

func (a Agent) String() string {
	return [...]string{"None", "Red", "Blue", "DoubleAgent", "Bystander", "Assassin"}[a]
}

func (t Team) String() string {
	return [...]string{"Red", "Blue"}[t]
}

func CreateGame(board *Board, startingTeam Team) *Game {
	game := Game{
		Board:        board,
		StartingTeam: startingTeam,
	}

	return &game
}

func GenerateBoard(startingTeam Team) *Board {

	var codenames [5][5]Codename
	var board Board
	var agents [25]Agent = [25]Agent{
		Red, Red, Red, Red, Red, Red, Red, Red,
		Blue, Blue, Blue, Blue, Blue, Blue, Blue, Blue,
		DoubleAgent,
		Bystander, Bystander, Bystander, Bystander, Bystander, Bystander, Bystander,
		Assassin}

	f, err := os.Open("codenames/codenames.csv")
	if err != nil {
		log.Fatal("Unable to read input file codenames.csv", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	r, err := csvReader.Read()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] }) //Shuffle input array with Fisher-Yates function
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(agents), func(i, j int) { agents[i], agents[j] = agents[j], agents[i] }) //Shuffle agents array with Fisher-Yates function

	counter := 0

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if agents[counter].String() == "DoubleAgent" {
				codenames[i][j] = Codename{Word: r[counter], Agent: startingTeam.String(), Position: [2]int{i, j}}
			} else {
				codenames[i][j] = Codename{Word: r[counter], Agent: agents[counter].String(), Position: [2]int{i, j}}
			}
			counter++
		}
	}

	if err != nil {
		log.Fatal("Unable to parse file as CSV for codenames.csv", err)
	}

	board.Codenames = codenames

	return &board
}

func PublishGameState() {

	board := GenerateBoard(TeamRed)

	b, err := json.Marshal(board.Codenames)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(b))

	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"codenames", // name
		"fanout",    // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	//body := "test"
	err = ch.Publish(
		"codenames", // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(b),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", b)

}

func PublishGameStateMQTT() {

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
