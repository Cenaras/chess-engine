package game

const pawnValue int = 100
const knightValue int = 300
const bishopValue int = 320
const rookValue int = 500
const queenValue int = 900

const (
	Infinity = 100_000_000
)

func EvaluatePosition(position *Position) int {
	materialCount := getMaterialCount(position)
	return materialCount
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

	// Negamax wants the evaluation from the perspective
	// of the player to move.
	if position.PlayerToMove == BLACK.Player() {
		score = -score
	}
	return score
}
