package board

import (
	"sort"

	"engine/evaluation/board/bitboards"
)

var (
	centralSquares = bitboards.BitBoard(0x3C3C3C3C0000)
	edgesOfBoard   = bitboards.BitBoard(0x8181818181818181)
)

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
		return (moves[i].MoveType == CastleKingside || moves[i].MoveType == CastleQueenside) && (moves[j].MoveType != CastleKingside || moves[j].MoveType != CastleQueenside)
	})

	return moves
}

func ControlCentre(board *Board) []Move {
	moves := board.LegalMoves()

	sort.Slice(moves, func(i, j int) bool {
		return (bitboards.New(moves[i].Destination)) < (bitboards.New(moves[j].Destination))
	})

	return moves
}

func (board Board) LegalMoves() []Move {
	moveMap := make(map[bitboards.BitBoard][]bitboards.BitBoard)

	allMoves := board.RegularMoves()
	if board.TurnBlack {
		attacks := board.WhiteAttacksMinimal()
		if board.IsAttacked(board.BlackKing.BitBoard(), attacks) {
			moveMap[board.BlackKing.BitBoard()] = (board.BlackKing.Moves(board.EmptySquares, board.WhitePieces) &^ attacks).Split()

			for piece, move := range StopCheck(board, allMoves) {
				moveMap[piece] = move
			}
			return board.BitBoardMapToMove(board.WhitePieces, moveMap)
		}
		return board.BitBoardMapToMove(board.WhitePieces, allMoves)
	} else {
		attacks := board.BlackAttacksMinimal()
		if board.IsAttacked(board.WhiteKing.BitBoard(), attacks) {
			moveMap[board.WhiteKing.BitBoard()] = (board.WhiteKing.Moves(board.EmptySquares, board.BlackPieces) &^ attacks).Split()

			for piece, move := range StopCheck(board, allMoves) {
				moveMap[piece] = move
			}

			return board.BitBoardMapToMove(board.BlackPieces, moveMap)
		}
		return board.BitBoardMapToMove(board.BlackPieces, allMoves)
	}
}
