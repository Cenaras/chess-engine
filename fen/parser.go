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
	position := game.Position{
		PossibleEnPassantCapture: game.NoSquare,
	}

	parts := strings.Fields(fen)

	if len(parts) != 6 {
		panic(fmt.Sprintf("invalid FEN: expected 6 fields, got %d", len(parts)))
	}

	// ------------------------------------------------------------
	// Piece placement
	// A1 = 0, H8 = 63
	// ------------------------------------------------------------

	rank := 7
	file := 0

	for _, ch := range parts[0] {
		switch {
		case ch == '/':
			if file != 8 {
				panic("invalid FEN: rank does not contain 8 squares")
			}

			rank--
			file = 0

		case ch >= '1' && ch <= '8':
			file += int(ch - '0')

			if file > 8 {
				panic("invalid FEN: too many squares in rank")
			}

		default:
			var piece game.Piece

			switch ch {
			case 'r':
				piece = game.ROOK | game.BLACK
			case 'n':
				piece = game.KNIGHT | game.BLACK
			case 'b':
				piece = game.BISHOP | game.BLACK
			case 'q':
				piece = game.QUEEN | game.BLACK
			case 'k':
				piece = game.KING | game.BLACK
				position.BlackKingSquare = game.Square(rank*8 + file)
			case 'p':
				piece = game.PAWN | game.BLACK

			case 'R':
				piece = game.ROOK | game.WHITE
			case 'N':
				piece = game.KNIGHT | game.WHITE
			case 'B':
				piece = game.BISHOP | game.WHITE
			case 'Q':
				piece = game.QUEEN | game.WHITE
			case 'K':
				piece = game.KING | game.WHITE
				position.WhiteKingSquare = game.Square(rank*8 + file)
			case 'P':
				piece = game.PAWN | game.WHITE

			default:
				panic(fmt.Sprintf(
					"invalid FEN: couldn't parse positional character %c",
					ch,
				))
			}

			if file >= 8 || rank < 0 {
				panic("invalid FEN: invalid piece placement")
			}

			square := game.Square(rank*8 + file)
			position.Board[square] = piece
			file++
		}
	}

	if rank != 0 || file != 8 {
		panic("invalid FEN: piece placement does not describe an 8x8 board")
	}

	// ------------------------------------------------------------
	// Player to move
	// ------------------------------------------------------------

	switch parts[1] {
	case "w":
		position.PlayerToMove = game.Player(game.WHITE)
	case "b":
		position.PlayerToMove = game.Player(game.BLACK)
	default:
		panic(`invalid FEN: player to move must be either "w" or "b"`)
	}

	// ------------------------------------------------------------
	// Castling rights
	// ------------------------------------------------------------

	for _, ch := range parts[2] {
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
			// No castling rights.
		default:
			panic(fmt.Sprintf(
				"invalid FEN castling character: %c",
				ch,
			))
		}
	}

	// ------------------------------------------------------------
	// En passant target
	// ------------------------------------------------------------

	enPassantTarget := parts[3]

	if enPassantTarget != "-" {
		if len(enPassantTarget) != 2 {
			panic(fmt.Sprintf(
				"invalid FEN en passant square: %s",
				enPassantTarget,
			))
		}

		fileChar := enPassantTarget[0]
		rankChar := enPassantTarget[1]

		if fileChar < 'a' || fileChar > 'h' ||
			rankChar < '1' || rankChar > '8' {
			panic(fmt.Sprintf(
				"invalid FEN en passant square: %s",
				enPassantTarget,
			))
		}

		fileIndex := int(fileChar - 'a')
		rankIndex := int(rankChar - '1')

		position.PossibleEnPassantCapture =
			game.Square(rankIndex*8 + fileIndex)
	}

	// TODO: halfmove/fullmove clocks
	// parts[4]
	// parts[5]

	return position
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}
