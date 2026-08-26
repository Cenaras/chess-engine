package game

const pawnValue int = 100
const knightValue int = 320
const bishopValue int = 330
const rookValue int = 500
const queenValue int = 900

const (
	Infinity  = 100_000_000
	MateScore = 100_000
)

// Due to limited depth + evaluation only considering piece count / position
// we cannot convert endgames.

// Static Evaluation: Does not consider 3-fold repetition or 50 move rule -- see search.go#search
func EvaluatePosition(position *Position) int {
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
		if piece.Type() == NONE {
			continue
		}

		sign := 1
		if piece.Player() == BLACK.Player() {
			sign = -1
		}

		pieceScore := GetPieceScore(piece) * sign
		score += pieceScore
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

func GetPieceScore(piece Piece) int {
	pieceType := piece.Type()
	if piece.Type() == NONE {
		panic("GetPieceScore should not be called on a NONE piece")
	}
	switch pieceType {
	case PAWN:
		return pawnValue
	case KNIGHT:
		return knightValue
	case BISHOP:
		return bishopValue
	case ROOK:
		return rookValue
	case QUEEN:
		return queenValue
	}
	// The King's value is not relevant, since he can never be captured
	return 0
}
