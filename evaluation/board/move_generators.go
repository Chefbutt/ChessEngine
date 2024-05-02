package board

import (
	"sort"

	"engine/evaluation/board/bitboards"
)

var centralSquares = bitboards.BitBoard(0x3C3C3C3C0000)

func Captures(board *Board) []Move {
	moves := board.LegalMoves()

	var filtered []Move
	for _, move := range moves {
		if board.PieceAt(move.Destination) != -1 {
			filtered = append(filtered, move)
		}
	}

	return filtered
}

func OrderedMoves(board *Board) []Move {
	moves := board.LegalMoves()

	for id, move := range moves {
		moves[id].CapturedPiece = board.PieceAt(move.Destination)
		undo, err := board.makeMove(move)
		if err != nil {
			panic(err)
		}
		// board.Display()
		if board.IsAttacked(board.KingInPlayAndOpponentAttacks()) {
			moves[id].IsCheck = true
		}
		board.UndoMove(undo)
	}

	sort.SliceStable(moves, func(i, j int) bool {
		// board.PieceAt(moves[i].Destination)

		// First, prioritize capturing moves
		if moves[i].CapturedPiece != -1 && moves[j].CapturedPiece == -1 {
			return true
		} else if moves[i].CapturedPiece == -1 && moves[j].CapturedPiece != -1 {
			return false
		}

		// Then, prioritize check moves
		if moves[i].IsCheck && !moves[j].IsCheck {
			return true
		} else if !moves[i].IsCheck && moves[j].IsCheck {
			return false
		}

		// Finally, prioritize quiet moves
		return false
	})

	return moves
}

func ControlCentre(board *Board) []Move {
	moves := board.LegalMoves()

	sort.Slice(moves, func(i, j int) bool {
		return bitboards.New(moves[i].Destination) < bitboards.New(moves[j].Destination)
	})

	return moves
}

func (board Board) LegalMoves() []Move {
	if board.TurnBlack {
		return board.AvailableBlackMoves()
	} else {
		return board.AvailableWhiteMoves()
	}
}
