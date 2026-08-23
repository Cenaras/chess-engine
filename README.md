# Chess Engine
This project contains a chess engine written in Go.

# Testing
This project has two different types of tests: Unit tests and Perft tests.

### Unit tests
This tests assert that expected behaviour is maintained for specific pieces, and that common pitfalls are avoided. Most of the tests do not enumerate an enormous amount of moves, but rather focus on specific pieces. Unit test coverage should therefore *not* be considered as an indicator for complete correctness, but rather as a useful metric for quickly testing general logic.

Run unit tests with:
`go test -short ./...`


### Perft tests
Perft tests (performance tests, move path enumeration) is a useful debugging function to enumerate all possible moves in a position. We use such tests for two things:
 - 1) Enumerate all possible positions and compare with known results for various depths, to assert implementation correctness
 - 2) Run performance tests to inspect the efficiency of the current implementation

The test suite contains common perft positions and inspects them at various depths (typically up to 5) to detect bugs and investigate performance. 
Use these tests to verify the overall implementation and performance results after major engine rewrites.

Run all perft tests with:
`go test ./game/perft_test.go`
Note: This may take a while. 


# TODO:
 - Refactor + improve code quality
 - Implement Iterative Deepening to better evaluate search impact
 - Then better move ordering
 - More extensive testing: A lot of code isn't testing
 - Better scripting for benchmark: Ideally python script that runs the fastchess 
    tournament and parses the result nicely
 



