package board

import (
	"errors"
	"fmt"

	"engine/evaluation/board/bitboards"
)

func (board *Board) MakeMove(move string) error {
	parsedMove := board.UCItoMove(move)

	err := board.makeMove(parsedMove)
	if err != nil {
		return err
	}

	fmt.Println()

	board.Display()

	return nil
}

func (board *Board) MakeNativeMove(move Move) error {
	if move.IsValid() {
		return errors.New("invalid move")
	}

	err := board.makeMove(move)
	if err != nil {
		return err
	}

	fmt.Println()
	board.Display()

	return nil
}

// Optimize to know the colour beforehand
func (board *Board) makeMove(move Move) error {
	if board.EmptySquares&board.OccupiedSquares != 0 {
		panic("invalid board state")
	}
	sourceBit := bitboards.New(move.Source)
	destBit := bitboards.New(move.Destination)

	*board.pieceBitboard(move.Piece) &= ^sourceBit
	*board.pieceBitboard(move.Piece) |= destBit

	switch move.MoveType {
	case Capture:
		*board.pieceBitboard(move.CapturedPiece) &= ^destBit
		board.updateAggregateBitboards()
		board.TurnBlack = !board.TurnBlack
		return nil
	case EnPassant:
		if move.Piece == WhitePawn {
			capturedPawnBit := bitboards.New(move.Destination - 8) // or +8 depending on direction
			*board.pieceBitboard(WhitePawn) &= ^capturedPawnBit    // or WhitePawn
		}

		if move.Piece == BlackPawn {
			capturedPawnBit := bitboards.New(move.Destination + 8) // or +8 depending on direction
			*board.pieceBitboard(BlackPawn) &= ^capturedPawnBit    // or WhitePawn
		}
	case CastleKingside:
		if move.Piece == WhiteKing && board.CastleWhite {
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(7) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(5)  // new rook position for kingside
			board.CastleWhite = false
		}
		if move.Piece == BlackKing && board.CastleBlack {
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(63) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(61)  // new rook position for kingside
			board.CastleBlack = false
		}
	case CastleQueenside:
		if move.Piece == WhiteKing && board.CastleWhite {
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(0) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(3)  // new rook position for kingside
			board.CastleWhite = false
		}
		if move.Piece == BlackKing && board.CastleBlack {
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(55) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(58)  // new rook position for kingside
			board.CastleBlack = false
		}
	case Promotion:
		*board.pieceBitboard(move.Piece) &= ^destBit         // Remove pawn from destination
		*board.pieceBitboard(move.PromotionPiece) |= destBit // Add queen to destination
	case NormalMove:
	default:
		panic("")
	}

	board.updateAggregateBitboards()
	board.TurnBlack = !board.TurnBlack
	board.HalfTurn++

	return nil
}
