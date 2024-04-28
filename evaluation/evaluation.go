package evaluation

import (
	"fmt"
	"math/bits"

	"engine/evaluation/board"
)

type Evaluation struct{}

func Evaluate(b board.Board) (board.Move, error) {
	// fmt.Print(b.BlackPawns.Moves(b.EmptySquares))
	// b.EmptySquares.Display()
	// b.BlackPawns.BitBoard().Display()
	moveSliceAllPieces := make(map[int]string)
	for _, move := range b.LegalMoves() {
		if old, ok := moveSliceAllPieces[move.Source]; ok {
			moveSliceAllPieces[move.Source] = old + ", " + board.IndexToPosition(uint64(move.Destination))
		} else {
			moveSliceAllPieces[move.Source] = fmt.Sprint(board.PieceSymbols[b.PieceAt(bits.TrailingZeros64(uint64(move.Source)))], "(", board.IndexToPosition(uint64(move.Source)), "): ", board.IndexToPosition(uint64(move.Destination)))
		}
	}

	for _, move := range moveSliceAllPieces {
		fmt.Print(move, " ")
	}

	fmt.Println()

	return board.Move{}, nil
}
