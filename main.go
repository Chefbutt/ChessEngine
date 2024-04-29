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
	b := board.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter moves in standard chess notation (e.g., 'e2e4'), type 'exit' to quit:")

	for {
		fmt.Print("Enter move: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input, please try again.")
			continue
		}

		// Convert CRLF to LF and trim any surrounding whitespace
		text = strings.TrimSpace(strings.Replace(text, "\r\n", "", -1))

		if text == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		// Validate move before making it

		// Make the move
		b.MakeHumanMove(text)
		fmt.Println("Move made:", text)
		b.Display() // Assuming there's a function to display the board state

		b.MakeMove()
		b.Display()
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
