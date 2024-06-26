package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"engine/evaluation/board"
	"engine/evaluation/board/bitboards"
	"engine/evaluation/library"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No mode specified")
		fmt.Println("Usage: go run main.go [engine-vs-engine | engine-vs-human] [debug | no-debug] [depth]")
		os.Exit(1)
	}

	if len(os.Args) < 4 {
		fmt.Println("Missing arguments. Usage: go run main.go [engine-vs-engine | engine-vs-human] [debug | no-debug] [depth]")
		os.Exit(1)
	}

	mode := os.Args[1]

	debug := os.Args[2]

	depth, _ := strconv.Atoi(os.Args[3])

	depth = depth - 2
	if depth == 0 {
		depth = 4
	}

	switch mode {
	case "engine-vs-engine":
		playEngineVsEngine(debug, depth)
	case "engine-vs-human":
		playEngineVsHuman(debug, depth)
	default:
		fmt.Println("Invalid mode specified")
		fmt.Println("Usage: go run main.go [engine-vs-engine | engine-vs-human] [debug | no-debug] [depth]")
		os.Exit(1)
	}
}

func playEngineVsEngine(debug string, depth int) {
	bitboards.InitBitboards()
	board.TranspositionTable = make(map[uint64]board.TranspositionEntry)
	board.InitZobristTable()
	b := board.New()

	if debug == "debug" {
		b.Debug = true
	}

	for {
		err := b.MakeMove(depth)
		if err != nil {
			break
		}
		b.Display()
	}
}

func playEngineVsHuman(debug string, depth int) {
	b := board.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter moves in standard chess notation (e.g., 'e2e4'), type 'exit' to quit:")
	bitboards.InitBitboards()
	board.TranspositionTable = make(map[uint64]board.TranspositionEntry)
	board.InitZobristTable()

	if debug == "debug" {
		b.Debug = true
	}

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
		err = b.MakeHumanMove(text)
		if err != nil {
			fmt.Println("Illegal move, please try again.")
			continue
		}

		fmt.Println("Move made:", text)
		// b.Display() // Assuming there's a function to display the board state

		b.MakeMove(depth)
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
