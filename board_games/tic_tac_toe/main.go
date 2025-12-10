//to do
//undo move, game history
//observer pattern
//optimise check win

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//What Makes Something a Model?
//1. Do we need to store that data
//2. Do we need to represent what happened

// Game Service/Component
type GameService interface {
	CreateGame(player1, player2 *Player) (*Game, error)
	Start(gameId string) error
	MakeMove(gameId string, move Move) error
	GetGame(gameId string) (*Game, error)
	GetBoard(gameId string) (*Board, error)
	GetCurrentPlayer(gameId string) (*Player, error)
	GetState(gameId string) (GameState, error)
	IsGameOver(gameId string) (bool, error)
}

type GameServiceV1 struct {
	Games         map[string]*Game //(game_id, game object)
	MoveValidator MoveValidatorService
	WinChecker    WinCheckerService
	mu            sync.RWMutex
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

func (s *GameServiceV1) CreateGame(player1, player2 *Player) (*Game, error) {
	game := &Game{
		Id:            "312321",
		Board:         &Board{Cells: [3][3]Symbol{}, Size: 3},
		Players:       [2]*Player{player1, player2},
		CurrentPlayer: 0,
		State:         NotStarted,
		CreatedAt:     time.Now(),
	}

	s.Games[game.Id] = game
	return game, nil
}
func (s *GameServiceV1) Start(gameId string) error {
	game, err := s.GetGame(gameId)
	if err != nil {
		return err
	}
	game.State = InProgress
	return nil
}

func (s *GameServiceV1) IsGameOver(gameId string) (bool, error) {
	gameState, err := s.GetState(gameId)
	if err != nil {
		return false, nil
	}
	if gameState == Won || gameState == Draw {
		return true, nil
	}
	return false, nil
}

func (s *GameServiceV1) MakeMove(gameId string, move Move) error {
	board, err := s.GetBoard(gameId)
	if err != nil {
		return err
	}
	game, err := s.GetGame(gameId)
	if err != nil {
		return err
	}
	// Check if game is in progress
	if game.State != InProgress {
		return errors.New("game not in progress")
	}

	// Check if it's the correct player's turn
	currentPlayer, err := s.GetCurrentPlayer(gameId)
	if err != nil {
		return err
	}
	if move.Player != currentPlayer {
		return errors.New("not this player's turn")
	}
	isValidMove := s.MoveValidator.IsValid(board, move)
	if !isValidMove {
		return errors.New("Move is invalid")
	}

	//place symbol on the board
	board.Cells[move.Position.Row][move.Position.Col] = move.Player.Symbol

	//check win condtion
	hasWon := s.WinChecker.CheckWin(board, move.Player)
	if hasWon {
		game.State = Won
		game.Winner = move.Player
		return nil
	}

	//check draw condtion
	isDraw := s.isBoardFull(board)
	if isDraw {
		game.State = Draw
	}

	//switch player turn
	game.CurrentPlayer = (game.CurrentPlayer + 1) % 2
	return nil
}

func (s *GameServiceV1) isBoardFull(board *Board) bool {
	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			if board.Cells[i][j] == Empty {
				return false
			}
		}
	}
	return true
}

func (s *GameServiceV1) GetBoard(gameId string) (*Board, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, err
	}
	return game.Board, nil
}
func (s *GameServiceV1) GetCurrentPlayer(gameId string) (*Player, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return nil, err
	}
	return game.Players[game.CurrentPlayer], nil
}

func (s *GameServiceV1) GetState(gameId string) (GameState, error) {
	game, err := s.GetGame(gameId)
	if err != nil {
		return "", err
	}
	return game.State, nil
}

// move validator service/component
type MoveValidatorService interface {
	IsValid(board *Board, move Move) bool
}

type MoveValidatorServiceV1 struct {
}

func (m *MoveValidatorServiceV1) IsValid(board *Board, move Move) bool {
	// Check if row is within bounds
	if move.Position.Row < 0 || move.Position.Row >= board.Size {
		return false
	}

	// Check if col is within bounds
	if move.Position.Col < 0 || move.Position.Col >= board.Size {
		return false
	}

	// Check if cell is empty
	if board.Cells[move.Position.Row][move.Position.Col] != Empty {
		return false
	}

	return true
}

// winchecker interface
type WinCheckerService interface {
	CheckWin(board *Board, player *Player) bool
}

type RowWinChecker struct { //class
}

func (r *RowWinChecker) CheckWin(board *Board, player *Player) bool {
	for row := 0; row < board.Size; row++ {
		allMatch := true
		for col := 0; col < board.Size; col++ {
			if board.Cells[row][col] != player.Symbol {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	return false
}

type RowColumnnWinChecker struct { //class
}

func (r *RowColumnnWinChecker) CheckWin(board *Board, player *Player) bool {
	for row := 0; row < board.Size; row++ {
		allMatch := true
		for col := 0; col < board.Size; col++ {
			if board.Cells[row][col] != player.Symbol {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	return false
}

// Column Win Checker
type ColWinChecker struct{}

func (c *ColWinChecker) CheckWin(board *Board, player *Player) bool {
	for col := 0; col < board.Size; col++ {
		allMatch := true
		for row := 0; row < board.Size; row++ {
			if board.Cells[row][col] != player.Symbol {
				allMatch = false
				break
			}
		}
		if allMatch {
			return true
		}
	}
	return false
}

// Diagonal Win Checker
type DiagonalWinChecker struct{}

func (d *DiagonalWinChecker) CheckWin(board *Board, player *Player) bool { //optimised
	// Check main diagonal (top-left to bottom-right)
	allMatch := true
	for i := 0; i < board.Size; i++ {
		if board.Cells[i][i] != player.Symbol {
			allMatch = false
			break
		}
	}
	if allMatch {
		return true
	}

	// Check anti-diagonal (top-right to bottom-left)
	allMatch = true
	for i := 0; i < board.Size; i++ {
		if board.Cells[i][board.Size-1-i] != player.Symbol {
			allMatch = false
			break
		}
	}
	return allMatch
}

// Composite Win Checker (combines all strategies)
type CompositeWinChecker struct {
	Checkers []WinCheckerService
}

func (c *CompositeWinChecker) CheckWin(board *Board, player *Player) bool {
	for _, checker := range c.Checkers {
		if checker.CheckWin(board, player) {
			return true
		}
	}
	return false
}

// Move represents a player's action
type Move struct {
	Position Position
	Player   *Player
}

type Board struct {
	Cells [3][3]Symbol //3x3 grid
	Size  int          //3 (could make it NxN later)
}

type Game struct {
	Id            string
	Board         *Board
	Players       [2]*Player
	CurrentPlayer int
	State         GameState
	CreatedAt     time.Time
	Winner        *Player
}

// Game state - current state of game
type GameState string

const (
	NotStarted GameState = "NOT_STARTED"
	InProgress GameState = "IN_PROGRESS"
	Won        GameState = "WON"
	Draw       GameState = "DRAW"
)

// 2. Symbol - What goes in a cell (X, O, or Empty)
type Symbol string

const (
	Empty Symbol = ""
	X     Symbol = "X"
	O     Symbol = "O"
)

type Player struct {
	Id     string
	Name   string
	Symbol Symbol
}

type Position struct { //We need to represent "what happened"
	Row int
	Col int
}

func main() {
	// Create validator
	var validator MoveValidatorService = &MoveValidatorServiceV1{}

	// Create win checker with composite pattern (Strategy Pattern!)
	var winChecker WinCheckerService = &CompositeWinChecker{
		Checkers: []WinCheckerService{
			&RowWinChecker{},
			&ColWinChecker{},
			&DiagonalWinChecker{},
		},
	}

	// Create game service
	var gameService GameService = &GameServiceV1{
		Games:         make(map[string]*Game),
		MoveValidator: validator,
		WinChecker:    winChecker,
	}

	player1 := &Player{
		Id:     "P1",
		Name:   "Girish",
		Symbol: X,
	}

	player2 := &Player{
		Id:     "P2",
		Name:   "Rohit",
		Symbol: O,
	}
	game, err := gameService.CreateGame(player1, player2)
	if err != nil {
		fmt.Printf("err: %v ", err)
	}
	gameService.Start(game.Id)

	// Hardcoded sequence of moves for demo
	moves := []Position{
		{0, 0}, // Girish
		{1, 1}, // Rohit
		{0, 1}, // Girish
		{2, 2}, // Rohit
		{0, 2}, // Girish -> wins row 0
	}
	moveIndex := 0

	for {
		currentPlayer, _ := gameService.GetCurrentPlayer(game.Id)
		pos := moves[moveIndex]
		moveIndex++

		move := Move{
			Position: pos,
			Player:   currentPlayer,
		}

		err := gameService.MakeMove(game.Id, move)
		if err != nil {
			fmt.Println("Error:", err)
		}

		printBoard(game.Board)
		gameOver, _ := gameService.IsGameOver(game.Id)
		if gameOver {
			state, _ := gameService.GetState(game.Id)
			if state == Won {
				fmt.Printf("Player %s (%s) wins!\n", currentPlayer.Name, currentPlayer.Symbol)
			} else {
				fmt.Println("It's a draw!")
			}
			break
		}
	}
}

func printBoard(board *Board) {
	fmt.Println("\nCurrent Board:")
	for i := 0; i < board.Size; i++ {
		for j := 0; j < board.Size; j++ {
			cell := board.Cells[i][j]
			if cell == Empty {
				fmt.Print("_ ")
			} else {
				fmt.Printf("%s ", cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
