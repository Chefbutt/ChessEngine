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

	// 1. e4 e5 2. Nf3 Nc6 3. Bb5 Nf6 4. O-O Bc5 5. Re1 O-O 6. c3 a6 7. Ba4 b5 8. Bc2 Re8 9. d4 exd4 10. cxd4 Bb6 11. e5
	moves := []string{"e2e4", "e7e5", "g1f3", "b8c6", "f1b5", "g8f6", "e1g1", "f8c5", "f1e1", "e8g8"}

	for _, move := range moves {
		err := b.MakeMove(move)
		if err != nil {
			panic(err)
		}
	}

	// 1 - baltų ėjimas, -1 - juodų
	fmt.Println("Eval: ", b.Evaluate(0))
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
