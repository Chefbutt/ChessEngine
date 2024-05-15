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
	PreviousAggregateBitboards   AggregateBitboards
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

// func (board *Board) Quiesce(depth int, alpha int, beta int, maximizingPlayer bool) Evaluation {
// 	standPat := toEval(board.Evaluate())
// 	if maximizingPlayer {
// 		if standPat.Sum() >= beta {
// 			return Evaluation{float64(beta), 0, 0, 0, 0, 0}
// 		}
// 		if alpha < standPat.Sum() {
// 			alpha = standPat.Sum()
// 		}
// 	} else {
// 		if standPat.Sum() <= alpha {
// 			return Evaluation{float64(alpha), 0, 0, 0, 0, 0}
// 		}
// 		if beta > standPat.Sum() {
// 			beta = standPat.Sum()
// 		}
// 	}

// 	if depth == 0 {
// 		return toEval(board.Evaluate())
// 	}

// 	legalMoves := Captures(board) // This should be modified to return only capture moves
// 	for _, move := range legalMoves {
// 		if board.PieceAt(move.Destination) == -1 { // Assuming you have a way to determine if the move is a capture
// 			continue
// 		}
// 		undo, err := board.MakeNativeMove(move)
// 		if err != nil {
// 			panic(err) // Handle errors appropriately
// 		}
// 		score := -board.Quiesce(depth-1, -beta, -alpha, !maximizingPlayer).Sum()
// 		board.UndoMove(undo)

// 		if maximizingPlayer {
// 			if score >= beta {
// 				return Evaluation{float64(beta), 0, 0, 0, 0, 0}
// 			}
// 			if score > alpha {
// 				alpha = score
// 			}
// 		} else {
// 			if score <= alpha {
// 				return Evaluation{float64(alpha), 0, 0, 0, 0, 0}
// 			}
// 			if score < beta {
// 				beta = score
// 			}
// 		}
// 	}
// 	if maximizingPlayer {
// 		return Evaluation{float64(alpha), 0, 0, 0, 0, 0}
// 	}
// 	return Evaluation{float64(beta), 0, 0, 0, 0, 0}
// }

type MoveEvaluation struct {
	Move  Move
	Score int
}

type TranspositionEntry struct {
	Depth    int
	Score    int
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

func (board *Board) BestMove(depth int, strategy func(Board) []Move, materialModifier, mobilityModifier, centreModifier, penaltyModifier int) (Move, int) {
	InitZobristTable()
	legalMoves := strategy(*board)
	if len(legalMoves) == 0 {
		return Move{}, 0 // or appropriate error handling
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
			score := MiniMax(tmpBoard, depth, -9999, 9999, false, strategy, materialModifier, mobilityModifier, centreModifier, penaltyModifier)
			tmpBoard.UndoMove(undo)
			results <- MoveEvaluation{Move: move, Score: score}
		}(move)
	}

	// Find the best move based on evaluations
	bestMove := Move{}
	bestScore := -9999

	for range legalMoves {
		result := <-results
		if board.Debug {
			fmt.Println(PieceSymbols[board.PieceAt(int(result.Move.Source))], "(", IndexToPosition(uint64(result.Move.Destination)), ") :", result.Score)
		}
		if result.Score > bestScore {
			bestScore = result.Score
			bestMove = result.Move
		}
	}

	return bestMove, bestScore
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MiniMax(board Board, depth int, alpha, beta int, maximizingPlayer bool, strategy func(Board) []Move, materialModifier, mobilityModifier, centreModifier, penaltyModifier int) int {
	if depth == 0 {
		return board.Evaluate(materialModifier, mobilityModifier, centreModifier, penaltyModifier)
	}
	hashKey := board.hash()
	if entry, exists := getTranspositionEntry(hashKey); exists && entry.Depth >= depth {
		switch entry.Flag {
		case exact:
			return entry.Score
		case lowerBound:
			alpha = max(alpha, entry.Score)
		case upperBound:
			beta = min(beta, entry.Score)
		}

		if alpha >= beta {
			return entry.Score
		}
	}
	legalMoves := strategy(board)
	if len(legalMoves) == 0 && board.IsCheckMate() {
		if maximizingPlayer {
			return -9999
		} else {
			return 9999
		}
	}

	if maximizingPlayer {
		maxEval := -9999
		var bestMove Move
		for _, move := range legalMoves {
			tmpBoard := board
			undo, err := tmpBoard.MakeNativeMove(move)
			if err != nil {
				panic(err) // Handle the error appropriately.
			}

			eval := MiniMax(tmpBoard, depth-1, alpha, beta, false, strategy, materialModifier, mobilityModifier, centreModifier, penaltyModifier)
			tmpBoard.UndoMove(undo)

			maxEval = max(maxEval, eval)
			alpha = max(alpha, eval)

			if beta <= alpha {
				break // alpha cut-off
			}
		}
		if len(legalMoves) == 0 {
			board.Display()
			fmt.Print("f")
		}
		setTranspositionEntry(hashKey, TranspositionEntry{Depth: depth, Score: maxEval, Flag: exact, BestMove: bestMove})
		return maxEval
	} else {
		minEval := 9999
		var bestMove Move
		for _, move := range legalMoves {
			tmpBoard := board
			undo, err := tmpBoard.MakeNativeMove(move)
			if err != nil {
				panic(err) // Handle the error appropriately.
			}

			eval := MiniMax(tmpBoard, depth-1, alpha, beta, true, strategy, materialModifier, mobilityModifier, centreModifier, penaltyModifier)
			tmpBoard.UndoMove(undo)

			minEval = min(eval, minEval)
			beta = min(beta, eval)

			if beta <= alpha {
				break // beta cut-off
			}
		}
		setTranspositionEntry(hashKey, TranspositionEntry{Depth: depth, Score: minEval, Flag: exact, BestMove: bestMove})
		return minEval
	}
}

func (b *Board) QuiescentSearch(alpha, beta, materialModifier, mobilityModifier, centreModifier, penaltyModifier int) int {
	standPat := b.Evaluate(materialModifier, mobilityModifier, centreModifier, penaltyModifier)
	if standPat >= beta {
		return beta
	}
	if alpha < standPat {
		alpha = standPat
	}

	captures := b.Captures()
	for _, move := range captures {
		b.makeMove(move)
		score := -b.QuiescentSearch(-beta, -alpha, materialModifier, mobilityModifier, centreModifier, penaltyModifier)

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
