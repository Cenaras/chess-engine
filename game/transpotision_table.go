package game

import (
	"unsafe"
)

// Transposition table: Stores the evaluation of a previously encounted position
type TranspositionTableEntry struct {
	Hash     uint64 // the hash of the position
	BestMove Move   // Best/Refutation move in this position
	Depth    int    // The depth we searched to
	Score    int    // the evaluation of the position
	NodeType        // type of node: PV gives exact score, all-node gives upper bound, cut-node gives lower bound
}

type TranspositionTable struct {
	entries []TranspositionTableEntry
}

// TODO: We need a bit better understanding of alpha and beta, but we can use those as well in our
// transpotion table: For instance, say we encounter a move in a position where bestScore > beta.
// At this point, we don't know that our current most is the best, but it doesn't matter -- while we
// may be able to do even better, our opponent has already found a move that would be better for him

type NodeType int

const (
	EXACT       NodeType = 0
	UPPER_BOUND NodeType = 1
	LOWER_BOUND NodeType = 2
)

func CreateTranspositionTable(megabytes int) TranspositionTable {
	bytes := megabytes * 1024 * 1024
	size := int(unsafe.Sizeof(TranspositionTableEntry{}))
	entries := bytes / size

	return TranspositionTable{
		entries: make([]TranspositionTableEntry, entries),
	}
}

var TT TranspositionTable = CreateTranspositionTable(64)

func (tt *TranspositionTable) index(hash uint64) int {
	return int(hash % uint64(len(tt.entries)))
}

func (tt *TranspositionTable) Store(entry TranspositionTableEntry) {
	index := tt.index(entry.Hash)
	tt.entries[index] = entry
}

func (tt *TranspositionTable) Lookup(hash uint64) (TranspositionTableEntry, bool) {
	index := tt.index(hash)
	entry := tt.entries[index]

	// ensures that it wasn't a miss
	if entry.Hash != hash {
		return TranspositionTableEntry{}, false
	}
	return entry, true
}

// During MoveGen we maintain the zobrist hash
// During search, we obvs make moves. That updates the current hash
// If, during search, we encounter a position where we already have the result in our tt,
// then we can stop searching further -- we already know from the TT what the best move is etc.
