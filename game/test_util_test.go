package game_test

import (
	"chess/game"
	"testing"
)

func sq(s string) game.Square {
	file := int(s[0] - 'a')
	rank := int(s[1] - '1')
	return game.Square(rank*8 + file)
}

func hasMove(
	moves []game.Move,
	from string,
	to string,
	flag game.MoveFlag,
) bool {
	fromSq := sq(from)
	toSq := sq(to)

	for _, move := range moves {
		if move.From == fromSq &&
			move.To == toSq &&
			move.Flag == flag {
			return true
		}
	}

	return false
}

func requireMove(
	t *testing.T,
	moves []game.Move,
	from string,
	to string,
	flag game.MoveFlag,
) {
	t.Helper()

	if !hasMove(moves, from, to, flag) {
		t.Fatalf(
			"expected move %s -> %s with flag %v",
			from,
			to,
			flag,
		)
	}
}

func requireNoMove(
	t *testing.T,
	moves []game.Move,
	from string,
	to string,
	flag game.MoveFlag,
) {
	t.Helper()

	if hasMove(moves, from, to, flag) {
		t.Fatalf(
			"did not expect move %s -> %s with flag %v",
			from,
			to,
			flag,
		)
	}
}
