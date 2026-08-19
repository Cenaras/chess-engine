package game

import "fmt"

const pawnValue int = 100
const knightValue int = 320
const bishopValue int = 330
const rookValue int = 500
const queenValue int = 900

const (
	Infinity = 100_000_000
)

// Due to limited depth + evaluation only considering piece count / position
// we cannot convert endgames.

func EvaluatePosition(position *Position) int {

	// 100-ply without capture or pawn advancement is a draw
	if position.HalfMoveClock == 100 {
		return 0
	}

	fmt.Println(position.HalfMoveClock)

	score := 0
	// Count pieces
	score += getMaterialCount(position)
	// attribute extra score to pieces on "good" squares
	score += evaluatePieceTables(position)

	// Negamax wants the evaluation from the perspective
	// of the player to move.
	if position.PlayerToMove == BLACK.Player() {
		score = -score
	}

	return score
}

// TOOD: Make an EvalConfig that allows us to enable/disable parts of our eval
// to test if its better or not!

// TODO: We should really store where all pieces are instead of iterating!
// NOTE: This function must evaluate the position relative the side to move,
// for NegaMax (search) to work
func getMaterialCount(position *Position) int {
	score := 0

	for square := range TOTAL_SQUARES {
		piece := position.GetPieceAt(Square(square))
		pieceType := piece.Type()

		sign := 1
		if piece.Player() == BLACK.Player() {
			sign = -1
		}

		switch pieceType {
		case PAWN:
			score += pawnValue * sign
		case KNIGHT:
			score += knightValue * sign
		case BISHOP:
			score += bishopValue * sign
		case ROOK:
			score += rookValue * sign
		case QUEEN:
			score += queenValue * sign
		}
	}

	return score
}

// TODO: The next thing we do is update to use bitboards or at least
// know where all pieces are at all times to avoid this traversal!
func evaluatePieceTables(position *Position) int {
	score := 0
	for square := range TOTAL_SQUARES {
		piece := position.GetPieceAt(Square(square))
		sign := 1
		if piece.Player() == BLACK.Player() {
			sign = -1
		}

		switch piece.Type() {
		case PAWN:
			score += pieceTableValue(pawnPieceTable, square, piece.Player()) * sign
		case KNIGHT:
			score += pieceTableValue(knightPieceTable, square, piece.Player()) * sign
		case BISHOP:
			score += pieceTableValue(bishopPieceTable, square, piece.Player()) * sign
		case ROOK:
			score += pieceTableValue(rookPieceTable, square, piece.Player()) * sign
		case QUEEN:
			score += pieceTableValue(queenPieceTable, square, piece.Player()) * sign
		}
	}
	return score
}
