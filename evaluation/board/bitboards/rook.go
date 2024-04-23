package bitboards

type RookBitboard BitBoard

func (b RookBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b *RookBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

// VerticalMoves calculates the vertical movement possibilities for a rook.

func (r RookBitboard) Attacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return (r.VerticalMoves(sameColorOccupancy, oppositeColorOccupancy) | r.HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy)) & oppositeColorOccupancy
}

func (r RookBitboard) Moves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return r.VerticalMoves(sameColorOccupancy, oppositeColorOccupancy) | r.HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy)
}

func (b RookBitboard) MovesByPiece(sameColorOccupancy, oppositeColorOccupancy BitBoard) map[BitBoard]BitBoard {
	rooks := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, rook := range rooks {
		moves[rook] = RookBitboard(rook).Moves(sameColorOccupancy, oppositeColorOccupancy)
	}

	return moves
}

// VerticalMoves calculates the vertical movement possibilities for a rook.
func (r RookBitboard) VerticalMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return r.verticalUpMoves(sameColorOccupancy, oppositeColorOccupancy) | r.verticalDownMoves(sameColorOccupancy, oppositeColorOccupancy)
}

// HorizontalMoves calculates the horizontal movement possibilities for a rook.
func (r RookBitboard) HorizontalMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return r.horizontalLeftMoves(sameColorOccupancy, oppositeColorOccupancy) | r.horizontalRightMoves(sameColorOccupancy, oppositeColorOccupancy)
}

// verticalUpMoves calculates all upwards moves from the rook's position until blocked.
func (r RookBitboard) verticalUpMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.northOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		if pos&sameColorOccupancy != 0 { // Check if there is a piece blocking further movement
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
		moves |= pos
	}
	return moves
}

// verticalDownMoves calculates all downwards moves from the rook's position until blocked.
func (r RookBitboard) verticalDownMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.southOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		if pos&sameColorOccupancy != 0 {
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
		moves |= pos
	}
	return moves
}

// horizontalRightMoves calculates all rightward moves from the rook's position until blocked.
func (r RookBitboard) horizontalRightMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.eastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		if pos&sameColorOccupancy != 0 {
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
		moves |= pos
	}
	return moves
}

// horizontalLeftMoves calculates all leftward moves from the rook's position until blocked.
func (r RookBitboard) horizontalLeftMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.westOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		if pos&sameColorOccupancy != 0 {
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
		moves |= pos
	}
	return moves
}
