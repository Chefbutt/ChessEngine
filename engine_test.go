package main

import (
	"bufio"
	"engine/evaluation/board"
	"engine/evaluation/board/bitboards"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgainstItself(t *testing.T) {
	bitboards.InitBitboards()
	board.TranspositionTable = make(map[uint64]board.TranspositionEntry)
	board.InitZobristTable()
	b := board.New()

	for {
		err := b.MakeMove()
		if err != nil {
			break
		}
		b.Display()
	}

	assert.True(t, true)
}

func TestAgainstHuman(t *testing.T) {
	b := board.New()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter moves in standard chess notation (e.g., 'e2e4'), type 'exit' to quit:")
	bitboards.InitBitboards()
	board.TranspositionTable = make(map[uint64]board.TranspositionEntry)
	board.InitZobristTable()

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

		b.MakeMove()
		b.Display()
	}
}
