package board

type RookBitboard BitBoard

// VerticalMoves calculates the vertical movement possibilities for a rook.
func (r RookBitboard) Moves(occupancy BitBoard) BitBoard {
	return r.VerticalMoves(occupancy) | r.HorizontalMoves(occupancy)
}

// VerticalMoves calculates the vertical movement possibilities for a rook.
func (r RookBitboard) VerticalMoves(occupancy BitBoard) BitBoard {
	return r.verticalUpMoves(occupancy) | r.verticalDownMoves(occupancy)
}

// HorizontalMoves calculates the horizontal movement possibilities for a rook.
func (r RookBitboard) HorizontalMoves(occupancy BitBoard) BitBoard {
	return r.horizontalLeftMoves(occupancy) | r.horizontalRightMoves(occupancy)
}

// verticalUpMoves calculates all upwards moves from the rook's position until blocked.
func (r RookBitboard) verticalUpMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.northOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		moves |= pos
		if pos&occupancy != 0 { // Check if there is a piece blocking further movement
			break
		}
	}
	return moves
}

// verticalDownMoves calculates all downwards moves from the rook's position until blocked.
func (r RookBitboard) verticalDownMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.southOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		moves |= pos
		if pos&occupancy != 0 {
			break
		}
	}
	return moves
}

// horizontalRightMoves calculates all rightward moves from the rook's position until blocked.
func (r RookBitboard) horizontalRightMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.eastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		moves |= pos
		if pos&occupancy != 0 {
			break
		}
	}
	return moves
}

// horizontalLeftMoves calculates all leftward moves from the rook's position until blocked.
func (r RookBitboard) horizontalLeftMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(r)
	for {
		if newPos := pos.westOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		moves |= pos
		if pos&occupancy != 0 {
			break
		}
	}
	return moves
}
