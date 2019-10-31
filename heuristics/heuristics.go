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


func removePiece(permutation *[]int, piece int) {
	idxPiece := -1
	for idx, currentPiece := range *permutation {
		if currentPiece == piece {
			idxPiece = idx
		}
	}

	if idxPiece == -1 {
		panic("removeElem error")
	}

	*permutation = append((*permutation)[:idxPiece], (*permutation)[idxPiece+1:]...)
}

func Search(task taskPackage.Task) solutionPackage.Solution {
	solution := solutionPackage.MakeEmptySolution(task.GetMaterialLength())
	permutation := task.GetAllPiecesByProperty("descending")
	allCountPiece := task.GetCountPieces()
	materialLength := task.GetMaterialLength()


	for countUsedMaterial := 0;; countUsedMaterial++ {
		possibleSolutions := make([][]int, 1)

		//the first element in the solution is remaining the free length
		possibleSolutions[0] = []int{materialLength}
		for _, piece := range permutation {
			for solutionNumber := 0; solutionNumber < len(possibleSolutions); solutionNumber++ {
				materialFreeLength := possibleSolutions[solutionNumber][0]
				pieceLength := task.GetPieceLength(piece)

				if pieceLength <= materialFreeLength {
					//decrease free length
					possibleSolutions[solutionNumber][0] -= pieceLength
					possibleSolutions[solutionNumber] = append(possibleSolutions[solutionNumber], piece)
				}
			}
			//add new solution
			possibleSolutions = append(possibleSolutions, []int{materialLength})
		}

		bestSolution := 0
		minFreeLength := materialLength
		minCountPieces := allCountPiece

		for idxSolution, solution := range possibleSolutions {
			freeLength := solution[0]
			countPieces := len(solution) - 1
			if (freeLength < minFreeLength) || ((freeLength == minFreeLength) && (countPieces < minCountPieces)) {
				minFreeLength = freeLength
				minCountPieces = countPieces
				bestSolution = idxSolution
			}
		}

		firstPiece := possibleSolutions[bestSolution][1]
		removePiece(&permutation, firstPiece)
		solution.CutDetailFromNewMaterial(firstPiece, task.GetPieceLength(firstPiece))
		for _, piece := range possibleSolutions[bestSolution][2:] {
			solution.CutDetail(countUsedMaterial, piece, task.GetPieceLength(piece))
			removePiece(&permutation, piece)
		}

		if len(permutation) == 0 {
			break
		}
	}

	return solution
}
