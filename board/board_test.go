package board

import (
	"fmt"
	"testing"

	"engine/board/bitboards"
)

func TestKingAttacks(t *testing.T) {
	// Example of a board with pieces and calculating king's attacks
	// board := Board{
	// 	// Pieces:    make(map[Square]Piece),
	// 	TurnBlack: false,
	// 	Move:      1,
	// }

	// King at position e1 (the 5th bit in the lowest row)
	var kingPos bitboards.BitBoard = 0x10
	var king bitboards.KingBitboard = 5 // King piece at index 5, which corresponds to e1

	// board.Pieces[kingPos] = Piece{PieceType: byte(king)}
	attacks := king.Moves(kingPos)

	fmt.Printf("Possible attacks for a king at e1: %064b\n", attacks)
}
