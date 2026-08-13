package game

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

type Square int // 0..64

type Position struct {
	Board                    [64]Piece
	PlayerToMove             Player
	CastleRights             CastleRights
	PossibleEnPassantCapture Square // todo: representable using 4bits. Using the playerToMove to indicate the rank and just store the file
}

// PieceType, Color
func (p Position) GetPieceAt(square Square) (Piece, Player) {
	piece := p.Board[square]
	return piece.Type(), piece.Player()
}

type Move struct {
	From Square
	To   Square
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

func IsSameColor(c1 Player, c2 Player) bool {
	return c1 == c2
}
