import argparse
import chess
import chess.engine


def play_game(white, black, movetime, game_id):
    board = chess.Board()

    while not board.is_game_over(claim_draw=True):
        engine = white if board.turn == chess.WHITE else black

        result = engine.play(
            board,
            chess.engine.Limit(time=movetime),
            game=game_id,
        )

        if result.move is None:
            raise RuntimeError("Engine returned no move in a non-terminal position")

        board.push(result.move)

    return board.result(claim_draw=True), board.ply()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("engine1", help="Path to first UCI engine")
    parser.add_argument("engine2", help="Path to second UCI engine")
    parser.add_argument("games", type=int, help="Number of games to play")
    parser.add_argument(
        "--movetime-ms",
        type=int,
        default=100,
        help="Time per move in milliseconds",
    )

    args = parser.parse_args()

    movetime = args.movetime_ms / 1000.0

    engine1 = chess.engine.SimpleEngine.popen_uci(args.engine1)
    engine2 = chess.engine.SimpleEngine.popen_uci(args.engine2)

    score = [0.0, 0.0]

    try:
        for game_number in range(args.games):
            # Alternate colors every game.
            if game_number % 2 == 0:
                white = engine1
                black = engine2
                white_index = 0
                black_index = 1
            else:
                white = engine2
                black = engine1
                white_index = 1
                black_index = 0

            # python-chess uses this to know that a new UCI game started.
            game_id = object()

            result, plies = play_game(
                white,
                black,
                movetime,
                game_id,
            )

            if result == "1-0":
                score[white_index] += 1
            elif result == "0-1":
                score[black_index] += 1
            else:
                score[0] += 0.5
                score[1] += 0.5

            print(
                f"Game {game_number + 1}: "
                f"{result} ({plies} plies)"
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