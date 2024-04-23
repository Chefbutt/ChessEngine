package library

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"engine/evaluation/library/json_converter"
)

type StoredPosition struct {
	Eval int
	Line string
}

func scaleMate(mate int) int {
	if mate > 0 {
		return 999
	}
	return -999
}

func JsonToStoredPosition(jsonPosition json_converter.JsonPosition) (string, StoredPosition) {
	if jsonPosition.Evals[0].Variation[0].Mate != 0 {
		return jsonPosition.Fen, StoredPosition{Eval: scaleMate(jsonPosition.Evals[0].Variation[0].Mate), Line: jsonPosition.Evals[0].Variation[0].Line}
	}

	return jsonPosition.Fen, StoredPosition{Eval: jsonPosition.Evals[0].Variation[0].Evaluation, Line: jsonPosition.Evals[0].Variation[0].Line}
}

// 1 to 6 are white piece ids
const (
	King uint8 = 1 + iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
)

const (
	Empty byte = 15
)

func formPiece(character rune) uint8 {
	var piece uint8

	if unicode.IsUpper(character) {
		piece = piece + 1 // Make piece black
		character = unicode.ToLower(character)
	}

	switch character {
	case 'p':
		piece = piece + Pawn<<1
	case 'n':
		piece = piece + Knight<<1
	case 'b':
		piece = piece + Bishop<<1
	case 'q':
		piece = piece + Queen<<1
	case 'r':
		piece = piece + Rook<<1
	case 'k':
		piece = piece + King<<1
	}

	return piece
}

func FormExtraneous(extraneous string) uint32 {
	var boardState uint32
	for _, character := range extraneous {
		switch character {
		case 'w':
			boardState = boardState + 16
		case 'K':
			boardState = boardState + 8
		case 'Q':
			boardState = boardState + 4
		case 'k':
			boardState = boardState + 2
		case 'q':
			boardState = boardState + 1
		default:
		}
	}
	return boardState
}

func FormFile(file string) uint32 {
	// Actually only takes up 32 bits
	var formedFile uint32

	var elementId int
	for _, character := range file {

		formedFile = formedFile << 4
		if unicode.IsDigit(character) {
			times, err := strconv.ParseInt(string(character), 10, 64)
			if err != nil {
				panic(err)
			}
			for time := range times {
				elementId++
				formedFile = formedFile + 15
				if times > 1 && time != times-1 {
					formedFile = formedFile << 4
				}
			}
			continue
		} else {
			elementId++
			// 4 bits
			piece := formPiece(character)

			formedFile = formedFile + uint32(piece)
		}
	}
	return formedFile
}

func FormBoardState(board string) [9]uint32 {
	// rnbqkbnr/pp1ppppp/8/2p5/4P3/8/PPPP1PPP/RNBQKBNR w KQkq c6 0 2
	parts := strings.SplitN(board, " ", 2)

	files := strings.SplitN(parts[0], "/", 8)

	var filesInBinary [9]uint32

	for id, file := range files {
		formedFile := FormFile(file)

		filesInBinary[id] = formedFile
	}

	filesInBinary[8] = FormExtraneous(parts[1])

	return filesInBinary
}

func ReverseFile(file uint32) string {
	mapping := map[uint32]rune{
		0b1011: 'R',
		0b0111: 'N',
		0b1001: 'B',
		0b1101: 'Q',
		0b0011: 'K',
		0b0101: 'P',
		0b1010: 'r',
		0b0110: 'n',
		0b1000: 'b',
		0b1100: 'q',
		0b0010: 'k',
		0b0100: 'p',
		0b1111: ' ',
	}

	var result string

	for i := 0; i < 8; i++ {
		segment := (file >> (4 * i)) & 0xF

		if char, exists := mapping[segment]; exists {
			result = string(char) + result
		}
	}

	return result
}

func ReverseExtraneous(state uint32) string {
	var result string
	if state&16 != 0 {
		result += "w "
	} else {
		result += "b "
	}
	if state&8 != 0 {
		result += "K"
	}
	if state&4 != 0 {
		result += "Q"
	}
	if state&2 != 0 {
		result += "k"
	}
	if state&1 != 0 {
		result += "q"
	}
	return result
}

func ReverseBoardState(boardState [9]uint32) string {
	var board strings.Builder

	// Reverse each file and join them with '/'
	for i, file := range boardState {
		if file == 0 {
			continue
		}
		if i < len(boardState)-1 { // The last one is for extraneous information
			if row := ReverseFile(file); row != "" {
				spaceRegexp := regexp.MustCompile(` +`)

				row := spaceRegexp.ReplaceAllStringFunc(row, func(match string) string {
					return strconv.Itoa(len(match))
				})
				board.WriteString(row)
			}
		}
		if i < 7 {
			board.WriteString("/")
		}
	}

	// Append extraneous information
	extraneous := ReverseExtraneous(boardState[len(boardState)-1])
	board.WriteString(" ")
	board.WriteString(extraneous)

	return board.String()
}

func ConvertEval(eval int) uint8 {
	return uint8((eval + 999) * 255 / 1998)
}

func ReverseConvertEval(val uint8) int {
	newVal := uint(val)
	return int(newVal*1998/255 - 999)
}

func convertPair(pair string) uint16 {
	// Character to numeric mapping
	charMap := map[rune]int{
		'a': 7,
		'b': 6,
		'c': 5,
		'd': 4,
		'e': 3,
		'f': 2,
		'g': 1,
		'h': 0,
	}

	// Convert character-number pair
	char := charMap[rune(pair[0])]
	number, _ := strconv.Atoi(string(pair[1]))
	value := char + (number-1)*8

	return uint16(value)
}

func ConvertMoves(input string) [10]uint16 {
	addCharMap := map[rune]uint16{
		'r': 0b1101,
		'n': 0b1011,
		'b': 0b1100,
		'q': 0b1110,
	}

	var output [10]uint16
	moves := strings.Split(input, " ")

	for id, move := range moves {
		var bits uint16
		var addChar rune
		if len(move) > 4 {
			addChar = rune(move[4])
			move = move[:4]
			bits = addCharMap[addChar] << 12
		}

		firstPairValue := convertPair(move[:2])
		secondPairValue := convertPair(move[2:])

		bits = bits + firstPairValue<<6
		bits = bits + secondPairValue

		output[id] = bits
	}

	// delimiter := uint16(16383) // "000011111111111"
	// output = append(output, delimiter)

	return output
}

func reverseConvertPair(value uint16) string {
	fileMap := map[uint16]rune{
		7: 'a', 6: 'b', 5: 'c', 4: 'd',
		3: 'e', 2: 'f', 1: 'g', 0: 'h',
	}

	charPart := value % 8
	numberPart := value/8 + 1

	char := fileMap[charPart]
	numberStr := strconv.Itoa(int(numberPart))

	return string(char) + numberStr
}

func ReverseConvertMoves(encodedMoves [10]uint16) string {
	reverseAddCharMap := map[uint16]rune{
		0b1101: 'r', 0b1011: 'n', 0b1100: 'b', 0b1110: 'q',
	}

	var moves []string

	for _, encodedMove := range encodedMoves {
		addCharBits := encodedMove >> 12
		firstPairValue := (encodedMove >> 6) & 0x3F // First 6 bits
		secondPairValue := encodedMove & 0x3F       // Last 6 bits

		firstPairStr := reverseConvertPair(firstPairValue)
		secondPairStr := reverseConvertPair(secondPairValue)

		if firstPairStr == secondPairStr {
			break
		}

		moveStr := firstPairStr + secondPairStr

		if char, exists := reverseAddCharMap[addCharBits]; exists {
			moveStr += string(char)
		}

		moves = append(moves, moveStr)
	}

	return strings.Join(moves, " ")
}

type positionLine struct {
	Eval  uint8
	Moves [10]uint16
}

type BinaryPosition struct {
	FEN  [9]uint32
	Line positionLine
}

func ReadPositionFromFile(filename string, index int64) (*BinaryPosition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Calculate the offset in the file
	offset := int64(binary.Size(BinaryPosition{})) * index
	if _, err := file.Seek(offset, 0); err != nil {
		return nil, err
	}

	var pos BinaryPosition
	if err := binary.Read(file, binary.LittleEndian, &pos); err != nil {
		return nil, err
	}

	return &pos, nil
}

func FindFen(filename string, fen string) (int64, error) {
	return binarySearch(filename, FormBoardState(fen))
}

func binarySearch(filename string, target [9]uint32) (int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	size, err := file.Stat()
	if err != nil {
		return -1, err
	}

	left, right := int64(0), int64(int(size.Size())/binary.Size(BinaryPosition{}))-1
	for left <= right {
		mid := left + (right-left)/2

		midPos, err := ReadPositionFromFile(filename, mid)
		if err != nil {
			return -1, err
		}

		switch compareFEN(midPos.FEN, target) {
		case -1:
			left = mid + 1
		case 1:
			right = mid - 1
		case 0:
			return mid, nil
		}
	}

	return -1, nil // not found
}

func compareFEN(a, b [9]uint32) int {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

func EncodeToFile(pos json_converter.JsonPosition) BinaryPosition {
	name, gotPos := JsonToStoredPosition(pos)

	var binaryPosition BinaryPosition

	binaryPosition.FEN = FormBoardState(name)

	var binaryLine positionLine

	binaryLine.Eval = ConvertEval(int(gotPos.Eval))

	binaryLine.Moves = ConvertMoves(gotPos.Line)

	binaryPosition.Line = binaryLine

	return binaryPosition
}

func DecodeFile(fileName string) []json_converter.JsonPosition {
	var positions []json_converter.JsonPosition

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var position json_converter.JsonPosition

	var uint32s [9]uint32

	// Read []uint32
	for i := range uint32s {
		var uint32 uint32
		err = binary.Read(reader, binary.LittleEndian, &uint32)
		if err != nil && err != io.EOF {
			panic(err)
		}

		uint32s[i] = uint32

		if err == io.EOF {
			break
		}
	}

	// Print read values to verify
	position.Fen = ReverseBoardState(uint32s)

	var evals json_converter.Evals

	var singleEval json_converter.JsonVariation

	var evaluation uint8

	err = binary.Read(reader, binary.LittleEndian, &evaluation)
	if err != nil && err != io.EOF {
		panic(err)
	}

	singleEval.Evaluation = ReverseConvertEval(evaluation)

	moves, err := readMoves(reader)
	if err != nil {
		panic(err)
	}

	singleEval.Line = moves

	evals.Variation = append(evals.Variation, singleEval)

	position.Evals = append(position.Evals, evals)

	positions = append(positions, position)

	return positions
}

func readMoves(reader io.Reader) (moves string, err error) {
	var uint16s [10]uint16
	for id := range uint16s {
		var num uint16
		err = binary.Read(reader, binary.LittleEndian, &num)
		if err != nil {
			panic(err) // Handle EOF or other reading errors appropriately
		}
		uint16s[id] = num

	}

	moves = ReverseConvertMoves(uint16s)

	return moves, err
}

func EncodeAllPositions(fromFile, toFile string) {
	file, err := os.Create(toFile + ".dat")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)

	var binaryPositions []BinaryPosition
	json_converter.UseLinesFromFiles(fromFile, func(s []byte) {
		pos := json_converter.UnmarshallPosition(s)

		binaryPositions = append(binaryPositions, EncodeToFile(pos))
	})

	sort.Slice(binaryPositions, func(i, j int) bool {
		switch compareFEN(binaryPositions[i].FEN, binaryPositions[j].FEN) {
		case -1:
			return true
		case 1:
			return false
		case 0:
			return false
		}
		return false
	})

	for _, binaryPosition := range binaryPositions {
		err := binary.Write(writer, binary.LittleEndian, binaryPosition.FEN)
		if err != nil {
			panic(err)
		}

		err = binary.Write(writer, binary.LittleEndian, binaryPosition.Line.Eval)
		if err != nil {
			panic(err)
		}
		err = binary.Write(writer, binary.LittleEndian, binaryPosition.Line.Moves)
		if err != nil {
			panic(err)
		}

		writer.Flush()
	}

	file.Close()
}
