import argparse
import json
import random

import chess
import chess.engine

DEBUG = True


def load_openings(path):
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)

    openings = []

    for entry in data:
        moves = entry.get("moves", [])

        if not moves:
            continue

        openings.append({
            "name": entry.get("name", ""),
            "moves": moves,
            "evaluation_cp": entry.get("evaluation_cp"),
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


def play_game(white, black, movetime, game_id, opening):
    board = setup_board_with_opening(opening)
    starting_ply = board.ply()

    while not board.is_game_over(claim_draw=True):
        engine = white if board.turn == chess.WHITE else black
        engine_name = "WHITE" if board.turn == chess.WHITE else "BLACK"

        if DEBUG:
            print()
            print(f"=== Ply {board.ply() + 1} ===")
            print(f"Engine: {engine_name}")
            print(f"FEN: {board.fen()}")
            # print(f"Moves: {' '.join(m.uci() for m in board.move_stack)}")

        try:
            result = engine.play(
                board,
                chess.engine.Limit(time=movetime),
                game=game_id,
            )
        except Exception:
            print()
            print("ENGINE FAILED")
            print(f"Side: {engine_name}")
            print(f"FEN: {board.fen()}")
            print(
                "Moves:",
                " ".join(m.uci() for m in board.move_stack),
            )
            raise

        if result.move is None:
            raise RuntimeError(
                "Engine returned no move in a non-terminal position"
            )

        board.push(result.move)

    played_plies = board.ply() - starting_ply

    return board.result(claim_draw=True), played_plies


def update_score(result, white_index, black_index, score):
    if result == "1-0":
        score[white_index] += 1
    elif result == "0-1":
        score[black_index] += 1
    else:
        score[0] += 0.5
        score[1] += 0.5


def main():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "engine1",
        help="Path to first UCI engine",
    )

    parser.add_argument(
        "engine2",
        help="Path to second UCI engine",
    )

    parser.add_argument(
        "games",
        type=int,
        help="Number of games to play",
    )

    parser.add_argument(
        "openings",
        help="Path to JSON opening database",
    )

    parser.add_argument(
        "--movetime-ms",
        type=int,
        default=100,
        help="Time per move in milliseconds",
    )

    parser.add_argument(
        "--seed",
        type=int,
        default=0,
        help="seed for randomness"
    )

    args = parser.parse_args()
    random.seed(args.seed)

    if args.games % 2 != 0:
        raise ValueError(
            "Number of games must be even because openings are played in pairs"
        )

    movetime = args.movetime_ms / 1000.0

    openings = load_openings(args.openings)

    if not openings:
        raise ValueError("No valid openings found")

    engine1 = chess.engine.SimpleEngine.popen_uci(args.engine1)
    engine2 = chess.engine.SimpleEngine.popen_uci(args.engine2)

    score = [0.0, 0.0]

    try:
        game_number = 0

        # Every opening produces two games:
        #   Game A: Engine 1 = White
        #   Game B: Engine 2 = White
        while game_number < args.games:
            opening = random.choice(openings)

            print()
            print(
                f'Opening: {opening["name"]} '
                f'({" ".join(opening["moves"])})'
            )

            # -------------------------
            # Game 1: Engine 1 as White
            # -------------------------

            game_id = object()

            result, plies = play_game(
                white=engine1,
                black=engine2,
                movetime=movetime,
                game_id=game_id,
                opening=opening,
            )

            update_score(
                result=result,
                white_index=0,
                black_index=1,
                score=score,
            )

            game_number += 1

            print(
                f"Game {game_number}: "
                f"{result} ({plies} played plies)"
            )

            # -------------------------
            # Game 2: Engine 2 as White
            # -------------------------

            game_id = object()

            result, plies = play_game(
                white=engine2,
                black=engine1,
                movetime=movetime,
                game_id=game_id,
                opening=opening,
            )

            update_score(
                result=result,
                white_index=1,
                black_index=0,
                score=score,
            )

            game_number += 1

            print(
                f"Game {game_number}: "
                f"{result} ({plies} played plies)"
            )

    finally:
        engine1.quit()
        engine2.quit()

    print()
    print("Final result")
    print(f"Engine 1: {score[0]}")
    print(f"Engine 2: {score[1]}")


if __name__ == "__main__":
    main()