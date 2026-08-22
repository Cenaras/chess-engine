package game

import (
	"context"
	"fmt"
	"time"
)

// negamax algorithm: https://chessprogramming.org/Negamax
type SearchOptions struct {
	WhiteTime      time.Duration
	BlackTime      time.Duration
	WhiteIncrement time.Duration
	BlackIncrement time.Duration

	Depth         int
	MaxSearchTime time.Duration
}

var nodesSearchedForIteration int = 0

// Find best move in the current position. This is the root of our search.
func FindBestMove(position *Position, options SearchOptions, ctx context.Context) Move {
	// TODO: This is sort of messy! We do root search and normal search differently?
	start := time.Now()
	moves := GenerateMoves(position)

	depth := options.Depth
	bestScore := -Infinity

	// TODO: May be uninitialized if search is terminated before it is updated
	var bestMove Move
	// TODO: Iterative deepening!
	// For now, if depth is 0 it just means "keep going until time limit"
	// as Time limit isn't supplied -- always assume 4...
	if depth == 0 {
		depth = 4
	}
	nodesSearchedForIteration = 0
	for _, move := range moves {
		undo := MakeMove(position, move)
		childScore, completed := search(position, depth-1, -Infinity, Infinity, ctx)
		UnmakeMove(position, move, undo)

		// Search was terminated; quit iterating and report best result so far
		if !completed {
			fmt.Println("info search was canelled")
			break
		}

		score := -childScore

		// At this point we know that the previous search was completed.
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	fmt.Printf(
		"info depth %d score cp %d time %d nodes %d\n",
		depth,
		bestScore,
		time.Since(start).Milliseconds(),
		nodesSearchedForIteration,
	)
	return bestMove
}

// Simple implementation of NegaMax search
// If the search is terminated, we report the best score that was fully
// evaluated
func search(position *Position, depth int, alpha int, beta int, ctx context.Context) (int, bool) {
	nodesSearchedForIteration++

	// Check for terminal draw rules
	if IsThreefoldRepetition(position) {
		return 0, true
	}
	if position.HalfMoveClock >= 100 {
		return 0, true
	}

	moves := GenerateMoves(position)
	// Terminal positions
	if len(moves) == 0 {
		if IsKingInCheck(position, position.PlayerToMove) {
			return -Infinity, true // checkmate
		}
		return 0, true // stalemate
	}

	// If non-terminal position, evaluate the position
	if depth == 0 {
		return EvaluatePosition(position), true
	}
	bestScore := -Infinity

	for _, move := range moves {
		// Search was termianted before we saw the final leaves.
		// Report that we cannot trust this evaluation.
		if ctx.Err() != nil {
			return 0, false
		}

		undo := MakeMove(position, move)
		// See above explanation for why the sign is negative
		childScore, completed := search(position, depth-1, -beta, -alpha, ctx)
		UnmakeMove(position, move, undo)

		if !completed {
			return 0, false
		}

		// Negamax requires us to invert the resul -- cannot do that at the return
		// since we are returning multiple values...
		score := -childScore
		if score > bestScore {
			bestScore = score
			if score > alpha {
				alpha = score
			}
		}
		if score >= beta {
			return bestScore, true
		}
	}
	return bestScore, true
}

func IsThreefoldRepetition(position *Position) bool {
	history := position.History
	current := position.GetCurrentPositionHash()

	occurrences := 1
	// Why -3?: -1 is the current position. -2 is opponents turn (always diff from current ) as
	// zobrist considers the playerToMove
	for i := len(history) - 3; i >= 0; i -= 2 {
		if history[i] == current {
			occurrences++
			if occurrences == 3 {
				return true
			}
		}
	}
	return false
}
