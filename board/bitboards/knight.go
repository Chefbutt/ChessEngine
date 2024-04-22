package bitboards

type KnightBitboard BitBoard

func (b *KnightBitboard) BitBoardPointer() *BitBoard {
	return (*BitBoard)(b)
}

func (b KnightBitboard) BitBoard() BitBoard {
	return BitBoard(b)
}

func (k KnightBitboard) Attacks(otherColorPieces, empty BitBoard) BitBoard {
	return k.Moves(empty) & otherColorPieces
}

func (k KnightBitboard) Moves(empty BitBoard) BitBoard {
	pos := BitBoard(k)
	northNorthEast := (pos << 17) & notAFile // two up, one right
	northEastEast := (pos << 10) & notABFile // one up, two right
	southEastEast := (pos >> 6) & notABFile  // one down, two right
	southSouthEast := (pos >> 15) & notAFile // two down, one right
	northNorthWest := (pos << 15) & notHFile // two up, one left
	northWestWest := (pos << 6) & notGHFile  // one up, two left
	southWestWest := (pos >> 10) & notGHFile // one down, two left
	southSouthWest := (pos >> 17) & notHFile // two down, one left

	return (northNorthEast | northEastEast | southEastEast | southSouthEast | northNorthWest | northWestWest | southWestWest | southSouthWest) & empty
}

// Board edges that knights can't wrap around
const (
	notAFile  BitBoard = 0xfefefefefefefefe // 11111110...
	notABFile BitBoard = 0xfcfcfcfcfcfcfcfc // 11111100...
	notHFile  BitBoard = 0x7f7f7f7f7f7f7f7f // 01111111...
	notGHFile BitBoard = 0x3f3f3f3f3f3f3f3f // 00111111...
)
