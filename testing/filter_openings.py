import argparse
import json

import chess
import chess.engine
import chess.pgn


MAX_CP = 50
STOCKFISH_NODES = 200_000


def load_pgn_openings(path):
    openings = []

    with open(path, "r", encoding="utf-8") as pgn:
        while True:
            game = chess.pgn.read_game(pgn)

            if game is None:
                break

            board = game.board()
            moves = []

            for move in game.mainline_moves():
                moves.append(move.uci())
                board.push(move)

            openings.append({
                "name": game.headers.get("Opening", ""),
                "moves": moves,
            })

    return openings


def setup_board_with_opening(opening):
    board = chess.Board()

    for move_str in opening["moves"]:
        move = chess.Move.from_uci(move_str)

        if move not in board.legal_moves:
            raise ValueError(
                f'Opening "{opening["name"]}" contained illegal move: {move}'
            )

        board.push(move)

    return board


def evaluate_opening(stockfish, opening):
    board = setup_board_with_opening(opening)

    if board.is_game_over(claim_draw=True):
        return None

    info = stockfish.analyse(
        board,
        chess.engine.Limit(nodes=STOCKFISH_NODES),
    )

    score = info["score"].pov(chess.WHITE)

    # We don't want openings where Stockfish sees a forced mate.
    if score.is_mate():
        return None

    return score.score()


def main():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "input",
        help="Path to 8moves_v3.pgn",
    )

    parser.add_argument(
        "stockfish",
        help="Path to Stockfish executable",
    )

    parser.add_argument(
        "output",
        help="Path for filtered JSON output",
    )

    args = parser.parse_args()

    openings = load_pgn_openings(args.input)

    print(f"Loaded {len(openings)} openings")

    stockfish = chess.engine.SimpleEngine.popen_uci(args.stockfish)

    balanced = []

    try:
        for i, opening in enumerate(openings):
            cp = evaluate_opening(stockfish, opening)

            if cp is not None and abs(cp) <= MAX_CP:
                opening["evaluation_cp"] = cp
                balanced.append(opening)

            if (i + 1) % 100 == 0:
                print(
                    f"{i + 1}/{len(openings)} processed, "
                    f"{len(balanced)} accepted"
                )

    finally:
        stockfish.quit()

    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(
            balanced,
            f,
            indent=2,
        )

    print()
    print(f"Accepted {len(balanced)} / {len(openings)} openings")
    print(f"Saved to {args.output}")


if __name__ == "__main__":
    main()