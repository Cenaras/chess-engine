package game

import "fmt"

type Piece uint8

const (
	NONE   Piece = 0  // 0b00000
	PAWN   Piece = 1  // 0b00001
	KNIGHT Piece = 2  // 0b00010
	BISHOP Piece = 3  // 0b00011
	ROOK   Piece = 4  // 0b00100
	QUEEN  Piece = 5  // 0b00101
	KING   Piece = 6  // 0b00110
	WHITE  Piece = 8  // 0b01000
	BLACK  Piece = 16 // 0b10000
)

// Use masks to extract piece type/piece color from a Piece
const (
	TypeMask  Piece = 0b00111
	ColorMask Piece = 0b11000
)

type Player uint8

func (p Piece) Player() Player {
	return Player(p & ColorMask)
}

func (p Piece) Type() Piece {
	return Piece(p & TypeMask)
}

type CastleRights uint8

const (
	WhiteKingSide  CastleRights = 0b0001
	WhiteQueenSide CastleRights = 0b0010
	BlackKingSide  CastleRights = 0b0100
	BlackQueenSide CastleRights = 0b1000
)

type Rank uint8

const (
	FILE_A Rank = iota + 1
	FILE_B
	FILE_C
	FILE_D
	FILE_E
	FILE_F
	FILE_G
	FILE_H
)

type File uint8

const (
	RANK_1 File = iota + 1
	RANK_2
	RANK_3
	RANK_4
	RANK_5
	RANK_6
	RANK_7
	RANK_8
)

type Square uint8 // 0..64
const NoSquare Square = 64

// TODO: Encapsulate to avoid exposing datatype?
type Position struct {
	Board                    [64]Piece
	PlayerToMove             Player
	CastleRights             CastleRights // 0001
	PossibleEnPassantCapture Square       // todo: representable using 4bits. Using the playerToMove to indicate the rank and just store the file
}

// func (p Position) GetPieceAt(square Square) Piece {
// 	return piece := p.Board[square]
// 	return piece.Type(), piece.Player()
// }

// PieceType, Color
func (p Position) GetPieceAt(square Square) (Piece, Player) {
	piece := p.Board[square]
	return piece.Type(), piece.Player()
}

func (p Position) SetPieceAt(piece Piece, square Square) {
	p.Board[square] = piece
}

type MoveFlag uint8

const (
	NormalMove MoveFlag = iota
	DoublePawnPush
	KingCastle
	QueenCastle
	EnPassantCapture
	PromoteKnight
	PromoteBishop
	PromoteRook
	PromoteQueen
)

type Move struct {
	From Square
	To   Square
	Flag MoveFlag
}

type UndoMoveState struct {
	CapturedPiece            Piece
	CastleRights             CastleRights
	PossibleEnPassantCapture Square
}

type Direction struct {
	Rank int
	File int
}

func RankFileToSquare(rank int, file int) Square {
	return Square(rank*8 + file)
}

func SquareToRankFile(square Square) (int, int) {
	return int(square) / 8, int(square) % 8
}

func IsLegalRank(rank int) bool {
	return rank >= 0 && rank < 8
}

func IsLegalFile(file int) bool {
	return file >= 0 && file < 8
}

func IsStartPawnRank(square Square, color Player) bool {
	if color == WHITE.Player() {
		return square/8 == 1
	}
	if color == BLACK.Player() {
		return square/8 == 6
	}
	panic(fmt.Sprintf("Invalid color for player: %d", color))
}

func MoveDirection(square Square, direction Direction) (Square, error) {
	startRank, startFile := SquareToRankFile(square)
	newRank := startRank + direction.Rank
	newFile := startFile + direction.File

	// Illegal direction for this piece
	if !IsLegalRank(newRank) || !IsLegalFile(newFile) {
		return 0, fmt.Errorf(
			"moving square %d by direction %+v goes off the board",
			square,
			direction)
	}
	// Target square is within the board
	return RankFileToSquare(newRank, newFile), nil
}

func IsSameColor(c1 Player, c2 Player) bool {
	return c1 == c2
}
