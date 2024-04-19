package board

type BishopBitboard BitBoard

// DiagonalMoves calculates all diagonal movement possibilities for a bishop.
func (b BishopBitboard) Moves(occupancy BitBoard) BitBoard {
	return b.diagonalNorthEastMoves(occupancy) | b.diagonalNorthWestMoves(occupancy) |
		b.diagonalSouthEastMoves(occupancy) | b.diagonalSouthWestMoves(occupancy)
}

// diagonalNorthEastMoves calculates all northeast diagonal moves until blocked.
func (b BishopBitboard) diagonalNorthEastMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&occupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
		moves |= pos
	}
	return moves
}

// diagonalNorthWestMoves calculates all northwest diagonal moves until blocked.
func (b BishopBitboard) diagonalNorthWestMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&occupancy != 0 {
			moves |= pos
			break
		}
		moves |= pos
	}
	return moves
}

// diagonalSouthEastMoves calculates all southeast diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthEastMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}
		if pos&occupancy != 0 {
			moves |= pos
			break
		}
		moves |= pos
	}
	return moves
}

// diagonalSouthWestMoves calculates all southwest diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthWestMoves(occupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&occupancy != 0 {
			moves |= pos
			break
		}
		moves |= pos
	}
	return moves
}
