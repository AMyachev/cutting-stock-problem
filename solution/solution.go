package solution

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Solution interface {
	GetCountUsedMaterials() int
	GetFreeLength(materialNumber int) int
	GetMaterialLength() int
	CutDetail(materialNumber int, detailNumber int, length int) error
	CutDetailFromNewMaterial(detailNumber int, length int)
	fixCut(materialNumber int, detailNumber int)
}

//Facade
func MakeEmptySolution(materialLength int) Solution {
	return &SolutionImpl{
		materialLength: materialLength,
		materialsFreeLength: []int{},
		piecesNumbersForMaterials: [][]int{},
	}
}


type SolutionImpl struct {
	materialLength int
	materialsFreeLength []int
	piecesNumbersForMaterials [][]int
}


func (solution *SolutionImpl) GetCountUsedMaterials() int {
	return len(solution.materialsFreeLength)
}

func (solution *SolutionImpl) GetFreeLength(materialNumber int) int {
	countUsedMaterials := solution.GetCountUsedMaterials()
	if materialNumber >= countUsedMaterials {
		log.WithFields(log.Fields{
			"materialNumber": materialNumber,
			"countUsedMaterials": countUsedMaterials,
		}).Fatal("GetFreeLength error")
	}

	return solution.materialsFreeLength[materialNumber]
}

func (solution *SolutionImpl) GetMaterialLength() int {
	return solution.materialLength
}

func (solution *SolutionImpl) CutDetail(materialNumber int, detailNumber int, length int) error {
	freeLength := solution.GetFreeLength(materialNumber)
	if length > freeLength {
		return fmt.Errorf("Error in CutDetail func: length > freeLength")
	}

	solution.fixCut(materialNumber, detailNumber)
	solution.materialsFreeLength[materialNumber] -= length
	return nil
}

func (solution *SolutionImpl) CutDetailFromNewMaterial(detailNumber int, length int) {
	solution.fixCut(solution.GetCountUsedMaterials(), detailNumber)
	solution.materialsFreeLength = append(solution.materialsFreeLength, solution.GetMaterialLength() - length)
}

func (solution *SolutionImpl) fixCut(materialNumber int, detailNumber int) {
	countUsedMaterials := solution.GetCountUsedMaterials()
	if countUsedMaterials == materialNumber {
		solution.piecesNumbersForMaterials = append(solution.piecesNumbersForMaterials, []int{detailNumber})
	} else if countUsedMaterials > materialNumber {
		solution.piecesNumbersForMaterials[materialNumber] = append(solution.piecesNumbersForMaterials[materialNumber], detailNumber)
	} else {
		log.WithFields(log.Fields{
			"materialNumber": materialNumber,
			"countUsedMaterials": countUsedMaterials,
		}).Fatal("fixCut error")
	}
}
