package bitboards

type QueenBitboard BitBoard

func (b *QueenBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (b QueenBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

// Moves calculates all possible movements for a queen combining both straight and diagonal paths.
func (q QueenBitboard) Moves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return q.VerticalMoves(sameColorOccupancy, oppositeColorOccupancy) | q.HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy) | q.DiagonalMoves(sameColorOccupancy, oppositeColorOccupancy)
}

func (q QueenBitboard) Attacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return (q.VerticalMoves(sameColorOccupancy, oppositeColorOccupancy) | q.HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy) | q.DiagonalMoves(sameColorOccupancy, oppositeColorOccupancy)) & oppositeColorOccupancy
}

func (b QueenBitboard) MovesByPiece(sameColorOccupancy, oppositeColorOccupancy BitBoard) map[BitBoard]BitBoard {
	queens := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, queen := range queens {
		moves[queen] = QueenBitboard(queen).Moves(sameColorOccupancy, oppositeColorOccupancy)
	}

	return moves
}

// VerticalMoves calculates the vertical movement possibilities for a queen.
func (q QueenBitboard) VerticalMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	rook := RookBitboard(q)
	return rook.verticalUpMoves(sameColorOccupancy, oppositeColorOccupancy) | rook.verticalDownMoves(sameColorOccupancy, oppositeColorOccupancy)
}

// HorizontalMoves calculates the horizontal movement possibilities for a queen.
func (q QueenBitboard) HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	rook := RookBitboard(q)
	return rook.horizontalRightMoves(sameColorOccupancy, oppositeColorOccupancy) | rook.horizontalLeftMoves(sameColorOccupancy, oppositeColorOccupancy)
}

// DiagonalMoves calculates the diagonal movement possibilities for a queen.
func (q QueenBitboard) DiagonalMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	bishop := BishopBitboard(q)
	return bishop.diagonalNorthEastMoves(sameColorOccupancy, oppositeColorOccupancy) | bishop.diagonalNorthWestMoves(sameColorOccupancy, oppositeColorOccupancy) |
		bishop.diagonalSouthEastMoves(sameColorOccupancy, oppositeColorOccupancy) | bishop.diagonalSouthWestMoves(sameColorOccupancy, oppositeColorOccupancy)
}
