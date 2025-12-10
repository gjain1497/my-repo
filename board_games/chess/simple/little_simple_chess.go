package chess

import (
	"fmt"
	"time"
)

//8*8 board
//6 piece type king, queen, rrok, bishop, knight, pawn
//2 players (black, white)
//winner -> checkmate opponent's king

//move validation -> different for different pieces
//check -> king is under attack
//checkmate -> king is under attack and cannot escape
//castling -> special king + rook move
//promotion - pawn reached end -> becomes queen/rook/bishop/knight

type GameService interface {
	CreateGame(player1, player2 *Player) (*Game, error)
	Start(gameId string) error
	GetGame(gameId string) (*Game, error)
	MakeMove(gameId string, move Move) error
	GetCurrentPlayer(gameId string) (*Player, error)
	GetState(gameId string) (GameState, error)
	IsGameOver(gameId string) (bool, error)
}

type GameServiceV1 struct {
	Games        map[string]*Game
	PieceFactory *PieceFactory
	// MoveExecutionService MoveExecutionService
	BoardService BoardService
}
type PieceFactory struct {
	Pieces map[PieceType]PieceService
}

func NewPiece() *PieceFactory {
	return &PieceFactory{
		Pieces: map[PieceType]PieceService{
			KING:   &King{},
			QUEEN:  &Queen{},
			ROOK:   &Rook{},
			BISHOP: &Bishop{},
			PAWN:   &Pawn{},
			KNIGHT: &Knight{},
		},
	}
}

func (f *PieceFactory) GetPieceValidationService(pieceType PieceType) (PieceService, error) { //kingValidation queenValidation etc
	pieceValidationService, exists := f.Pieces[pieceType]
	if !exists {
		return nil, fmt.Errorf("pieceService %s not found", pieceType)
	}
	return pieceValidationService, nil
}

type Game struct {
	ID            string
	Board         *Board
	Players       [2]*Player
	CurrentPlayer int
	State         GameState
	CreatedAt     time.Time
	Winner        *Player
	MoveHistory   []Move
	IsCheck       bool
}

type GameState string

const (
	NotStarted GameState = "NOT_STARTED"
	InProgress GameState = "IN_PROGRESS"
	Checkmate  GameState = "CHECKMATE"
	Stalemate  GameState = "STALEMATE"
	Draw       GameState = "DRAW"
)

type BoardService interface {
	InitializeBoard() *Board
	GetPieceAt(board *Board, pos Position) *Piece
	SetPieceAt(board *Board, pos Position, piece *Piece) error
	RemovePieceAt(board *Board, pos Position) error
	IsPositionEmpty(board *Board, pos Position) bool
	IsPositionValid(pos Position) bool
}
type BoardServiceV1 struct {
}

type Board struct {
	Cells [8][8]*Piece
}

type Piece struct {
	Color    Color
	Position Position
	Type     PieceType
	HasMoved bool //not sure why this is required
}

type PieceType string

const (
	KING   PieceType = "KING"
	QUEEN  PieceType = "QUEEN"
	ROOK   PieceType = "ROOK"
	BISHOP PieceType = "BISHOP"
	KNIGHT PieceType = "KNIGHT"
	PAWN   PieceType = "PAWN"
)

type Color string

const (
	White Color = "WHITE"
	Black Color = "BLACK"
)

type Player struct {
	Id    string
	Name  string
	Color Color
}

type PieceService interface {
	ValidateMove(move *Move, board *Board) (bool, error)
	GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error)
	// CLAUDE, cmiiw, As ExecuteMove will be common for all pieces. So if keep in this interface
	// every piece has to implement it which is  not ideal
	// So here we can apply ISP and break this interface
	// where ExecuteMove can be part of a seperate interface
	ExecuteMove(game *Game, move *Move) error
}

type King struct {
}

func (k *King) ValidateMove(move *Move, board *Board) (bool, error) {

}

func (k *King) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

func (k *King) ExecuteMove(game *Game, move *Move) error {

}

type Bishop struct {
}

func (b *Bishop) ValidateMove(move *Move, board *Board) (bool, error) {

}

func (b *Bishop) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}
func (b *Bishop) ExecuteMove(game *Game, move *Move) error {

}

type Queen struct {
}

func (q *Queen) ValidateMove(move *Move, board *Board) (bool, error) {

}
func (q *Queen) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}
func (q *Queen) ExecuteMove(game *Game, move *Move) error {

}

type Rook struct {
}

func (r *Rook) ValidateMove(move *Move, board *Board) (bool, error) {

}

func (r *Rook) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}
func (r *Rook) ExecuteMove(game *Game, move *Move) error {

}

type Pawn struct {
}

func (p *Pawn) ValidateMove(move *Move, board *Board) (bool, error) {

}

func (p *Pawn) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}
func (p *Pawn) ExecuteMove(game *Game, move *Move) error {

}

type Knight struct {
}

func (k *Knight) ValidateMove(move *Move, board *Board) (bool, error) {

}

func (k *Knight) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}
func (k *Knight) ExecuteMove(game *Game, move *Move) error {

}

// type MoveExecutionService interface {
// 	ExecuteMove(game *Game, move *Move) error
// }

// type MoveExecutionServiceV1 struct {
// 	BoardService BoardService
// }

// func (s *MoveExecutionServiceV1) ExecuteMove(game *Game, move *Move) error {

// }

//Instead of doing this we will do factory
//because this is what piece service has to decide
//based on type of piece
// type PieceServiceV1 struct{
// 	Type PieceMoverType
// }

type Move struct {
	From          Position
	To            Position
	Piece         *Piece
	Player        *Player
	CapturedPiece *Piece
	IsSpecialMove bool
	PromotionTo   *PieceType
}

type Position struct {
	Row int
	Col int
}
