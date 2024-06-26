package bitboards

// KingBitboard represents the position of a king on the board using an index.
type KingBitboard uint64

func (b KingBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b *KingBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (k KingBitboard) Attacks(oppositeColorOccupancy BitBoard) BitBoard {
	kingPos := BitBoard(k)
	moves := kingPos.eastOne() | kingPos.westOne() | kingPos.northOne() | kingPos.southOne() | kingPos.southEastOne() | kingPos.southWestOne() | kingPos.northEastOne() | kingPos.northWestOne()
	return moves & oppositeColorOccupancy
}

// Moves calculates all possible attack positions for a king from a given position.
func (k KingBitboard) Moves(empty, oppositeColorOccupancy BitBoard) BitBoard {
	kingPos := BitBoard(k)
	moves := kingPos.eastOne() | kingPos.westOne() | kingPos.northOne() | kingPos.southOne() | kingPos.southEastOne() | kingPos.southWestOne() | kingPos.northEastOne() | kingPos.northWestOne()
	return moves & (oppositeColorOccupancy | empty)
}

func (b KingBitboard) MovesByPiece(empty, oppositeColorOccupancy BitBoard) map[BitBoard]BitBoard {
	bishops := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, king := range bishops {
		moves[king] = KingBitboard(king).Moves(empty, oppositeColorOccupancy)
	}

	return moves
}
