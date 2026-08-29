package testing

import (
	"chess/fen"
	"chess/game"
	"context"
	"testing"
	"time"
)

func TestFixedDepthPosition(t *testing.T) {
	position := fen.LoadFenPosition(
		"r2q1rk1/1bp2pp1/2pp3p/p2n2b1/Q1N5/2PP1NB1/PP3PPP/R3R1K1 b - - 1 17",
	)

	options := game.SearchOptions{
		WhiteTime:      200000 * time.Millisecond,
		BlackTime:      200000 * time.Millisecond,
		WhiteIncrement: 0 * time.Millisecond,
		BlackIncrement: 0 * time.Millisecond,
		Depth:          6,
	}

	ctx, _ := context.WithTimeout(context.Background(), 200000*time.Millisecond)
	game.FindBestMove(&position, options, ctx)
}
