package board

import (
	"fmt"
	"strings"

	"engine/board/bitboards"
)

// Board holds the state of the game, including piece positions, turn, and move count.
type Board struct {
	WhitePawns   bitboards.WhitePawnBitboard // Bitboard for white pawns
	BlackPawns   bitboards.BlackPawnBitboard // Bitboard for black pawns
	WhiteKnights bitboards.KnightBitboard    // Bitboard for white knights
	BlackKnights bitboards.KnightBitboard    // Bitboard for black knights
	WhiteBishops bitboards.BishopBitboard    // Bitboard for white bishops
	BlackBishops bitboards.BishopBitboard    // Bitboard for black bishops
	WhiteRooks   bitboards.RookBitboard      // Bitboard for white rooks
	BlackRooks   bitboards.RookBitboard      // Bitboard for black rooks
	WhiteQueens  bitboards.QueenBitboard     // Bitboard for white queens
	BlackQueens  bitboards.QueenBitboard     // Bitboard for black queens
	WhiteKing    bitboards.KingBitboard      // Bitboard for the white king
	BlackKing    bitboards.KingBitboard      // Bitboard for the black king

	OccupiedSquares bitboards.BitBoard // Bitboard for all occupied squares
	EmptySquares    bitboards.BitBoard // Bitboard for all empty squares

	WhitePieces bitboards.BitBoard // Bitboard for all white pieces
	BlackPieces bitboards.BitBoard // Bitboard for all black pieces

	EnPassantTarget bitboards.BitBoard // Bitboard for possible en passant capture squares

	CastlingRights uint8 // Flags for castling rights, encoded as bits

	TurnBlack bool // Flag to indicate if it's black's turn to move

	Turn int // What is the turn
}

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
