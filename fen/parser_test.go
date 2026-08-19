package fen

import (
	"chess/game"
	"testing"
)

func TestFenBoardOrientation(t *testing.T) {
	position := parseFen(
		"8/8/8/8/8/8/2n5/8 b - - 0 1",
	)

	// c2 = 1*8 + 2 = 10
	expected := game.KNIGHT | game.BLACK

	if position.GetPieceAt(10) != expected {
		t.Fatalf(
			"expected black knight on c2, got %v",
			position.GetPieceAt(10),
		)
	}
}

func TestHalfMoveClock(t *testing.T) {
	position := parseFen("8/5k2/3p4/1p1Pp2p/pP2Pp1P/P4P1K/8/8 b - - 99 50")
	if position.HalfMoveClock != 99 {
		t.Fatalf(
			"expected half move clock to be 99, but was %d",
			position.HalfMoveClock)
	}
}

func TestFenEnPassantOrientation(t *testing.T) {
	position := parseFen(
		"8/8/8/8/8/8/8/8 w - e3 0 1",
	)

	// e3 = 2*8 + 4 = 20
	if position.PossibleEnPassantCapture != 20 {
		t.Fatalf(
			"expected e3 to map to square 20, got %d",
			position.PossibleEnPassantCapture,
		)
	}
}
