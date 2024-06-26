package board

import (
	"fmt"
	"strings"

	"engine/evaluation/board/bitboards"
)

var PieceSymbols = map[int]string{
	0:  "♙",  // WhitePawn
	1:  "♟︎", // BlackPawn
	2:  "♘",  // WhiteKnight
	3:  "♞",  // BlackKnight
	4:  "♗",  // WhiteBishop
	5:  "♝",  // BlackBishop
	6:  "♖",  // WhiteRook
	7:  "♜",  // BlackRook
	8:  "♕",  // WhiteQueen
	9:  "♛",  // BlackQueen
	10: "♔",  // WhiteKing
	11: "♚",  // BlackKing
}

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

	CastleBlackKingside  bool
	CastleBlackQueenside bool
	CastleWhiteKingside  bool
	CastleWhiteQueenside bool

	BlackCastled bool
	WhiteCastled bool

	TurnBlack bool // Flag to indicate if it's black's turn to move

	HalfTurn int // What is the half turn
	Debug    bool
}

func New() Board {
	whitePawns := bitboards.WhitePawnBitboard(0x000000000000FF00)
	blackPawns := bitboards.BlackPawnBitboard(0x00FF000000000000)

	whiteKnights := bitboards.KnightBitboard(0x0000000000000042)
	blackKnights := bitboards.KnightBitboard(0x4200000000000000)

	whiteBishops := bitboards.BishopBitboard(0x0000000000000024)
	blackBishops := bitboards.BishopBitboard(0x2400000000000000)

	whiteRooks := bitboards.RookBitboard(0x0000000000000081)
	blackRooks := bitboards.RookBitboard(0x8100000000000000)

	whiteQueens := bitboards.QueenBitboard(0x0000000000000008)
	blackQueens := bitboards.QueenBitboard(0x0800000000000000)

	whiteKing := bitboards.KingBitboard(0x0000000000000010)
	blackKing := bitboards.KingBitboard(0x1000000000000000)

	occupiedSquares := bitboards.BitBoard(0xFFFF00000000FFFF)
	emptySquares := ^occupiedSquares

	whitePieces := bitboards.BitBoard(0x000000000000FFFF)
	blackPieces := bitboards.BitBoard(0xFFFF000000000000)

	enPassantTarget := bitboards.BitBoard(0)

	turnBlack := false

	return Board{
		WhitePawns:           whitePawns,
		BlackPawns:           blackPawns,
		WhiteKnights:         whiteKnights,
		BlackKnights:         blackKnights,
		WhiteBishops:         whiteBishops,
		BlackBishops:         blackBishops,
		WhiteRooks:           whiteRooks,
		BlackRooks:           blackRooks,
		WhiteQueens:          whiteQueens,
		BlackQueens:          blackQueens,
		WhiteKing:            whiteKing,
		BlackKing:            blackKing,
		OccupiedSquares:      occupiedSquares,
		EmptySquares:         emptySquares,
		WhitePieces:          whitePieces,
		BlackPieces:          blackPieces,
		EnPassantTarget:      enPassantTarget,
		CastleBlackKingside:  true,
		CastleBlackQueenside: true,
		CastleWhiteKingside:  true,
		CastleWhiteQueenside: true,
		TurnBlack:            turnBlack,
	}
}

func (board Board) Display() {
	fmt.Println("-------------------")
	for row := 0; row < 8; row++ {
		for col := 0; col < 9; col++ {
			if col == 8 {
				fmt.Print("|", 8-row)
				break
			}
			position := 8*(7-row) + col
			piece := board.PieceAt(position)
			if piece != -1 {
				fmt.Print(PieceSymbols[piece], " ")
			} else {
				fmt.Print("▢ ")
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------")
	fmt.Println("a b c d e f g h")
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

func (board Board) PieceAt(index int) int {
	indexMask := bitboards.New(index)

	if (board.WhitePawns.BitBoard() & indexMask) > 0 {
		return WhitePawn
	} else if (board.BlackPawns.BitBoard() & indexMask) > 0 {
		return BlackPawn
	} else if (board.WhiteKnights.BitBoard() & indexMask) > 0 {
		return WhiteKnight
	} else if (board.BlackKnights.BitBoard() & indexMask) > 0 {
		return BlackKnight
	} else if (board.WhiteBishops.BitBoard() & indexMask) > 0 {
		return WhiteBishop
	} else if (board.BlackBishops.BitBoard() & indexMask) > 0 {
		return BlackBishop
	} else if (board.WhiteRooks.BitBoard() & indexMask) > 0 {
		return WhiteRook
	} else if (board.BlackRooks.BitBoard() & indexMask) > 0 {
		return BlackRook
	} else if (board.WhiteQueens.BitBoard() & indexMask) > 0 {
		return WhiteQueen
	} else if (board.BlackQueens.BitBoard() & indexMask) > 0 {
		return BlackQueen
	} else if (board.WhiteKing.BitBoard() & indexMask) > 0 {
		return WhiteKing
	} else if (board.BlackKing.BitBoard() & indexMask) > 0 {
		return BlackKing
	}

	return -1
}

func (board Board) ToFEN() string {
	rankStrings := []string{}
	for rank := 7; rank >= 0; rank-- {
		emptyCount := 0
		rankString := ""

		for file := 0; file < 8; file++ {
			index := rank*8 + file
			piece := board.PieceAt(index)

			if piece == -1 {
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

	if board.TurnBlack {
		return strings.Join(rankStrings, "/") + " b - -"
	} else {
		return strings.Join(rankStrings, "/") + " w - -"
	}
}

func (b *Board) updateAggregateBitboards() {
	b.WhitePieces = b.WhitePawns.BitBoard() | b.WhiteKnights.BitBoard() | b.WhiteBishops.BitBoard() | b.WhiteRooks.BitBoard() | b.WhiteQueens.BitBoard() | b.WhiteKing.BitBoard()
	b.BlackPieces = b.BlackPawns.BitBoard() | b.BlackKnights.BitBoard() | b.BlackBishops.BitBoard() | b.BlackRooks.BitBoard() | b.BlackQueens.BitBoard() | b.BlackKing.BitBoard()
	b.OccupiedSquares = b.WhitePieces | b.BlackPieces
	b.EnPassantTarget = bitboards.BitBoard(0)
	b.EmptySquares = ^b.OccupiedSquares
}

func (b *Board) pieceBitboard(piece int) *bitboards.BitBoard {
	switch piece {
	case WhitePawn:
		return b.WhitePawns.BitBoardPointer()
	case BlackPawn:
		return b.BlackPawns.BitBoardPointer()
	case WhiteKnight:
		return b.WhiteKnights.BitBoardPointer()
	case BlackKnight:
		return b.BlackKnights.BitBoardPointer()
	case WhiteBishop:
		return b.WhiteBishops.BitBoardPointer()
	case BlackBishop:
		return b.BlackBishops.BitBoardPointer()
	case WhiteRook:
		return b.WhiteRooks.BitBoardPointer()
	case BlackRook:
		return b.BlackRooks.BitBoardPointer()
	case WhiteQueen:
		return b.WhiteQueens.BitBoardPointer()
	case BlackQueen:
		return b.BlackQueens.BitBoardPointer()
	case WhiteKing:
		return b.WhiteKing.BitBoardPointer()
	case BlackKing:
		return b.BlackKing.BitBoardPointer()
	default:
		panic("Invalid piece type")
	}
}

func (b Board) isOccupied(pos int) bool {
	occupiedMask := bitboards.New(pos)
	return (b.OccupiedSquares & occupiedMask) != 0
}
