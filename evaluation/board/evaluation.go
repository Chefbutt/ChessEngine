package board

import (
	"engine/evaluation/board/bitboards"
)

func (board *Board) whiteMaterial() int8 {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	wQ := board.WhiteQueens.BitBoard().PopCount()
	wR := board.WhiteRooks.BitBoard().PopCount()
	wN := board.WhiteKnights.BitBoard().PopCount()
	wB := board.WhiteBishops.BitBoard().PopCount()
	wP := board.WhitePawns.BitBoard().PopCount()

	return int8(queenWt*wQ + rookWt*wR + knightWt*wN + bishopWt*wB + pawnWt*wP)
}

func (board *Board) blackMaterial() int8 {
	queenWt, rookWt, knightWt, bishopWt, pawnWt := 9, 5, 3, 3, 1

	bQ := board.BlackQueens.BitBoard().PopCount()
	bR := board.BlackRooks.BitBoard().PopCount()
	bN := board.BlackKnights.BitBoard().PopCount()
	bB := board.BlackBishops.BitBoard().PopCount()
	bP := board.BlackPawns.BitBoard().PopCount()

	return int8(queenWt*bQ + rookWt*bR + knightWt*bN + bishopWt*bB + pawnWt*bP)
}

func (board *Board) calculateDoubledPawns() int8 {
	var doubled int8
	for file := 0; file < 8; file++ {
		fileMask := bitboards.FileMask(file)
		pawnsInFile := board.WhitePawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled += int8(pawnsInFile.PopCount())
		}

		pawnsInFile = board.BlackPawns.BitBoard() & fileMask
		if pawnsInFile.PopCount() > 1 {
			doubled -= int8(pawnsInFile.PopCount())
		}
	}

	return doubled
}

func (board *Board) calculateBlockedPawns() int8 {
	occupied := bitboards.BitBoard(board.WhitePawns) | bitboards.BitBoard(board.BlackPawns)
	blockedWhite := (bitboards.BitBoard(board.WhitePawns) << 8) & occupied
	blockedBlack := (bitboards.BitBoard(board.BlackPawns) >> 8) & occupied

	return int8(blockedWhite.PopCount() - blockedBlack.PopCount())
}

func (board *Board) calculateIsolatedPawns() int8 {
	isolatedWhite := board.WhitePawns
	isolatedBlack := board.BlackPawns

	// Check neighbors
	neighborFilesWhite := (board.WhitePawns << 1) | (board.WhitePawns >> 1)
	neighborFilesBlack := (board.BlackPawns << 1) | (board.BlackPawns >> 1)

	isolatedWhite &= ^neighborFilesWhite
	isolatedBlack &= ^neighborFilesBlack

	return int8(isolatedWhite.BitBoard().PopCount() - isolatedBlack.BitBoard().PopCount())
}

// Calculate mobility score (stub)
func (board *Board) mobilityScore() int8 {
	totalMoves := 0

	totalMoves += board.WhitePawns.Moves(board.EmptySquares, board.BlackPieces, board.EnPassantTarget).PopCount()
	totalMoves += board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces).PopCount()
	totalMoves += board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	totalMoves += board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	// totalMoves += board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces).PopCount()
	// totalMoves -= board.WhiteKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	totalMoves -= board.BlackPawns.Moves(board.EmptySquares, board.WhitePieces, board.EnPassantTarget).PopCount()
	totalMoves -= board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces).PopCount()
	totalMoves -= board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	totalMoves -= board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	// totalMoves -= board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces).PopCount()
	// totalMoves += board.BlackKing.Moves(board.EmptySquares, board.BlackPieces).PopCount()

	return int8(totalMoves)
}

func (board *Board) knightsOnRim() int8 {
	knightsOnRim := (board.WhiteKnights.BitBoard() & edgesOfBoard).PopCount()
	knightsOnRim = knightsOnRim - (board.BlackKnights.BitBoard() & edgesOfBoard).PopCount()
	return int8(-knightsOnRim)
}

func (board *Board) piecesInCentre() int8 {
	pieces := (board.WhitePawns.BitBoard() & centralSquares).PopCount()
	pieces = pieces - (board.BlackPawns.BitBoard() & centralSquares).PopCount()

	return int8(pieces)
}

// IsAttacked

func (board *Board) IsCheckMate() bool {
	if board.TurnBlack {
		attacks := board.WhiteAttacksMinimal()
		if board.IsAttacked(board.BlackKing.BitBoard(), attacks) {
			legalMoves := board.LegalMoves()
			return len(legalMoves) == 0
		}
	} else {
		attacks := board.BlackAttacksMinimal()
		if board.IsAttacked(board.WhiteKing.BitBoard(), attacks) {
			legalMoves := board.LegalMoves()
			return len(legalMoves) == 0
		}
	}
	return false
}

func (board *Board) IsStaleMate() bool {
	if !board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
		legalMoves := board.LegalMoves()
		return len(legalMoves) == 0
	}

	return false
}

func (board *Board) KingSafetyBonus() int8 {
	bonus := 0
	if board.CastleWhiteKingside || board.CastleWhiteQueenside {
		bonus += 30 // Add 30 points for white castling
	}
	if board.CastleBlackKingside || board.CastleBlackQueenside {
		bonus -= 30 // Subtract 30 points when black castles, as lower is better for black
	}
	return int8(bonus)
}

func (board *Board) GamePhase() float64 {
	totalPieces := board.WhitePawns.BitBoard().PopCount() + board.BlackPawns.BitBoard().PopCount() +
		board.WhiteKnights.BitBoard().PopCount() + board.BlackKnights.BitBoard().PopCount() +
		board.WhiteBishops.BitBoard().PopCount() + board.BlackBishops.BitBoard().PopCount() +
		board.WhiteRooks.BitBoard().PopCount() + board.BlackRooks.BitBoard().PopCount() +
		board.WhiteQueens.BitBoard().PopCount() + board.BlackQueens.BitBoard().PopCount()

	return 1 - float64(totalPieces)/32.0
}

func (board *Board) DynamicKingSafety() int8 {
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

func (board *Board) Evaluate() Evaluation {
	if board.IsCheckMate() {
		if board.TurnBlack {
			return Evaluation{-128, -128, -128, -128, -128, -128}
		} else {
			return Evaluation{127, 127, 127, 127, 127, 127}
		}
	} else if board.IsStaleMate() {
		return Evaluation{0, 0, 0, 0, 0, 0}
	}

	material := (board.whiteMaterial() - board.blackMaterial()) * -8

	doubled := board.calculateDoubledPawns()
	blocked := board.calculateBlockedPawns()
	isolated := board.calculateIsolatedPawns()
	mobility := board.mobilityScore()
	centre := board.piecesInCentre()
	kingSafety := board.DynamicKingSafety()
	misplacedKnights := board.knightsOnRim()

	pawnPenalties := doubled + blocked + isolated
	mobilityBonus := -1 * mobility
	centreBonus := -4 * centre
	knightBonus := -5 * misplacedKnights
	safety := -5 * kingSafety

	if board.TurnBlack {
		return Evaluation{material, pawnPenalties, mobilityBonus, centreBonus, safety, knightBonus}
	} else {
		return Evaluation{-material, -pawnPenalties, -mobilityBonus, -centreBonus, -safety, -knightBonus}
	}
}
