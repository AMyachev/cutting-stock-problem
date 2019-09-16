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


func GreedyAlgorithmForAscending(task taskPackage.Task) solutionPackage.Solution {
	permutation := task.GetAllPiecesByProperty("ascending")

	log.WithFields(log.Fields{
		"task": task,
		"permutation": permutation,
	}).Info("Computing solution by GreedyAlgorithmForAscending ...")

	return GreedyAlgorithm(task, permutation)
}

func GreedyAlgorithmForDescending(task taskPackage.Task) solutionPackage.Solution {
	permutation := task.GetAllPiecesByProperty("descending")

	log.WithFields(log.Fields{
		"task": task,
		"permutation": permutation,
	}).Info("Computing solution by GreedyAlgorithmForDescending ...")

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
	computeCriterion := func(task *taskPackage.Task, permutation *[]int) (int, int) {
		solution := GreedyAlgorithm(*task, *permutation)
		return solution.GetCountUsedMaterials(), solution.GetAllFreeLength()
	}

	bestCriterion, bestAllFreeLength := computeCriterion(&task, &permutation)

	for idx := 0; idx < task.GetCountPieces() - 1; idx++ {
		pairExchange(&permutation, idx, idx + 1)
		tempCriterion, allFreeLength := computeCriterion(&task, &permutation)
		if tempCriterion < bestCriterion || (tempCriterion == bestCriterion) && (allFreeLength < bestAllFreeLength) {
			bestCriterion = tempCriterion
			bestAllFreeLength = allFreeLength
			continue
		}
		pairExchange(&permutation, idx, idx + 1)
	}

	return permutation, GreedyAlgorithm(task, permutation)
}


func firstNotAssigned(permutationIndexes []int) (int, int) {
	for position, idx := range permutationIndexes {
		if idx != -1 {
			return position, idx
		}
	}
}

func assign(permutationIndexes []int, position int) int {
	permutationIndexes[position] = -1
}

func Search(task taskPackage.Task) ([]int, solutionPackage.Solution) {
	solution := solutionPackage.MakeEmptySolution(task.GetMaterialLength())
	permutation := task.GetAllPiecesByProperty("descending")
	assigned_pieces := 0
	all_pieces := task.GetCountPieces()
	materialLength := task.GetMaterialLength()


	for countUsedMaterial := 0;; countUsedMaterial++ {
		position, idxPiece := firstNotAssigned(permutation)
		solution.CutDetailFromNewMaterial(idxPiece, task.GetPieceLength(idxPiece))
		assign(permutation, position)
		assigned_pieces++
		if assigned_pieces == all_pieces {
			break
		}

		materialFreeLength := solution.GetFreeLength(countUsedMaterial)
		for position, idxPiece := range permutation {
			if idxPiece == -1 {
				idxPieceLength := task.GetPieceLength(idxPiece)
				if idxPieceLength > materialFreeLength {
					continue
				} else if idxPieceLength == materialFreeLength {
					solution.CutDetail(countUsedMaterial, idxPiece, task.GetPieceLength(idxPiece))
					assign(permutation, position)
					assigned_pieces++
					break
				}
			}
		}

		if assigned_pieces == all_pieces {
			break
		}

	}
}
