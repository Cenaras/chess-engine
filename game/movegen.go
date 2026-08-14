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

// var pawnMoveOffsets = [...]Direction{
// 	// White Advance, White Capture left, White Capture Right
// 	{1,0},  {1, -1}, {1, 1},
// 	// Black Advance, Black Capture left, Black Capture Right
// 	{1,0},  {1, -1}, {1, 1},
// }

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
			fmt.Printf("Pseudo-legal moves after knight: %d\n", len(moves))
		}

		if pieceType == PAWN {
			// TODO
			moves = genPawnMoves(startSquare, pieceColor, moves, position)
			fmt.Printf("Pseudo-legal moves after pawn: %d\n", len(moves))
		}
	}

	return moves
}

func genPawnMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	var forwards, doubleForwards Direction
	if color == WHITE.Player() {
		forwards = Direction{1, 0}
		doubleForwards = Direction{2, 0}
	} else {
		forwards = Direction{-1, 0}
		doubleForwards = Direction{-2, 0}
	}

	// Single forwards move
	singleMoveSquare, err := MoveDirection(startSquare, forwards)
	if err != nil {
		panic(err) // should be impossible: pawn would have been promoted
	}
	if singleMoveTargetPiece, _ := position.GetPieceAt(singleMoveSquare); singleMoveTargetPiece == NONE {
		moves = append(moves, Move{startSquare, singleMoveSquare})
		// Also look for double pawn moves here, since no piece was infront
		if IsStartPawnRank(startSquare, color) {
			doubleMoveSquare, err := MoveDirection(startSquare, doubleForwards)
			if err != nil {
				panic(err) // should never be out of bounds for a 2x move. Just in case...
			}
			if doubleMoveTargetPiece, _ := position.GetPieceAt(doubleMoveSquare); doubleMoveTargetPiece == NONE {
				moves = append(moves, Move{startSquare, doubleMoveSquare})
			}
		}

	}
	return moves
}

// Appends all pseudo-legal knight moves to the move list
func genKnightMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	// ensure that the knight doesn't move off the board
	for _, direction := range knightOffsets {

		targetSquare, err := MoveDirection(startSquare, direction)
		if err != nil {
			continue
		}
		_, targetColor := position.GetPieceAt(targetSquare)
		if !IsSameColor(color, targetColor) {
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
	fmt.Printf("Pseudo-legal moves: %d\n", len(pseudo))
}
