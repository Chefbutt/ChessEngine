package bitboards

// rank4 and rank5 constants represent the 4th and 5th rows on a chess board.
const (
	rank4 BitBoard = 0x00000000FF000000
	rank5 BitBoard = 0x000000FF00000000
)

type WhitePawnBitboard BitBoard

type BlackPawnBitboard BitBoard

func (b BlackPawnBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b WhitePawnBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (b *BlackPawnBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (b *WhitePawnBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (p WhitePawnBitboard) Attacks() BitBoard {
	return p.whitePawnEastAttacks() | p.whitePawnWestAttacks()
}

func (p WhitePawnBitboard) Moves(empty, oppositeColorOccupancy, enPassantTarget BitBoard) BitBoard {
	return p.SinglePushTargets(empty) | p.DoublePushTargets(empty) | ((p.whitePawnEastAttacks() | p.whitePawnWestAttacks()) & (oppositeColorOccupancy | enPassantTarget))
}

func (b WhitePawnBitboard) MovesByPiece(empty, oppositeColorOccupancy, enPassantTarget BitBoard) map[BitBoard]BitBoard {
	pawns := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, pawn := range pawns {
		moves[pawn] = WhitePawnBitboard(pawn).Moves(empty, oppositeColorOccupancy, enPassantTarget)
	}

	return moves
}

func (b BlackPawnBitboard) MovesByPiece(empty, oppositeColorOccupancy, enPassantTarget BitBoard) map[BitBoard]BitBoard {
	pawns := b.BitBoard().Split()
	moves := make(map[BitBoard]BitBoard)

	for _, pawn := range pawns {
		moves[pawn] = BlackPawnBitboard(pawn).Moves(empty, oppositeColorOccupancy, enPassantTarget)
	}

	return moves
}

func (p BlackPawnBitboard) Attacks() BitBoard {
	return p.eastAttacks() | p.westAttacks()
}

// whiteSinglePushTargets calculates the targets for white pawns that can move forward one square.
func (p WhitePawnBitboard) SinglePushTargets(empty BitBoard) BitBoard {
	return BitBoard(p).northOne() & empty
}

// whiteDoublePushTargets calculates the targets for white pawns that can move forward two squares.
func (p WhitePawnBitboard) DoublePushTargets(empty BitBoard) BitBoard {
	singlePushes := p.SinglePushTargets(empty)
	return singlePushes.northOne() & empty & rank4
}

// whitePawnEastAttacks computes the eastward attacks for white pawns.
func (p WhitePawnBitboard) whitePawnEastAttacks() BitBoard {
	return BitBoard(p).northEastOne()
}

// wPawnWestAttacks computes the westward attacks for white pawns.
func (p WhitePawnBitboard) whitePawnWestAttacks() BitBoard {
	return BitBoard(p).northWestOne()
}

func (p BlackPawnBitboard) Moves(empty, oppositeColorOccupancy, enPassantTarget BitBoard) BitBoard {
	return p.SinglePushTargets(empty) | p.DoublePushTargets(empty) | ((p.eastAttacks() | p.westAttacks()) & (oppositeColorOccupancy | enPassantTarget))
}

// whitePawnEastAttacks computes the eastward attacks for white pawns.
func (p BlackPawnBitboard) eastAttacks() BitBoard {
	return BitBoard(p).southEastOne()
}

// wPawnWestAttacks computes the westward attacks for white pawns.
func (p BlackPawnBitboard) westAttacks() BitBoard {
	return BitBoard(p).southWestOne()
}

// blackSinglePushTargets calculates the targets for black pawns that can move forward one square.
func (p BlackPawnBitboard) SinglePushTargets(empty BitBoard) BitBoard {
	return BitBoard(p).southOne() & empty
}

// blackDoublePushTargets calculates the targets for black pawns that can move forward two squares.
func (p BlackPawnBitboard) DoublePushTargets(empty BitBoard) BitBoard {
	singlePushes := p.SinglePushTargets(empty)
	return singlePushes.southOne() & empty & rank5
}

// blackPawnEastAttacks computes the eastward attacks for black pawns.
func (p BlackPawnBitboard) PawnEastAttacks() BitBoard {
	return BitBoard(p).southEastOne()
}

// blackPawnWestAttacks computes the westward attacks for black pawns.
func (p BlackPawnBitboard) PawnWestAttacks() BitBoard {
	return BitBoard(p).southWestOne()
}
