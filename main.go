package main

import (
	"chess/fen"
	"chess/game"
)

func main() {
	position := fen.LoadFenPosition(fen.StartingFEN)
	game.GenerateMoves(&position)
}
