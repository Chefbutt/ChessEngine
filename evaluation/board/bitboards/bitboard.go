package bitboards

import "fmt"

type BitBoard uint64

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

func (b *BitBoard) PopLSB() uint8 {
	if *b == 0 {
		return 65
	}

	lsb := *b & -*b

	*b &= *b - 1

	var index uint8
	for lsb != 1 {
		lsb >>= 1
		index++
	}

	return index
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
