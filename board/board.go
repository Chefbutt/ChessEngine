package board

import (
	"fmt"
	"strings"
)

// Board holds the state of the game, including piece positions, turn, and move count.
type Board struct {
	WhitePawns   WhitePawnBitboard // Bitboard for white pawns
	BlackPawns   BlackPawnBitboard // Bitboard for black pawns
	WhiteKnights KnightBitboard    // Bitboard for white knights
	BlackKnights KnightBitboard    // Bitboard for black knights
	WhiteBishops BishopBitboard    // Bitboard for white bishops
	BlackBishops BishopBitboard    // Bitboard for black bishops
	WhiteRooks   RookBitboard      // Bitboard for white rooks
	BlackRooks   RookBitboard      // Bitboard for black rooks
	WhiteQueens  QueenBitboard     // Bitboard for white queens
	BlackQueens  QueenBitboard     // Bitboard for black queens
	WhiteKing    KingBitboard      // Bitboard for the white king
	BlackKing    KingBitboard      // Bitboard for the black king

	OccupiedSquares BitBoard // Bitboard for all occupied squares
	EmptySquares    BitBoard // Bitboard for all empty squares

	WhitePieces BitBoard // Bitboard for all white pieces
	BlackPieces BitBoard // Bitboard for all black pieces

	EnPassantTarget BitBoard // Bitboard for possible en passant capture squares

	CastlingRights uint8 // Flags for castling rights, encoded as bits

	TurnBlack bool // Flag to indicate if it's black's turn to move

	Turn int // What is the turn
}

// BitBoard represents a position on a chess board using a 64-bit integer.
type BitBoard uint64

func pieceToFEN(piece int) string {
	switch piece {
	case WhitePawn:
		return "P"
	case BlackPawn:
		return "p"
	case WhiteKnight:
		return "N"
	case BlackKnight:
		return "n"
	case WhiteBishop:
		return "B"
	case BlackBishop:
		return "b"
	case WhiteRook:
		return "R"
	case BlackRook:
		return "r"
	case WhiteQueen:
		return "Q"
	case BlackQueen:
		return "q"
	case WhiteKing:
		return "K"
	case BlackKing:
		return "k"
	default:
		return "?"
	}
}

func (board *Board) ToFEN() string {
	rankStrings := []string{}
	for rank := 7; rank >= 0; rank-- { // From rank 8 to 1
		emptyCount := 0
		rankString := ""

		for file := 0; file < 8; file++ { // From file a to h
			index := rank*8 + file
			piece := identifyPieceAt(board, index)

			if piece == -1 { // No piece at this position
				emptyCount++
			} else {
				if emptyCount > 0 {
					rankString += fmt.Sprintf("%d", emptyCount)
					emptyCount = 0
				}
				rankString += pieceToFEN(piece)
			}
		}
		if emptyCount > 0 {
			rankString += fmt.Sprintf("%d", emptyCount)
		}
		rankStrings = append(rankStrings, rankString)
	}

	// Assemble the full FEN for the piece placement
	if board.TurnBlack {
		return strings.Join(rankStrings, "/") + " b - -"
	} else {
		return strings.Join(rankStrings, "/") + " w - -"
	}
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
