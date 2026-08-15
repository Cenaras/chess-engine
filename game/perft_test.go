package game_test

import (
	"chess/fen"
	"chess/game"
	"testing"
)

func TestStartingFen(t *testing.T) {
	performPerft(t, fen.StartingFEN, 0, 1)
	performPerft(t, fen.StartingFEN, 1, 20)

}

func performPerft(t *testing.T, testFen string, depth int, expectedNodes uint64) {
	position := fen.LoadFenPosition(testFen)
	actualNodes := perft(&position, depth)
	if actualNodes != expectedNodes {
		t.Errorf("Expected %d moves, but got %d for position: \n%s", expectedNodes, actualNodes, testFen)
	}
}

func perft(pos *game.Position, depth int) uint64 {
	var nodes uint64 = 0
	if depth == 0 {
		return 1
	}

	moves := game.GenerateMoves(pos)
	for _, move := range moves {
		game.MakeMove(move)
		nodes += perft(pos, depth-1)
		game.UnmakeMove(move)
	}
	return nodes
}
