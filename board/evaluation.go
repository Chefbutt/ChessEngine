package board

import "engine/board/bitboards"

func (board *Board) materialScore() int {
	kingWt, queenWt, rookWt, knightWt, bishopWt, pawnWt := 200, 9, 5, 3, 3, 1

	wK := board.WhiteKing.BitBoard().PopCount()
	bK := board.BlackKing.BitBoard().PopCount()
	wQ := board.WhiteQueens.BitBoard().PopCount()
	bQ := board.BlackQueens.BitBoard().PopCount()
	wR := board.WhiteRooks.BitBoard().PopCount()
	bR := board.BlackRooks.BitBoard().PopCount()
	wN := board.WhiteKnights.BitBoard().PopCount()
	bN := board.BlackKnights.BitBoard().PopCount()
	wB := board.WhiteBishops.BitBoard().PopCount()
	bB := board.BlackBishops.BitBoard().PopCount()
	wP := board.WhitePawns.BitBoard().PopCount()
	bP := board.BlackPawns.BitBoard().PopCount()

	return kingWt*(wK-bK) + queenWt*(wQ-bQ) + rookWt*(wR-bR) + knightWt*(wN-bN) + bishopWt*(wB-bB) + pawnWt*(wP-bP)
}

// Stub functions for calculating doubled, blocked, and isolated pawns
func (board *Board) calculateDoubledPawns() int {
	// Example for white pawns, same logic applies for black with relevant bitboard
	doubled := 0
	for file := 0; file < 8; file++ {
		fileMask := bitboards.FileMask(file)
		pawnsInFile := board.WhitePawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled += pawnsInFile.PopCount() - 1
		}

		pawnsInFile = board.BlackPawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled += pawnsInFile.PopCount() + 1
		}
	}
	return doubled
}

func (board *Board) calculateBlockedPawns() int {
	// Shift white pawns one rank up and check for overlap with occupied squares (both colors)
	occupied := bitboards.BitBoard(board.WhitePawns) | bitboards.BitBoard(board.BlackPawns)
	blockedWhite := (bitboards.BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (bitboards.BitBoard(board.BlackPawns) >> 8) & occupied

	return blockedWhite.PopCount() + blockedBlack.PopCount()
}

func (board *Board) calculateIsolatedPawns() int {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return isolatedWhite.BitBoard().PopCount() + isolatedBlack.BitBoard().PopCount()
}

// Calculate mobility score (stub)
func (board *Board) mobilityScore() int {
	totalMoves := 0

	totalMoves += board.WhitePawns.Moves(board.EmptySquares).PopCount()
	totalMoves += board.WhiteKnights.Moves(board.EmptySquares).PopCount()
	totalMoves += board.WhiteBishops.Moves(board.OccupiedSquares).PopCount()
	totalMoves += board.WhiteRooks.Moves(board.OccupiedSquares).PopCount()
	totalMoves += board.WhiteQueens.Moves(board.OccupiedSquares).PopCount()
	totalMoves += board.WhiteKing.Moves(board.EmptySquares).PopCount()

	totalMoves -= board.BlackPawns.Moves(board.EmptySquares).PopCount()
	totalMoves -= board.BlackKnights.Moves(board.EmptySquares).PopCount()
	totalMoves -= board.BlackBishops.Moves(board.OccupiedSquares).PopCount()
	totalMoves -= board.BlackRooks.Moves(board.OccupiedSquares).PopCount()
	totalMoves -= board.BlackQueens.Moves(board.OccupiedSquares).PopCount()
	totalMoves -= board.BlackKing.Moves(board.EmptySquares).PopCount()

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
