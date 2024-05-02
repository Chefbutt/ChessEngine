package board

import (
	"math/bits"

	"engine/evaluation/board/bitboards"
)

func (board Board) whiteMaterial() int8 {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	wQ := bits.OnesCount64(uint64(board.WhiteQueens))
	wR := bits.OnesCount64(uint64(board.WhiteRooks))
	wN := bits.OnesCount64(uint64(board.WhiteKnights))
	wB := bits.OnesCount64(uint64(board.WhiteBishops))
	wP := bits.OnesCount64(uint64(board.WhitePawns))

	return int8(queenWt*wQ + rookWt*wR + knightWt*wN + bishopWt*wB + pawnWt*wP)
}

func (board Board) blackMaterial() int8 {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	bQ := bits.OnesCount64(uint64(board.BlackQueens))
	bR := bits.OnesCount64(uint64(board.BlackRooks))
	bN := bits.OnesCount64(uint64(board.BlackKnights))
	bB := bits.OnesCount64(uint64(board.BlackBishops))
	bP := bits.OnesCount64(uint64(board.BlackPawns))

	return int8(queenWt*bQ + rookWt*bR + knightWt*bN + bishopWt*bB + pawnWt*bP)
}

func (board Board) calculateDoubledPawns() int8 {
	var doubled int8
	for file := 0; file < 8; file++ {
		fileMask := bitboards.FileMask(file)
		pawnsInFile := board.WhitePawns.BitBoard() & fileMask
		if bits.OnesCount64(uint64(pawnsInFile)) > 1 {
			doubled += int8(bits.OnesCount64(uint64(pawnsInFile)))
		}

		pawnsInFile = board.BlackPawns.BitBoard() & fileMask
		if bits.OnesCount64(uint64(pawnsInFile)) > 1 {
			doubled -= int8(bits.OnesCount64(uint64(pawnsInFile)))
		}
	}

	return doubled
}

func (board Board) calculateBlockedPawns() int8 {
	occupied := bitboards.BitBoard(board.WhitePawns) | bitboards.BitBoard(board.BlackPawns)
	blockedWhite := (bitboards.BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (bitboards.BitBoard(board.BlackPawns) >> 8) & occupied

	return int8(bits.OnesCount64(uint64(blockedWhite)) - bits.OnesCount64(uint64(blockedBlack)))
}

func (board Board) calculateIsolatedPawns() int8 {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return int8(bits.OnesCount64(uint64(isolatedWhite)) - bits.OnesCount64(uint64(isolatedBlack)))
}

// Calculate mobility score (stub)
func (board Board) mobilityScore() int8 {
	totalMoves := 0

	totalMoves += bits.OnesCount64(uint64(board.WhitePawns.Moves(board.EmptySquares, board.BlackPieces, board.EnPassantTarget)))
	totalMoves += bits.OnesCount64(uint64(board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces)))
	totalMoves += bits.OnesCount64(uint64(board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces)))
	totalMoves += bits.OnesCount64(uint64(board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces)))
	// totalMoves += bits.OnesCount64(uint64(board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces)))
	// totalMoves -= board.WhiteKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	totalMoves -= bits.OnesCount64(uint64(board.BlackPawns.Moves(board.EmptySquares, board.WhitePieces, board.EnPassantTarget)))
	totalMoves -= bits.OnesCount64(uint64(board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces)))
	totalMoves -= bits.OnesCount64(uint64(board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces)))
	totalMoves -= bits.OnesCount64(uint64(board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces)))
	// totalMoves -= bits.OnesCount64(uint64(board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces)))
	// totalMoves += board.BlackKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	return int8(totalMoves)
}

var edgesOfBoard = bitboards.BitBoard(0x8181818181818181)

func (board Board) knightsOnRim() int8 {
	knightsOnRim := (board.WhiteKnights.BitBoard() & edgesOfBoard).PopCount() - (board.BlackKnights.BitBoard() & edgesOfBoard).PopCount()
	return int8(knightsOnRim)
}

func (board Board) piecesInCentre() int8 {
	pieces := bits.OnesCount64(uint64(board.WhitePieces & centralSquares))
	pieces = pieces - bits.OnesCount64(uint64(board.BlackPieces&centralSquares))

	return int8(pieces)
}

// IsAttacked

func (board Board) IsStaleMate() bool {
	if !board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
		legalMoves := board.AvailableBlackMoves()
		return len(legalMoves) == 0
	}

	return false
}

func (board Board) KingSafetyBonus() int8 {
	bonus := 0
	if board.WhiteCastled {
		bonus += 15 // Add 30 points for white castling
	}
	if board.BlackCastled {
		bonus -= 15 // Subtract 30 points when black castles, as lower is better for black
	}
	return int8(bonus)
}

func (board Board) GamePhase() float64 {
	totalPieces := bits.OnesCount64(uint64(board.WhitePawns)) + bits.OnesCount64(uint64(board.BlackPawns)) +
		bits.OnesCount64(uint64(board.WhiteKnights)) + bits.OnesCount64(uint64(board.BlackKnights)) +
		bits.OnesCount64(uint64(board.WhiteBishops)) + bits.OnesCount64(uint64(board.BlackBishops)) +
		bits.OnesCount64(uint64(board.WhiteRooks)) + bits.OnesCount64(uint64(board.BlackRooks)) +
		bits.OnesCount64(uint64(board.WhiteQueens)) + bits.OnesCount64(uint64(board.BlackQueens))

	return 1 - float64(totalPieces)/32.0
}

func (board Board) DynamicKingSafety() int8 {
	phase := board.GamePhase()
	return int8(float64(board.KingSafetyBonus()) * (0.5 + 0.5*phase))
}

type Evaluation struct {
	material      int8
	pawnPenalties int8
	mobilityBonus int8
	centreBonus   int8
	safety        int8
	knightBonus   int8
}

func (e Evaluation) Sum() int16 {
	var sum int16
	sum = sum + int16(e.material) + int16(e.pawnPenalties) + int16(e.mobilityBonus) + int16(e.centreBonus) + int16(e.safety) + int16(e.knightBonus)
	return sum
}

func (board Board) IsCheckMate() bool {
	// Check if the current player's king is in check
	var kingInCheck bool
	var kingMoves []Move

	if board.TurnBlack {
		kingInCheck = board.isKingInCheck(board.BlackKing, false)
		kingMoves = board.AvailableBlackMoves()
	} else {
		kingInCheck = board.isKingInCheck(board.WhiteKing, true)
		kingMoves = board.AvailableWhiteMoves()
	}

	// If the king is in check and no legal moves are available to get out of check
	return kingInCheck && len(kingMoves) == 0
}

// isKingInCheck checks if a given king is in check
func (board Board) isKingInCheck(king bitboards.KingBitboard, opponentBlack bool) bool {
	kingPosition := king.BitBoard()

	// Determine if the king is attacked by any opponent piece
	attacks := board.generateAttacks(opponentBlack)
	return attacks&kingPosition != 0
}

func (board Board) generateAttacks(opponentBlack bool) bitboards.BitBoard {
	if opponentBlack {
		return board.AvailableBlackAttacks()
	}

	return board.AvailableWhiteAttacks()
}

func (board Board) Evaluate() Evaluation {
	if board.IsCheckMate() {
		// Checkmate detection
		if board.TurnBlack {
			return Evaluation{-128, -128, -128, -128, -128, -128}
		} else {
			return Evaluation{127, 127, 127, 127, 127, 127}
		}
	}

	if board.IsStaleMate() {
		// Stalemate detection
		return Evaluation{0, 0, 0, 0, 0, 0}
	}

	material := (board.whiteMaterial() - board.blackMaterial()) * 8

	doubled := board.calculateDoubledPawns()
	blocked := board.calculateBlockedPawns()
	isolated := board.calculateIsolatedPawns()
	mobility := board.mobilityScore()
	centre := board.piecesInCentre()
	kingSafety := board.DynamicKingSafety()
	// misplacedKnights := board.knightsOnRim()

	pawnPenalties := doubled + blocked + isolated
	mobilityBonus := -1 * mobility
	centreBonus := -4 * centre
	// knightBonus := 0

	if board.TurnBlack {
		return Evaluation{material, pawnPenalties, mobilityBonus, -centreBonus, kingSafety, 0}
	} else {
		return Evaluation{-material, -pawnPenalties, -mobilityBonus, centreBonus, -kingSafety, 0}
	}
}
