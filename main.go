package main

import (
	"chess/fen"
	"chess/game"
)

const USE_DEFAULT_FEN = true

func main() {
	position := fen.LoadFenPosition(USE_DEFAULT_FEN)
	game.GenerateMove(&position)
}

// TODO: FEN Parser and Board representation
// Board: 64-sized array of numbers
// Pieces: numbers (enums)
// Move: Struct of (from, to)

// FEN parser: sets up the game position.

// Add testing system that tests number of legal moves for a certain depth.
// this can be precomputed and serves as a valuable test.
