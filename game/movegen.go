package game

var rookDirections = []Direction{
	// S, N, W, E
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
}
var bishopDirections = []Direction{
	// NW, NE,  SE, SW
	{1, -1}, {1, 1}, {-1, 1}, {-1, -1},
}

// TODO: Precompute target squares for every square
var knightOffsets = [...]Direction{
	// NNW, NNE,
	{2, -1}, {2, 1},
	// NWW, NEE,
	{1, -2}, {1, 2},
	// SWW, SEE,
	{-1, -2}, {-1, 2},
	// SSW, SSE
	{-2, -1}, {-2, 1},
}

// Generate all pseudo legal moves for the current position
func pseudoLegalmove(position *Position) []Move {
	// loop every single piece, compute all its pseudolegal moves
	var moves []Move
	board := position.Board
	for idx := range board {
		startSquare := Square(idx)
		pieceType, pieceColor := position.GetPieceAt(startSquare)

		// No piece or wrong piece color
		if pieceType == NONE || pieceColor != position.PlayerToMove {
			continue
		}

		switch pieceType {
		case KNIGHT:
			moves = genKnightMoves(startSquare, pieceColor, moves, position)
		case PAWN:
			moves = genPawnMoves(startSquare, pieceColor, moves, position)
		case KING:
			moves = genKingMoves(startSquare, pieceColor, moves, position)
		case BISHOP:
			moves = genSlidingPieceMoves(startSquare, pieceColor, moves, position, bishopDirections)
		case ROOK:
			moves = genSlidingPieceMoves(startSquare, pieceColor, moves, position, rookDirections)
		case QUEEN:
			moves = genSlidingPieceMoves(startSquare, pieceColor, moves, position, rookDirections)
			moves = genSlidingPieceMoves(startSquare, pieceColor, moves, position, bishopDirections)
		}
	}
	return moves
}

func genSlidingPieceMoves(startSquare Square, color Player, moves []Move, position *Position, directions []Direction) []Move {
	for _, direction := range directions {
		currentSquare := startSquare
		for {
			targetSquare, err := MoveDirection(currentSquare, direction)
			// Moving further in this direction leaves the board
			if err != nil {
				break
			}

			// Get piece at targetSquare
			pieceType, pieceColor := position.GetPieceAt(targetSquare)
			// Friendly piece blocks us
			if pieceType != NONE && pieceColor == color {
				break
			}

			// Square is empty
			moves = append(moves, Move{currentSquare, targetSquare, NormalMove})

			// If we captured a piece, break as well
			if pieceType != NONE {
				break
			}

			// Continue search from new square
			currentSquare = targetSquare
		}
	}
	return moves
}

func genKingMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	var directions = [...]Direction{
		{1, 0}, {1, 1}, {1, -1},
		{0, 1}, {0, -1},
		{-1, 0}, {-1, 1}, {-1, -1},
	}

	// Regular moves
	for _, direction := range directions {
		targetSquare, err := MoveDirection(startSquare, direction)
		if err != nil {
			continue
		}
		targetPiece, targetColor := position.GetPieceAt(targetSquare)
		// Empty or opponent piece
		if targetPiece == NONE || !IsSameColor(targetColor, color) {
			moves = append(moves, Move{startSquare, targetSquare, NormalMove})
		}

		// Castle moves:
		// TODO, indicate in the Move that it is a castle, or check
		//	whenever king makes a move?
		kingSide := WhiteKingSide
		queenSide := WhiteQueenSide

		if color == BLACK.Player() {
			kingSide = BlackKingSide
			queenSide = BlackQueenSide
		}

		addCastleMove := func(side CastleRights, direction int, sideFlag MoveFlag) {
			if position.CastleRights&side != 0 {
				// since castling is valid we dont need to check if OOB
				moves = append(moves, Move{startSquare, startSquare + Square(direction), sideFlag})
			}
		}
		addCastleMove(kingSide, 2, KingCastle)
		addCastleMove(queenSide, -2, QueenCastle)
	}
	return moves
}

func genPawnMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	var forwards, doubleForwards Direction
	if color == WHITE.Player() {
		forwards = Direction{1, 0}
		doubleForwards = Direction{2, 0}
	} else {
		forwards = Direction{-1, 0}
		doubleForwards = Direction{-2, 0}
	}

	// Single forwards move
	singleMoveSquare, err := MoveDirection(startSquare, forwards)
	if err != nil {
		panic(err) // should be impossible: pawn would have been promoted
	}
	if singleMoveTargetPiece, _ := position.GetPieceAt(singleMoveSquare); singleMoveTargetPiece == NONE {
		moves = append(moves, Move{startSquare, singleMoveSquare, NormalMove})
		// Also look for double pawn moves here, since no piece was infront
		if IsStartPawnRank(startSquare, color) {
			doubleMoveSquare, err := MoveDirection(startSquare, doubleForwards)
			if err != nil {
				panic(err) // should never be out of bounds for a 2x move. Just in case...
			}
			if doubleMoveTargetPiece, _ := position.GetPieceAt(doubleMoveSquare); doubleMoveTargetPiece == NONE {
				moves = append(moves, Move{startSquare, doubleMoveSquare, DoublePawnPush})
			}
		}
	}
	var attackLeft, attackRight Direction
	if color == WHITE.Player() {
		attackLeft = Direction{1, -1}
		attackRight = Direction{1, -1}
	} else {
		attackLeft = Direction{-1, -1}
		attackRight = Direction{-1, 1}
	}

	attackSide := func(square Square, direction Direction) []Move {
		squareToAttack, err := MoveDirection(square, direction)
		if err != nil {
			// check for normal attack
			pieceToAttack, pieceColor := position.GetPieceAt(squareToAttack)
			if pieceToAttack != NONE && !IsSameColor(color, pieceColor) {
				moves = append(moves, Move{startSquare, squareToAttack, NormalMove})
			}
			// check for en-passant
			if pieceToAttack == NONE && squareToAttack == position.PossibleEnPassantCapture {
				moves = append(moves, Move{startSquare, squareToAttack, EnPassantCapture})
			}
		}
		return moves
	}

	moves = attackSide(startSquare, attackLeft)
	moves = attackSide(startSquare, attackRight)
	return moves
}

// Appends all pseudo-legal knight moves to the move list
func genKnightMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	// ensure that the knight doesn't move off the board
	for _, direction := range knightOffsets {

		targetSquare, err := MoveDirection(startSquare, direction)
		if err != nil {
			continue
		}
		_, targetColor := position.GetPieceAt(targetSquare)
		if !IsSameColor(color, targetColor) {
			moves = append(moves, Move{startSquare, targetSquare, NormalMove})
		}

	}
	return moves
}

// Generate all legal moves for the position
func GenerateMoves(position *Position) []Move {
	// generate all pseudo legal moves, play them and check if our king is capturable.
	// optimize later: attack maps etc...
	pseudo := pseudoLegalmove(position)
	return pseudo
}

func MakeMove(p *Position, move Move) UndoMoveState {
	// TODO: make the move

	from := move.From
	to := move.To
	moveFlag := move.Flag
	piece, _ := p.GetPieceAt(from)
	// Move pieces around
	p.SetPieceAt(piece, to)
	p.SetPieceAt(NONE, from)

	// Update the position state
	switch moveFlag {
	case DoublePawnPush:
	case EnPassantCapture:
	case KingCastle:
	case QueenCastle:
	case PromoteBishop:
	case PromoteKnight:
	case PromoteRook:
	case PromoteQueen:
	}

	// TODO: ...
	return UndoMoveState{}

}

func UnmakeMove(p *Position, move Move, undo UndoMoveState) {
	// TODO: unmake the move.
	// Remember to undo any side-effects making the move had
	// such as en-passant squares etc
}
