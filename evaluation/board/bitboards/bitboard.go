package bitboards

import (
	"fmt"
	"math/bits"
)

type BitBoard uint64

var (
	SquareBB        [65]BitBoard
	BlackPawnPushes map[BitBoard]BitBoard
	WhitePawnPushes map[BitBoard]BitBoard
	KnightMoves     map[BitBoard]BitBoard
	KingMoves       map[BitBoard]BitBoard
)

func InitBitboards() {
	var sq uint8
	for sq = 0; sq < 65; sq++ {
		SquareBB[sq] = 0x8000000000000000 << sq
	}

	WhitePawnPushes = make(map[BitBoard]BitBoard)
	BlackPawnPushes = make(map[BitBoard]BitBoard)
	KnightMoves = make(map[BitBoard]BitBoard)
	KingMoves = make(map[BitBoard]BitBoard)

	for sq = 0; sq < 64; sq++ {
		squareBB := New(int(sq))

		// White Pawn Pushes
		WhitePawnPushes[squareBB] = squareBB << 8

		BlackPawnPushes[squareBB] = squareBB >> 8

		// Knight Moves
		KnightMoves[squareBB] = (squareBB<<17)&notAFile | (squareBB<<10)&notABFile | (squareBB>>6)&notABFile | (squareBB>>15)&notAFile | (squareBB<<15)&notHFile | (squareBB<<6)&notGHFile | (squareBB>>10)&notGHFile | (squareBB>>17)&notHFile

		// King Moves
		KingMoves[squareBB] = squareBB.eastOne() | squareBB.westOne() | squareBB.northOne() | squareBB.southOne() | squareBB.southEastOne() | squareBB.southWestOne() | squareBB.northEastOne() | squareBB.northWestOne()
	}
}

func (s BitBoard) eastOne() BitBoard {
	return (s >> 1) & 0x7f7f7f7f7f7f7f7f
}

func (s BitBoard) westOne() BitBoard {
	return (s << 1) & 0xfefefefefefefefe
}

func (s BitBoard) northOne() BitBoard {
	return s << 8
}

func (s BitBoard) southOne() BitBoard {
	return s >> 8
}

func (b BitBoard) northEastOne() BitBoard {
	return (b << 9) & ^BitBoard(0x0101010101010101)
}

func (b BitBoard) northWestOne() BitBoard {
	return (b << 7) & ^BitBoard(0x8080808080808080)
}

func (b BitBoard) southEastOne() BitBoard {
	return (b >> 7) & ^BitBoard(0x0101010101010101)
}

func (b BitBoard) southWestOne() BitBoard {
	return (b >> 9) & ^BitBoard(0x8080808080808080)
}

func (b BitBoard) Display() {
	fmt.Println()
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			position := 8*(7-row) + col
			if (b & (1 << position)) != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

func (b BitBoard) PopCount() int {
	count := 0
	for b != 0 {
		count++
		b &= b - 1
	}
	return count
}

func (bitboard *BitBoard) SetBit(sq uint8) {
	*bitboard |= SquareBB[sq]
}

// Clear the bit at given square.
func (bitboard *BitBoard) ClearBit(sq uint8) {
	*bitboard &= ^SquareBB[sq]
}

// Test whether the bit of the given bitbord at the given
// position is set.
func (bb BitBoard) BitSet(sq uint8) bool {
	return (bb & SquareBB[sq]) != 0
}

// Get the position of the MSB of the given bitboard.
func (bitboard BitBoard) Lsb() uint8 {
	return uint8(bits.TrailingZeros64(uint64(bitboard)))
}

// Get the position of the MSB of the given bitboard,
// and clear the MSB.
func (bitboard *BitBoard) PopBit() uint8 {
	sq := bitboard.Lsb()
	bitboard.ClearBit(sq)
	return sq
}

// Count the bits in a given bitboard using the SWAR-popcount
// algorithm for 64-bit integers.
func (bitboard BitBoard) CountBits() int {
	return bits.OnesCount64(uint64(bitboard))
}

func (b *BitBoard) PopLSB() uint8 {
	lsb := b.Lsb()

	*b &= *b - 1

	return lsb
}

func (b BitBoard) Split() []BitBoard {
	var bitboards []BitBoard
	for i := 0; i < 64; i++ {
		if b&(1<<i) != 0 {
			bitboards = append(bitboards, 1<<i)
		}
	}
	return bitboards
}

func NewFull() BitBoard {
	return BitBoard(0)
}

func New(mask int) BitBoard {
	return BitBoard(1) << mask
}

func FileMask(file int) BitBoard {
	return 0x0101010101010101 << file
}
