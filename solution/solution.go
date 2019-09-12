package solution

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Solution interface {
	GetCountUsedMaterials() int
	GetFreeLength(materialNumber int) int
	GetMaterialLength() int
	CutDetail(materialNumber int, length int) error
	CutDetailFromNewMaterial(length int)
	fixCut(materialNumber int, length int)
}

type SolutionImpl struct {
	materialLength int
	materialsFreeLength []int
	piecesMaterials [][]int
}

//Constructors
func MakeEmptySolution(materialLength int) Solution {
	return &SolutionImpl{
		materialLength: materialLength,
		materialsFreeLength: []int{},
		piecesMaterials: [][]int{},
	}
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

func (solution *SolutionImpl) CutDetail(materialNumber int, length int) error {
	freeLength := solution.GetFreeLength(materialNumber)
	if length > freeLength {
		return fmt.Errorf("Error in CutDetail func: length > freeLength")
	}

	solution.fixCut(materialNumber, length)
	solution.materialsFreeLength[materialNumber] -= length
	return nil
}

func (solution *SolutionImpl) CutDetailFromNewMaterial(length int) {
	solution.fixCut(solution.GetCountUsedMaterials(), length)
	solution.materialsFreeLength = append(solution.materialsFreeLength, solution.GetMaterialLength() - length)
}

func (solution *SolutionImpl) fixCut(materialNumber int, length int) {
	countUsedMaterials := solution.GetCountUsedMaterials()
	if countUsedMaterials == materialNumber {
		solution.piecesMaterials = append(solution.piecesMaterials, []int{length})
	} else if countUsedMaterials > materialNumber {
		solution.piecesMaterials[materialNumber] = append(solution.piecesMaterials[materialNumber], length)
	} else {
		log.WithFields(log.Fields{
			"materialNumber": materialNumber,
			"countUsedMaterials": countUsedMaterials,
		}).Fatal("fixCut error")
	}
}
