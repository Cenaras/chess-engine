package fen

import (
	"chess/game"
	"fmt"
	"os"
	"strings"
)

const StartingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func LoadFenPosition(fen string) game.Position {
	position := parseFen(fen)
	return position
}

func parseFen(fen string) game.Position {
	position := game.Position{}

	// Parse FEN into parts, separated by spaces
	parts := strings.Fields(fen)

	// Set up position from parts[0]
	index := 63
	for _, ch := range parts[0] {
		switch ch {
		case 'r':
			position.Board[index] = game.ROOK | game.BLACK
			index--
		case 'n':
			position.Board[index] = game.KNIGHT | game.BLACK
			index--
		case 'b':
			position.Board[index] = game.BISHOP | game.BLACK
			index--
		case 'q':
			position.Board[index] = game.QUEEN | game.BLACK
			index--
		case 'k':
			position.Board[index] = game.KING | game.BLACK
			index--
		case 'p':
			position.Board[index] = game.PAWN | game.BLACK
			index--
		case 'R':
			position.Board[index] = game.ROOK | game.WHITE
			index--
		case 'N':
			position.Board[index] = game.KNIGHT | game.WHITE
			index--
		case 'B':
			position.Board[index] = game.BISHOP | game.WHITE
			index--
		case 'Q':
			position.Board[index] = game.QUEEN | game.WHITE
			index--
		case 'K':
			position.Board[index] = game.KING | game.WHITE
			index--
		case 'P':
			position.Board[index] = game.PAWN | game.WHITE
			index--
		case '/':
			continue // NO OP as we are 1D array
		case '1', '2', '3', '4', '5', '6', '7', '8':
			emptySquares := int(ch - '0') // 'a' - '0' = a
			index -= emptySquares
		default:
			panic(fmt.Sprintf("Invalid FEN: couldn't parse positional character %c", ch))

		}
	}

	// player in turn from parts[1]
	activeColor := parts[1]
	switch activeColor {
	case "w":
		position.PlayerToMove = game.Player(game.WHITE)
	case "b":
		position.PlayerToMove = game.Player(game.BLACK)
	default:
		panic("Invalid FEN: Player to move must be either \"w\" or \"b\"")
	}

	// Where can both players castle to?
	castlingRights := parts[2]
	for _, ch := range castlingRights {
		switch ch {
		case 'K':
			position.CastleRights |= game.WhiteKingSide
		case 'Q':
			position.CastleRights |= game.WhiteQueenSide
		case 'k':
			position.CastleRights |= game.BlackKingSide
		case 'q':
			position.CastleRights |= game.BlackQueenSide
		case '-':
			position.CastleRights = 0
		}
	}

	// What pawn can be captured en-passant?
	enPassantTarget := parts[3]
	if enPassantTarget == "-" {
		position.PossibleEnPassantCapture = 0
	} else if len(enPassantTarget) != 2 {
		panic(fmt.Sprintf("Invalid FEN en passant square: %s", fen))
	} else {
		file := enPassantTarget[0]
		rank := enPassantTarget[1]
		if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
			panic(fmt.Sprintf("invalid fen en passant square: %s", fen))
		}
		fileIndex := int(file - 'a')
		rankIndex := int('8' - rank)
		index := rankIndex*8 + fileIndex
		position.PossibleEnPassantCapture = game.Square(index)
	}

	// TODO: Fix half and full move clocks
	// halfMoveClock := parts[4]
	// fullMoveClock := parts[5]

	return position
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}
