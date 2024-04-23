package board

import "engine/evaluation/board/bitboards"

const (
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

type Move struct {
	Source         int
	Destination    int
	Piece          int // Use constants to define pieces
	CapturedPiece  int
	MoveType       int // Use constants to define move types
	PromotionPiece int // Added for promotion clarity
}

func (m Move) IsValid() bool {
	return m.Source != m.Destination && m.Piece != -1 && m.MoveType != -1
}

func (board Board) InCheck() bool {
	if board.TurnBlack {
		return (board.WhiteAttacks() | bitboards.BitBoard(board.BlackKing)) != 0
	} else {
		return (board.BlackAttacks() | bitboards.BitBoard(board.BlackKing)) != 0
	}
}

func (board Board) WhiteAttacks() bitboards.BitBoard {
	return board.WhiteBishops.Attacks(board.WhitePieces, board.BlackPieces) | board.WhiteKing.Attacks(board.BlackPieces) | board.WhiteKnights.Attacks(board.BlackPieces) | board.WhitePawns.Attacks(board.BlackPieces, board.EnPassantTarget) | board.WhiteQueens.Attacks(board.WhitePieces, board.BlackPieces) | board.WhiteRooks.Attacks(board.WhitePieces, board.BlackPieces)
}

func (board Board) BlackAttacks() bitboards.BitBoard {
	return board.BlackBishops.Attacks(board.BlackPieces, board.WhitePieces) | board.BlackKing.Attacks(board.WhitePieces) | board.BlackKnights.Attacks(board.WhitePieces) | board.BlackPawns.Attacks(board.WhitePieces, board.EnPassantTarget) | board.BlackQueens.Attacks(board.BlackPieces, board.WhitePieces) | board.BlackRooks.Attacks(board.BlackPieces, board.WhitePieces)
}

func (board Board) RegularMoves() map[bitboards.BitBoard][]bitboards.BitBoard {
	moveMap := make(map[bitboards.BitBoard][]bitboards.BitBoard)
	if board.TurnBlack {
		for src, pawnMoves := range board.BlackPawns.MovesByPiece(board.EmptySquares, board.WhitePieces, board.EnPassantTarget) {
			moveBitboards := pawnMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, bishopMoves := range board.BlackBishops.MovesByPiece(board.BlackPieces, board.WhitePieces) {
			moveBitboards := bishopMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, knightMoves := range board.BlackKnights.MovesByPiece(board.EmptySquares, board.WhitePieces) {
			moveBitboards := knightMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, rookMoves := range board.BlackRooks.MovesByPiece(board.BlackPieces, board.WhitePieces) {
			moveBitboards := rookMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, queenMoves := range board.BlackQueens.MovesByPiece(board.BlackPieces, board.WhitePieces) {
			moveBitboards := queenMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, kingMoves := range board.BlackKing.MovesByPiece(board.EmptySquares, board.WhitePieces) {
			moveBitboards := kingMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
	} else {
		for src, pawnMoves := range board.WhitePawns.MovesByPiece(board.EmptySquares, board.BlackPieces, board.EnPassantTarget) {
			moveBitboards := pawnMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, bishopMoves := range board.WhiteBishops.MovesByPiece(board.WhitePieces, board.BlackPieces) {
			moveBitboards := bishopMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, knightMoves := range board.WhiteKnights.MovesByPiece(board.EmptySquares, board.BlackPieces) {
			moveBitboards := knightMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, rookMoves := range board.WhiteRooks.MovesByPiece(board.WhitePieces, board.BlackPieces) {
			moveBitboards := rookMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, queenMoves := range board.WhiteQueens.MovesByPiece(board.WhitePieces, board.BlackPieces) {
			moveBitboards := queenMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
		for src, kingMoves := range board.WhiteKing.MovesByPiece(board.EmptySquares, board.BlackPieces) {
			moveBitboards := kingMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = moveBitboards
			}
		}
	}
	return moveMap
}

func (board *Board) UCItoMove(uci string) Move {
	if len(uci) < 4 || len(uci) > 5 {
		panic("Invalid UCI string")
	}

	source := positionToIndex(uci[0:2])
	destination := positionToIndex(uci[2:4])

	piece := board.PieceAt(source)

	// Determine if it's a capture
	isCapture := board.isOccupied(destination)
	capturedPiece := -1
	moveType := NormalMove
	if isCapture {
		moveType = Capture
		capturedPiece = board.PieceAt(destination)
	}

	if (piece == WhiteKing || piece == BlackKing) && (source+2 == destination) {
		moveType = CastleKingside
	}

	if (piece == WhiteKing || piece == BlackKing) && (source-3 == destination) {
		moveType = CastleQueenside
	}

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

func positionToIndex(pos string) int {
	file := pos[0] - 'a'
	rank := pos[1] - '1'
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
