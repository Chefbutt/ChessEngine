package bitboards

// KingBitboard represents the position of a king on the board using an index.
type KingBitboard uint64

func (b KingBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b *KingBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

// Moves calculates all possible attack positions for a king from a given position.
func (k KingBitboard) Moves(empty BitBoard) BitBoard {
	kingPos := BitBoard(k)
	attacks := kingPos.eastOne() | kingPos.westOne() | kingPos.northOne() | kingPos.southOne()
	return attacks & empty
}
