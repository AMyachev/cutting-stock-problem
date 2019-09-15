package heuristics

import (
	"fmt"
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
			if err := solution.CutDetail(idxMaterial, idxDetail, detailLength); err != nil {
				continue	
			}
			detailCuttedOff = true
			break
		}
		if !detailCuttedOff {
			solution.CutDetailFromNewMaterial(idxDetail, detailLength)
		}
	}

	return solution
}

func CopySliceInts(slice []int) []int {
	newSlice := make([]int, 0)
	return append(newSlice, slice...)
}

func GreedyAlgorithmByAscending(task taskPackage.Task) solutionPackage.Solution {
	permutation := task.GetAllPiecesByProperty("ascending")

	log.WithFields(log.Fields{
		"task": task,
		"permutation": permutation,
	}).Info("Computing solution by GreedyAlgorithm ...")

	return GreedyAlgorithm(task, permutation)
}


func pairExchange(permutation *[]int, i int, j int) {
	if i >= len(*permutation) || j >= len(*permutation) {
		log.WithFields(log.Fields{
			"first_index": i,
			"secong_index": j,
			"permutation": permutation,
		}).Fatal("pairExchange error ...")
	}

	temp := (*permutation)[i]
	(*permutation)[i] = (*permutation)[j]
	(*permutation)[j] = temp
}


func LocalOptimization(task taskPackage.Task, permutation []int) ([]int, solutionPackage.Solution) {
	computeCriterion := func(task taskPackage.Task, permutation []int) (int, int) {
		solution := GreedyAlgorithm(task, permutation)
		return solution.GetCountUsedMaterials(), solution.GetAllFreeLength()
	}

	bestCriterion, bestAllFreeLength := computeCriterion(task, permutation)

	for idx := 0; idx < task.GetCountPieces() - 1; idx++ {
		pairExchange(&permutation, idx, idx + 1)
		tempCriterion, allFreeLength := computeCriterion(task, permutation)
		if tempCriterion < bestCriterion || (tempCriterion == bestCriterion) && (allFreeLength < bestAllFreeLength) {
			bestCriterion = tempCriterion
			bestAllFreeLength = allFreeLength
			fmt.Println("better")
			continue
		}
		pairExchange(&permutation, idx, idx + 1)
	}

	return permutation, GreedyAlgorithm(task, permutation)
}
