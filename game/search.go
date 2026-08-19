package game

import (
	"fmt"
)

// negamax algorithm: https://chessprogramming.org/Negamax

/*Chess search is a minimax problem. The idea is as follows:
The algorithm is used to determine the score in a zero-sum game with two players:
 - The max player: Trying to maximize the score
 - The min player: Trying to minimize the score

In a one-ply search, simply generate all moves, look at the one with the best
score, select that.

For a two-ply search, assuming perfect play, the min player selects the worst move for us.
Therefore, the score of each move is the score of the worst that the opponent can do.

For instance:
 - MAX has 3 moves available, MIN has 1 response for each
 - MAX's move scores are 2, 3, 5
 - MIN's move scores are -1, 1, -3

The question is: What is MAX's best move?
 - Clearly it is move 2! That gives the best move for him.
 - That is, the best move for MAX is the worst move for MIN

Now, we could do two searches: int MAX(int depth) and int MIN(int depth)
And then compute MAX as: score = mini(depth-1); if score > max then best = score
So, max tries to maximize the score, min tries to minimize

NegaMax simplifies this: based on the fact that max(a, b) == -min(-a, -b)
So we can use a single method to compute both MAX and MIN's score

NOTICE: For this to work, the evaluation function must return a score relative
to the side being evaluated -- since min requires signs flipped,
black is also evaluated as positive on their turn
*/

// Find best move in the current position. This is the root of our search.
func FindBestMove(position *Position, depth int) Move {
	moves := GenerateMoves(position)
	bestScore := -Infinity
	var bestMove Move
	for _, move := range moves {
		undo := MakeMove(position, move)
		score := -search(position, depth-1)
		UnmakeMove(position, move, undo)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	fmt.Printf("Score for bestMove: %d\n", bestScore)
	return bestMove
}

// Simple implementation of NegaMax search
func search(position *Position, depth int) int {

	// Even at depth 0, if we are in check we must calculate if
	// we are checkmated or not
	if depth == 0 {
		if IsKingInCheck(position, position.PlayerToMove) {
			moves := GenerateMoves(position)
			if len(moves) == 0 {
				return -Infinity
			}
		}
		// If not, simply evaluate the position here
		return EvaluatePosition(position)
	}

	moves := GenerateMoves(position)
	if len(moves) == 0 {
		// Checkmake
		if IsKingInCheck(position, position.PlayerToMove) {
			return -Infinity
		}
		// Stalemate
		return 0
	}
	bestScore := -Infinity

	for _, move := range moves {
		undo := MakeMove(position, move)
		// See above explanation for why the sign is negative
		score := -search(position, depth-1)
		if score > bestScore {
			bestScore = score
		}
		UnmakeMove(position, move, undo)
	}
	return bestScore
}
