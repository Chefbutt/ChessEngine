package evaluation

import (
	"fmt"
	"math/bits"
	"strings"

	"engine/evaluation/board"
)

type Evaluation struct{}

func indexToPosition(bitboard uint64) string {
	index := bits.TrailingZeros64(bitboard)
	rank := index % 8
	file := index / 8
	return string('a'+rune(rank)) + string('1'+rune(file))
}

func Evaluate(b board.Board) (board.Move, error) {
	// fmt.Print(b.BlackPawns.Moves(b.EmptySquares))
	// b.EmptySquares.Display()
	// b.BlackPawns.BitBoard().Display()
	var moveSliceAllPieces []string
	for src, moves := range b.RegularMoves() {
		var moveSlice []string
		for _, move := range moves {
			moveSlice = append(moveSlice, indexToPosition(uint64(move)))
		}
		moveSliceAllPieces = append(moveSliceAllPieces, fmt.Sprint(board.PieceSymbols[b.PieceAt(bits.TrailingZeros64(uint64(src)))], "(", indexToPosition(uint64(src)), "): ", strings.Join(moveSlice, ", ")))
	}

	fmt.Print(strings.Join(moveSliceAllPieces, "; "))
	return board.Move{}, nil
}
