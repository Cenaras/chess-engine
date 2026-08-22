package game

import (
	"math/rand/v2"
)

// See:https://en.wikipedia.org/wiki/Zobrist_hashing

// indices for zobrist hashing
const (
	whitePawn int = iota
	whiteKnight
	whiteBishop
	whiteRook
	whiteQueen
	whiteKing
	blackPawn
	blackKnight
	blackBishop
	blackRook
	blackQueen
	blackKing
)

type zobristTable struct {
	Table         [64][12]uint64
	BlackToMove   uint64
	Castling      [4]uint64
	EnPassantFile [8]uint64
}

var Zobrist zobristTable = newZobristTable()

// Fills the zobrist table with random bitstrings for each piece
func newZobristTable() zobristTable {
	var zobrist zobristTable
	for idx := range TOTAL_SQUARES {
		for piece_idx := range 12 {
			zobrist.Table[idx][piece_idx] = genBitstring()
		}
	}
	for i := range 4 {
		zobrist.Castling[i] = genBitstring()
	}
	zobrist.BlackToMove = genBitstring()
	for file := range 8 {
		zobrist.EnPassantFile[file] = genBitstring()
	}
	return zobrist
}

// For initializing the boards zobrist hash. game.MakeMove is responsoble for maintaining the new hash
func SetupZobristHash(position *Position) uint64 {
	var hash uint64 = 0
	if position.PlayerToMove == BLACK.Player() {
		hash ^= Zobrist.BlackToMove
	}
	for square := range TOTAL_SQUARES {
		piece := position.GetPieceAt(square)
		pieceType := piece.Type()
		if pieceType == NONE {
			continue
		}
		idx := ZobristPieceIndex(pieceType, Piece(piece.Player()))
		hash ^= Zobrist.Table[square][idx]
	}

	rights := position.CastleRights
	for i := range 4 {
		// 0001, 0010, 0100, 1000 compactly iterated
		if rights&(1<<i) != 0 {
			hash ^= Zobrist.Castling[i]
		}
	}

	// FIXME: A position should only consider en-passant in the repetition,
	// if the capture is possible! This logic will technically cause our engine
	// to miss the fact that it could claim a draw (likely doesn't matter!)
	// Fix later with attack bitboards etc...
	if position.PossibleEnPassantCapture != NO_SQUARE {
		file := position.PossibleEnPassantCapture % 8
		hash ^= Zobrist.EnPassantFile[file]
	}

	return hash
}

// For debugging, use seeded rng
var rng = rand.New(rand.NewPCG(
	0x123456789abcdef0,
	0xfedcba9876543210,
))

// TODO: consider custom algorithm for performance
func genBitstring() uint64 {
	return rand.Uint64()
}

// Important: This assumes the order of the zobrist indices are consistent
// with the piece IDs!
func ZobristPieceIndex(piece Piece, player Piece) int {
	if piece < PAWN || piece > KING {
		panic("invalid piece type")
	}

	index := int(piece) - 1

	if player == BLACK {
		index += 6
	} else if player != WHITE {
		panic("invalid player")
	}

	return index
}
