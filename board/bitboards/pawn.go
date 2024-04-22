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

func (p WhitePawnBitboard) Moves(empty BitBoard) BitBoard {
	return p.SinglePushTargets(empty) & p.DoublePushTargets(empty)
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

func (p BlackPawnBitboard) Moves(empty BitBoard) BitBoard {
	return p.SinglePushTargets(empty) & p.DoublePushTargets(empty)
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
