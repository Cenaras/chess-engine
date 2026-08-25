# How the testing framework works
We use `8moves_v3` for random openings. For each such opening, `filter_openings.py` filters away inequal openings. To do this we play the opening, invoke stockfish and evaluate the position.

**Generate openings**
```sh
python filter_openings.py 8moves_v3.pgn /path/to/stockfish.exe /path/to/output.json
```
Example:
```sh
python testing/filter_openings.py testing/8moves_v3.pgn ../../chess/stockfish/stockfish/stockfish-windows-x86-64-avx2.exe testing/openings.json
```

**Run the testing framework**
You need a binary for each of the two engines you wish to compare.
```sh
python testing/compare.py /path/to/engine1.exe /path/to/engine2.exe no-of-games /path/to/openings.json [--movetime-ms 100 ] [--seed 0]
```

Invoke the `compare.py` script to perform the testing.
Example -- 100 games (50 black 50 white):
```sh
python testing/compare.py /path/to/engine1.exe /path/to/engine2.exe 100 testing/openings.json --movetime-ms 100
```

To reproduce a run, copy the reported `SEED:` from the execution and use `--seed`