package board

import (
	"sort"

	"engine/evaluation/board/bitboards"
)

var centralSquares = bitboards.BitBoard(0x3C3C000000)

func OrderedMoves(board Board) []Move {
	moves := board.LegalMoves()

	for id, move := range moves {
		board := board
		moves[id].CapturedPiece = board.PieceAt(move.Destination)
		undo, err := board.makeMove(move)
		if err != nil {
			panic(err)
		}
		// board.Display()
		if board.TurnBlack {
			if board.isKingInCheck(board.BlackKing, false) {
				moves[id].IsCheck = true
			}
		} else {
			if board.isKingInCheck(board.WhiteKing, true) {
				moves[id].IsCheck = true
			}
		}

		board.UndoMove(undo)
	}

	sort.SliceStable(moves, func(i, j int) bool {
		// board.PieceAt(moves[i].Destination)

		// Then, prioritize check moves
		if moves[i].IsCheck && !moves[j].IsCheck {
			return true
		} else if !moves[i].IsCheck && moves[j].IsCheck {
			return false
		}

		// First, prioritize capturing moves
		if moves[i].CapturedPiece != -1 && moves[j].CapturedPiece == -1 {
			return true
		} else if moves[i].CapturedPiece == -1 && moves[j].CapturedPiece != -1 {
			return false
		}

		// Finally, prioritize quiet moves
		return false
	})

	return moves
}

func (board Board) LegalMoves() []Move {
	if board.TurnBlack {
		if board.BlackKing == 0 {
			return nil
		}
		return board.AvailableBlackMoves()
	} else {
		if board.WhiteKing == 0 {
			return nil
		}
		return board.AvailableWhiteMoves()
	}
}

func (board Board) Captures() []Move {
	if board.WhiteKing == 0 || board.BlackKing == 0 {
		return nil
	}

	if board.TurnBlack {
		if board.BlackKing == 0 {
			return nil
		}
		return board.AvailableBlackCaptures()
	} else {
		if board.WhiteKing == 0 {
			return nil
		}
		return board.AvailableWhiteCaptures()
	}
}
