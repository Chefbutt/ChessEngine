package board

// KingBitboard represents the position of a king on the board using an index.
type KingBitboard uint64

// Moves calculates all possible attack positions for a king from a given position.
func (k KingBitboard) Moves(empty BitBoard) BitBoard {
	kingPos := BitBoard(k)
	attacks := kingPos.eastOne() | kingPos.westOne() | kingPos.northOne() | kingPos.southOne()
	return attacks & empty
}
