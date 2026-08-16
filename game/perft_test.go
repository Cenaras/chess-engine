package game_test

// See: https://chessprogramming.org/Perft_Results

import (
	"chess/fen"
	"chess/game"
	"fmt"
	"testing"
	"time"
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
		{Depth: 1, Nodes: 20},
		{Depth: 2, Nodes: 400},
		{Depth: 3, Nodes: 8902},
		{Depth: 4, Nodes: 197281},
		{Depth: 5, Nodes: 4865609},
		{Depth: 6, Nodes: 119060324},
	},
}

var position3Perft = ExpectedPerftTable{
	FEN: "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
	Results: []ExpectedPerft{
		{Depth: 1, Nodes: 14},
		{Depth: 2, Nodes: 191},
		{Depth: 3, Nodes: 2812},
		{Depth: 4, Nodes: 43238},
		{Depth: 5, Nodes: 674624},
		{Depth: 6, Nodes: 11030083},
	},
}

// Only 4 nodes as depth 5 and 6 explode. Return and extend with larger tables once
// implementation is faster
var position4Perft = ExpectedPerftTable{
	FEN: "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
	Results: []ExpectedPerft{
		{Depth: 1, Nodes: 6},
		{Depth: 2, Nodes: 264},
		{Depth: 3, Nodes: 9467},
		{Depth: 4, Nodes: 422333},
		{Depth: 5, Nodes: 15833292},
		{Depth: 6, Nodes: 706045033},
	},
}

var position5Perft = ExpectedPerftTable{
	FEN: "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
	Results: []ExpectedPerft{
		{Depth: 1, Nodes: 44},
		{Depth: 2, Nodes: 1486},
		{Depth: 3, Nodes: 62379},
		{Depth: 4, Nodes: 2103487},
		{Depth: 5, Nodes: 89941194},
	},
}

func TestStartingFen(t *testing.T) {
	performPerft(t, startingPosPerft, 6)
}

func TestPosition3(t *testing.T) {
	performPerft(t, position3Perft, 6)
}

func TestPosition4(t *testing.T) {
	performPerft(t, position4Perft, 6)
}

func TestPosition5(t *testing.T) {
	performPerft(t, position5Perft, 5)
}

func TestPosition5CanCastleKingSide(t *testing.T) {
	position := fen.LoadFenPosition(
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
	)

	moves := game.GenerateMoves(&position)

	found := false

	for _, move := range moves {
		if move.Flag == game.KingCastle {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected white kingside castling to be legal")
	}
}

func performPerft(t *testing.T, table ExpectedPerftTable, maxDepth int) {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping pert tests in short mode")
	}

	for _, expected := range table.Results {
		if expected.Depth > maxDepth {
			continue
		}
		t.Run(fmt.Sprintf("depth_%d", expected.Depth), func(t *testing.T) {

			fmt.Printf(
				"Starting depth %d (expected %d nodes)...\n",
				expected.Depth,
				expected.Nodes,
			)
			position := fen.LoadFenPosition(table.FEN)

			start := time.Now()
			actualNodes := bulkPerft(&position, expected.Depth)
			duration := time.Since(start)

			fmt.Printf(
				"Completed depth %d: %d nodes in %s\n",
				expected.Depth,
				actualNodes,
				duration,
			)

			if actualNodes != expected.Nodes {
				t.Errorf(
					"depth %d: expected %d nodes, got %d\nFEN :%s",
					expected.Depth,
					expected.Nodes,
					actualNodes,
					table.FEN,
				)
			}
		})
	}
}

func bulkPerft(pos *game.Position, depth int) uint64 {
	if depth < 1 {
		panic("bulk perft must have depth >= 1")
	}
	var nodes uint64
	moves := game.GenerateMoves(pos)
	numberOfMoves := len(moves)

	if depth == 1 {
		return uint64(numberOfMoves)
	}
	for _, move := range moves {
		undo := game.MakeMove(pos, move)
		nodes += bulkPerft(pos, depth-1)
		game.UnmakeMove(pos, move, undo)
	}
	return nodes
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
