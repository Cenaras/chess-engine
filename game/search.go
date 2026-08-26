package game

import (
	"cmp"
	"context"
	"fmt"
	"slices"
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
	start := time.Now()
	maxDepth := options.Depth
	if maxDepth == 0 {
		// If no max depth, keep searching until time limit
		maxDepth = Infinity
	}
	ttEntry, hasTTMove := TT.Lookup(position.Hash)

	moves := GenerateMoves(position)
	orderedMoves := OrderMoves(moves, position, hasTTMove, ttEntry.BestMove)
	var bestMove Move
	if len(orderedMoves) == 0 {
		return Move{} // Will error
	}

	// To ensure valid legal move in case of early termination
	bestMove = orderedMoves[0]

	// Iteratively search the root to increasing depths, starting at depth 1
	for iteration := 1; iteration <= maxDepth; iteration++ {
		// The best move from last iteration is likely still a good move: so explore that first
		if iteration > 1 {
			moveToFront(moves, bestMove)
		}

		nodesSearchedForIteration = 0
		// lower/upper bounds are reset per iteration
		alpha := -Infinity
		beta := Infinity

		var bestMoveForIteration Move
		bestScoreForIteration := -Infinity
		// Consider ever possible move
		for _, move := range orderedMoves {
			undo := MakeMove(position, move)
			childScore, completed := search(position, iteration-1, -beta, -alpha, ctx)
			UnmakeMove(position, move, undo)

			// When search is terminated, return the best move we've seen for this depth.
			if !completed {
				return bestMove
			}
			// Update best move.
			scoreForIteration := -childScore
			if scoreForIteration > bestScoreForIteration {
				bestMoveForIteration = move
				bestScoreForIteration = scoreForIteration
			}
			// Update alpha as the best score so far from other branch
			if scoreForIteration > alpha {
				alpha = scoreForIteration
			}
		}
		// Once the iteration completed, update the currently found best move for this
		// terminated iteration
		bestMove = bestMoveForIteration
		fmt.Printf(
			"info depth %d score cp %d time %d nodes %d\n",
			iteration,
			bestScoreForIteration,
			time.Since(start).Milliseconds(),
			nodesSearchedForIteration,
		)
	}
	return bestMove
}

// Simple implementation of NegaMax search
// If the search is terminated, we report the best score that was fully evaluated.
/*
	Alpha: a lower bound on the score MAX can already guarantee.
	If a MIN node finds a reply with value <= alpha, this branch
	cannot improve MAX's result, so the remaining replies can be pruned.

	Beta: an upper bound on the score MIN can already guarantee MAX will get.
	If a MAX node finds a continuation with value >= beta, MIN would
	never choose the ancestor branch leading here, so the remaining
	continuations can be pruned.
*/
func search(position *Position, depth int, alpha int, beta int, ctx context.Context) (int, bool) {
	nodesSearchedForIteration++

	// Search was termianted before we saw the final leaves.
	// Report that we cannot trust this evaluation.
	if ctx.Err() != nil {
		return 0, false
	}

	// Check for terminal draw rules
	if IsThreefoldRepetition(position) {
		return 0, true
	}
	if position.HalfMoveClock >= 100 {
		return 0, true
	}

	originalAlpha := alpha
	originalBeta := beta

	// TODO: Look in TT to see if the current position has a known result
	ttEntry, ttFound := TT.Lookup(position.Hash)
	// We have already seen this position before: Make sure it is at a sufficiently deep depth,
	// otherwise the result isn't trustworthy
	if ttFound && ttEntry.Depth >= depth {
		switch ttEntry.NodeType {
		// if the score is exact, we can definitely use it
		case EXACT:
			return ttEntry.Score, true
		// if the score was a lower bound on the position value, and that lower bound is greated than beta,
		// our opponent will never pick this move, so don't continue the search
		case LOWER_BOUND:
			alpha = max(alpha, ttEntry.Score)
		// if the score was an upper bound, and we have a better move in alpha already, quit the search
		case UPPER_BOUND:
			beta = min(beta, ttEntry.Score)
		}

		if alpha >= beta {
			return ttEntry.Score, true
		}

		// If none of these were hits, we can use the lower bound and upper bound to tighten
		// the search window.
		// Say we have alpha = 20, beta = 80, and we found LOWER_BOUND = 40
		// We cannot report any score, since it *may* be 50, 100, 10000
		// However, we now know that we have another move that gives us a score of 40, so we can
		// tighten alpha and beta:
		// - If the value is a LOWER_BOUND smaller than beta, then our opponent will still consider this move, but is is better than alpha, so it is a new lower bound
		// - If the value is an UPPER_BOUND greater than alpha, then we can tighten beta with it
		//		For instance before: 20 <= value < 80. Now 40 is an upper bound, so
		//		20 <= trueValue <= 40
	}

	// Order moves based on heuristic
	moves := GenerateMoves(position)
	// We can still attempt the best TT move even though the depth was less than what we are currently searching
	orderedMoves := OrderMoves(moves, position, ttFound, ttEntry.BestMove)

	// Terminal positions
	if len(orderedMoves) == 0 {
		if IsKingInCheck(position, position.PlayerToMove) {
			return -MateScore, true // checkmate // TODO: +ply-to-mate
		}
		return 0, true // stalemate
	}

	// If non-terminal position, evaluate the position
	if depth == 0 {
		return EvaluatePosition(position), true
	}
	bestScore := -Infinity
	var bestMove Move = orderedMoves[0]

	for _, move := range orderedMoves {
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
			bestMove = move
		}
		// Tighten lower bound
		alpha = max(alpha, score)
		// Parent will never choose the line leading to this node.
		// MIN will prefer the branch leading to beta, so no point looking at further positions
		if alpha >= beta {
			break
		}
	}
	nodeType := EXACT
	// If bestScore <= originalAlpha, no move in this search managed to increase our lower bound.
	// Therefore the true evaluation of the current position is at most bestScore
	if bestScore <= originalAlpha {
		nodeType = UPPER_BOUND
	}
	// If bestScore >= beta, we found a move that is at least beta and stopped searching.
	// The true value may be higher
	if bestScore >= originalBeta {
		nodeType = LOWER_BOUND
	}
	ttEntry = TranspositionTableEntry{
		Hash:     position.Hash,
		BestMove: bestMove,
		Depth:    depth,
		Score:    bestScore,
		NodeType: nodeType,
	}
	TT.Store(ttEntry)
	return bestScore, true
}

func moveToFront(moves []Move, move Move) {
	for i := range moves {
		if moves[i] == move {
			moves[0], moves[i] = moves[i], moves[0]
			return
		}
	}
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

// The PV mode from previous will be added to the front in the iterative deepening
func OrderMoves(legalMoves []Move, position *Position, ttMoveFound bool, ttMove Move) []Move {

	// Ideas for move ordering:
	// Promotions

	moveScore := func(move Move, ttMove Move) int {
		// Transposition Table move gets highest priority (PV mode handled by iterative deepening)
		if ttMoveFound && move == ttMove {
			// highest priority
			return 1_000_000
		}

		movingPiece := position.GetPieceAt(move.From)
		capturedPiece := position.GetPieceAt(move.To)
		if capturedPiece.Type() == NONE {
			// down-prioritize non-capture moves for now
			return -1
		}
		gain := GetPieceScore(capturedPiece) - GetPieceScore(movingPiece)

		// Evaluate moves based on material gain
		if gain >= 0 {
			return 100_000 + gain
		}
		return 0
	}

	// Sort moves based on their assigned score
	slices.SortStableFunc(legalMoves, func(a, b Move) int {
		return cmp.Compare(
			moveScore(b, ttMove),
			moveScore(a, ttMove),
		)
	})
	return legalMoves
}
