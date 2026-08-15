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
		piece := position.GetPieceAt(startSquare)
		pieceType := piece.Type()
		pieceColor := piece.Player()

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
			piece := position.GetPieceAt(targetSquare)
			pieceType := piece.Type()
			pieceColor := piece.Player()
			// Friendly piece blocks us
			if pieceType != NONE && pieceColor == color {
				break
			}

			// Square is empty
			moves = append(moves, Move{startSquare, targetSquare, NormalMove})

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
		targetPiece := position.GetPieceAt(targetSquare)
		targetColor := targetPiece.Player()
		// Empty or opponent piece
		if targetPiece == NONE || !IsSameColor(targetColor, color) {
			moves = append(moves, Move{startSquare, targetSquare, NormalMove})
		}
	}

	// Castle moves:
	kingSide := WhiteKingSide
	queenSide := WhiteQueenSide

	if color == BLACK.Player() {
		kingSide = BlackKingSide
		queenSide = BlackQueenSide
	}

	addCastleMove := func(side CastleRights, direction int, sideFlag MoveFlag, emptyOffsets ...int) {
		// Check if we are allowed to castle
		if position.CastleRights&side == 0 {
			return
		}
		// Ensure no piece is blocking the castle
		for _, offset := range emptyOffsets {
			if position.GetPieceAt(Square(int(startSquare)+offset)) != NONE {
				return
			}
		}
		// since castling is valid we dont need to check if OOB
		targetSquare := Square(int(startSquare) + direction)
		moves = append(moves, Move{startSquare, targetSquare, sideFlag})
	}
	addCastleMove(kingSide, 2, KingCastle, 1, 2)
	addCastleMove(queenSide, -2, QueenCastle, -1, -2, -3)
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

	singleMoveTargetPiece := position.GetPieceAt(singleMoveSquare)
	if singleMoveTargetPiece == NONE {
		moves = append(moves, Move{startSquare, singleMoveSquare, NormalMove})
		// Also look for double pawn moves here, since no piece was infront
		if IsStartPawnRank(startSquare, color) {
			doubleMoveSquare, err := MoveDirection(startSquare, doubleForwards)
			if err != nil {
				panic(err) // should never be out of bounds for a 2x move. Just in case...
			}
			doubleMoveTargetPiece := position.GetPieceAt(doubleMoveSquare)
			if doubleMoveTargetPiece == NONE {
				moves = append(moves, Move{startSquare, doubleMoveSquare, DoublePawnPush})
			}
		}
	}
	var attackLeft, attackRight Direction
	if color == WHITE.Player() {
		attackLeft = Direction{1, -1}
		attackRight = Direction{1, 1}
	} else {
		attackLeft = Direction{-1, -1}
		attackRight = Direction{-1, 1}
	}

	attackSide := func(square Square, direction Direction) []Move {
		squareToAttack, err := MoveDirection(square, direction)
		if err != nil {
			return moves
		}
		// check for normal attack
		piece := position.GetPieceAt(squareToAttack)
		pieceColor := piece.Player()
		if piece != NONE && !IsSameColor(color, pieceColor) {
			moves = append(moves, Move{startSquare, squareToAttack, NormalMove})
		}
		// check for en-passant
		if piece == NONE && squareToAttack == position.PossibleEnPassantCapture {
			moves = append(moves, Move{startSquare, squareToAttack, EnPassantCapture})
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
		piece := position.GetPieceAt(targetSquare)
		targetColor := piece.Player()
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
	from := move.From
	to := move.To
	moveFlag := move.Flag

	movingPiece := p.GetPieceAt(from)
	capturedPiece := p.GetPieceAt(to)

	// capture the current board state so UnmakeMove can restore it
	undoMoveState := UndoMoveState{
		CapturedPiece:            capturedPiece,
		CastleRights:             p.CastleRights,
		PossibleEnPassantCapture: p.PossibleEnPassantCapture,
	}

	// Move pieces around
	p.SetPieceAt(movingPiece, to)
	p.SetPieceAt(NONE, from)

	// TODO: handle special operations -- first rook/king move invalidates castle etc

	// Update the position state
	// Clear the en-passant move (will be reset if the move was double pawn)
	p.PossibleEnPassantCapture = NoSquare
	switch moveFlag {
	case DoublePawnPush:
		// the square behind the pawn is now a possible en-passant capture
		enPassantCaptureSquare := to - 8
		if movingPiece.Player() == BLACK.Player() {
			enPassantCaptureSquare = to + 8
		}
		p.PossibleEnPassantCapture = enPassantCaptureSquare
	case EnPassantCapture:
		// if capturing en-passant, remove the piece behind the destination square...
		// NOTE: UnmakeMove is responsible for reinstating the captured pawn
		capturedEnPassantSquare := to - 8
		if movingPiece.Player() == BLACK.Player() {
			capturedEnPassantSquare = to + 8
		}
		p.SetPieceAt(NONE, capturedEnPassantSquare)

	// Move the rook -- king has already moved: (flags are updated below)
	case KingCastle:
		rookFrom := to + 1
		rookTo := to - 1
		rook := p.GetPieceAt(rookFrom)
		p.SetPieceAt(rook, rookTo)
		p.SetPieceAt(NONE, rookFrom)

	case QueenCastle:
		rookFrom := to - 2
		rookTo := to + 1
		rook := p.GetPieceAt(rookFrom)
		p.SetPieceAt(rook, rookTo)
		p.SetPieceAt(NONE, rookFrom)

	// TODO: Allow promotions as well!
	case PromoteBishop:
	case PromoteKnight:
	case PromoteRook:
	case PromoteQueen:
	}

	// Update castling rights when king/rook moves
	switch movingPiece.Type() {
	case KING:
		if movingPiece.Player() == WHITE.Player() {
			p.CastleRights &^= WhiteKingSide | WhiteQueenSide
		} else {
			p.CastleRights &^= BlackKingSide | BlackQueenSide
		}
	case ROOK:
		switch from {
		// A1
		case 0:
			p.CastleRights &^= WhiteQueenSide
		// H1
		case 7:
			p.CastleRights &^= WhiteKingSide
		// A8
		case 56:
			p.CastleRights &^= BlackQueenSide
		// H8
		case 63:
			p.CastleRights &^= BlackKingSide
		}
	}

	if capturedPiece.Type() == ROOK {
		switch to {
		// A1
		case 0:
			p.CastleRights &^= WhiteQueenSide
		// H1
		case 7:
			p.CastleRights &^= WhiteKingSide
		// A8
		case 56:
			p.CastleRights &^= BlackQueenSide
		// H8
		case 63:
			p.CastleRights &^= BlackKingSide
		}
	}

	// TODO: Flip player to move

	return undoMoveState

}

func UnmakeMove(p *Position, move Move, undo UndoMoveState) {

	// TODO: unmake the move.
	// Remember to undo any side-effects making the move had
	// such as en-passant squares etc
	// in particular, if the capture is an EnPassantCapture, reinstate the captured piece
}
