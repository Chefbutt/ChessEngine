package bitboards

import "fmt"

// BitBoard represents a position on a chess board using a 64-bit integer.
type BitBoard uint64

// eastOne shifts the square one position to the right (East) on the board.
func (s BitBoard) eastOne() BitBoard {
	return (s >> 1) & 0x7f7f7f7f7f7f7f7f
}

// westOne shifts the square one position to the left (West) on the board.
func (s BitBoard) westOne() BitBoard {
	return (s << 1) & 0xfefefefefefefefe
}

// northOne shifts the square eight positions up (North) on the board.
func (s BitBoard) northOne() BitBoard {
	return s << 8
}

// southOne shifts the square eight positions down (South) on the board.
func (s BitBoard) southOne() BitBoard {
	return s >> 8
}

// northEastOne shifts north-east for capturing moves.
func (b BitBoard) northEastOne() BitBoard {
	return (b << 9) & ^BitBoard(0x0101010101010101) // Exclude a-file wraparounds
}

// northWestOne shifts north-west for capturing moves.
func (b BitBoard) northWestOne() BitBoard {
	return (b << 7) & ^BitBoard(0x8080808080808080) // Exclude h-file wraparounds
}

// southEastOne shifts south-east for capturing moves.
func (b BitBoard) southEastOne() BitBoard {
	return (b >> 7) & ^BitBoard(0x0101010101010101) // Exclude a-file wraparounds
}

// southWestOne shifts south-west for capturing moves.
func (b BitBoard) southWestOne() BitBoard {
	return (b >> 9) & ^BitBoard(0x8080808080808080) // Exclude h-file wraparounds
}

func (b BitBoard) Display() {
	// Iterate over each row
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			// Calculate the position of the bit to check
			position := 8*(7-row) + col // bit position from the top left
			if (b & (1 << position)) != 0 {
				fmt.Print("1 ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println() // New line after each row
	}
}

func (b BitBoard) PopCount() int {
	count := 0
	for b != 0 {
		count++
		b &= b - 1 // reset least significant bit
	}
	return count
}

func New(mask int) BitBoard {
	return BitBoard(1) << mask
}

func FileMask(file int) BitBoard {
	return 0x0101010101010101 << file
}
