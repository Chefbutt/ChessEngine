package position_test

import (
	"fmt"
	"strconv"
	"testing"

	"engine/position"
	"engine/position/json_converter"

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

var jsonPositionTestTwo = json_converter.JsonPosition{
	Fen: "2bq1rk1/pr3ppn/1p2p3/7P/2pP1B1P/2P5/PPQ2PB1/R3R1K1 w - -",
	Evals: []json_converter.Evals{{
		Variation: []json_converter.JsonVariation{
			{
				Evaluation: 311,
				Line:       "g2e4 f7f5 e4b7 c8b7 f2f3 b7f3 e1e6 d8h4 c2h2 h4g4",
			},
			{
				Evaluation: 292,
				Line:       "f4g3 f7f5 e1e5 d8f6 a1e1 b7f7 g2c6 f8d8 d4d5 e6d5",
			},
		},
	}},
}

var morePositionsTest = []json_converter.JsonPosition{
	jsonPositionTest,
	jsonPositionTestTwo,
	jsonPositionTest,
	jsonPositionTestTwo,
	jsonPositionTest,
	jsonPositionTestTwo,
}

func TestConvertPosition(t *testing.T) {
	name, gotPos := position.JsonToStoredPosition(jsonPositionTest)
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
	file := position.FormFile("rnbqkbnr")
	assert.Equal(t, strconv.FormatInt(int64(file), 2), "10100110100011001110100001101010")
}

func TestFormExtraneous(t *testing.T) {
	//  w KQkq c6 0 2 should be 01111
	assert.Equal(t, strconv.FormatInt(int64(position.FormExtraneous(" w KQkq c6 0 2")), 2), "1111")
}

func TestFormBoardState(t *testing.T) {
	//  w KQkq c6 0 2 should be 01111

	// 10100110100011001110100001101010
	// 01000100111101000100010001000100
	// 10011001100110011001100110010000
	//  100100110011001100110010000
	pos := position.FormBoardState("2bq1rk1/pr3ppn/1p2p3/7P/2pP1B1P/2P5/PPQ2PB1/R3R1K1 w - -")
	for _, positionFragment := range pos {
		fmt.Print(positionFragment, " ")
	}

	// assert.Equal(t, strconv.FormatInt(int64(, 2), "1111")
}

func TestReverseFile(t *testing.T) {
	pos := position.FormBoardState("8/p2b3p/2pPp3/P7/1p6/3B4/kQ5P/1R5K b - - 8 42")
	for _, positionFragment := range pos {
		fmt.Println(position.ReverseFile(positionFragment))
	}
}

func TestConvertMoves(t *testing.T) {
	moves := position.ConvertMoves("c3a4 c3a4q")
	for _, move := range moves {
		fmt.Printf("%016b\n", move)
	}
}

func TestReverseConvertMoves(t *testing.T) {
	moves := position.ConvertMoves("c3a4 c3a4q")

	decoded := position.ReverseConvertMoves(moves)

	assert.Equal(t, decoded, "c3a4 c3a4q")
}
