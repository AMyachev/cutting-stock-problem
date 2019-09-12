package heuristics

import (
	taskPackage "csp/task"
	solutionPackage "csp/solution"

	log "github.com/sirupsen/logrus"
)


func GreedyAlgorithm(task taskPackage.Task, permutation []int) solutionPackage.Solution {
	solution := solutionPackage.MakeEmptySolution(task.GetMaterialLength())
	
	for _, idxDetail := range permutation {
		detailCuttedOff := false
		detailLength := task.GetPieceLength(idxDetail)
		for idxMaterial := 0; idxMaterial < solution.GetCountUsedMaterials(); idxMaterial++ {
			if err := solution.CutDetail(idxMaterial, detailLength); err != nil {
				continue	
			}
			detailCuttedOff = true
		}
		if !detailCuttedOff {
			solution.CutDetailFromNewMaterial(detailLength)
		}
	}

	return solution
}

func CopySliceInts(slice []int) []int {
	newSlice := make([]int, 0)
	return append(newSlice, slice...)
}

func GreedyAlgorithmByAscending(task taskPackage.Task) solutionPackage.Solution {
	permutation := task.GetAscendingPermutation()

	log.WithFields(log.Fields{
		"task": task,
		"permutation": permutation,
	}).Info("Computing solution by GreedyAlgorithm ...")

	return GreedyAlgorithm(task, permutation)
}
