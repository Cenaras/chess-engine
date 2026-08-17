package game

func parsePieceTable(values [64]int) [64]int {
	var table [64]int

	// Input is written in chess-board display order:
	//
	//   a8 b8 ... h8
	//   a7 b7 ... h7
	//   ...
	//   a1 b1 ... h1
	//
	// Our internal square representation is:
	//
	//   A1 = 0, B1 = 1, ... H8 = 63
	//
	// Therefore we only need to vertically flip the ranks.
	for i, value := range values {
		inputRank := i / 8 // 0 = rank 8, 7 = rank 1
		file := i % 8

		internalRank := 7 - inputRank // 0 = rank 1, 7 = rank 8

		square := internalRank*8 + file
		table[square] = value
	}

	return table
}

var pawnPieceTable = parsePieceTable([64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
})

var bishopPieceTable = parsePieceTable([64]int{
	-20, -10, -10, -10, -10, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 5, 5, 10, 10, 5, 5, -10,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-10, 10, 10, 10, 10, 10, 10, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-20, -10, -10, -10, -10, -10, -10, -20,
})

// From whites point of view
var knightPieceTable = parsePieceTable([64]int{
	-50, -40, -30, -30, -30, -30, -40, -50, // rank 8
	-40, -20, 0, 0, 0, 0, -20, -40, // rank 7
	-30, 0, 10, 15, 15, 10, 0, -30, // rank 6
	-30, 5, 15, 20, 20, 15, 5, -30, // rank 5
	-30, 0, 15, 20, 20, 15, 0, -30, // rank 4
	-30, 5, 10, 15, 15, 10, 5, -30, // rank 3
	-40, -20, 0, 5, 5, 0, -20, -40, // rank 2
	-50, -40, -30, -30, -30, -30, -40, -50, // rank 1
})

var rookPieceTable = parsePieceTable([64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 10, 10, 10, 10, 10, 5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	0, 0, 0, 5, 5, 0, 0, 0,
})

var queenPieceTable = parsePieceTable([64]int{
	-20, -10, -10, -5, -5, -10, -10, -20,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-5, 0, 5, 5, 5, 5, 0, -5,
	0, 0, 5, 5, 5, 5, 0, -5,
	-10, 5, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 0, 0, 0, 0, -10,
	-20, -10, -10, -5, -5, -10, -10, -20,
})

// TODO Nothing for king yet :) Also maybe update these later on in the game?

// Notice that ^56 flips the ranks horizontally
func mirrorSquare(square Square) Square {
	return square ^ 56
}

func pieceTableValue(table [64]int, square Square, player Player) int {
	if player == BLACK.Player() {
		square = mirrorSquare(square)
	}
	return table[square]
}
