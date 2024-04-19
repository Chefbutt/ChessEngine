package json_converter

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type JsonVariation struct {
	Evaluation int    `json:"cp"`
	Mate       int    `json:"mate"`
	Line       string `json:"line"`
}
type Evals struct {
	Variation []JsonVariation `json:"pvs"`
}

type JsonPosition struct {
	Fen   string  `json:"fen"`
	Evals []Evals `json:"evals"`
}

func UnmarshallPosition(opening []byte) JsonPosition {
	var dat JsonPosition

	if err := json.Unmarshal(opening, &dat); err != nil {
		panic(err)
	}

	return dat
}

func UseLinesFromFiles(fileName string, useLine func([]byte)) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		useLine(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
