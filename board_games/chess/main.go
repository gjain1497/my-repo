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
	Games                 map[string]*Game
	MoveValidationService MoveValidationService
	MoveExecutionService  MoveExecutionService
	BoardService          BoardService
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
	King   PieceType = "KING"
	Queen  PieceType = "QUEEN"
	Rook   PieceType = "ROOK"
	Bishop PieceType = "BISHOP"
	Knight PieceType = "KNIGHT"
	Pawn   PieceType = "PAWN"
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

type MoveValidationService interface {
	ValidateMove(move *Move, board *Board) (bool, error)
	GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error)
}

type MoveValidationServiceV1 struct {
	//PieceMovement
	PieceMovementFactory *PieceMovementFactory
}

func (m *MoveValidationServiceV1) ValidateMove(move *Move, board *Board) (bool, error) {

	movement, err := m.PieceMovementFactory.GetPieceMovement(move.Piece.Type)
	if err != nil {
		return false, err
	}
	isValid, err := movement.IsValidMove(move.From, move.To, move.Piece, board)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	return true, nil
}

func (m *MoveValidationServiceV1) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {
	movement, err := m.PieceMovementFactory.GetPieceMovement(piece.Type)
	if err != nil {
		return nil, err
	}
	postitions, err := movement.GetValidMovesForPiece(piece, board)
	if err != nil {
		return nil, err
	}
	return postitions, nil
}

type PieceMovementFactory struct {
	Movements map[PieceType]PieceMovement
}

func NewPieceMovementFactory() *PieceMovementFactory {
	return &PieceMovementFactory{
		Movements: map[PieceType]PieceMovement{
			King:   &KingMovement{},
			Queen:  &QueenMovement{},
			Rook:   &RookMovement{},
			Bishop: &BishopMovement{},
			Pawn:   &PawnMovement{},
			Knight: &KnightMovement{},
		},
	}
}

func (f *PieceMovementFactory) GetPieceMovement(pieceType PieceType) (PieceMovement, error) {
	pieceMovement, exists := f.Movements[pieceType]
	if !exists {
		return nil, fmt.Errorf("processor %s not found", pieceType)
	}
	return pieceMovement, nil
}

type PieceMovement interface {
	IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error)
	GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error)
}

type KingMovement struct {
}

func (k *KingMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}

func (k *KingMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type BishopMovement struct {
}

func (b *BishopMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}

func (b *BishopMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type QueenMovement struct {
}

func (q *QueenMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}
func (q *QueenMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type RookMovement struct {
}

func (r *RookMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}

func (r *RookMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type PawnMovement struct {
}

func (p *PawnMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}

func (p *PawnMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type KnightMovement struct {
}

func (k *KnightMovement) IsValidMove(from, to Position, piece *Piece, board *Board) (bool, error) {

}

func (k *KnightMovement) GetValidMovesForPiece(piece *Piece, board *Board) ([]Position, error) {

}

type MoveExecutionService interface {
	ExecuteMove(game *Game, move *Move) error
}

type MoveExecutionServiceV1 struct {
	BoardService BoardService
}

func (s *MoveExecutionServiceV1) ExecuteMove(game *Game, move *Move) error {

}

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
