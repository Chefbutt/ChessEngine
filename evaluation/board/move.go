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

func (originalBoard Board) AvailableBlackAttacks() bitboards.BitBoard {
	board := originalBoard

	var attack bitboards.BitBoard

	movesList := board.BlackKnights.Moves(board.EmptySquares, board.WhitePieces)
	attack = attack | movesList

	movesList = board.BlackBishops.Moves(board.BlackPieces, board.WhitePieces)
	attack = attack | movesList

	movesList = board.BlackRooks.Moves(board.BlackPieces, board.WhitePieces)
	attack = attack | movesList

	movesList = board.BlackQueens.Moves(board.BlackPieces, board.WhitePieces)
	attack = attack | movesList

	movesList = board.BlackPawns.Moves(board.EmptySquares, board.WhitePieces, board.EnPassantTarget)
	attack = attack | movesList

	return attack
}

func (originalBoard Board) AvailableWhiteAttacks() bitboards.BitBoard {
	board := originalBoard

	var attack bitboards.BitBoard

	movesList := board.WhiteKnights.Moves(board.EmptySquares, board.BlackPieces)
	attack = attack | movesList

	movesList = board.WhiteBishops.Moves(board.WhitePieces, board.BlackPieces)
	attack = attack | movesList

	movesList = board.WhiteRooks.Moves(board.WhitePieces, board.BlackPieces)
	attack = attack | movesList

	movesList = board.WhiteQueens.Moves(board.WhitePieces, board.BlackPieces)
	attack = attack | movesList

	movesList = board.WhitePawns.Moves(board.EmptySquares, board.BlackPieces, board.EnPassantTarget)
	attack = attack | movesList

	return attack
}

func (originalBoard Board) AvailableBlackMoves() []Move {
	board := originalBoard
	var moves []Move

	attacks := board.WhiteAttacksMinimal()
	if board.CastleBlackQueenside && !board.isOccupied(58) && !board.isOccupied(59) && !board.isOccupied(57) && board.PieceAt(56) == BlackRook && !board.IsAttacked((board.BlackKing>>1).BitBoard(), attacks) && !board.IsAttacked((board.BlackKing>>2).BitBoard(), attacks) {
		moves = append(moves, Move{Source: 60, Destination: 58, MoveType: CastleQueenside, Piece: BlackKing})
	}
	if board.CastleBlackKingside && !board.isOccupied(61) && !board.isOccupied(62) && !board.IsAttacked((board.BlackKing<<1).BitBoard(), attacks) && board.PieceAt(63) == BlackRook && !board.IsAttacked((board.BlackKing<<2).BitBoard(), attacks) {
		moves = append(moves, Move{Source: 60, Destination: 62, MoveType: CastleKingside, Piece: BlackKing})
	}

	knights := board.BlackKnights
	for knights != 0 {
		from := knights.BitBoardPointer().PopLSB()

		movesList := bitboards.KnightBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.WhitePieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: BlackKnight})
		}
	}

	bishops := board.BlackBishops
	for bishops != 0 {
		from := bishops.BitBoardPointer().PopLSB()

		movesList := bitboards.BishopBitboard(bitboards.New(int(from))).Moves(board.BlackPieces, board.WhitePieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: BlackBishop})
		}
	}

	rooks := board.BlackRooks
	for rooks != 0 {
		from := rooks.BitBoardPointer().PopLSB()

		movesList := bitboards.RookBitboard(bitboards.New(int(from))).Moves(board.BlackPieces, board.WhitePieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: BlackRook})
		}
	}

	queens := board.BlackQueens
	for queens != 0 {
		from := queens.BitBoardPointer().PopLSB()

		movesList := bitboards.QueenBitboard(bitboards.New(int(from))).Moves(board.BlackPieces, board.WhitePieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: BlackQueen})
		}
	}

	pawns := board.BlackPawns
	for pawns != 0 {
		from := pawns.BitBoardPointer().PopLSB()

		movesList := bitboards.BlackPawnBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.WhitePieces, board.EnPassantTarget)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: BlackPawn})
		}
	}

	kings := board.BlackKing
	for kings != 0 {
		from := kings.BitBoardPointer().PopLSB()

		movesList := bitboards.KingBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.WhitePieces)
		for movesList != 0 {
			board := board
			to := movesList.PopLSB()

			move := Move{Source: bits.TrailingZeros64(uint64(board.BlackKing)), Destination: int(to), Piece: BlackKing}
			board.makeMove(move)
			if !board.isKingInCheck(board.BlackKing, false) && move.Source != move.Destination {
				moves = append(moves, move)
			}
		}
	}
	// Remove attacked

	if board.isKingInCheck(board.BlackKing, false) {
		var legalMoves []Move
		for _, move := range moves {
			board := board
			_, err := board.makeMove(move)
			if err != nil {
				panic(err)
			}
			if !board.isKingInCheck(board.BlackKing, false) {
				legalMoves = append(legalMoves, move)
			}
		}
		return legalMoves
	}

	return moves
}

func (originalBoard Board) AvailableWhiteMoves() []Move {
	board := originalBoard
	var moves []Move

	attacks := board.BlackAttacksMinimal()
	if board.CastleWhiteQueenside && board.isOccupied(1) && !board.isOccupied(2) && !board.isOccupied(3) && board.PieceAt(4) == WhiteRook && !board.IsAttacked((board.WhiteKing>>1).BitBoard(), attacks) && !board.IsAttacked((board.WhiteKing>>2).BitBoard(), attacks) {
		moves = append(moves, Move{Source: 4, Destination: 1, MoveType: CastleQueenside, Piece: WhiteKing})
	}
	if board.CastleWhiteKingside && !board.isOccupied(5) && !board.isOccupied(6) && !board.IsAttacked((board.WhiteKing<<1).BitBoard(), attacks) && !board.IsAttacked((board.WhiteKing<<2).BitBoard(), attacks) {
		moves = append(moves, Move{Source: 4, Destination: 6, MoveType: CastleKingside, Piece: WhiteKing})
	}

	knights := board.WhiteKnights
	for knights != 0 {
		from := knights.BitBoardPointer().PopLSB()

		movesList := bitboards.KnightBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.BlackPieces)
		for movesList != 0 {
			to := movesList.PopLSB()
			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: WhiteKnight})
		}
	}

	bishops := board.WhiteBishops
	for bishops != 0 {
		from := bishops.BitBoardPointer().PopLSB()

		movesList := bitboards.BishopBitboard(bitboards.New(int(from))).Moves(board.WhitePieces, board.BlackPieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: WhiteBishop})
		}
	}

	rooks := board.WhiteRooks
	for rooks != 0 {
		from := rooks.BitBoardPointer().PopLSB()

		movesList := bitboards.RookBitboard(bitboards.New(int(from))).Moves(board.WhitePieces, board.BlackPieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: WhiteRook})
		}
	}

	queens := board.WhiteQueens
	for queens != 0 {
		from := queens.BitBoardPointer().PopLSB()

		movesList := bitboards.QueenBitboard(bitboards.New(int(from))).Moves(board.WhitePieces, board.BlackPieces)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: WhiteQueen})
		}
	}

	pawns := board.WhitePawns
	for pawns != 0 {
		from := pawns.BitBoardPointer().PopLSB()

		movesList := bitboards.WhitePawnBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.BlackPieces, board.EnPassantTarget)
		for movesList != 0 {
			to := movesList.PopLSB()

			moves = append(moves, Move{Source: int(from), Destination: int(to), Piece: WhitePawn})
		}
	}

	// Remove attacked
	kings := board.WhiteKing
	for kings != 0 {
		from := kings.BitBoardPointer().PopLSB()

		movesList := bitboards.KingBitboard(bitboards.New(int(from))).Moves(board.EmptySquares, board.BlackPieces)
		for movesList != 0 {
			board := board
			to := movesList.PopLSB()

			move := Move{Source: bits.TrailingZeros64(uint64(board.WhiteKing)), Destination: int(to), Piece: WhiteKing}
			board.makeMove(move)
			if !board.isKingInCheck(board.WhiteKing, true) && move.Source != move.Destination {
				moves = append(moves, move)
			}
		}
	}

	if board.isKingInCheck(board.WhiteKing, true) {
		var legalMoves []Move
		for _, move := range moves {
			board := board
			_, err := board.makeMove(move)
			if err != nil {
				panic(err)
			}
			if !board.isKingInCheck(board.WhiteKing, true) {
				legalMoves = append(legalMoves, move)
			}
		}
		return legalMoves
	}

	return moves
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

func (board Board) UCItoMove(uci string) Move {
	if len(uci) < 4 || len(uci) > 5 {
		return Move{}
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
