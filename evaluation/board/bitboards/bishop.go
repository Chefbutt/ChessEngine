package bitboards

type BishopBitboard BitBoard

func (b BishopBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b *BishopBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (b BishopBitboard) Attacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return b.diagonalNorthEastAttacks(sameColorOccupancy, oppositeColorOccupancy) | b.diagonalNorthWestAttacks(sameColorOccupancy, oppositeColorOccupancy) |
		b.diagonalSouthEastAttacks(sameColorOccupancy, oppositeColorOccupancy) | b.diagonalSouthWestAttacks(sameColorOccupancy, oppositeColorOccupancy)
}

func (b BishopBitboard) diagonalNorthEastAttacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
	}
	return moves
}

// diagonalNorthWestMoves calculates all northwest diagonal moves until blocked.
func (b BishopBitboard) diagonalNorthWestAttacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
	}
	return moves
}

// diagonalSouthEastMoves calculates all southeast diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthEastAttacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			break
		}
		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
	}
	return moves
}

// diagonalSouthWestMoves calculates all southwest diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthWestAttacks(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&oppositeColorOccupancy != 0 { // Check for blockage
			moves |= pos // Include the blocking piece's square
			break
		}
	}
	return moves
}

// DiagonalMoves calculates all diagonal movement possibilities for a bishop.
func (b BishopBitboard) Moves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	return b.diagonalNorthEastMoves(sameColorOccupancy, oppositeColorOccupancy) | b.diagonalNorthWestMoves(sameColorOccupancy, oppositeColorOccupancy) |
		b.diagonalSouthEastMoves(sameColorOccupancy, oppositeColorOccupancy) | b.diagonalSouthWestMoves(sameColorOccupancy, oppositeColorOccupancy)
}

func (b BishopBitboard) MovesByPiece(sameColorOccupancy, oppositeColorOccupancy BitBoard) map[BitBoard]BitBoard {
	bishops := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, bishop := range bishops {
		moves[bishop] = BishopBitboard(bishop).Moves(sameColorOccupancy, oppositeColorOccupancy)
	}

	return moves
}

// diagonalNorthEastMoves calculates all northeast diagonal moves until blocked.
func (b BishopBitboard) diagonalNorthEastMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			// moves |= pos // Include the blocking piece's square
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

// diagonalNorthWestMoves calculates all northwest diagonal moves until blocked.
func (b BishopBitboard) diagonalNorthWestMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.northWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			// moves |= pos // Include the blocking piece's square
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

// diagonalSouthEastMoves calculates all southeast diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthEastMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southEastOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			// moves |= pos // Include the blocking piece's square
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

// diagonalSouthWestMoves calculates all southwest diagonal moves until blocked.
func (b BishopBitboard) diagonalSouthWestMoves(sameColorOccupancy, oppositeColorOccupancy BitBoard) BitBoard {
	moves := BitBoard(0)
	pos := BitBoard(b)
	for {
		if newPos := pos.southWestOne(); pos == newPos {
			break
		} else {
			pos = newPos
		}

		if pos&sameColorOccupancy != 0 { // Check for blockage
			// moves |= pos // Include the blocking piece's square
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
