package main

import (
	"fmt"
	"log"
	"time"

	"engine/board"
	"engine/position"
)

func main() {
	// position.EncodeAllPositions("./position/sources/lichess_db_eval.jsonl", "./position/library/lichess_db_eval")

	// position.DecodeFile("./position/library/lichess_db_eval.dat")
	// getPos()

	b := board.NewBoard()

	b.OccupiedSquares.Display()

	uci := "e2e4"
	move := b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "e7e5"

	getPos(b.ToFEN())

	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "g1f3"

	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "b8c6"
	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "f1b5"
	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "g8f6"
	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	uci = "e1g1"
	move = b.UCItoMove(uci)

	b.MakeMove(move)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	b.OccupiedSquares.Display()

	fmt.Println(b.ToFEN())
	// getPos()

	fmt.Println()
	fmt.Println()

	// 1 - baltų ėjimas, -1 - juodų
	fmt.Println("Eval: ", b.Evaluate(0))

	b.OccupiedSquares.Display()
}

func getPos(fen string) {
	start := time.Now()

	idx, _ := position.FindFen("./position/library/lichess_db_eval.dat", fen)

	if idx == -1 {
		fmt.Println("Not in db :()")
		return
	}

	pos, _ := position.ReadPositionFromFile("./position/library/lichess_db_eval.dat", idx)

	fmt.Print(position.ReverseBoardState(pos.FEN))
	fmt.Print(position.ReverseConvertEval(pos.Line.Eval), " ")
	fmt.Println(position.ReverseConvertMoves(pos.Line.Moves))

	elapsed := time.Since(start)
	log.Printf("Lookup took %s", elapsed)
}
