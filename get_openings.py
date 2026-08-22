# Downloded from https://github.com/Disservin/fastchess/blob/master/tools/get_lichess_openings.py
# and modified to our needs

# Script that downloads bin/gen.py from lichess-org/chess-openings, runs it to generate combined data, and updates 'openings_data.cpp' file

import csv
import io
import json
import sys
from pathlib import Path

try:
    import chess
    import chess.pgn
except ImportError:
    print("Missing dependency: python-chess")
    print("Install it with: pip install chess")
    sys.exit(1)

try:
    import requests
except ImportError:
    print("Missing dependency: requests")
    print("Install it with: pip install requests")
    sys.exit(1)


REPO_BASE_URL = (
    "https://raw.githubusercontent.com/"
    "lichess-org/chess-openings/refs/heads/master"
)

TSV_FILES = [
    "a.tsv",
    "b.tsv",
    "c.tsv",
    "d.tsv",
    "e.tsv",
]

SCRIPT_DIR = Path(__file__).resolve().parent
OUTPUT_DIR = SCRIPT_DIR / "openings"

ALL_TSV_PATH = OUTPUT_DIR / "all.tsv"
PGN_PATH = OUTPUT_DIR / "openings.pgn"
JSON_PATH = OUTPUT_DIR / "openings.json"


def download_tsv(filename: str) -> str:
    url = f"{REPO_BASE_URL}/{filename}"

    print(f"Downloading {filename}...")

    response = requests.get(url, timeout=30)
    response.raise_for_status()

    return response.text


def parse_opening(eco: str, name: str, pgn: str) -> dict:
    """
    Convert one lichess opening entry into a richer representation.
    """

    game = chess.pgn.read_game(io.StringIO(pgn))

    if game is None:
        raise ValueError(f"Could not parse PGN: {pgn}")

    board = game.board()
    uci_moves = []

    for move in game.mainline_moves():
        uci_moves.append(move.uci())
        board.push(move)

    return {
        "eco": eco,
        "name": name,
        "pgn": pgn,
        "uci": " ".join(uci_moves),
        "epd": board.epd(),
    }


def load_openings() -> list[dict]:
    openings = []

    for filename in TSV_FILES:
        text = download_tsv(filename)

        reader = csv.DictReader(
            io.StringIO(text),
            delimiter="\t",
        )

        for row in reader:
            try:
                opening = parse_opening(
                    row["eco"].strip(),
                    row["name"].strip(),
                    row["pgn"].strip(),
                )
            except Exception as e:
                print(
                    f"Failed parsing {filename}: "
                    f"{row.get('name', '?')}: {e}"
                )
                raise

            openings.append(opening)

    return openings


def write_tsv(openings: list[dict]) -> None:
    print(f"Writing {ALL_TSV_PATH}...")

    with open(
        ALL_TSV_PATH,
        "w",
        encoding="utf-8",
        newline="",
    ) as f:
        writer = csv.DictWriter(
            f,
            fieldnames=[
                "eco",
                "name",
                "pgn",
                "uci",
                "epd",
            ],
            delimiter="\t",
        )

        writer.writeheader()
        writer.writerows(openings)


def write_json(openings: list[dict]) -> None:
    print(f"Writing {JSON_PATH}...")

    with open(
        JSON_PATH,
        "w",
        encoding="utf-8",
    ) as f:
        json.dump(
            openings,
            f,
            ensure_ascii=False,
            indent=2,
        )


def write_pgn(openings: list[dict]) -> None:
    print(f"Writing {PGN_PATH}...")

    with open(
        PGN_PATH,
        "w",
        encoding="utf-8",
        newline="\n",
    ) as f:
        for opening in openings:
            game = chess.pgn.read_game(
                io.StringIO(opening["pgn"])
            )

            if game is None:
                continue

            game.headers.clear()

            game.headers["Event"] = opening["name"]
            game.headers["ECO"] = opening["eco"]
            game.headers["Result"] = "*"

            exporter = chess.pgn.StringExporter(
                headers=True,
                variations=False,
                comments=False,
            )

            f.write(game.accept(exporter))
            f.write("\n\n")


def main() -> None:
    OUTPUT_DIR.mkdir(
        parents=True,
        exist_ok=True,
    )

    openings = load_openings()

    print(f"Loaded {len(openings)} opening lines.")

    write_tsv(openings)
    write_json(openings)
    write_pgn(openings)

    print()
    print("Done!")
    print(f"TSV:  {ALL_TSV_PATH}")
    print(f"PGN:  {PGN_PATH}")
    print(f"JSON: {JSON_PATH}")


if __name__ == "__main__":
    main()