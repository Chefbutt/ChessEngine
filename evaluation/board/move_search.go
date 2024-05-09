package board

import (
	"fmt"
	"math/rand"
	"sync"

	"engine/evaluation/board/bitboards"
)

func (board Board) AggregateBitboards() AggregateBitboards {
	return AggregateBitboards{
		WhitePieces:     board.WhitePieces,
		BlackPieces:     board.BlackPieces,
		OccupiedSquares: board.OccupiedSquares,
		EnPassantTarget: board.EnPassantTarget,
		EmptySquares:    board.EmptySquares,
	}
}

type MoveUndo struct {
	Source                       int
	Destination                  int
	Piece                        int
	CapturedPiece                int
	PromotionPiece               int
	MoveType                     int
	PreviousCastleWhiteKingside  bool
	PreviousCastleWhiteQueenside bool
	PreviousCastleBlackKingside  bool
	PreviousCastleBlackQueenside bool
	PreviousTurnBlack            bool
	PreviousHalfTurn             int
	PreviousAggregateBitboards   AggregateBitboards // Assume this type holds all necessary board state
}

func (board *Board) UndoMove(undo *MoveUndo) {
	*board.pieceBitboard(undo.Piece) &= ^bitboards.New(undo.Destination)
	*board.pieceBitboard(undo.Piece) |= bitboards.New(undo.Source)

	if undo.CapturedPiece != -1 {
		*board.pieceBitboard(undo.CapturedPiece) |= bitboards.New(undo.Destination)
	}

	if undo.MoveType == Promotion {
		*board.pieceBitboard(undo.PromotionPiece) &= ^bitboards.New(undo.Destination)
		*board.pieceBitboard(undo.Piece) |= bitboards.New(undo.Destination)
	}

	// Restore castling rights
	board.CastleWhiteKingside = undo.PreviousCastleWhiteKingside
	board.CastleWhiteQueenside = undo.PreviousCastleWhiteQueenside
	board.CastleBlackKingside = undo.PreviousCastleBlackKingside
	board.CastleBlackQueenside = undo.PreviousCastleBlackQueenside

	// Restore turn and half turn counters
	board.TurnBlack = undo.PreviousTurnBlack
	board.HalfTurn = undo.PreviousHalfTurn

	// Restore aggregate bitboards if needed
	board.BlackPieces = undo.PreviousAggregateBitboards.BlackPieces
	board.EmptySquares = undo.PreviousAggregateBitboards.EmptySquares
	board.EnPassantTarget = undo.PreviousAggregateBitboards.EnPassantTarget
	board.OccupiedSquares = undo.PreviousAggregateBitboards.OccupiedSquares
	board.WhitePieces = undo.PreviousAggregateBitboards.WhitePieces

	// Ensure the board's internal state is consistent
	board.updateAggregateBitboards()
}

type MoveEvaluation struct {
	Move  Move
	Score Evaluation
}

type TranspositionEntry struct {
	Depth    int
	Score    Evaluation
	Flag     int
	BestMove Move
}

const (
	exact      = 0
	lowerBound = 1
	upperBound = 2
)

var zobristTable [64][12]uint64 // 64 squares, 12 possible pieces (6 white, 6 black)

func InitZobristTable() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 12; j++ {
			zobristTable[i][j] = uint64(rand.Uint32())<<32 + uint64(rand.Uint32()) // Assume randUint64() generates a random 64-bit number
		}
	}
}

func (board *Board) hash() uint64 {
	var h uint64
	for i := 0; i < 64; i++ {
		if piece := board.PieceAt(i); piece != -1 {
			h ^= zobristTable[i][piece]
		}
	}
	return h
}

var (
	TranspositionTable = make(map[uint64]TranspositionEntry)
	tableLock          = sync.RWMutex{} // Mutex to protect map access
)

func getTranspositionEntry(hashKey uint64) (TranspositionEntry, bool) {
	tableLock.RLock()
	entry, exists := TranspositionTable[hashKey]
	tableLock.RUnlock()
	return entry, exists
}

func setTranspositionEntry(hashKey uint64, entry TranspositionEntry) {
	tableLock.Lock()
	TranspositionTable[hashKey] = entry
	tableLock.Unlock()
}

func (board *Board) BestMove(depth int, strategy func(Board) []Move) (Move, Evaluation) {
	InitZobristTable()
	legalMoves := strategy(*board)
	if len(legalMoves) == 0 {
		return Move{}, Evaluation{} // or appropriate error handling
	}

	results := make(chan MoveEvaluation, len(legalMoves))
	defer close(results)

	for _, move := range legalMoves {
		go func(move Move) {
			tmpBoard := *board
			undo, err := tmpBoard.MakeNativeMove(move)
			if err != nil {
				panic(err)
			}
			score := tmpBoard.MiniMax(depth, -9999, 9999, false, strategy)
			tmpBoard.UndoMove(undo)
			results <- MoveEvaluation{Move: move, Score: score}
		}(move)
	}

	// Find the best move based on evaluations
	bestMove := Move{}
	bestScore := Evaluation{-128, -128, -128, -128, -128, -128}

	for range legalMoves {
		result := <-results
		fmt.Println(PieceSymbols[board.PieceAt(int(result.Move.Source))], "(", IndexToPosition(uint64(result.Move.Destination)), ") material: ", result.Score.material, ", centre bonus: ", result.Score.centreBonus, ", mobility bonus: ", result.Score.mobilityBonus, ", pawn structure bonus: ", result.Score.pawnPenalties, ", knight placement bonus: ", result.Score.knightBonus, ", king safety bonus: ", result.Score.safety)
		if result.Score.Sum() > bestScore.Sum() {
			bestScore = result.Score
			bestMove = result.Move
		}
	}

	return bestMove, bestScore
}

func max(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func (board *Board) MiniMax(depth int, alpha, beta int16, maximizingPlayer bool, strategy func(Board) []Move) Evaluation {
	if depth == 0 {
		return board.Evaluate()
	}
	hashKey := board.hash()
	if entry, exists := getTranspositionEntry(hashKey); exists && entry.Depth >= depth {
		switch entry.Flag {
		case exact:
			return entry.Score
		case lowerBound:
			alpha = max(alpha, entry.Score.Sum())
		case upperBound:
			beta = min(beta, entry.Score.Sum())
		}

		if alpha >= beta {
			return entry.Score
		}
	}
	legalMoves := strategy(*board)
	if len(legalMoves) == 0 {
		return board.Evaluate()
	}

	if maximizingPlayer {
		maxEval := Evaluation{material: -128, pawnPenalties: -128, mobilityBonus: -128, centreBonus: -128, safety: -128, knightBonus: -128}
		var bestMove Move
		for _, move := range legalMoves {
			tmpBoard := *board
			undo, err := tmpBoard.MakeNativeMove(move)
			if err != nil {
				panic(err) // Handle the error appropriately.
			}
			eval := tmpBoard.MiniMax(depth-1, -beta, -alpha, false, strategy)
			tmpBoard.UndoMove(undo)

			if eval.Sum() > maxEval.Sum() {
				maxEval = eval
			}
			alpha = max(alpha, eval.Sum())
			if beta <= alpha {
				break // alpha cut-off
			}
		}
		setTranspositionEntry(hashKey, TranspositionEntry{Depth: depth, Score: maxEval, Flag: exact, BestMove: bestMove})
		return maxEval
	} else {
		minEval := Evaluation{material: 127, pawnPenalties: 127, mobilityBonus: 127, centreBonus: 127, safety: 127, knightBonus: 127}
		var bestMove Move
		for _, move := range legalMoves {
			tmpBoard := *board
			undo, err := tmpBoard.MakeNativeMove(move)
			if err != nil {
				panic(err) // Handle the error appropriately.
			}
			eval := tmpBoard.MiniMax(depth-1, -beta, -alpha, true, strategy)
			tmpBoard.UndoMove(undo)

			if eval.Sum() < minEval.Sum() {
				minEval = eval
			}
			beta = min(beta, eval.Sum())
			if alpha >= beta {
				break // beta cut-off
			}
		}
		setTranspositionEntry(hashKey, TranspositionEntry{Depth: depth, Score: minEval, Flag: exact, BestMove: bestMove})
		return minEval
	}
}
