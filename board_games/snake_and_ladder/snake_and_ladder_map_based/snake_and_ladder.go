package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//snake and ladder

type GameService interface {
	CreateGame(snakes map[int]int, ladders map[int]int) (*Game, error)
	StartGame(gameId string) (bool, error)
	MakeMove(gameId string, diceRollResult int) error
	GetGameState(gameId string) (GameState, error)
	GetCurrentPlayer(gameId string) (*Player, error)
	GetBoard(gameId string) (*Board, error)
	GetGame(gameId string) (*Game, error)
	IsGameOver(gameId string) (bool, error)
	RollDice()
}

type GameServiceV1 struct {
	Games map[string]*Game
	//MoveValidatorService MoveValidatorService no validation required here as in tic tac toe we validate BEFORE making the move:
	// - Is position within bounds?
	// - Is cell empty?
	// - Is it player's turn?
	// User does NOT provide position!
	// User just rolls dice → System calculates position

	//WinCheckerService    WinCheckerService //dont need this as its very simple

	//In tic tac toe we didn't require seperate movehandler as it was very simple over there
	// 2. Make the move (ONE LINE!) -> board.Cells[move.Position.Row][move.Position.Col] = move.Player.Symbol
	//but here ladder and snake complex move logic
	MoveHandlerService MoveHandlerService
	DiceService        DiceService
	mu                 sync.RWMutex
}

func NewGameServiceV1(moveHandlerService MoveHandlerService, diceService DiceService) *GameServiceV1 {
	return &GameServiceV1{
		Games:              make(map[string]*Game),
		MoveHandlerService: moveHandlerService,
		DiceService:        diceService,
	}
}

func (s *GameServiceV1) CreateGame(snakes map[int]int, ladders map[int]int, numPlayers int, players []*Player) (*Game, error) {
	log.Println("players here: ", players)
	game := &Game{
		Id: "ds212132",
		Board: &Board{
			Snakes:  snakes,
			Ladders: ladders,
			Size:    100,
		},
		Players:         players,
		CurrentPlayer:   0,
		NumberOfPlayers: numPlayers,
		State:           NotStarted,
		CreatedAt:       time.Now(),
	}
	s.Games[game.Id] = game
	return game, nil
}
func (s *GameServiceV1) StartGame(gameId string) (bool, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return false, errors.New("Game not found to start")
	}
	game.State = InProgress
	return true, nil
}
func (s *GameServiceV1) MakeMove(gameId string, diceRollResult int) error {
	board, err := s.GetBoard(gameId)
	if err != nil {
		return errors.New("Board not found for this game")
	}

	game, err := s.GetGame(gameId)
	if err != nil {
		return errors.New("Game not found to start")
	}
	currPlayer, err := s.GetCurrentPlayer(gameId)
	if err != nil {
		return err
	}
	log.Println("currPlayer ", currPlayer)
	newPos := currPlayer.Position + diceRollResult
	log.Println("newPos ", newPos)
	if newPos >= 100 {
		game.State = Won
		game.Winner = currPlayer
		currPlayer.Position = 100
		return nil
	}

	newPos, err = s.MoveHandlerService.HandleMove(newPos, board)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	currPlayer.Position = newPos
	//switch player turn
	game.CurrentPlayer = (game.CurrentPlayer + 1) % game.NumberOfPlayers
	return nil
}
func (s *GameServiceV1) GetGameState(gameId string) (GameState, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return "", errors.New("Game not found to start")
	}
	return game.State, nil
}
func (s *GameServiceV1) GetCurrentPlayer(gameId string) (*Player, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, errors.New("Game not found to start")
	}
	fmt.Println("game is: ", game)
	fmt.Println("Index of current player ", game.CurrentPlayer)
	return game.Players[game.CurrentPlayer], nil
}
func (s *GameServiceV1) GetBoard(gameId string) (*Board, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, errors.New("Game not found to start")
	}
	return game.Board, nil
}
func (s *GameServiceV1) GetGame(gameId string) (*Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	game, ok := s.Games[gameId]
	if !ok {
		return nil, errors.New("Game not found")
	}
	return game, nil
}
func (s *GameServiceV1) IsGameOver(gameId string) (bool, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return false, errors.New("Game not found to start")
	}
	return game.State == Won, nil
}
func (s *GameServiceV1) RollDice() int {
	return s.DiceService.Roll()
}

type Game struct {
	Id              string
	Board           *Board
	Players         []*Player //n players
	CurrentPlayer   int
	NumberOfPlayers int
	State           GameState
	CreatedAt       time.Time
	Winner          *Player
}

// Game state - current state of game
type GameState string

const (
	NotStarted GameState = "NOT_STARTED"
	InProgress GameState = "IN_PROGRESS"
	Won        GameState = "WON"
)

type Snake struct {
	Head int //snake starts(higher number)
	Tail int //snake ends(lower number)
}

type Ladder struct {
	Bottom int //ladder starts (lower number)
	Top    int //ladder ends (higher number)
}

type Board struct {
	//Snakes  []Snake -> can be optimised to use map instead
	Snakes  map[int]int //map[head]tail -> map[98]7
	Ladders map[int]int //map[bottom]top -> map[4]6
	Size    int
}

type Player struct {
	Id       string
	Name     string
	Position int
}

type DiceService interface {
	Roll() int //return 1-6
}

type DiceServiceV1 struct {
}

func (s *DiceServiceV1) Roll() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dice := r.Intn(6) + 1
	fmt.Println("Dice roll:", dice)
	return dice
}

type MoveHandlerService interface {
	HandleMove(position int, board *Board) (int, error) //return 1-6
}

type MoveHandlerServiceV1 struct {
}

func (s *MoveHandlerServiceV1) HandleMove(position int, board *Board) (int, error) {
	snakes := board.Snakes
	ladders := board.Ladders
	newPos := position

	for {
		moved := false
		snakeTail, ok := snakes[newPos] //check if this position has snake top
		if ok {
			fmt.Printf("Snake at %v! Sliding down to %v\n", newPos, snakeTail)
			newPos = snakeTail
			moved = true
		}
		ladderTop, ok := ladders[newPos] //check if this position has ladder bottom
		if ok {
			fmt.Printf("Ladder at %v! Climbing up to %v\n", newPos, ladderTop)
			newPos = ladderTop
			moved = true
		}
		if !moved {
			break
		}
	}

	return newPos, nil
}

func main() {
	//create game
	//get current player
	//roll dice -> will give number
	//make move with that number
	//check isGameOVer
	//if won print player details
	moveHandlerServiceV1 := &MoveHandlerServiceV1{}
	dicServiceV1 := &DiceServiceV1{}

	gameService := NewGameServiceV1(moveHandlerServiceV1, dicServiceV1)

	snake1 := &Snake{
		Head: 99,
		Tail: 1,
	}
	snake2 := &Snake{
		Head: 87,
		Tail: 2,
	}
	snake3 := &Snake{
		Head: 43,
		Tail: 27,
	}
	snake4 := &Snake{
		Head: 78,
		Tail: 32,
	}
	snake5 := &Snake{
		Head: 22,
		Tail: 11,
	}

	ladder1 := &Ladder{
		Bottom: 3,
		Top:    97,
	}
	ladder2 := &Ladder{
		Bottom: 18,
		Top:    29,
	}
	ladder3 := &Ladder{
		Bottom: 27,
		Top:    81,
	}
	ladder4 := &Ladder{
		Bottom: 21,
		Top:    31,
	}
	ladder5 := &Ladder{
		Bottom: 13,
		Top:    49,
	}
	ladders := make(map[int]int)
	ladders[ladder1.Bottom] = ladder1.Top
	ladders[ladder2.Bottom] = ladder2.Top
	ladders[ladder3.Bottom] = ladder3.Top
	ladders[ladder4.Bottom] = ladder4.Top
	ladders[ladder5.Bottom] = ladder5.Top

	snakes := make(map[int]int)
	snakes[snake1.Head] = snake1.Tail
	snakes[snake2.Head] = snake2.Tail
	snakes[snake3.Head] = snake3.Tail
	snakes[snake4.Head] = snake4.Tail
	snakes[snake5.Head] = snake5.Tail

	numPlayers := 5

	player1 := &Player{
		Id:   "0",
		Name: "Girish",
	}

	player2 := &Player{
		Id:   "1",
		Name: "Nipun",
	}
	player3 := &Player{
		Id:   "2",
		Name: "Akash",
	}
	player4 := &Player{
		Id:   "3",
		Name: "Ramya",
	}
	player5 := &Player{
		Id:   "4",
		Name: "Vaibhav",
	}

	players := []*Player{player1, player2, player3, player4, player5}
	log.Printf("playere before: %v", players)

	game, err := gameService.CreateGame(snakes, ladders, numPlayers, players)
	if err != nil {
		fmt.Printf("Not able to create the game %v ", err.Error())
	}

	for {
		currPlayer, err := gameService.GetCurrentPlayer(game.Id)
		if err != nil {
			fmt.Printf("Not able to find the current player %v ", err.Error())
		}

		dice := gameService.RollDice()
		log.Println("dice ", dice)

		err = gameService.MakeMove(game.Id, dice)
		if err != nil {
			fmt.Printf("Not able to make the move %v ", err.Error())
		}
		gameOver, _ := gameService.IsGameOver(game.Id)
		if gameOver {
			state, _ := gameService.GetGameState(game.Id)
			if state == Won {
				fmt.Printf("Player %s (%s) wins!\n", currPlayer.Name, currPlayer.Id)
			}
			break
		}
	}

}
