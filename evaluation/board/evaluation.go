package board

import (
	"engine/evaluation/board/bitboards"
)

func (board *Board) materialScore() int {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

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

	return queenWt*(wQ-bQ) + rookWt*(wR-bR) + knightWt*(wN-bN) + bishopWt*(wB-bB) + pawnWt*(wP-bP)
}

func (board *Board) calculateDoubledPawns() int {
	doubled := 0
	for file := 0; file < 8; file++ {
		fileMask := bitboards.FileMask(file)
		pawnsInFile := board.WhitePawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled += pawnsInFile.PopCount()
		}

		pawnsInFile = board.BlackPawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled -= pawnsInFile.PopCount()
		}
	}

	return doubled
}

func (board *Board) calculateBlockedPawns() int {
	occupied := bitboards.BitBoard(board.WhitePawns) | bitboards.BitBoard(board.BlackPawns)
	blockedWhite := (bitboards.BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (bitboards.BitBoard(board.BlackPawns) >> 8) & occupied

	return blockedWhite.PopCount() - blockedBlack.PopCount()
}

func (board *Board) calculateIsolatedPawns() int {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return isolatedWhite.BitBoard().PopCount() - isolatedBlack.BitBoard().PopCount()
}

// Calculate mobility score (stub)
func (board *Board) mobilityScore() int {
	totalMoves := 0

	totalMoves += board.WhitePawns.Moves(board.EmptySquares, board.BlackPieces, board.EnPassantTarget).PopCount()
	totalMoves += board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces).PopCount()
	totalMoves += board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	totalMoves += board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	totalMoves += board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	// totalMoves -= board.WhiteKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	totalMoves -= board.BlackPawns.Moves(board.EmptySquares, board.WhitePieces, board.EnPassantTarget).PopCount()
	totalMoves -= board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces).PopCount()
	totalMoves -= board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	totalMoves -= board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	totalMoves -= board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	// totalMoves += board.BlackKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	return totalMoves
}

func (board *Board) knightsOnRim() int {
	knightsOnRim := (board.WhiteKnights.BitBoard() & edgesOfBoard).PopCount()
	knightsOnRim = knightsOnRim - (board.BlackKnights.BitBoard() & edgesOfBoard).PopCount()
	return -knightsOnRim
}

func (board *Board) piecesInCentre() int {
	pieces := (board.WhitePieces & centralSquares).PopCount()
	pieces = pieces - (board.BlackPieces & centralSquares).PopCount()

	return pieces
}

// IsAttacked

func (board *Board) IsCheckMate() bool {
	if !board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
		return false
	}

	legalMoves := board.LegalMoves()

	return len(legalMoves) == 0 // If in check and no legal moves, it's checkmate
}

func (board *Board) IsStaleMate() bool {
	if !board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
		legalMoves := board.LegalMoves()
		return len(legalMoves) == 0
	}

	return false
}

func (board *Board) KingSafetyBonus() int {
	bonus := 0
	if board.CastleWhiteKingside || board.CastleWhiteQueenside {
		bonus += 30 // Add 30 points for white castling
	}
	if board.CastleBlackKingside || board.CastleBlackQueenside {
		bonus -= 30 // Subtract 30 points when black castles, as lower is better for black
	}
	return bonus
}

func (board *Board) GamePhase() float64 {
	// Simplistic calculation: count the number of minor and major pieces
	totalPieces := board.WhitePawns.BitBoard().PopCount() + board.BlackPawns.BitBoard().PopCount() +
		board.WhiteKnights.BitBoard().PopCount() + board.BlackKnights.BitBoard().PopCount() +
		board.WhiteBishops.BitBoard().PopCount() + board.BlackBishops.BitBoard().PopCount() +
		board.WhiteRooks.BitBoard().PopCount() + board.BlackRooks.BitBoard().PopCount() +
		board.WhiteQueens.BitBoard().PopCount() + board.BlackQueens.BitBoard().PopCount()

	// Assuming a typical game starts with 32 pieces
	return 1 - float64(totalPieces)/32.0 // This gives us a range from 0 (start) to 1 (few pieces left)
}

func (board *Board) DynamicKingSafety() int {
	phase := board.GamePhase()
	// Scale the king safety bonus based on the game phase, increasing as the game progresses
	return int(float64(board.KingSafetyBonus()) * (0.5 + 0.5*phase)) // Weight more heavily towards the endgame
}

// Evaluation function
func (board *Board) Evaluate() int {
	// Check for game-over conditions
	if board.IsCheckMate() {
		if board.TurnBlack {
			return -9999 // Checkmate against black
		} else {
			return 9999 // Checkmate against white
		}
	} else if board.IsStaleMate() {
		return 0 // Stalemate considered a draw
	}

	material := board.materialScore() * 30
	doubled := board.calculateDoubledPawns()
	blocked := board.calculateBlockedPawns()
	isolated := board.calculateIsolatedPawns()
	mobility := board.mobilityScore()
	centre := board.piecesInCentre()
	kingSafety := board.DynamicKingSafety()
	misplacedKnights := board.knightsOnRim()

	// Additional factors like doubled, blocked, and isolated pawns
	pawnPenalties := 6 * float64(doubled+blocked+isolated)
	mobilityBonus := 0.2 * float64(mobility)
	centreBonus := 29 * float64(centre)
	knightBonus := 5 * float64(misplacedKnights)
	safety := 4.5 * float64(kingSafety)

	score := float64(material) + pawnPenalties + mobilityBonus + centreBonus + safety + knightBonus

	return int(score)
}
