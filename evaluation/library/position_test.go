package library

import (
	"fmt"
	"strconv"
	"testing"

	"engine/evaluation/library/json_converter"

	"github.com/stretchr/testify/assert"
)

var jsonPositionTest = json_converter.JsonPosition{
	Fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
	Evals: []json_converter.Evals{{
		Variation: []json_converter.JsonVariation{
			{
				Evaluation: 0,
				Line:       "f4g3 f7f5 e1e5 f5f4 g2e4 h7f6 e4b7 c8b7 g3f4 f6g4",
			},
			{
				Evaluation: 10,
				Line:       "f4g3 f7f5 e1e5 f5f4 g2e4 h7f6 e4b7 c8b7 g3f4 f6g4",
			},
		},
	}},
}

func TestConvertPosition(t *testing.T) {
	name, gotPos := JsonToStoredPosition(jsonPositionTest)
	fmt.Print(name, gotPos)
}

// func TestEncodePosition(t *testing.T) {
// 	for _, pos := range morePositionsTest {
// 		position.EncodeToFile(pos, "output.dat")
// 	}

// 	position.DecodeFile("output.dat")
// }

func TestFormFile(t *testing.T) {
	// rnbqkbnr should be 10100110100011001110100001101010
	// fmt.Println(strconv.FormatInt(position.FormFile("rnbqkbnr"), 2))
	file := FormFile("1rb1k2r")
	assert.Equal(t, strconv.FormatInt(int64(file), 2), "10100110100011001110100001101010")
}

func TestFormExtraneous(t *testing.T) {
	//  w KQkq c6 0 2 should be 01111
	assert.Equal(t, strconv.FormatInt(int64(FormExtraneous(" w KQkq c6 0 2")), 2), "1111")
}

func TestFormBoardState(t *testing.T) {
	//  w KQkq c6 0 2 should be 01111

	// 101001101000110011101000011010100100010011110100010001000100010010011001100110011001100110010000100100110011001100110010000
	pos := FormBoardState("1rb1k2r/p2p1ppp/2p1p3/4P1N1/1b1q1P2/2Bn4/PP1K2PP/2Q2B2 w - - 0 1")
	for _, positionFragment := range pos {
		fmt.Printf("%b\n", positionFragment)
	}

	// assert.Equal(t, strconv.FormatInt(int64(, 2), "1111")
}

func TestReverseFile(t *testing.T) {
	pos := FormBoardState("8/p2b3p/2pPp3/P7/1p6/3B4/kQ5P/1R5K b - - 8 42")
	for _, positionFragment := range pos {
		fmt.Println(ReverseFile(positionFragment))
	}
}

func TestConvertMoves(t *testing.T) {
	moves := ConvertMoves("c3a4 c3a4q")
	for _, move := range moves {
		fmt.Printf("%016b\n", move)
	}
}

func TestReverseConvertMoves(t *testing.T) {
	moves := ConvertMoves("c3a4 c3a4q")

	decoded := ReverseConvertMoves(moves)

	assert.Equal(t, decoded, "c3a4 c3a4q")
}
