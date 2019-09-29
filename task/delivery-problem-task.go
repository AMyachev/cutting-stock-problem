package task

import (
	"io/ioutil"
	"github.com/emirpasic/gods/sets/treeset"
	log "github.com/sirupsen/logrus"
	"strings"
	"strconv"
	"fmt"
)

type Vertex interface {
	ID() int

	LowerBound() int
	SetLowerBound(int)

	UpperBound() int
	SetUpperBound(int)

	NextVertexes(task TaskBnB) (Vertex, []Vertex)

	Permutation() []int
	SetPermutation(permutation []int)

	RemainingPoints() []int
}

type TaskBnB interface {
	CreateInitVertexSet() *treeset.Set
	Dimension() int
	DeliveryTime(firstPoint int, secondPoint int) int
	DirectiveTime(point int) int

	LowerBound(vertex Vertex) int
	UpperBound(vertex Vertex, modificate bool) int

	IncreaseCountTraversedVertexes(n int)
	CountTraversedVertexes() int

	Criterion([]int) (int, int)
}

type taskBnB struct {
	dimension int

	// 1 x dimension
	directiveTimes []int

	// (dimension + 1) x (dimension + 1)
	deliveryTimes [][]int

	branchingStrategy string
	lowerBoundStrategy string
	upperBoundStrategy string

	countTraversedVertexes int
}

func MakeDeliveryProblemFromFile(taskFile, branchingStrategy, lowerBoundStrategy, upperBoundStrategy string) TaskBnB {
	fmt.Println(taskFile)
	content, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}
  
	contentLines := strings.Split(string(content), "\n")
	
	dimension, err := strconv.Atoi(strings.TrimSpace(contentLines[0]))
	if err != nil {
		log.Fatal(err)
	}

	directiveTimesString := strings.Split(strings.TrimSpace(contentLines[1]), " ")
	directiveTimesInt := make([]int, dimension)
	for pos, value := range directiveTimesString {
		directiveTimesInt[pos], err = strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	deliveryTimes := make([][]int, dimension + 1)
	for i := 0; i < dimension + 1; i++ {
		deliveryTimesString := strings.Split(strings.TrimSpace(contentLines[2 + i]), "\t")
		deliveryTimesInt := make([]int, dimension + 1)
		for pos, value := range deliveryTimesString {
			deliveryTimesInt[pos], err = strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
		}
		deliveryTimes[i] = deliveryTimesInt
	}

	log.WithFields(log.Fields{
		"dimension": dimension,
		"directiveTimes": directiveTimesInt,
		"deliveryTimes": deliveryTimes,

		"branchingStrategy": branchingStrategy,
		"lowerBoundStrategy": lowerBoundStrategy,
		"upperBoundStrategy": upperBoundStrategy,
		"countTraversedVertexes": 0,
	}).Info("creating taskBnB ...")

	return &taskBnB{
		dimension: dimension,
		directiveTimes: directiveTimesInt,
		deliveryTimes: deliveryTimes,

		branchingStrategy: branchingStrategy,
		lowerBoundStrategy: lowerBoundStrategy,
		upperBoundStrategy: upperBoundStrategy,
		countTraversedVertexes: 0,
	}
}

func (task *taskBnB) optimisticBranching() func(a, b interface{}) int {
	return func(a, b interface{}) int {
		aVert := a.(Vertex)
		bVert := b.(Vertex)
		//aUB := task.UpperBound(aVert)
		//bUB := task.UpperBound(bVert)

		if aVert.ID() == bVert.ID() {
			return 0
		}

		if aVert.ID() > bVert.ID() {
			return -1
		} else {
			return 1
		}
	}
}

func (task *taskBnB) CreateInitVertexSet() *treeset.Set {
	var comparator func(a, b interface{}) int

	switch task.branchingStrategy {
	case "optimistic":
		comparator = task.optimisticBranching()
	default:
		panic("not implemented")
	}

	//Dimension() exclude zero point
	countRemainingPoints := task.Dimension()
	remainingPoints := make([]int, countRemainingPoints)
	for i := 0; i < countRemainingPoints; i++ {
		remainingPoints[i] = i + 1
	}

	id := 0
	vertex := MakeVertex([]int{0}, remainingPoints, id)
	task.IncreaseCountTraversedVertexes(1)
	return treeset.NewWith(comparator, vertex)
}

func (task *taskBnB) Dimension() int {
	return task.dimension
}

func (task *taskBnB) DeliveryTime(firstPoint int, secondPoint int) int {
	return task.deliveryTimes[firstPoint][secondPoint]
}

func (task *taskBnB) DirectiveTime(point int) int {
	return task.directiveTimes[point]
}

func (task *taskBnB) lowerBoundDefault(vertex Vertex) int {
	permutation := vertex.Permutation()
	result, endTime := task.Criterion(permutation)
	lastPoint := permutation[len(permutation) - 1]

	for _, remainPoint := range vertex.RemainingPoints() {
		currentTime := endTime + task.DeliveryTime(lastPoint, remainPoint)
		//exclude zero point: -1
		if currentTime > task.DirectiveTime(remainPoint - 1) {
			result++
		}
	}

	return result
}

func (task *taskBnB) LowerBound(vertex Vertex) int {
	var lowerBound int

	// Already compute
	if currentLowerBound := vertex.LowerBound(); currentLowerBound != -1 {
		return currentLowerBound
	}

	switch task.lowerBoundStrategy {
	case "default":
		lowerBound = task.lowerBoundDefault(vertex)
	default:
		panic("not implemented")
	}

	vertex.SetLowerBound(lowerBound)
	return lowerBound
}

func (task *taskBnB) upperBoundDefault(vertex Vertex, modificate bool) int {
	var currentTime int
	var posPoint int
	countRemainingPoints := len(vertex.RemainingPoints())
	newPermutation := vertex.Permutation()
	newRemainingPoints := CopySliceInts(vertex.RemainingPoints(), countRemainingPoints)

	for i := 0; i < countRemainingPoints; i++ {
		lastPoint := newPermutation[len(newPermutation) - 1]
		minTime := 10000000
		for pos, remainingPoint := range newRemainingPoints {
			currentTime = task.DeliveryTime(lastPoint, remainingPoint)
			if currentTime < minTime {
				minTime = currentTime
				posPoint = pos
			}
		}
		newPermutation = append(newPermutation, newRemainingPoints[posPoint])

		if posPoint == len(newRemainingPoints) - 1 {
			//Removed last elem
			newRemainingPoints = newRemainingPoints[:posPoint]
		} else {
			newRemainingPoints = append(newRemainingPoints[:posPoint], newRemainingPoints[posPoint+1:]...)
		}
	}

	crit, _ := task.Criterion(newPermutation)
	if modificate {
		vertex.SetPermutation(newPermutation)
	}
	return crit
}

func (task *taskBnB) UpperBound(vertex Vertex, modificate bool) int {
	currentUpperBound := vertex.UpperBound()
	if modificate || currentUpperBound == -1 {
		switch task.upperBoundStrategy {
		case "default":
			currentUpperBound = task.upperBoundDefault(vertex, modificate)
		default:
			panic("not implemented")
		}
	}
	vertex.SetUpperBound(currentUpperBound)

	return currentUpperBound
}

func (task *taskBnB) IncreaseCountTraversedVertexes(n int) {
	task.countTraversedVertexes += n
}

func (task *taskBnB) CountTraversedVertexes() int {
	return task.countTraversedVertexes
}

func (task *taskBnB) Criterion(permutation []int) (int, int) {
	result := 0
	currentTime := 0
	for i := 0; i < len(permutation) - 1; i++ {
		currentTime += task.DeliveryTime(permutation[i], permutation[i+1])
		if currentTime > task.DirectiveTime(permutation[i+1] - 1) {
			//late
			result++
		}
	}
	return result, currentTime
}

type vertexImpl struct {
	id int
	permutation []int
	remainingPoints []int

	upperBound int
	lowerBound int
}

func MakeVertex(permutation []int, remainingPoints []int, id int) Vertex {
	return &vertexImpl{
		id: id,
		permutation: permutation,
		remainingPoints: remainingPoints,

		upperBound: -1,
		lowerBound: -1,
	}
}

func (vert *vertexImpl) ID() int {
	return vert.id
}

func (vert *vertexImpl) LowerBound() int {
	return vert.lowerBound
}

func (vert *vertexImpl) SetLowerBound(lowerBound int) {
	vert.lowerBound = lowerBound
}

func (vert *vertexImpl) UpperBound() int {
	return vert.upperBound
}

func (vert *vertexImpl) SetUpperBound(upperBound int) {
	vert.upperBound = upperBound
}

func (vert *vertexImpl) NextVertexes(task TaskBnB) (Vertex, []Vertex) {
	var bestVertex Vertex

	countRemainingPoints := len(vert.remainingPoints)
	countPermutationPoints := len(vert.permutation)
	resultVertexes := make([]Vertex, countRemainingPoints)

	countVertex := task.CountTraversedVertexes()
	task.IncreaseCountTraversedVertexes(countRemainingPoints)

	for i := 0; i < countRemainingPoints; i++ {
		remainingPoints := CopySliceInts(vert.remainingPoints, countRemainingPoints)

		permutation := CopySliceInts(vert.permutation, countPermutationPoints + 1)
		permutation[countPermutationPoints] = remainingPoints[i]

		if i < (countRemainingPoints - 1) {
			remainingPoints = append(remainingPoints[:i], remainingPoints[i+1:]...)
		} else {
			remainingPoints = remainingPoints[:countRemainingPoints-1]
		}

		id := countVertex + i
		resultVertexes[i] = MakeVertex(permutation, remainingPoints, id)

		if bestVertex == nil {
			//init
			bestVertex = resultVertexes[0]
		} else if task.UpperBound(resultVertexes[i], false) <= task.UpperBound(bestVertex, false) {
			bestVertex = resultVertexes[i]
		}
	}

	return bestVertex, resultVertexes
}

func (vert *vertexImpl) Permutation() []int {
	return vert.permutation
}

func (vert *vertexImpl) SetPermutation(permutation []int) {
	vert.permutation = permutation
	vert.remainingPoints = nil
}

func (vert *vertexImpl) RemainingPoints() []int {
	return vert.remainingPoints
}



func CopySliceInts(slice []int, cap int) []int {
	newSlice := make([]int, cap)
	for i, value := range slice {
		newSlice[i] = value
	}
	return newSlice
}
