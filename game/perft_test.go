package game_test

import (
	"chess/fen"
	"chess/game"
	"fmt"
	"testing"
)

type ExpectedPerft struct {
	Depth int
	Nodes uint64
}

type ExpectedPerftTable struct {
	FEN     string
	Results []ExpectedPerft
}

var startingPosPerft = ExpectedPerftTable{
	FEN: fen.StartingFEN,
	Results: []ExpectedPerft{
		{Depth: 0, Nodes: 1},
		{Depth: 1, Nodes: 20},
		{Depth: 2, Nodes: 400},
		{Depth: 3, Nodes: 8902},
		{Depth: 4, Nodes: 197281},
		{Depth: 5, Nodes: 4865609},
	},
}

var position5FEN = "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"
var position5Perft = ExpectedPerftTable{
	FEN: position5FEN,
	Results: []ExpectedPerft{
		{Depth: 1, Nodes: 44},
		{Depth: 2, Nodes: 1486},
		{Depth: 3, Nodes: 62379},
		{Depth: 4, Nodes: 2103487},
		{Depth: 5, Nodes: 89941194},
	},
}

func TestStartingFen(t *testing.T) {
	performPerft(t, startingPosPerft, 5)
}

func TestPosition5(t *testing.T) {
	performPerft(t, position5Perft, 5)
}

func performPerft(t *testing.T, table ExpectedPerftTable, maxDepth int) {
	t.Helper()

	for _, expected := range table.Results {
		if expected.Depth > maxDepth {
			continue
		}
		t.Run(fmt.Sprintf("depth_%d", expected.Depth), func(t *testing.T) {
			position := fen.LoadFenPosition(table.FEN)
			actualNodes := perft(&position, expected.Depth)

			if actualNodes != expected.Nodes {
				t.Errorf(
					"depth %d: expected %d nodes, got %d\nFEN :%s",
					expected.Depth,
					expected.Nodes,
					actualNodes,
					table.FEN,
				)
			}
			fmt.Sprintf("Success for depth %d", expected.Depth)
		})
	}
}

func perft(pos *game.Position, depth int) uint64 {
	var nodes uint64 = 0
	if depth == 0 {
		return 1
	}

	moves := game.GenerateMoves(pos)
	for _, move := range moves {
		undo := game.MakeMove(pos, move)
		nodes += perft(pos, depth-1)
		game.UnmakeMove(pos, move, undo)
	}
	return nodes
}
