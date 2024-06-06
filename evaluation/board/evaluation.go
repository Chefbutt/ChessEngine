package board

import (
	"math"
	"math/bits"

	"engine/evaluation/board/bitboards"
)

func (board Board) whiteMaterial() int {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	wQ := bits.OnesCount64(uint64(board.WhiteQueens))
	wR := bits.OnesCount64(uint64(board.WhiteRooks))
	wN := bits.OnesCount64(uint64(board.WhiteKnights))
	wB := bits.OnesCount64(uint64(board.WhiteBishops))
	wP := bits.OnesCount64(uint64(board.WhitePawns))

	return queenWt*wQ + rookWt*wR + knightWt*wN + bishopWt*wB + pawnWt*wP
}

func (board Board) blackMaterial() int {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	bQ := bits.OnesCount64(uint64(board.BlackQueens))
	bR := bits.OnesCount64(uint64(board.BlackRooks))
	bN := bits.OnesCount64(uint64(board.BlackKnights))
	bB := bits.OnesCount64(uint64(board.BlackBishops))
	bP := bits.OnesCount64(uint64(board.BlackPawns))

	return queenWt*bQ + rookWt*bR + knightWt*bN + bishopWt*bB + pawnWt*bP
}

func (board Board) calculateDoubledPawns() int {
	var doubled int
	for file := 0; file < 8; file++ {
		fileMask := bitboards.FileMask(file)
		pawnsInFile := board.WhitePawns.BitBoard() & fileMask
		if bits.OnesCount64(uint64(pawnsInFile)) > 1 {
			doubled += bits.OnesCount64(uint64(pawnsInFile))
		}

		pawnsInFile = board.BlackPawns.BitBoard() & fileMask
		if bits.OnesCount64(uint64(pawnsInFile)) > 1 {
			doubled -= bits.OnesCount64(uint64(pawnsInFile))
		}
	}

	return doubled
}

func (board Board) calculateBlockedPawns() int {
	occupied := bitboards.BitBoard(board.WhitePawns) | bitboards.BitBoard(board.BlackPawns)
	blockedWhite := (bitboards.BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (bitboards.BitBoard(board.BlackPawns) >> 8) & occupied

	return bits.OnesCount64(uint64(blockedWhite)) - bits.OnesCount64(uint64(blockedBlack))
}

func (board Board) calculateIsolatedPawns() int {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return bits.OnesCount64(uint64(isolatedWhite)) - bits.OnesCount64(uint64(isolatedBlack))
}

// Calculate mobility score (stub)
func (board Board) mobilityScore() int {
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

	return totalMoves
}

var edgesOfBoard = bitboards.BitBoard(0x8181818181818181)

func (board Board) knightsOnRim() int {
	knightsOnRim := bits.OnesCount64(uint64(board.WhiteKnights.BitBoard()&edgesOfBoard)) - bits.OnesCount64(uint64(board.BlackKnights.BitBoard()&edgesOfBoard))
	return knightsOnRim
}

func (board Board) piecesInCentre() int {
	pieces := bits.OnesCount64(uint64(board.WhitePieces & centralSquares))

	return pieces - bits.OnesCount64(uint64(board.BlackPieces&centralSquares))
}

// IsAttacked

func (board Board) IsStaleMate() bool {
	if !board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
		legalMoves := board.AvailableBlackMoves()
		return len(legalMoves) == 0
	}

	return false
}

func (board Board) KingSafetyBonus() int {
	bonus := 0
	if board.WhiteCastled {
		bonus += 15 // Add 30 points for white castling
	}
	if board.BlackCastled {
		bonus -= 15 // Subtract 30 points when black castles, as lower is better for black
	}
	return bonus
}

func (board Board) GamePhase() float64 {
	totalPieces := bits.OnesCount64(uint64(board.WhitePawns)) + bits.OnesCount64(uint64(board.BlackPawns)) +
		bits.OnesCount64(uint64(board.WhiteKnights)) + bits.OnesCount64(uint64(board.BlackKnights)) +
		bits.OnesCount64(uint64(board.WhiteBishops)) + bits.OnesCount64(uint64(board.BlackBishops)) +
		bits.OnesCount64(uint64(board.WhiteRooks)) + bits.OnesCount64(uint64(board.BlackRooks)) +
		bits.OnesCount64(uint64(board.WhiteQueens)) + bits.OnesCount64(uint64(board.BlackQueens))

	return 1 - float64(totalPieces)/32.0
}

func (board Board) DynamicKingSafety() float64 {
	phase := board.GamePhase()
	return (0.5 + 0.5*phase) * float64(board.KingSafetyBonus())
}

func (board Board) IsCheckMate() bool {
	// Check if the current player's king is in check
	var kingInCheck bool
	var kingMoves []Move

	if board.TurnBlack {
		if board.BlackKing == 0 {
			return true
		}
		kingInCheck = board.isKingInCheck(board.BlackKing, false)
		kingMoves = board.AvailableBlackMoves()
	} else {
		if board.WhiteKing == 0 {
			return true
		}
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

func (board Board) Evaluate(materialModifier, mobilityModifier, centreModifier, penaltyModifier int) int {
	if board.IsStaleMate() {
		// Stalemate detection
		return 0
	}

	material := (board.whiteMaterial() - board.blackMaterial()) * materialModifier // 8

	doubled := board.calculateDoubledPawns()
	blocked := board.calculateBlockedPawns()
	isolated := board.calculateIsolatedPawns()
	mobility := board.mobilityScore() / 2
	centre := board.piecesInCentre() * 2
	kingSafety := board.DynamicKingSafety()
	misplacedKnights := board.knightsOnRim()
	aggression := board.GamePhase()

	pawnPenalties := doubled + blocked + isolated
	mobilityBonus := mobilityModifier * mobility //-1
	centreBonus := centreModifier * centre
	// knightBonus := 0

	if board.TurnBlack {
		return (-material - penaltyModifier*pawnPenalties - mobilityBonus + centreBonus - int(math.Round(kingSafety*3)) - misplacedKnights - int(math.Round(aggression*3)))
	} else {
		return (material + penaltyModifier*pawnPenalties + mobilityBonus - centreBonus + int(math.Round(kingSafety*3)) + misplacedKnights + int(math.Round(aggression*3)))
	}
}
