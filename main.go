package main

import (
	"chess/fen"
	"chess/game"
	"fmt"
)

func main() {
	position := fen.LoadFenPosition(fen.StartingFEN)
	// position := fen.LoadFenPosition(qwe)
	depth := 4
	totalMoves := 0

	for {
		moves := game.GenerateMoves(&position)
		if len(moves) == 0 {
			if game.IsKingInCheck(&position, position.PlayerToMove) {
				fmt.Printf(
					"CHECKMATE: Winner is %s\n",
					game.PlayerToString(position.PlayerToMove.Opponent()),
				)
				break
			}
			fmt.Println("Stalemate")
			break
		}
		bestMove := game.FindBestMove(&position, depth)
		game.MakeMove(&position, bestMove)
		game.PrintMove(bestMove)

		totalMoves++

		// DEBUGGING
		if totalMoves > 250 {
			break
		}

		// TODO: Make a small test engine, plays x games counts W/L/D
		// Use opening book for random games
	}
}
