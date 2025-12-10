package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// ============================================
// MODELS
// ============================================

type CellType string

const (
	NormalCell  CellType = "NORMAL"
	SnakeHead   CellType = "SNAKE_HEAD"
	LadderStart CellType = "LADDER_START"
)

type Cell struct {
	Position    int
	Type        CellType
	Destination *int
}

type Board struct {
	Cells [101]*Cell
	Size  int
}

type Player struct {
	Id       string
	Name     string
	Position int
}

type Game struct {
	Id              string
	Board           *Board
	Players         []*Player
	CurrentPlayer   int
	NumberOfPlayers int
	State           GameState
	CreatedAt       time.Time
	Winner          *Player
}

type GameState string

const (
	NotStarted GameState = "NOT_STARTED"
	InProgress GameState = "IN_PROGRESS"
	Won        GameState = "WON"
)

// ============================================
// SERVICES
// ============================================

// --- BoardService --- ✅ NEW!
type BoardService interface {
	CreateBoard(snakes map[int]int, ladders map[int]int) *Board
}

type BoardServiceV1 struct{}

func (s *BoardServiceV1) CreateBoard(snakes map[int]int, ladders map[int]int) *Board {
	board := &Board{
		Cells: [101]*Cell{},
		Size:  100,
	}

	// Initialize all cells as normal
	for i := 1; i <= 100; i++ {
		board.Cells[i] = &Cell{
			Position:    i,
			Type:        NormalCell,
			Destination: nil,
		}
	}

	// Set snakes
	for head, tail := range snakes {
		board.Cells[head].Type = SnakeHead
		board.Cells[head].Destination = &tail
	}

	// Set ladders
	for bottom, top := range ladders {
		board.Cells[bottom].Type = LadderStart
		board.Cells[bottom].Destination = &top
	}

	return board
}

// --- GameService ---
type GameService interface {
	CreateGame(snakes map[int]int, ladders map[int]int, numPlayers int, players []*Player) (*Game, error)
	StartGame(gameId string) (bool, error)
	MakeMove(gameId string, diceRollResult int) error
	GetGameState(gameId string) (GameState, error)
	GetCurrentPlayer(gameId string) (*Player, error)
	GetBoard(gameId string) (*Board, error)
	GetGame(gameId string) (*Game, error)
	IsGameOver(gameId string) (bool, error)
	RollDice() int
}

type GameServiceV1 struct {
	Games              map[string]*Game
	BoardService       BoardService // ✅ Injected!
	MoveHandlerService MoveHandlerService
	DiceService        DiceService
	mu                 sync.RWMutex
}

func NewGameServiceV1(boardService BoardService, moveHandlerService MoveHandlerService, diceService DiceService) *GameServiceV1 {
	return &GameServiceV1{
		Games:              make(map[string]*Game),
		BoardService:       boardService, // ✅ Injected!
		MoveHandlerService: moveHandlerService,
		DiceService:        diceService,
	}
}

func (s *GameServiceV1) CreateGame(snakes map[int]int, ladders map[int]int, numPlayers int, players []*Player) (*Game, error) {
	log.Println("Creating game with players:", players)

	// ✅ Use BoardService to create board
	board := s.BoardService.CreateBoard(snakes, ladders)

	game := &Game{
		Id:              generateGameID(),
		Board:           board,
		Players:         players,
		CurrentPlayer:   0,
		NumberOfPlayers: numPlayers,
		State:           NotStarted,
		CreatedAt:       time.Now(),
	}

	s.mu.Lock()
	s.Games[game.Id] = game
	s.mu.Unlock()

	return game, nil
}

func (s *GameServiceV1) StartGame(gameId string) (bool, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return false, errors.New("game not found")
	}

	game.State = InProgress
	return true, nil
}

func (s *GameServiceV1) MakeMove(gameId string, diceRollResult int) error {
	game, err := s.GetGame(gameId)
	if err != nil {
		return err
	}

	board := game.Board
	currPlayer := game.Players[game.CurrentPlayer]

	log.Printf("Current player: %s at position %d", currPlayer.Name, currPlayer.Position)

	newPos := currPlayer.Position + diceRollResult
	log.Printf("After dice roll: %d", newPos)

	if newPos >= board.Size {
		game.State = Won
		game.Winner = currPlayer
		currPlayer.Position = board.Size
		log.Printf("🎉 %s wins!", currPlayer.Name)
		return nil
	}

	finalPos, err := s.MoveHandlerService.HandleMove(newPos, board)
	if err != nil {
		return err
	}

	currPlayer.Position = finalPos
	log.Printf("%s now at position %d\n", currPlayer.Name, finalPos)

	game.CurrentPlayer = (game.CurrentPlayer + 1) % game.NumberOfPlayers

	time.Sleep(500 * time.Millisecond)
	return nil
}

func (s *GameServiceV1) GetGameState(gameId string) (GameState, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return "", err
	}
	return game.State, nil
}

func (s *GameServiceV1) GetCurrentPlayer(gameId string) (*Player, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, err
	}
	return game.Players[game.CurrentPlayer], nil
}

func (s *GameServiceV1) GetBoard(gameId string) (*Board, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, err
	}
	return game.Board, nil
}

func (s *GameServiceV1) GetGame(gameId string) (*Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, ok := s.Games[gameId]
	if !ok {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func (s *GameServiceV1) IsGameOver(gameId string) (bool, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return false, err
	}
	return game.State == Won, nil
}

func (s *GameServiceV1) RollDice() int {
	return s.DiceService.Roll()
}

// --- DiceService ---
type DiceService interface {
	Roll() int
}

type DiceServiceV1 struct{}

func (s *DiceServiceV1) Roll() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	dice := r.Intn(6) + 1
	fmt.Printf("🎲 Dice roll: %d\n", dice)
	return dice
}

// --- MoveHandlerService ---
type MoveHandlerService interface {
	HandleMove(position int, board *Board) (int, error)
}

type MoveHandlerServiceV1 struct{}

func (s *MoveHandlerServiceV1) HandleMove(position int, board *Board) (int, error) {
	if position < 1 || position > board.Size {
		return position, errors.New("invalid position")
	}

	newPos := position

	for {
		cell := board.Cells[newPos]

		if cell.Destination != nil {
			if cell.Type == SnakeHead {
				fmt.Printf("🐍 Snake at %d! Sliding down to %d\n", newPos, *cell.Destination)
			} else if cell.Type == LadderStart {
				fmt.Printf("🪜 Ladder at %d! Climbing up to %d\n", newPos, *cell.Destination)
			}
			newPos = *cell.Destination
		} else {
			break
		}
	}

	return newPos, nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func generateGameID() string {
	return fmt.Sprintf("game_%d", time.Now().UnixNano())
}

// ============================================
// MAIN
// ============================================

func main() {
	// ✅ Initialize ALL services (including BoardService)
	boardService := &BoardServiceV1{}
	moveHandlerService := &MoveHandlerServiceV1{}
	diceService := &DiceServiceV1{}

	// ✅ Inject BoardService into GameService
	gameService := NewGameServiceV1(boardService, moveHandlerService, diceService)

	snakes := map[int]int{
		99: 1,
		87: 2,
		43: 27,
		78: 32,
		22: 11,
	}

	ladders := map[int]int{
		3:  97,
		18: 29,
		27: 81,
		21: 31,
		13: 49,
	}

	players := []*Player{
		{Id: "1", Name: "Girish", Position: 0},
		{Id: "2", Name: "Nipun", Position: 0},
		{Id: "3", Name: "Akash", Position: 0},
		{Id: "4", Name: "Ramya", Position: 0},
		{Id: "5", Name: "Vaibhav", Position: 0},
	}

	numPlayers := len(players)

	game, err := gameService.CreateGame(snakes, ladders, numPlayers, players)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}

	_, err = gameService.StartGame(game.Id)
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}

	fmt.Println("\n Game Started!")

	for {
		currPlayer, err := gameService.GetCurrentPlayer(game.Id)
		if err != nil {
			log.Printf("Error getting current player: %v", err)
			break
		}

		fmt.Printf("\n--- %s's turn ---\n", currPlayer.Name)

		dice := gameService.RollDice()

		err = gameService.MakeMove(game.Id, dice)
		if err != nil {
			log.Printf("Error making move: %v", err)
		}

		gameOver, _ := gameService.IsGameOver(game.Id)
		if gameOver {
			state, _ := gameService.GetGameState(game.Id)
			if state == Won {
				fmt.Printf("\n🏆 %s wins! 🏆\n", currPlayer.Name)
			}
			break
		}
	}
}
