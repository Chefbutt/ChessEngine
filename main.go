package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"engine/evaluation/board"
	"engine/evaluation/library"
)

func main() {
	// position.EncodeAllPositions("./position/sources/lichess_db_eval.jsonl", "./position/library/lichess_db_eval")

	// position.DecodeFile("./position/library/lichess_db_eval.dat")
	// getPos()

	b := board.New()
	reader := bufio.NewReader(os.Stdin)
	// 1. e4 e5 2. Nf3 Nc6 3. Bb5 Nf6 4. O-O Bc5 5. Re1 O-O 6. c3 a6 7. Ba4 b5 8. Bc2 Re8 9. d4 exd4 10. cxd4 Bb6 11. e5
	// moves := []string{"e2e4", "e7e5", "g1f3", "b8c6", "f1b5", "g8f6", "e1g1", "f8c5", "f1e1", "e8g8", "c2c3", "c5f2"}

	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		b.MakeHumanMove(text)
		b.MakeMove(text)
	}
}

func getPos(fen string) error {
	start := time.Now()

	idx, _ := library.FindFen("./position/library/lichess_db_eval.dat", fen)

	if idx == -1 {
		return errors.New("move not in DB")
	}

	pos, _ := library.ReadPositionFromFile("./position/library/lichess_db_eval.dat", idx)

	fmt.Print(library.ReverseBoardState(pos.FEN))
	fmt.Print(library.ReverseConvertEval(pos.Line.Eval), " ")
	fmt.Println(library.ReverseConvertMoves(pos.Line.Moves))

	elapsed := time.Since(start)
	log.Printf("Lookup took %s", elapsed)

	return nil
}
