package board

func popCount(b BitBoard) int {
	count := 0
	for b != 0 {
		count++
		b &= b - 1 // reset least significant bit
	}
	return count
}

func (board *Board) materialScore() int {
	kingWt, queenWt, rookWt, knightWt, bishopWt, pawnWt := 200, 9, 5, 3, 3, 1

	wK := popCount(BitBoard(board.WhiteKing))
	bK := popCount(BitBoard(board.BlackKing))
	wQ := popCount(BitBoard(board.WhiteQueens))
	bQ := popCount(BitBoard(board.BlackQueens))
	wR := popCount(BitBoard(board.WhiteRooks))
	bR := popCount(BitBoard(board.BlackRooks))
	wN := popCount(BitBoard(board.WhiteKnights))
	bN := popCount(BitBoard(board.BlackKnights))
	wB := popCount(BitBoard(board.WhiteBishops))
	bB := popCount(BitBoard(board.BlackBishops))
	wP := popCount(BitBoard(board.WhitePawns))
	bP := popCount(BitBoard(board.BlackPawns))

	return kingWt*(wK-bK) + queenWt*(wQ-bQ) + rookWt*(wR-bR) + knightWt*(wN-bN) + bishopWt*(wB-bB) + pawnWt*(wP-bP)
}

// Stub functions for calculating doubled, blocked, and isolated pawns
func (board *Board) calculateDoubledPawns() int {
	// Example for white pawns, same logic applies for black with relevant bitboard
	doubled := 0
	for file := 0; file < 8; file++ {
		fileMask := BitBoard(0x0101010101010101 << file)
		pawnsInFile := BitBoard(board.WhitePawns) & fileMask
		if popCount(pawnsInFile) > 1 {
			doubled += popCount(pawnsInFile) - 1
		}

		pawnsInFile = BitBoard(board.BlackPawns) & fileMask
		if popCount(pawnsInFile) > 1 {
			doubled += popCount(pawnsInFile) + 1
		}
	}
	return doubled
}

func (board *Board) calculateBlockedPawns() int {
	// Shift white pawns one rank up and check for overlap with occupied squares (both colors)
	occupied := BitBoard(board.WhitePawns) | BitBoard(board.BlackPawns)
	blockedWhite := (BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (BitBoard(board.BlackPawns) >> 8) & occupied

	return popCount(blockedWhite) + popCount(blockedBlack)
}

func (board *Board) calculateIsolatedPawns() int {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return popCount(BitBoard(isolatedWhite)) + popCount(BitBoard(isolatedBlack))
}

// Calculate mobility score (stub)
func (board *Board) mobilityScore() int {
	totalMoves := 0

	totalMoves += popCount(board.WhitePawns.Moves(board.EmptySquares))
	totalMoves += popCount(board.WhiteKnights.Moves(board.EmptySquares))
	totalMoves += popCount(board.WhiteBishops.Moves(board.OccupiedSquares))
	totalMoves += popCount(board.WhiteRooks.Moves(board.OccupiedSquares))
	totalMoves += popCount(board.WhiteQueens.Moves(board.OccupiedSquares))
	totalMoves += popCount(board.WhiteKing.Moves(board.EmptySquares))

	totalMoves -= popCount(board.BlackPawns.Moves(board.EmptySquares))
	totalMoves -= popCount(board.BlackKnights.Moves(board.EmptySquares))
	totalMoves -= popCount(board.BlackBishops.Moves(board.OccupiedSquares))
	totalMoves -= popCount(board.BlackRooks.Moves(board.OccupiedSquares))
	totalMoves -= popCount(board.BlackQueens.Moves(board.OccupiedSquares))
	totalMoves -= popCount(board.BlackKing.Moves(board.EmptySquares))

	return totalMoves
}

// Evaluation function
func (board *Board) Evaluate(who2Move int) int {
	material := board.materialScore()
	doubled := board.calculateDoubledPawns()
	blocked := board.calculateBlockedPawns()
	isolated := board.calculateIsolatedPawns()
	mobility := board.mobilityScore()

	// Additional factors like doubled, blocked, and isolated pawns
	pawnPenalties := -0.5 * float64(doubled+blocked+isolated)
	mobilityBonus := 0.1 * float64(mobility)

	score := float64(material) + pawnPenalties + mobilityBonus

	// Adjust score based on who's move it is
	return int(score) * who2Move
}
