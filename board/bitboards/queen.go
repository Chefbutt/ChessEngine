package bitboards

type QueenBitboard BitBoard

func (b *QueenBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (b QueenBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

// Moves calculates all possible movements for a queen combining both straight and diagonal paths.
func (q QueenBitboard) Moves(occupancy BitBoard) BitBoard {
	return q.VerticalMoves(occupancy) | q.HorizontalMoves(occupancy) | q.DiagonalMoves(occupancy)
}

// VerticalMoves calculates the vertical movement possibilities for a queen.
func (q QueenBitboard) VerticalMoves(occupancy BitBoard) BitBoard {
	rook := RookBitboard(q)
	return rook.verticalUpMoves(occupancy) | rook.verticalDownMoves(occupancy)
}

// HorizontalMoves calculates the horizontal movement possibilities for a queen.
func (q QueenBitboard) HorizontalMoves(occupancy BitBoard) BitBoard {
	rook := RookBitboard(q)
	return rook.horizontalRightMoves(occupancy) | rook.horizontalLeftMoves(occupancy)
}

// DiagonalMoves calculates the diagonal movement possibilities for a queen.
func (q QueenBitboard) DiagonalMoves(occupancy BitBoard) BitBoard {
	bishop := BishopBitboard(q)
	return bishop.diagonalNorthEastMoves(occupancy) | bishop.diagonalNorthWestMoves(occupancy) |
		bishop.diagonalSouthEastMoves(occupancy) | bishop.diagonalSouthWestMoves(occupancy)
}
