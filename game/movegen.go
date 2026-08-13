package game

import (
	"fmt"
)

var slidingDirectionOffsets = [...]Direction{
	// S, N, W, E
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	// NW, NE,  SE, SW
	{1, -1}, {1, 1}, {-1, 1}, {-1, -1},
}

// TODO: Precompute target squares for every square
var knightOffsets = [...]Direction{
	// NNW, NNE,
	{2, -1}, {2, 1},
	// NWW, NEE,
	{1, -2}, {1, 2},
	// SWW, SEE,
	{-1, -2}, {-1, 2},
	// SSW, SSE
	{-2, -1}, {-2, 1},
}

// Generate all pseudo legal moves for the current position
func pseudoLegalmove(position *Position) []Move {
	// loop every single piece, compute all its pseudolegal moves
	var moves []Move
	board := position.Board
	for idx := range board {
		startSquare := Square(idx)
		pieceType, pieceColor := position.GetPieceAt(startSquare)

		// No piece or wrong piece color
		if pieceType == NONE || pieceColor != position.PlayerToMove {
			continue
		}

		if pieceType == KNIGHT {
			moves = genKnightMoves(startSquare, pieceColor, moves, position)
		}
	}

	return moves
}

// Appends all pseudo-legal knight moves to the move list
func genKnightMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	startRank, startFile := SquareToRankFile(startSquare)

	fmt.Printf("startRank, startFile: (%d, %d)\n", startRank, startFile)

	// ensure that the knight doesn't move off the board
	for _, move := range knightOffsets {
		newRank := startRank + move.Rank
		newFile := startFile + move.File

		// check that the piece is within bounds
		if !IsLegalRank(newRank) || !IsLegalFile(newFile) {
			continue
		}

		targetSquare := RankFileToSquare(newRank, newFile)
		_, targetColor := position.GetPieceAt(targetSquare)
		if !IsSameColor(color, targetColor) {
			fmt.Printf("newRank, newFile: (%d, %d)\n", newRank, newFile)
			moves = append(moves, Move{startSquare, targetSquare})
		}

	}
	return moves
}

// Generate all legal moves for the position
func GenerateMove(position *Position) {
	// generate all pseudo legal moves, play them and check if our king is capturable.
	// optimize later: attack maps etc...
	pseudo := pseudoLegalmove(position)
	fmt.Printf("Pseudo-legal knight moves: %d", len(pseudo))
}
