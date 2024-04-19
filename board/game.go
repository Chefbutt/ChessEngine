package board

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
	occupiedMask := BitBoard(1) << pos
	return (b.OccupiedSquares & occupiedMask) != 0
}

func identifyPieceAt(board *Board, index int) int {
	indexMask := BitBoard(1) << index

	if (BitBoard(board.WhitePawns) & indexMask) != 0 {
		return WhitePawn
	} else if (BitBoard(board.BlackPawns) & indexMask) != 0 {
		return BlackPawn
	} else if (BitBoard(board.WhiteKnights) & indexMask) != 0 {
		return WhiteKnight
	} else if (BitBoard(board.BlackKnights) & indexMask) != 0 {
		return BlackKnight
	} else if (BitBoard(board.WhiteBishops) & indexMask) != 0 {
		return WhiteBishop
	} else if (BitBoard(board.BlackBishops) & indexMask) != 0 {
		return BlackBishop
	} else if (BitBoard(board.WhiteRooks) & indexMask) != 0 {
		return WhiteRook
	} else if (BitBoard(board.BlackRooks) & indexMask) != 0 {
		return BlackRook
	} else if (BitBoard(board.WhiteQueens) & indexMask) != 0 {
		return WhiteQueen
	} else if (BitBoard(board.BlackQueens) & indexMask) != 0 {
		return BlackQueen
	} else if (BitBoard(board.WhiteKing) & indexMask) != 0 {
		return WhiteKing
	} else if (BitBoard(board.BlackKing) & indexMask) != 0 {
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

func (b *Board) pieceBitboard(piece int) *BitBoard {
	var bitboard *BitBoard
	switch piece {
	case WhitePawn:
		bitboard = (*BitBoard)(&b.WhitePawns)
	case BlackPawn:
		bitboard = (*BitBoard)(&b.BlackPawns)
	case WhiteKnight:
		bitboard = (*BitBoard)(&b.WhiteKnights)
	case BlackKnight:
		bitboard = (*BitBoard)(&b.BlackKnights)
	case WhiteBishop:
		bitboard = (*BitBoard)(&b.WhiteBishops)
	case BlackBishop:
		bitboard = (*BitBoard)(&b.BlackBishops)
	case WhiteRook:
		bitboard = (*BitBoard)(&b.WhiteRooks)
	case BlackRook:
		bitboard = (*BitBoard)(&b.BlackRooks)
	case WhiteQueen:
		bitboard = (*BitBoard)(&b.WhiteQueens)
	case BlackQueen:
		bitboard = (*BitBoard)(&b.BlackQueens)
	case WhiteKing:
		bitboard = (*BitBoard)(&b.WhiteKing)
	case BlackKing:
		bitboard = (*BitBoard)(&b.BlackKing)
	default:
		panic("Invalid piece type")
	}
	return bitboard
}

func (b *Board) updateAggregateBitboards() {
	b.WhitePieces = BitBoard(b.WhitePawns) | BitBoard(b.WhiteKnights) | BitBoard(b.WhiteBishops) | BitBoard(b.WhiteRooks) | BitBoard(b.WhiteQueens) | BitBoard(b.WhiteKing)
	b.BlackPieces = BitBoard(b.BlackPawns) | BitBoard(b.BlackKnights) | BitBoard(b.BlackBishops) | BitBoard(b.BlackRooks) | BitBoard(b.BlackQueens) | BitBoard(b.BlackKing)
	b.OccupiedSquares = b.WhitePieces | b.BlackPieces
	b.EmptySquares = ^b.OccupiedSquares
}

func (board *Board) MakeMove(move Move) error {
	sourceBit := BitBoard(1) << move.Source
	destBit := BitBoard(1) << move.Destination

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
		capturedPawnBit := BitBoard(1) << (move.Destination - 8) // or +8 depending on direction
		*board.pieceBitboard(BlackPawn) &= ^capturedPawnBit      // or WhitePawn
	case CastleKingside:
		// Update both king and rook positions, assuming kingside castle for example
		if move.Piece == WhiteKing {
			*board.pieceBitboard(WhiteRook) &= ^(BitBoard(1) << 7) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= BitBoard(1) << 5    // new rook position for kingside
		}
		if move.Piece == BlackKing {
			*board.pieceBitboard(BlackRook) &= ^(BitBoard(1) << 7) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= BitBoard(1) << 5    // new rook position for kingside
		}
	case CastleQueenside:
		if move.Piece == WhiteKing {
			*board.pieceBitboard(WhiteRook) &= ^(BitBoard(1) << 7) // original rook position for kingside
			*board.pieceBitboard(WhiteRook) |= BitBoard(1) << 5    // new rook position for kingside
		}
		if move.Piece == BlackKing {
			*board.pieceBitboard(BlackRook) &= ^(BitBoard(1) << 7) // original rook position for kingside
			*board.pieceBitboard(BlackRook) |= BitBoard(1) << 5    // new rook position for kingside
		}
	case Promotion:
		// Assume promoting to a queen for simplicity
		*board.pieceBitboard(move.Piece) &= ^destBit // Remove pawn from destination
		*board.pieceBitboard(WhiteQueen) |= destBit  // Add queen to destination
	}

	// Update occupied, empty, and aggregate piece bitboards
	board.updateAggregateBitboards()

	// Toggle turn
	board.TurnBlack = !board.TurnBlack

	return nil
}

func NewBoard() Board {
	whitePawns := WhitePawnBitboard(0x000000000000FF00)
	blackPawns := BlackPawnBitboard(0x00FF000000000000)

	whiteKnights := KnightBitboard(0x0000000000000042)
	blackKnights := KnightBitboard(0x4200000000000000)

	whiteBishops := BishopBitboard(0x0000000000000024)
	blackBishops := BishopBitboard(0x2400000000000000)

	whiteRooks := RookBitboard(0x0000000000000081)
	blackRooks := RookBitboard(0x8100000000000000)

	whiteQueens := QueenBitboard(0x0000000000000008)
	blackQueens := QueenBitboard(0x0800000000000000)

	whiteKing := KingBitboard(0x0000000000000010)
	blackKing := KingBitboard(0x1000000000000000)

	// Setting up the occupancy based on all pieces
	occupiedSquares := BitBoard(0xFFFF00000000FFFF)
	emptySquares := ^occupiedSquares // All squares not occupied

	whitePieces := BitBoard(0xFFFF000000000000)
	blackPieces := BitBoard(0x000000000000FFFF)

	enPassantTarget := BitBoard(0) // No en passant possible on first move
	castlingRights := uint8(0xF)   // All castling rights available initially (both kingside and queenside for both colors)

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
