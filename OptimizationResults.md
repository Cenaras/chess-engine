# Optimization Results

This small document holds some notes on how the performance has changed of perft tests, move generation etc...

### Run Profiling
Run the profiling tool with the following example command: (exclude -run to run all tests)
`go test ./game -run "^TestStartingFen$" -cpuprofile=cpu.prof -count=1`
Inspect the result:
`go tool pprof cpu`
And look at the top-most hotspots
`top` or `top -cum`

Or investigate in the browser:
`go tool pprof -http=:8080 cpu.prof`

Meaning of columns:
 - cum: cumulative time spend in this function (including sub-calls)
 - flat: time spend directly inside this method (excluding sub-calls)

To inspect a specific function further, do: 
`list [functiongame]`

### Perft tests
**Starting Position (5), Position 3 (3) , Position 5 (5)**
Additions are incremental, meaning row i+1 contains all optimizations from 0...i

Naive Make-All-Response-Moves:  TIMEOUT / N/A
IsSquareAttackedChecK:          366.77s
KingSquareInPosition:           368.57s
BooleanOverErrorForOOB:         6.98s   (Constructing an error after every illegal move is expensive!!!)

**Starting Position (6), Position 3 (6), Position 4 (6), Position 6 (5)**
BooleanOverErrorForOOB         63.2s