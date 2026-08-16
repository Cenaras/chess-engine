package game

var rookDirections = []Direction{
	// S, N, W, E
	{-1, 0}, {1, 0}, {0, -1}, {0, 1},
}
var bishopDirections = []Direction{
	// NW, NE,  SE, SW
	{1, -1}, {1, 1}, {-1, 1}, {-1, -1},
}

// Precomputes, for each square on the board, all possible target squares
// according to the direction offsets
func precomputeMoves(dirOffsets []Direction) [64][]Square {
	var moves [64][]Square
	for idx := 0; idx < 64; idx++ {
		square := Square(idx)
		for _, dir := range dirOffsets {
			targetSquare, success := MoveDirection(square, dir)
			if !success {
				continue
			}
			moves[idx] = append(moves[idx], targetSquare)
		}
	}
	return moves
}

var kingDirection = []Direction{
	{1, 0}, {1, 1}, {1, -1},
	{0, 1}, {0, -1},
	{-1, 0}, {-1, 1}, {-1, -1},
}
var kingMoves [64][]Square = precomputeMoves(kingDirection)
var knightOffsets = []Direction{
	// NNW, NNE,
	{2, -1}, {2, 1},
	// NWW, NEE,
	{1, -2}, {1, 2},
	// SWW, SEE,
	{-1, -2}, {-1, 2},
	// SSW, SSE
	{-2, -1}, {-2, 1},
}
var knightMoves [64][]Square = precomputeMoves(knightOffsets)

var whitePawnForwardMoves = precomputeMoves([]Direction{{1, 0}})
var blackPawnForwardMoves = precomputeMoves([]Direction{{-1, 0}})
var whitePawnDoubleMoves = precomputeMoves([]Direction{{2, 0}})
var blackPawnDoubleMoves = precomputeMoves([]Direction{{-2, 0}})
var whitePawnAttackMoves = precomputeMoves([]Direction{{1, -1}, {1, 1}})
var blackPawnAttackMoves = precomputeMoves([]Direction{{-1, -1}, {-1, 1}})

// Directions to check if a pawn attacks our square (used by isSquareAttacked)
var whitePawnAttackers = precomputeMoves([]Direction{{-1, -1}, {-1, 1}})
var blackPawnAttackers = precomputeMoves([]Direction{{1, -1}, {1, 1}})

// Generate all pseudo legal moves for the current position
func pseudoLegalMove(position *Position) []Move {
	// loop every single piece, compute all its pseudolegal moves
	moves := make([]Move, 0, 64) // preallocate 64
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
			targetSquare, success := MoveDirection(currentSquare, direction)
			// Moving further in this direction leaves the board
			if !success {
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
	// Regular moves
	for _, targetSquare := range kingMoves[startSquare] {
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

func appendPawnMove(moves []Move, from Square, to Square, color Player, flag MoveFlag) []Move {
	if IsPromotionRank(to, color) {
		return append(
			moves,
			Move{from, to, PromoteQueen},
			Move{from, to, PromoteRook},
			Move{from, to, PromoteBishop},
			Move{from, to, PromoteKnight},
		)
	}
	return append(moves, Move{from, to, flag})
}

func genPawnMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	// Get forwards precomputed moves based on color
	var forwardTargets []Square
	var doubleTargets []Square
	var attackTargets []Square

	if color == WHITE.Player() {
		forwardTargets = whitePawnForwardMoves[startSquare]
		doubleTargets = whitePawnDoubleMoves[startSquare]
		attackTargets = whitePawnAttackMoves[startSquare]
	} else {
		forwardTargets = blackPawnForwardMoves[startSquare]
		doubleTargets = blackPawnDoubleMoves[startSquare]
		attackTargets = blackPawnAttackMoves[startSquare]
	}

	// Single forward move.
	if len(forwardTargets) > 0 {
		targetSquare := forwardTargets[0]

		if position.GetPieceAt(targetSquare) == NONE {
			moves = appendPawnMove(moves, startSquare, targetSquare, color, NormalMove)

			if IsStartPawnRank(startSquare, color) &&
				len(doubleTargets) > 0 {

				doubleTarget := doubleTargets[0]

				if position.GetPieceAt(doubleTarget) == NONE {
					moves = append(moves, Move{startSquare, doubleTarget, DoublePawnPush})
				}
			}
		}
	}

	for _, targetSquare := range attackTargets {
		piece := position.GetPieceAt(targetSquare)

		if piece != NONE && piece.Player() != color {
			moves = appendPawnMove(moves, startSquare, targetSquare, color, NormalMove)
		}

		if piece == NONE &&
			targetSquare == position.PossibleEnPassantCapture {
			moves = append(moves, Move{startSquare, targetSquare, EnPassantCapture})
		}
	}
	return moves
}

func genKnightMoves(startSquare Square, color Player, moves []Move, position *Position) []Move {
	for _, targetSquare := range knightMoves[startSquare] {
		piece := position.GetPieceAt(targetSquare)
		targetColor := piece.Player()
		if !IsSameColor(color, targetColor) {
			moves = append(moves, Move{startSquare, targetSquare, NormalMove})
		}
	}
	return moves
}

/*
	Optimizations:
	 - Bitboards and attack masks
*/

// Generate all legal moves for the position
func GenerateMoves(position *Position) []Move {
	// TODO: Preallocate a 256 sized buffer and reuse that instead of allocating on the heap
	pseudo := pseudoLegalMove(position)
	legal := make([]Move, 0, len(pseudo))

	player := position.PlayerToMove
	opponent := player.Opponent()

	for _, move := range pseudo {
		if move.Flag == KingCastle || move.Flag == QueenCastle {
			// Cannot castle out of check.
			if isSquareAttackedBy(position, move.From, opponent) {
				continue
			}
			// Cannot castle through check.
			offset := 1
			if move.Flag == QueenCastle {
				offset = -1
			}
			throughSquare := Square(int(move.From) + offset)
			if isSquareAttackedBy(position, throughSquare, opponent) {
				continue
			}
		}

		undo := MakeMove(position, move)
		kingSquare := position.FindKing(player)

		// Handles normal king safety and the final castling square.
		kingAttacked := isSquareAttackedBy(
			position,
			kingSquare,
			opponent,
		)
		UnmakeMove(position, move, undo)
		if kingAttacked {
			continue
		}
		legal = append(legal, move)
	}
	return legal
}

// Calculate from the current position, if square is attacked by player
func isSquareAttackedBy(position *Position, square Square, attacker Player) bool {
	// TODO!!! This is pretty much the same logic as checking moves for sliding pieces -- REFACTOR

	scanRay := func(rays []Direction, rayType Piece) bool {
		for _, dir := range rays {
			currentSquare := square
			for {
				targetSquare, success := MoveDirection(currentSquare, dir)
				if !success {
					break
				}
				piece := position.GetPieceAt(targetSquare)
				pieceType := piece.Type()

				// Empty square, search outwards
				if pieceType == NONE {
					currentSquare = targetSquare
					continue
				}
				// found bishop/rook or queen owned by attacker that can see the square
				if piece.Player() == attacker && (pieceType == rayType || pieceType == QUEEN) {
					return true
				}
				// A piece that isn't attacking us blocks the ray.
				break
			}
		}
		return false
	}
	// Scan outwards for sliding pieces
	if scanRay(rookDirections, ROOK) || scanRay(bishopDirections, BISHOP) {
		return true
	}
	// Check knights
	for _, targetSquare := range knightMoves[square] {
		piece := position.GetPieceAt(targetSquare)
		if piece.Player() == attacker && piece.Type() == KNIGHT {
			return true
		}
	}

	// Check kings (also probably duplicate of king movement)
	for _, targetSquare := range kingMoves[square] {
		piece := position.GetPieceAt(targetSquare)
		if piece.Player() == attacker && piece.Type() == KING {
			return true
		}
	}

	// Check pawns
	var pawnAttackers []Square

	if attacker == WHITE.Player() {
		pawnAttackers = whitePawnAttackers[square]
	} else {
		pawnAttackers = blackPawnAttackers[square]
	}

	for _, attackerSquare := range pawnAttackers {
		piece := position.GetPieceAt(attackerSquare)

		if piece.Player() == attacker && piece.Type() == PAWN {
			return true
		}
	}
	return false
}

func MakeMove(p *Position, move Move) UndoMoveState {
	from := move.From
	to := move.To
	moveFlag := move.Flag

	movingPiece := p.GetPieceAt(from)
	capturedPiece := p.GetPieceAt(to)
	oldKingSquare := p.WhiteKingSquare
	if movingPiece.Player() == BLACK.Player() {
		oldKingSquare = p.BlackKingSquare
	}

	// capture the current board state so UnmakeMove can restore it
	undoMoveState := UndoMoveState{
		CapturedPiece:            capturedPiece,
		CastleRights:             p.CastleRights,
		PossibleEnPassantCapture: p.PossibleEnPassantCapture,
		OldKingSquare:            oldKingSquare,
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

	case PromoteBishop:
		p.SetPieceAt(BISHOP|Piece(movingPiece.Player()), to)
	case PromoteKnight:
		p.SetPieceAt(KNIGHT|Piece(movingPiece.Player()), to)
	case PromoteRook:
		p.SetPieceAt(ROOK|Piece(movingPiece.Player()), to)
	case PromoteQueen:
		p.SetPieceAt(QUEEN|Piece(movingPiece.Player()), to)
	}

	// Update castling rights when king/rook moves and track new king position
	switch movingPiece.Type() {
	case KING:
		if movingPiece.Player() == WHITE.Player() {
			p.CastleRights &^= WhiteKingSide | WhiteQueenSide
			p.WhiteKingSquare = move.To
		} else {
			p.CastleRights &^= BlackKingSide | BlackQueenSide
			p.BlackKingSquare = move.To
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

	// Flip player to move
	p.PlayerToMove = p.PlayerToMove.Opponent()
	return undoMoveState

}

func UnmakeMove(p *Position, move Move, undo UndoMoveState) {
	moveFlag := move.Flag
	// Move pieces back / reinstate captured piece
	movedPiece := p.GetPieceAt(move.To)
	pieceToRestore := movedPiece
	capturedPiece := undo.CapturedPiece

	// If capture was en-passant, restore captured pawn
	switch moveFlag {
	case EnPassantCapture:
		// Captured pawn is always 1 rank behind pawns move's to
		rowOffset := -8
		pawn := PAWN | BLACK
		if movedPiece.Player() == BLACK.Player() {
			rowOffset = 8
			pawn = PAWN | WHITE
		}
		capturedSquare := Square(int(move.To) + rowOffset)
		p.SetPieceAt(pawn, capturedSquare)

	case KingCastle:
		// Move the rook back
		rookNewSquare := move.To - 1
		rookOriginalSquare := move.To + 1
		rook := p.GetPieceAt(rookNewSquare)
		p.SetPieceAt(rook, rookOriginalSquare)
		p.SetPieceAt(NONE, rookNewSquare)
	case QueenCastle:
		// Move the rook back
		rookNewSquare := move.To + 1
		rookOriginalSquare := move.To - 2
		rook := p.GetPieceAt(rookNewSquare)
		p.SetPieceAt(rook, rookOriginalSquare)
		p.SetPieceAt(NONE, rookNewSquare)

	// Remove promoted piece, restore pawn. Since the piece at the To square
	// is now the promoted piece, we should ensure that a pawn is placed
	case PromoteKnight, PromoteBishop, PromoteRook, PromoteQueen:
		pieceToRestore = PAWN | Piece(movedPiece.Player())
	}
	// Restore the actual move
	p.SetPieceAt(pieceToRestore, move.From)
	p.SetPieceAt(capturedPiece, move.To)

	// Restore state
	p.CastleRights = undo.CastleRights
	p.PossibleEnPassantCapture = undo.PossibleEnPassantCapture
	if movedPiece.Player() == WHITE.Player() {
		p.WhiteKingSquare = undo.OldKingSquare
	} else {
		p.BlackKingSquare = undo.OldKingSquare
	}

	p.PlayerToMove = p.PlayerToMove.Opponent()
}

// TODO: unmake the move.
// Remember to undo any side-effects making the move had
// such as en-passant squares etc
// in particular, if the capture is an EnPassantCapture, reinstate the captured piece
