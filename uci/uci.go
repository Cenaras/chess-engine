package uci

// See commands: https://publish.obsidian.md/modern-uci-doc/UCI+Docs/Commands/Commands

import (
	"bufio"
	"chess/fen"
	"chess/game"
	"fmt"
	"os"
	"strings"
)

const EngineName string = "CenEngine"

func StartUCI() {

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if strings.HasPrefix(text, string(Position)) {
			inputPosition(text)
		}

		switch text {
		case string(UCI):
			inputUCI()
		case string(IsReady):
			inputIsReady()

		}
	}
}

func inputUCI() {
	fmt.Printf("id name %s\n", EngineName)
	fmt.Println("id author Cenaras")

}
func inputIsReady() {
	// TODO: Start precomputation's here! Wait for them to return.
	fmt.Println("readyok")

}
func inputPosition(input string) {
	// position [fen <fenstring> | startpos ] moves <move 1>, ..., <move n>
	cmd := strings.Split(input, " ")
	firstArg := cmd[1]
	var position game.Position

	switch firstArg {
	case "fen":
		panic("Not implemented")
		// fen := cmd[2]
		// parse the fen
	case "startpos":
		position = fen.LoadFenPosition(fen.StartingFEN)
	}
	if !(cmd[2] == "moves") {
		panic("arbitrary fen not supported")
	}

	// Parse all the moves in the list after "moves" command
	for i := 3; i < len(cmd); i++ {
		moveNotation := cmd[i]
		move, err := FindMoveFromUCI(&position, moveNotation)
		if err != nil {
			panic("invalid uci move")
		}
		game.MakeMove(&position, move)
	}

}
func inputGo() {}

// Finds the move from game, corresponding to the algebraic notation string.
// As all moves are valid, no verification needs to be done besides the notation structure
func FindMoveFromUCI(position *game.Position, notation string) (game.Move, error) {
	if len(notation) != 4 || len(notation) != 5 {
		return game.Move{}, fmt.Errorf("invalid UCI move %q", notation)
	}

	from, errFrom := AlgebraicToSquare(notation[0:2])
	if errFrom != nil {
		return game.Move{}, errFrom
	}

	to, errTo := AlgebraicToSquare(notation[2:4])
	if errTo != nil {
		return game.Move{}, errTo
	}
	moves := game.GenerateMoves(position)

	for _, move := range moves {
		if move.From != from || move.To != to {
			continue
		}

		// Let us validate for sanity check, that the move isn't a promotion move
		// assuming valid UCI input (e.g., not c7c8), this isn't needed.
		if len(notation) == 4 {
			// A non-promotion move.
			if !isPromotion(move.Flag) {
				return move, nil
			}

			continue
		}

		// Promotion move, e.g. c7c8q.
		if promotionMatches(move.Flag, notation[4]) {
			return move, nil
		}
	}

	return game.Move{}, fmt.Errorf("move %q is not legal in the current position", notation)
}
func isPromotion(flag game.MoveFlag) bool {
	switch flag {
	case game.PromoteKnight, game.PromoteBishop, game.PromoteRook, game.PromoteQueen:
		return true
	default:
		return false
	}
}

func promotionMatches(flag game.MoveFlag, promotion byte) bool {
	switch promotion {
	case 'n', 'N':
		return flag == game.PromoteKnight
	case 'b', 'B':
		return flag == game.PromoteBishop
	case 'r', 'R':
		return flag == game.PromoteRook
	case 'q', 'Q':
		return flag == game.PromoteQueen
	default:
		return false
	}
}

func AlgebraicToSquare(s string) (game.Square, error) {
	if len(s) != 2 {
		return game.NO_SQUARE, fmt.Errorf("invalid square: %q", s)
	}

	file := s[0]
	rank := s[1]

	if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
		return game.NO_SQUARE, fmt.Errorf("invalid square: %q", s)
	}

	return game.RankFileToSquare(
		int(rank-'1'),
		int(file-'a'),
	), nil
}
