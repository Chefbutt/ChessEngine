package board

import (
	"errors"
	"fmt"

	"engine/evaluation/board/bitboards"
)

func IndexToPosition(bitboard uint64) string {
	rank := bitboard % 8
	file := bitboard / 8
	return string('a'+rune(rank)) + string('1'+rune(file))
}

func (board *Board) MakeHumanMove(move string) error {
	parsedMove := board.UCItoMove(move)
	if parsedMove.Source == parsedMove.Destination {
		return errors.New("illegal move")
	}

	moves := board.AvailableWhiteMoves()

	canMake := false
	for _, move := range moves {
		if parsedMove.Source == move.Source && parsedMove.Destination == move.Destination {
			canMake = true
			break
		}
	}
	if !canMake {
		return errors.New("illegal move")
	}
	// bestMove, eval := board.BestMove(12, OrderedMoves)

	// fmt.Print(PieceSymbols[board.PieceAt(int(bestMove.Source))], "(", IndexToPosition(uint64(bestMove.Source)), "): ", IndexToPosition(uint64(bestMove.Destination)), " ", eval)

	if board.EnPassantTarget&bitboards.BitBoard(parsedMove.Destination) != 0 {
		parsedMove.MoveType = EnPassant
	}

	_, err := board.makeMove(parsedMove)
	if err != nil {
		return err
	}

	// board.EvaluationDetails()

	fmt.Println()

	board.Display()

	return nil
}

// Create a fixed number of workers

func (board *Board) MakeMove(depth int) error {
	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	log.Fatal("could not create CPU profile: ", err)
	// }
	// defer f.Close() // error handling omitted for example
	// if err := pprof.StartCPUProfile(f); err != nil {
	// 	log.Fatal("could not start CPU profile: ", err)
	// }
	// defer pprof.StopCPUProfile()

	// fM, err := os.Create("mem.prof")
	// if err != nil {
	// 	log.Fatal("could not create memory profile: ", err)
	// }
	// defer fM.Close() // error handling omitted for example
	// runtime.GC()     // get up-to-date statistics
	// if err := pprof.WriteHeapProfile(fM); err != nil {
	// 	log.Fatal("could not write memory profile: ", err)
	// }

	bestMove, eval := board.BestMove(depth, OrderedMoves, -8, 3, 8, 1)

	if board.TurnBlack {
		if board.LastMoveWhite != bestMove {
			board.Repetition = 0
		}
	} else {
		if board.LastMoveWhite != bestMove {
			board.Repetition = 0
		}
	}

	if board.TurnBlack {
		if board.LastMoveBlack == bestMove {
			board.Repetition = board.Repetition + 1
		}
		board.LastMoveBlack = bestMove
	} else {
		if board.LastMoveWhite == bestMove {
			board.Repetition = board.Repetition + 1
		}
		board.LastMoveWhite = bestMove
	}

	if board.Repetition == 3 {
		fmt.Println("Draw by threefold repetition")
		return errors.New("Repetition")
	}

	if board.Debug {
		fmt.Println(PieceSymbols[board.PieceAt(int(bestMove.Source))], "(", IndexToPosition(uint64(bestMove.Destination)), ") : ", eval)
	}

	if bestMove.Source == bestMove.Destination {
		if board.TurnBlack {
			fmt.Println("Black resigns")
		} else {
			fmt.Println("White resigns")
		}
		return errors.New("Resign")
	}

	_, err := board.makeMove(bestMove)
	if err != nil {
		return err
	}

	return nil
}

func (board *Board) MakeNativeMove(move Move) (*MoveUndo, error) {
	if move.IsValid() {
		return nil, errors.New("invalid move")
	}

	undo, err := board.makeMove(move)
	if err != nil {
		return undo, err
	}

	return undo, nil
}

type AggregateBitboards struct {
	WhitePieces     bitboards.BitBoard
	BlackPieces     bitboards.BitBoard
	OccupiedSquares bitboards.BitBoard
	EnPassantTarget bitboards.BitBoard
	EmptySquares    bitboards.BitBoard
}

// Optimize to know the colour beforehand
func (board *Board) makeMove(move Move) (*MoveUndo, error) {
	if board.EmptySquares&board.OccupiedSquares != 0 {
		panic("invalid board state")
	}

	if move.Source == move.Destination {
		panic("")
	}

	move.CapturedPiece = board.PieceAt(move.Destination)

	if move.Piece == -1 {
		move.Piece = board.PieceAt(move.Source)
	}

	if move.MoveType == 0 && move.CapturedPiece != -1 {
		move.MoveType = Capture
	}

	if move.MoveType == 0 && move.CapturedPiece == -1 {
		move.MoveType = NormalMove
	}

	undo := MoveUndo{
		Source:                       move.Source,
		Destination:                  move.Destination,
		Piece:                        move.Piece,
		CapturedPiece:                move.CapturedPiece,
		PromotionPiece:               move.PromotionPiece,
		MoveType:                     move.MoveType,
		PreviousCastleWhiteKingside:  board.CastleWhiteKingside,
		PreviousCastleWhiteQueenside: board.CastleWhiteQueenside,
		PreviousCastleBlackKingside:  board.CastleBlackKingside,
		PreviousCastleBlackQueenside: board.CastleBlackQueenside,
		PreviousTurnBlack:            board.TurnBlack,
		PreviousHalfTurn:             board.HalfTurn,
		PreviousAggregateBitboards:   board.AggregateBitboards(), // Example, assume this captures all necessary board pieces
	}

	sourceBit := bitboards.New(move.Source)
	destBit := bitboards.New(move.Destination)

	*board.pieceBitboard(move.Piece) &= ^sourceBit
	*board.pieceBitboard(move.Piece) |= destBit

	switch move.MoveType {
	case Capture:
		undo.CapturedPiece = move.CapturedPiece
		*board.pieceBitboard(move.CapturedPiece) &= ^destBit
		board.updateAggregateBitboards()
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
		if !board.CastleWhiteKingside && !board.CastleBlackKingside {
			panic("")
		}
		if move.Piece == WhiteKing && board.CastleWhiteKingside {
			board.WhiteCastled = true
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(7) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(5)  // new rook position for kingside
			board.CastleWhiteKingside = false
			board.CastleWhiteQueenside = false
		}
		if move.Piece == BlackKing && board.CastleBlackKingside {
			board.BlackCastled = true
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(63) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(61)  // new rook position for kingside
			board.CastleBlackKingside = false
			board.CastleBlackQueenside = false
		}
	case CastleQueenside:
		if !board.CastleWhiteQueenside && !board.CastleBlackQueenside {
			panic("?")
		}
		if move.Piece == WhiteKing && board.CastleWhiteQueenside {
			board.BlackCastled = true
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(0) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(3)  // new rook position for kingside
			board.CastleWhiteKingside = false
			board.CastleWhiteQueenside = false
		}
		if move.Piece == BlackKing && board.CastleBlackQueenside {
			board.WhiteCastled = true
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(56) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(59)  // new rook position for kingside
			board.CastleBlackKingside = false
			board.CastleBlackQueenside = false
		}
	case Promotion:
		*board.pieceBitboard(move.Piece) &= ^destBit         // Remove pawn from destination
		*board.pieceBitboard(move.PromotionPiece) |= destBit // Add queen to destination
	case NormalMove:
	default:
		panic("")
	}

	if move.Piece == BlackKing {
		board.CastleBlackKingside = false
		board.CastleBlackQueenside = false
	}

	if move.Piece == WhiteKing {
		board.CastleWhiteKingside = false
		board.CastleWhiteQueenside = false
	}

	if move.Source == 7 {
		board.CastleWhiteKingside = false
	}

	if move.Source == 0 {
		board.CastleWhiteQueenside = false
	}

	if move.Source == 63 {
		board.CastleBlackKingside = false
	}

	if move.Source == 55 {
		board.CastleBlackQueenside = false
	}

	board.updateAggregateBitboards()
	board.TurnBlack = !board.TurnBlack
	board.HalfTurn++

	return &undo, nil
}
