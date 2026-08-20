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
	table         [64][12]uint64
	blackToMove   uint64
	castling      [4]uint64
	enPassantFile [8]uint64
}

var zobrist zobristTable = newZobristTable()

// Fills the zobrist table with random bitstrings for each piece
func newZobristTable() zobristTable {
	var zobrist zobristTable
	for idx := range TOTAL_SQUARES {
		for piece_idx := range 12 {
			zobrist.table[idx][piece_idx] = genBitstring()
		}
	}
	for i := range 4 {
		zobrist.castling[i] = genBitstring()
	}
	zobrist.blackToMove = genBitstring()
	for file := range 8 {
		zobrist.enPassantFile[file] = genBitstring()
	}
	return zobrist
}

// TODO: Incremental hashing: Take the current position, take the move
// and move flag, "xor the pieces out corresponding to the move". That shou;d
// be a lot more efficient than recomputing the hash each time
func ZobristHash(position *Position) uint64 {
	var hash uint64 = 0
	if position.PlayerToMove == BLACK.Player() {
		hash = hash ^ zobrist.blackToMove
	}
	for square := range TOTAL_SQUARES {
		piece := position.GetPieceAt(square)
		pieceType := piece.Type()
		if pieceType == NONE {
			continue
		}
		idx := zobristPieceIndex(pieceType, Piece(piece.Player()))
		hash ^= zobrist.table[square][idx]
	}

	rights := position.CastleRights
	for i := range 4 {
		// 0001, 0010, 0100, 1000 compactly iterated
		if rights&(1<<i) != 0 {
			hash ^= zobrist.castling[i]
		}
	}

	// FIXME: A position should only consider en-passant in the repetition,
	// if the capture is possible! This logic will technically cause our engine
	// to miss the fact that it could claim a draw (likely doesn't matter!)
	// Fix later with attack bitboards etc...
	if position.PossibleEnPassantCapture != NO_SQUARE {
		file := position.PossibleEnPassantCapture % 8
		hash ^= zobrist.enPassantFile[file]
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
func zobristPieceIndex(piece Piece, player Piece) int {
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
