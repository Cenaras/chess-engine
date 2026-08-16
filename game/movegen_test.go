package game_test

import (
	"chess/fen"
	"chess/game"
	"testing"
)

// An LLM generated test suite for regression tests on common pitfalls

func TestPawnInitialMoves(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/8/8/8/4P3/4K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "e2", "e3", game.NormalMove)
	requireMove(t, moves, "e2", "e4", game.DoublePawnPush)
}
func TestPawnCannotMoveThroughPiece(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/8/8/4n3/4P3/4K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireNoMove(t, moves, "e2", "e3", game.NormalMove)
	requireNoMove(t, moves, "e2", "e4", game.DoublePawnPush)
}
func TestPawnCaptures(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/3n4/4P3/8/8/4K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "e4", "d5", game.NormalMove)
	requireNoMove(t, moves, "e4", "f5", game.NormalMove)
}
func TestKnightDoesNotWrapAroundBoard(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/8/8/8/8/N3K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "a1", "b3", game.NormalMove)
	requireMove(t, moves, "a1", "c2", game.NormalMove)

	count := 0
	for _, move := range moves {
		if move.From == sq("a1") {
			count++
		}
	}

	if count != 2 {
		t.Fatalf("expected knight on a1 to have 2 moves, got %d", count)
	}
}
func TestKingSideCastleAllowed(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/8/8/8/8/4K2R w K - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "e1", "g1", game.KingCastle)
}
func TestCannotCastleThroughCheck(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/8/2b5/8/8/4K2R w K - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireNoMove(t, moves, "e1", "g1", game.KingCastle)
}
func TestCannotCastleOutOfCheck(t *testing.T) {
	position := fen.LoadFenPosition(
		"k3r3/8/8/8/8/8/8/4K2R w K - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireNoMove(t, moves, "e1", "g1", game.KingCastle)
}
func TestPinnedPieceCannotExposeKing(t *testing.T) {
	position := fen.LoadFenPosition(
		"k3r3/8/8/8/8/8/4R3/4K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	// Moving sideways exposes the king on e1 to the black rook on e8.
	requireNoMove(t, moves, "e2", "d2", game.NormalMove)
}
func TestKingCannotMoveIntoCheck(t *testing.T) {
	position := fen.LoadFenPosition(
		"k2r4/8/8/8/8/8/8/4K3 w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	// Black rook on d8 attacks d1.
	requireNoMove(t, moves, "e1", "d1", game.NormalMove)
}
func TestIllegalEnPassantDiscoveredCheck(t *testing.T) {
	position := fen.LoadFenPosition(
		"k7/8/8/4KPpr/8/8/8/8 w - g6 0 1",
	)

	moves := game.GenerateMoves(&position)

	// White pawn f5 could geometrically capture g6 en passant.
	//
	// But removing both f5 and g5 would open:
	//
	// h5 rook -> e5 king
	//
	// so the move must be illegal.
	requireNoMove(
		t,
		moves,
		"f5",
		"g6",
		game.EnPassantCapture,
	)
}
func TestEnPassantMakeAndUnmake(t *testing.T) {
	position := fen.LoadFenPosition(
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
	)

	moves := game.GenerateMoves(&position)

	var epMove game.Move
	found := false

	for _, move := range moves {
		if move.From == sq("e5") &&
			move.To == sq("d6") &&
			move.Flag == game.EnPassantCapture {
			epMove = move
			found = true
			break
		}
	}

	if !found {
		t.Fatal("expected e5xd6 en passant")
	}

	undo := game.MakeMove(&position, epMove)

	if position.GetPieceAt(sq("d6")) != game.PAWN|game.WHITE {
		t.Fatal("expected white pawn on d6 after en passant")
	}

	if position.GetPieceAt(sq("d5")) != game.NONE {
		t.Fatal("expected captured pawn on d5 to be removed")
	}

	game.UnmakeMove(&position, epMove, undo)

	if position.GetPieceAt(sq("e5")) != game.PAWN|game.WHITE {
		t.Fatal("expected white pawn restored to e5")
	}

	if position.GetPieceAt(sq("d5")) != game.PAWN|game.BLACK {
		t.Fatal("expected black pawn restored to d5")
	}

	if position.GetPieceAt(sq("d6")) != game.NONE {
		t.Fatal("expected d6 to be empty after unmake")
	}
}
func TestPawnPromotionGeneratesFourMoves(t *testing.T) {
	position := fen.LoadFenPosition(
		"7k/P7/8/8/8/8/8/7K w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "a7", "a8", game.PromoteQueen)
	requireMove(t, moves, "a7", "a8", game.PromoteRook)
	requireMove(t, moves, "a7", "a8", game.PromoteBishop)
	requireMove(t, moves, "a7", "a8", game.PromoteKnight)
}
func TestPawnCapturePromotionGeneratesFourMoves(t *testing.T) {
	position := fen.LoadFenPosition(
		"r6k/1P6/8/8/8/8/8/7K w - - 0 1",
	)

	moves := game.GenerateMoves(&position)

	requireMove(t, moves, "b7", "a8", game.PromoteQueen)
	requireMove(t, moves, "b7", "a8", game.PromoteRook)
	requireMove(t, moves, "b7", "a8", game.PromoteBishop)
	requireMove(t, moves, "b7", "a8", game.PromoteKnight)
}
