package board

import (
	"math/bits"

	"engine/evaluation/board/bitboards"
)

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

	Check
)

type Move struct {
	Source         int
	Destination    int
	Piece          int // Use constants to define pieces
	CapturedPiece  int
	MoveType       int // Use constants to define move types
	IsCheck        bool
	PromotionPiece int // Added for promotion clarity
}

func (m Move) IsValid() bool {
	return m.Source != m.Destination && m.Piece == -1 && m.MoveType != -1
}

func (board Board) KingInPlayAndOpponentAttacks() (bitboards.BitBoard, bitboards.BitBoard) {
	if board.TurnBlack {
		return board.BlackKing.BitBoard(), board.WhiteAttacksMinimal()
	} else {
		return board.WhiteKing.BitBoard(), board.BlackAttacksMinimal()
	}
}

func (board Board) IsAttacked(square, attackedSquares bitboards.BitBoard) bool {
	if board.TurnBlack {
		return attackedSquares&square != 0
	} else {
		return attackedSquares&square != 0
	}
}

func (board Board) FromToToMove(from, to bitboards.BitBoard) Move {
	var enemySquares bitboards.BitBoard
	if board.TurnBlack {
		enemySquares = board.WhitePieces
	} else {
		enemySquares = board.BlackPieces
	}
	if enemySquares&to != 0 {
		return Move{Source: bits.TrailingZeros64(uint64(from)), Destination: bits.TrailingZeros64(uint64(to)), MoveType: Capture, Piece: board.PieceAt(bits.TrailingZeros64(uint64(from))), CapturedPiece: board.PieceAt(bits.TrailingZeros64(uint64(to)))}
	}
	return Move{Source: bits.TrailingZeros64(uint64(from)), Destination: bits.TrailingZeros64(uint64(to)), MoveType: NormalMove, Piece: board.PieceAt(bits.TrailingZeros64(uint64(from)))}
}

func StopCheck(board Board, allMoves map[bitboards.BitBoard][]bitboards.BitBoard) map[bitboards.BitBoard][]bitboards.BitBoard {
	viableBlocks := make(map[bitboards.BitBoard][]bitboards.BitBoard)
	if board.TurnBlack {
		for piecePos, moves := range allMoves {
			var pieceMoves []bitboards.BitBoard
			for _, move := range moves {
				tempBoard := board
				tempBoard.makeMove(tempBoard.FromToToMove(piecePos, move))
				attacked, _ := tempBoard.WhiteAttacks()
				if !tempBoard.IsAttacked(board.BlackKing.BitBoard(), attacked) {
					pieceMoves = append(pieceMoves, move)
				}
			}
			if len(pieceMoves) > 0 {
				viableBlocks[piecePos] = pieceMoves
			}
		}
	} else {
		for piecePos, moves := range allMoves {
			var pieceMoves []bitboards.BitBoard
			for _, move := range moves {
				tempBoard := board
				tempBoard.makeMove(tempBoard.FromToToMove(piecePos, move))
				attacked, _ := tempBoard.BlackAttacks()
				if !tempBoard.IsAttacked(board.WhiteKing.BitBoard(), attacked) {
					pieceMoves = append(pieceMoves, move)
				}
			}
			if len(pieceMoves) > 0 {
				viableBlocks[piecePos] = pieceMoves
			}
		}
	}
	return viableBlocks
}

func (board Board) BitBoardMapToMove(opposingColor bitboards.BitBoard, moves map[bitboards.BitBoard][]bitboards.BitBoard) []Move {
	var moveSlice []Move
	for from, toSlice := range moves {
		for _, to := range toSlice {
			moveSlice = append(moveSlice, Move{
				Source:      int(bits.TrailingZeros64(uint64(from))),
				Destination: int(bits.TrailingZeros64(uint64(to))),
			})
		}
	}
	return moveSlice
}

func (board Board) WhiteAttacks() (bitboards.BitBoard, int) {
	attacks := bitboards.NewFull()
	count := 0
	if bishop := board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.WhiteKing.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.WhitePawns.Attacks(); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	return attacks, count
}

func (board Board) WhiteAttacksMinimal() bitboards.BitBoard {
	attacks := bitboards.NewFull()
	if bishop := board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.WhiteKing.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.WhitePawns.Attacks(); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces); bishop != 0 {
		attacks = attacks | bishop
	}
	return attacks
}

func (board Board) BlackAttacks() (bitboards.BitBoard, int) {
	attacks := bitboards.New(0)
	count := 0
	if bishop := board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.BlackKing.Moves(board.EmptySquares, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.BlackPawns.Attacks(); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	if bishop := board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
		count++
	}
	return attacks, count
}

func (board Board) BlackAttacksMinimal() bitboards.BitBoard {
	attacks := bitboards.New(0)
	if bishop := board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.BlackKing.Moves(board.EmptySquares, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.BlackPawns.Attacks(); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
	}
	if bishop := board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces); bishop != 0 {
		attacks = attacks | bishop
	}
	return attacks
}

func (board Board) availableCastles() map[bitboards.BitBoard]bitboards.BitBoard {
	kingMoves := make(map[bitboards.BitBoard]bitboards.BitBoard)
	if board.TurnBlack {
		attacked, _ := board.WhiteAttacks()
		if board.CastleBlackKingside && !board.isOccupied(bits.TrailingZeros64(uint64(board.BlackKing>>1))) && !board.isOccupied(bits.TrailingZeros64(uint64(board.BlackKing>>2))) && !board.IsAttacked((board.BlackKing>>1).BitBoard(), attacked) && !board.IsAttacked((board.BlackKing>>2).BitBoard(), attacked) {
			kingMoves[board.BlackKing.BitBoard()] = (board.BlackKing >> 2).BitBoard()
		}
		if board.CastleBlackKingside && !board.isOccupied(bits.TrailingZeros64(uint64(board.BlackKing<<1))) && !board.isOccupied(bits.TrailingZeros64(uint64(board.BlackKing<<2))) && !board.IsAttacked((board.BlackKing<<1).BitBoard(), attacked) && !board.IsAttacked((board.BlackKing<<2).BitBoard(), attacked) {
			if _, ok := kingMoves[board.BlackKing.BitBoard()]; ok {
				kingMoves[board.BlackKing.BitBoard()] = kingMoves[board.BlackKing.BitBoard()] & (board.BlackKing << 2).BitBoard()
			} else {
				kingMoves[board.BlackKing.BitBoard()] = (board.BlackKing << 2).BitBoard()
			}
		}
	} else {
		attacked, _ := board.BlackAttacks()
		if board.CastleWhiteKingside && !board.isOccupied(bits.TrailingZeros64(uint64(board.WhiteKing>>1))) && !board.isOccupied(bits.TrailingZeros64(uint64(board.WhiteKing>>2))) && !board.IsAttacked((board.WhiteKing>>1).BitBoard(), attacked) && !board.IsAttacked((board.WhiteKing>>2).BitBoard(), attacked) {
			kingMoves[board.WhiteKing.BitBoard()] = (board.WhiteKing >> 2).BitBoard()
		}
		if board.CastleWhiteQueenside && !board.isOccupied(bits.TrailingZeros64(uint64(board.WhiteKing<<1))) && !board.isOccupied(bits.TrailingZeros64(uint64(board.WhiteKing<<2))) && !board.IsAttacked((board.WhiteKing<<1).BitBoard(), attacked) && !board.IsAttacked((board.WhiteKing<<2).BitBoard(), attacked) {
			if _, ok := kingMoves[board.WhiteKing.BitBoard()]; ok {
				kingMoves[board.WhiteKing.BitBoard()] = kingMoves[board.WhiteKing.BitBoard()] & (board.WhiteKing << 2).BitBoard()
			} else {
				kingMoves[board.WhiteKing.BitBoard()] = (board.WhiteKing << 2).BitBoard()
			}
		}
	}
	return kingMoves
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
		for src, kingMoves := range board.availableCastles() {
			moveBitboards := kingMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = append(moveMap[src], moveBitboards...)
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
		for src, kingMoves := range board.availableCastles() {
			moveBitboards := kingMoves.Split()
			if len(moveBitboards) != 0 {
				moveMap[src] = append(moveMap[src], moveBitboards...)
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
