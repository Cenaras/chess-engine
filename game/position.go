package game

import "fmt"

type Piece uint8

func (p Piece) Type() Piece {
	return p & TypeMask
}

func (p Piece) Player() Player {
	return Player(p & ColorMask)
}

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

func PlayerToString(player Player) string {
	if player == WHITE.Player() {
		return "White"
	}
	return "Black"
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
const NO_SQUARE Square = 64
const TOTAL_SQUARES Square = 64

// TODO: Encapsulate fields to avoid exposing datatype?
// TODO: Represent flags using more compact notation
type Position struct {
	board                    [64]Piece
	PlayerToMove             Player
	CastleRights             CastleRights
	PossibleEnPassantCapture Square
	WhiteKingSquare          Square
	BlackKingSquare          Square
	HalfMoveClock            int
}

func (p Player) Opponent() Player {
	if p == WHITE.Player() {
		return BLACK.Player()
	}
	return WHITE.Player()
}

func (p *Position) GetPieceAt(square Square) Piece {
	return p.board[square]
}

func (p *Position) SetPieceAt(piece Piece, square Square) {
	p.board[square] = piece
}

func (p *Position) FindKing(player Player) Square {
	if player == WHITE.Player() {
		return p.WhiteKingSquare
	}
	return p.BlackKingSquare
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

func PrintMove(move Move) {
	fmt.Printf(
		"%s -> %s\n",
		SquareToString(move.From),
		SquareToString(move.To),
	)
}
func SquareToString(square Square) string {
	file := byte('a' + square%8)
	rank := byte('1' + square/8)

	return string([]byte{file, rank})
}

type UndoMoveState struct {
	CapturedPiece            Piece
	CastleRights             CastleRights
	PossibleEnPassantCapture Square
	OldKingSquare            Square
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

func IsPromotionRank(square Square, color Player) bool {
	if color == WHITE.Player() {
		return square/8 == 7
	}
	return square/8 == 0
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

func MoveDirection(square Square, direction Direction) (Square, bool) {
	startRank, startFile := SquareToRankFile(square)
	newRank := startRank + direction.Rank
	newFile := startFile + direction.File

	// Illegal direction for this piece
	if !IsLegalRank(newRank) || !IsLegalFile(newFile) {
		return 0, false
	}
	// Target square is within the board
	return RankFileToSquare(newRank, newFile), true
}

func IsSameColor(c1 Player, c2 Player) bool {
	return c1 == c2
}
