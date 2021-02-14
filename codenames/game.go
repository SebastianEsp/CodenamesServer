package game

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"
)

type Team int

const (
	None Team = iota
	TeamRed
	TeamBlue
)

type Agent int

const (
	Red Agent = iota
	Blue
	DoubleAgent
	Bystander
	Assassin
)

type Board struct {
	Codenames [5][5]Codename
}

type Codename struct {
	Word  string
	Agent string
	PosX  int
	PosY  int
	State CodenameState
}

type CodenameState struct {
	Guessed   bool
	GuessedBy string
}

type Game struct {
	Board        *Board
	StartingTeam string
	GameState    GameState
}

type GameState struct {
	RedPoints   int
	BluePoints  int
	CurrentTeam string
	Ended       bool
	Winner      string
}

func (a Agent) String() string {
	return [...]string{"Red", "Blue", "DoubleAgent", "Bystander", "Assassin"}[a]
}

func (t Team) String() string {
	return [...]string{"None", "Red", "Blue"}[t]
}

func NewGame(startingTeam Team) *Game {
	game := new(Game)

	game.Board = game.GenerateBoard(startingTeam)
	game.StartingTeam = startingTeam.String()
	game.GameState = GameState{0, 0, startingTeam.String(), false, None.String()}

	return game
}

func (game Game) GenerateBoard(startingTeam Team) *Board {

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
				codenames[i][j] = Codename{Word: r[counter], Agent: startingTeam.String(), PosX: i, PosY: j, State: CodenameState{Guessed: false, GuessedBy: None.String()}}
			} else {
				codenames[i][j] = Codename{Word: r[counter], Agent: agents[counter].String(), PosX: i, PosY: j, State: CodenameState{Guessed: false, GuessedBy: None.String()}}
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
