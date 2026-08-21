package uci

// See commands: https://publish.obsidian.md/modern-uci-doc/UCI+Docs/Commands/Commands

import (
	"bufio"
	"chess/fen"
	"chess/game"
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// The engine holds the position, and the ability to cancel the search
type Engine struct {
	Position     game.Position
	cancelSearch context.CancelFunc
}

// TODO: Make these methods on the engine instead
var engine Engine = Engine{}

const EngineName string = "CenEngine"

func StartUCI() {

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, string(Position)) {
			handlePosition(text)
		}
		if strings.HasPrefix(text, string(Go)) {
			handleGo(text)
		}

		switch text {
		case string(UCI):
			handleUCI()
		case string(IsReady):
			handleIsReady()
		case string(Quit):
			handleQuit()

		}
	}
}

func handleUCI() {
	fmt.Printf("id name %s\n", EngineName)
	fmt.Println("id author Cenaras")
	// TODO: Any engine options?
	fmt.Println("uciok")

}
func handleIsReady() {
	// TODO: Start precomputation's here! Wait for them to return.
	fmt.Println("readyok")
}

func parseGoCommand(input string) game.SearchOptions {
	var options game.SearchOptions = game.SearchOptions{}
	cmd := strings.Fields(input)
	getIntOption := func(option string) int {
		idx := slices.Index(cmd, option)
		if idx == -1 {
			return 0
		}

		if idx+1 >= len(cmd) {
			panic(fmt.Sprintf("missing value for %s", option))
		}

		val, err := strconv.Atoi(cmd[idx+1])
		if err != nil {
			panic(fmt.Sprintf("invalid value for %s: %q", option, cmd[idx+1]))
		}

		return val
	}

	options.WhiteTime = time.Duration(getIntOption("wtime")) * time.Millisecond
	options.BlackTime = time.Duration(getIntOption("btime")) * time.Millisecond

	options.WhiteIncrement = time.Duration(getIntOption("winc")) * time.Millisecond
	options.WhiteIncrement = time.Duration(getIntOption("binc")) * time.Millisecond

	options.Depth = getIntOption("depth")
	return options
}

func handleGo(input string) {
	options := parseGoCommand(input)

	// TODO: more/less time / calulate time based on the remainind time
	var searchTime = 1000 * time.Millisecond

	// create a cancel context, to end the seach.
	ctx, cancel := context.WithTimeout(context.Background(), searchTime)
	engine.cancelSearch = cancel

	// Create a copy to avoid any nasty race conditions: Since
	// searching for the best move will modify the state
	position := engine.Position

	// Start search in a new thread to avoid blocking the UCI interpreter.
	go func() {
		defer cancel()
		bestMove := game.FindBestMove(&position, options, ctx)
		notation := MoveToAlgebraic(bestMove)
		fmt.Printf("bestmove %s\n", notation)
	}()
}

func handlePosition(input string) {
	// position [fen <fenstring> | startpos ] moves <move 1>, ..., <move n>
	cmd := strings.Fields(input)
	if len(cmd) < 2 {
		panic("invalid")
	}
	var next int

	switch cmd[1] {
	case "startpos":
		engine.Position = fen.LoadFenPosition(fen.StartingFEN)
		next = 2
	case "fen":
		if len(cmd) < 8 {
			panic("invalid")
		}
		engine.Position = fen.LoadFenPosition(strings.Join(cmd[2:8], " "))
		next = 8
	default:
		panic("invalid")
	}

	if next == len(cmd) {
		return
	}

	if cmd[next] != "moves" {
		panic("expected moves !")
	}

	// Parse all the moves in the list after "moves" command
	for _, notation := range cmd[next+1:] {
		move, err := FindMoveFromUCI(&engine.Position, notation)
		if err != nil {
			panic(fmt.Sprintf("invalid uci move %s", err.Error()))

		}
		game.MakeMove(&engine.Position, move)
	}
}

// Finds the move from game, corresponding to the algebraic notation string.
// As all moves are valid, no verification needs to be done besides the notation structure
func FindMoveFromUCI(position *game.Position, notation string) (game.Move, error) {
	if len(notation) != 4 && len(notation) != 5 {
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

func MoveToAlgebraic(move game.Move) string {
	fromRank, fromFile := game.SquareToRankFile(move.From)
	toRank, toFile := game.SquareToRankFile(move.To)

	notation := fmt.Sprintf(
		"%c%c%c%c",
		'a'+fromFile,
		'1'+fromRank,
		'a'+toFile,
		'1'+toRank,
	)

	switch move.Flag {
	case game.PromoteKnight:
		notation += "n"
	case game.PromoteBishop:
		notation += "b"
	case game.PromoteRook:
		notation += "r"
	case game.PromoteQueen:
		notation += "q"
	}

	return notation
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
func handleQuit() {
	// If the
	if engine.cancelSearch != nil {
		engine.cancelSearch()
	}
}
