package board

import (
	"fmt"

	"engine/board/bitboards"
)

type Move struct {
	Source         int
	Destination    int
	Piece          int // Use constants to define pieces
	CapturedPiece  int
	MoveType       int // Use constants to define move types
	PromotionPiece int // Added for promotion clarity
}

// Converts chess position e.g., "e2" to board index 0-63
func positionToIndex(pos string) int {
	file := pos[0] - 'a' // Convert 'a' to 'h' to 0-7
	rank := pos[1] - '1' // Convert '1' to '8' to 0-7
	return int(rank*8 + file)
}

func mapPromotionPiece(char byte, isBlackTurn bool) int {
	switch char {
	case 'q':
		if isBlackTurn {
			return BlackQueen
		}
		return WhiteQueen
	case 'r':
		if isBlackTurn {
			return BlackRook
		}
		return WhiteRook
	case 'b':
		if isBlackTurn {
			return BlackBishop
		}
		return WhiteBishop
	case 'n':
		if isBlackTurn {
			return BlackKnight
		}
		return WhiteKnight
	default:
		panic("Invalid promotion piece")
	}
}

func (b *Board) isOccupied(pos int) bool {
	occupiedMask := bitboards.New(pos)
	return (b.OccupiedSquares & occupiedMask) != 0
}

func identifyPieceAt(board *Board, index int) int {
	indexMask := bitboards.New(index)

	if (board.WhitePawns.BitBoard() & indexMask) != 0 {
		return WhitePawn
	} else if (board.BlackPawns.BitBoard() & indexMask) != 0 {
		return BlackPawn
	} else if (board.WhiteKnights.BitBoard() & indexMask) != 0 {
		return WhiteKnight
	} else if (board.BlackKnights.BitBoard() & indexMask) != 0 {
		return BlackKnight
	} else if (board.WhiteBishops.BitBoard() & indexMask) != 0 {
		return WhiteBishop
	} else if (board.BlackBishops.BitBoard() & indexMask) != 0 {
		return BlackBishop
	} else if (board.WhiteRooks.BitBoard() & indexMask) != 0 {
		return WhiteRook
	} else if (board.BlackRooks.BitBoard() & indexMask) != 0 {
		return BlackRook
	} else if (board.WhiteQueens.BitBoard() & indexMask) != 0 {
		return WhiteQueen
	} else if (board.BlackQueens.BitBoard() & indexMask) != 0 {
		return BlackQueen
	} else if (board.WhiteKing.BitBoard() & indexMask) != 0 {
		return WhiteKing
	} else if (board.BlackKing.BitBoard() & indexMask) != 0 {
		return BlackKing
	}

	return -1
}

func (board *Board) UCItoMove(uci string) Move {
	if len(uci) < 4 {
		panic("Invalid UCI string")
	}

	source := positionToIndex(uci[0:2])
	destination := positionToIndex(uci[2:4])

	piece := identifyPieceAt(board, source)
	capturedPiece := -1 // Use -1 when no piece is captured

	// Determine if it's a capture
	isCapture := board.isOccupied(destination)
	moveType := NormalMove
	if isCapture {
		moveType = Capture
		capturedPiece = identifyPieceAt(board, destination)
	}

	if (piece == WhiteKing || piece == BlackKing) && (source+2 == destination) {
		moveType = CastleKingside
	}

	if (piece == WhiteKing || piece == BlackKing) && (source-3 == destination) {
		moveType = CastleQueenside
	}

	// Check for promotion (indicated by a fifth character)
	promotionPiece := -1
	if len(uci) == 5 {
		moveType = Promotion
		promotionPiece = mapPromotionPiece(uci[4], board.TurnBlack)
	}

	return Move{
		Source:         source,
		Destination:    destination,
		Piece:          piece,
		CapturedPiece:  capturedPiece,
		MoveType:       moveType,
		PromotionPiece: promotionPiece,
	}
}

const (
	// Pieces
	WhitePawn = iota
	BlackPawn
	WhiteKnight
	BlackKnight
	WhiteBishop
	BlackBishop
	WhiteRook
	BlackRook
	WhiteQueen
	BlackQueen
	WhiteKing
	BlackKing

	// Move types
	NormalMove
	Capture
	EnPassant
	CastleKingside

	CastleQueenside
	Promotion
)

func (b *Board) pieceBitboard(piece int) *bitboards.BitBoard {
	switch piece {
	case WhitePawn:
		return b.WhitePawns.BitBoardPointer()
	case BlackPawn:
		return b.BlackPawns.BitBoardPointer()
	case WhiteKnight:
		return b.WhiteKnights.BitBoardPointer()
	case BlackKnight:
		return b.BlackKnights.BitBoardPointer()
	case WhiteBishop:
		return b.WhiteBishops.BitBoardPointer()
	case BlackBishop:
		return b.BlackBishops.BitBoardPointer()
	case WhiteRook:
		return b.WhiteRooks.BitBoardPointer()
	case BlackRook:
		return b.BlackRooks.BitBoardPointer()
	case WhiteQueen:
		return b.WhiteQueens.BitBoardPointer()
	case BlackQueen:
		return b.BlackQueens.BitBoardPointer()
	case WhiteKing:
		return b.WhiteKing.BitBoardPointer()
	case BlackKing:
		return b.BlackKing.BitBoardPointer()
	default:
		panic("Invalid piece type")
	}
	return nil
}

func (b *Board) updateAggregateBitboards() {
	b.WhitePieces = b.WhitePawns.BitBoard() | b.WhiteKnights.BitBoard() | b.WhiteBishops.BitBoard() | b.WhiteRooks.BitBoard() | b.WhiteQueens.BitBoard() | b.WhiteKing.BitBoard()
	b.BlackPieces = b.BlackPawns.BitBoard() | b.BlackKnights.BitBoard() | b.BlackBishops.BitBoard() | b.BlackRooks.BitBoard() | b.BlackQueens.BitBoard() | b.BlackKing.BitBoard()
	b.OccupiedSquares = b.WhitePieces | b.BlackPieces
	b.EmptySquares = ^b.OccupiedSquares
}

func (board *Board) MakeMove(move string) error {
	parsedMove := board.UCItoMove(move)

	err := board.makeMove(parsedMove)
	if err != nil {
		return err
	}

	fmt.Println()

	board.OccupiedSquares.Display()

	return nil
}

// Optimize to know the colour beforehand
func (board *Board) makeMove(move Move) error {
	sourceBit := bitboards.New(move.Source)
	destBit := bitboards.New(move.Destination)

	// Clear source and set destination for the moving piece
	*board.pieceBitboard(move.Piece) &= ^sourceBit
	*board.pieceBitboard(move.Piece) |= destBit

	// Handle capture
	if move.MoveType == Capture {
		*board.pieceBitboard(move.CapturedPiece) &= ^destBit
		board.updateAggregateBitboards()
		board.TurnBlack = !board.TurnBlack
		return nil
	}

	// Special moves handling
	switch move.MoveType {
	case EnPassant:
		// Clear the pawn in passing
		capturedPawnBit := bitboards.New(move.Destination - 8) // or +8 depending on direction
		*board.pieceBitboard(BlackPawn) &= ^capturedPawnBit    // or WhitePawn
	case CastleKingside:
		// Update both king and rook positions, assuming kingside castle for example
		if move.Piece == WhiteKing {
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(7) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(5)  // new rook position for kingside
		}
		if move.Piece == BlackKing {
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(63) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(61)  // new rook position for kingside
		}
	case CastleQueenside:
		if move.Piece == WhiteKing {
			*board.pieceBitboard(WhiteRook) &= ^bitboards.New(0) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= bitboards.New(3)  // new rook position for kingside
		}
		if move.Piece == BlackKing {
			*board.pieceBitboard(BlackRook) &= ^bitboards.New(55) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= bitboards.New(58)  // new rook position for kingside
		}
	case Promotion:
		// Assume promoting to a queen for simplicity
		*board.pieceBitboard(move.Piece) &= ^destBit // Remove pawn from destination
		*board.pieceBitboard(WhiteQueen) |= destBit  // Add queen to destination
	default:
		panic("")
	}

	// Update occupied, empty, and aggregate piece bitboards
	board.updateAggregateBitboards()

	// Toggle turn
	board.TurnBlack = !board.TurnBlack

	return nil
}

func NewBoard() Board {
	whitePawns := bitboards.WhitePawnBitboard(0x000000000000FF00)
	blackPawns := bitboards.BlackPawnBitboard(0x00FF000000000000)

	whiteKnights := bitboards.KnightBitboard(0x0000000000000042)
	blackKnights := bitboards.KnightBitboard(0x4200000000000000)

	whiteBishops := bitboards.BishopBitboard(0x0000000000000024)
	blackBishops := bitboards.BishopBitboard(0x2400000000000000)

	whiteRooks := bitboards.RookBitboard(0x0000000000000081)
	blackRooks := bitboards.RookBitboard(0x8100000000000000)

	whiteQueens := bitboards.QueenBitboard(0x0000000000000008)
	blackQueens := bitboards.QueenBitboard(0x0800000000000000)

	whiteKing := bitboards.KingBitboard(0x0000000000000010)
	blackKing := bitboards.KingBitboard(0x1000000000000000)

	// Setting up the occupancy based on all pieces
	occupiedSquares := bitboards.BitBoard(0xFFFF00000000FFFF)
	emptySquares := ^occupiedSquares // All squares not occupied

	whitePieces := bitboards.BitBoard(0xFFFF000000000000)
	blackPieces := bitboards.BitBoard(0x000000000000FFFF)

	enPassantTarget := bitboards.BitBoard(0) // No en passant possible on first move
	castlingRights := uint8(0xF)             // All castling rights available initially (both kingside and queenside for both colors)

	turnBlack := false // White moves first in chess

	return Board{
		WhitePawns:      whitePawns,
		BlackPawns:      blackPawns,
		WhiteKnights:    whiteKnights,
		BlackKnights:    blackKnights,
		WhiteBishops:    whiteBishops,
		BlackBishops:    blackBishops,
		WhiteRooks:      whiteRooks,
		BlackRooks:      blackRooks,
		WhiteQueens:     whiteQueens,
		BlackQueens:     blackQueens,
		WhiteKing:       whiteKing,
		BlackKing:       blackKing,
		OccupiedSquares: occupiedSquares,
		EmptySquares:    emptySquares,
		WhitePieces:     whitePieces,
		BlackPieces:     blackPieces,
		EnPassantTarget: enPassantTarget,
		CastlingRights:  castlingRights,
		TurnBlack:       turnBlack,
	}
}
