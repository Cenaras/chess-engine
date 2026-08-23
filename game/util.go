package game

import (
	"fmt"
)

func MoveToAlgebraic(move Move) string {
	fromRank, fromFile := SquareToRankFile(move.From)
	toRank, toFile := SquareToRankFile(move.To)

	notation := fmt.Sprintf(
		"%c%c%c%c",
		'a'+fromFile,
		'1'+fromRank,
		'a'+toFile,
		'1'+toRank,
	)

	if notation == "a1a1" {
		return "0000"
	}

	switch move.Flag {
	case PromoteKnight:
		notation += "n"
	case PromoteBishop:
		notation += "b"
	case PromoteRook:
		notation += "r"
	case PromoteQueen:
		notation += "q"
	}

	return notation
}

func AlgebraicToSquare(s string) (Square, error) {
	if len(s) != 2 {
		return NO_SQUARE, fmt.Errorf("invalid square: %q", s)
	}

	file := s[0]
	rank := s[1]

	if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
		return NO_SQUARE, fmt.Errorf("invalid square: %q", s)
	}

	return RankFileToSquare(
		int(rank-'1'),
		int(file-'a'),
	), nil
}
