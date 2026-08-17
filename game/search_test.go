package game_test

// LLM GENERATED TESTS

import (
	"testing"

	"chess/fen"
	"chess/game"
)

func searchSquare(square string) game.Square {
	file := int(square[0] - 'a')
	rank := int(square[1] - '1')

	return game.Square(rank*8 + file)
}

func requireBestMove(
	t *testing.T,
	fenString string,
	depth int,
	expectedFrom string,
	expectedTo string,
) {
	t.Helper()

	position := fen.LoadFenPosition(fenString)

	bestMove := game.FindBestMove(&position, depth)

	expectedFromSquare := searchSquare(expectedFrom)
	expectedToSquare := searchSquare(expectedTo)

	if bestMove.From != expectedFromSquare ||
		bestMove.To != expectedToSquare {

		t.Fatalf(
			"expected %s -> %s, got %s -> %s",
			expectedFrom,
			expectedTo,
			game.SquareToString(bestMove.From),
			game.SquareToString(bestMove.To),
		)
	}
}
func TestFindBestMoveWhiteCapturesFreeQueen(t *testing.T) {
	requireBestMove(
		t,
		"7k/q7/8/8/8/8/8/R6K w - - 0 1",
		1,
		"a1",
		"a7",
	)
}
func TestFindBestMoveBlackCapturesFreeQueen(t *testing.T) {
	requireBestMove(
		t,
		"r6k/8/8/8/8/8/Q7/7K b - - 0 1",
		1,
		"a8",
		"a2",
	)
}
func TestFindBestMoveLooksAheadBeforeCapturing(t *testing.T) {
	requireBestMove(
		t,
		"k2q4/3r4/8/7b/8/8/8/K2Q4 w - - 0 1",
		2,
		"d1",
		"h5",
	)
}
func TestFindBestMoveFindsMateInOne(t *testing.T) {
	requireBestMove(
		t,
		"k7/7R/8/8/8/8/8/K5R1 w - - 0 1",
		1,
		"g1",
		"g8",
	)
}
func TestFindBestMoveMateInTwo1(t *testing.T) {
	// 1. Re4-g4!
	//    ... Kf8
	// 2. Rh6-f6#
	requireBestMove(
		t,
		"3K4/5k2/7R/8/4R3/8/8/8 w - - 0 1",
		4,
		"e4",
		"g4",
	)
}

func TestFindBestMoveMateInTwo2(t *testing.T) {
	// 1. Ra7-g7!
	//    ... Kh2
	// 2. Re4-h4#
	requireBestMove(
		t,
		"8/R7/4K3/8/4R3/7k/8/8 w - - 0 1",
		4,
		"a7",
		"g7",
	)
}

func TestFindBestMoveMateInTwo3(t *testing.T) {
	// 1. Rh1-b1!
	//    ... Ka3
	// 2. Rc5-a5#
	requireBestMove(
		t,
		"8/8/8/2R5/k3K3/8/8/7R w - - 0 1",
		4,
		"h1",
		"b1",
	)
}

func TestFindBestMoveMateInTwo4(t *testing.T) {
	// 1. Rf6-g6!
	//    ... Kh1
	// 2. Re3-h3#
	requireBestMove(
		t,
		"8/8/5R2/8/8/2K1R3/7k/8 w - - 0 1",
		4,
		"f6",
		"g6",
	)
}

func TestFindBestMoveMateInTwo5(t *testing.T) {
	// 1. Rf4-e4!
	//    ... Kc8
	// 2. Re4-e8#
	requireBestMove(
		t,
		"3k4/5R2/8/8/5R2/4K3/8/8 w - - 0 1",
		4,
		"f4",
		"e4",
	)
}
