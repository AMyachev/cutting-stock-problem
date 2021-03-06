package task

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"math"
	"sort"

	log "github.com/sirupsen/logrus"
)

type Task interface {
	GetMaterialLength() int
	GetCountPieces() int
	GetPieceLength(numberPiece int) int
	ComputeLowerBound() int
	GetAllPiecesByProperty(property string) []int
}


// Facade
func MakeOneDimensionalCuttingStockProblem(materialLength int, piecesLength []int) (Task, error) {
	for _, pieceLength := range piecesLength {
		if materialLength < pieceLength {
			return nil, fmt.Errorf("materialLength < pieceLength: %d < %d", materialLength, pieceLength)
		}
	}

	return oneDimensionalCuttingStockProblem{materialLength: materialLength, countPieces: len(piecesLength), piecesLength: piecesLength}, nil
}

func MakeOneDimensionalCuttingStockProblemFromFile(fileName string) (Task, error) {
	contentBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	content := string(contentBytes)
	splittedContent := strings.Split(content, " ")

	splittedContentInts := make([]int, len(splittedContent))

	for idx, value := range splittedContent {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		splittedContentInts[idx] = intValue
	}

	materialLength := splittedContentInts[0]
	pieceLength := splittedContentInts[1:]

	return MakeOneDimensionalCuttingStockProblem(materialLength, pieceLength)
}


type oneDimensionalCuttingStockProblem struct {
	materialLength int
	countPieces int
	piecesLength []int
}

func (problem oneDimensionalCuttingStockProblem) GetMaterialLength() int {
	return problem.materialLength
}

func (problem oneDimensionalCuttingStockProblem) GetCountPieces() int {
	return problem.countPieces
}

func (problem oneDimensionalCuttingStockProblem) GetPieceLength(numberPiece int) int {
	countPieces := problem.GetCountPieces()
	if numberPiece > countPieces {
		log.WithFields(log.Fields{
			"numberPiece": numberPiece,
			"countPieces": countPieces,
		}).Fatal("GetPieceLength error")
	}

	return problem.piecesLength[numberPiece]
}

func (problem oneDimensionalCuttingStockProblem) ComputeLowerBound() int {
	result := 0
	for _, length := range problem.piecesLength {
		result += length
	}

	return int(math.Round(float64(result)/float64(problem.GetMaterialLength())))
}

func (problem oneDimensionalCuttingStockProblem) GetAllPiecesByProperty(property string) []int {
	//format idx, pieceLength
	piecesLengthWithIdx := make([][2]int, problem.GetCountPieces())
	for idx, pieceLength := range problem.piecesLength {
		piecesLengthWithIdx[idx][0] = idx
		piecesLengthWithIdx[idx][1] = pieceLength
	}

	if property == "ascending" {
		sort.Slice(piecesLengthWithIdx, func(i, j int) bool {
			return piecesLengthWithIdx[i][1] < piecesLengthWithIdx[j][1]
		})
	} else if property == "descending" {
		sort.Slice(piecesLengthWithIdx, func(i, j int) bool {
			return piecesLengthWithIdx[i][1] > piecesLengthWithIdx[j][1]
		})
	} else {
		panic("not implemented property")
	}

	resultSlice := make([]int, problem.GetCountPieces())
	for idx, pieceLengthWitholdIdx := range piecesLengthWithIdx {
		resultSlice[idx] = pieceLengthWitholdIdx[0]
	}

	return resultSlice
}
