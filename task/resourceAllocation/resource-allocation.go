package resourceAllocation

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const MaxInt = int(^uint(0) >> 1)

type resourceAllocationTask struct {
	// N
	countSuppliers int
	// M
	countCastomers int
	// T
	countTacts int

	// a1, a2, ..., aN
	volumeGoodsFromSuppliers []int

	// b11, b12, ..., b1T
	// ...
	// bN1, bN2, ..., bNT
	volumeGoodsFromSuppliersOnTact [][]int

	// c11, c12, ..., c1T
	// ...
	// cM1, cM2, ..., cMT
	volumeUsedGoodsByCastomersOnTact [][]int

	// D (D1, D2, ..., DM)
	setsSuppliersForCastomers [][]int
}

func MakeResourceAllocationTaskFromFile(taskFile string) *resourceAllocationTask {
	content, err := ioutil.ReadFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}
	contentLines := strings.Split(string(content), "\r\n")

	countSuppliers, err := strconv.Atoi(contentLines[0])
	if err != nil {
		log.Fatal(err)
	}

	countCastomers, err := strconv.Atoi(contentLines[1])
	if err != nil {
		log.Fatal(err)
	}

	countTacts, err := strconv.Atoi(contentLines[2])
	if err != nil {
		log.Fatal(err)
	}

	volumeGoodsFromSuppliersString := strings.Trim(contentLines[4], " ")
	volumeGoodsFromSuppliersStrings := strings.Split(volumeGoodsFromSuppliersString, " ")
	volumeGoodsFromSuppliers := make([]int, countSuppliers)
	for i, volumeGood := range volumeGoodsFromSuppliersStrings {
		volumeGoodsFromSuppliers[i], err = strconv.Atoi(volumeGood)
		if err != nil {
			log.Fatal(err)
		}
	}

	// allocation B
	volumeGoodsFromSuppliersOnTact := make([][]int, countSuppliers)
	for i := 0; i < countSuppliers; i++ {
		volumeGoodsFromSuppliersOnTact[i] = make([]int, countTacts)
	}

	// initialization B
	for i, contentLine := range contentLines[6 : 6+countSuppliers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedVolumesStrings := strings.Split(trimmedLine, " ")
		for j, volumeString := range splittedVolumesStrings {
			volumeGoodsFromSuppliersOnTact[i][j], err = strconv.Atoi(volumeString)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// allocation C
	volumeUsedGoodsByCastomersOnTact := make([][]int, countCastomers)
	for i := 0; i < countCastomers; i++ {
		volumeUsedGoodsByCastomersOnTact[i] = make([]int, countTacts)
	}

	// initialization C
	startCastomersLine := 6 + countSuppliers + 1
	for i, contentLine := range contentLines[startCastomersLine : startCastomersLine+countCastomers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedVolumesStrings := strings.Split(trimmedLine, " ")
		for j, volumeString := range splittedVolumesStrings {
			volumeUsedGoodsByCastomersOnTact[i][j], err = strconv.Atoi(volumeString)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// allocation D
	startSetsSuppliersLine := startCastomersLine + countCastomers + 1
	setsSuppliersForCastomers := make([][]int, countCastomers)
	for i := 0; i < countCastomers; i++ {
		setsSuppliersForCastomers[i] = []int{}
	}

	// initialization D
	for i, contentLine := range contentLines[startSetsSuppliersLine : startSetsSuppliersLine+countCastomers] {
		trimmedLine := strings.Trim(contentLine, " ")
		splittedSuppliersStrings := strings.Split(trimmedLine, " ")
		for _, supplierString := range splittedSuppliersStrings {
			supplier, err := strconv.Atoi(supplierString)
			if err != nil {
				log.Fatal(err)
			}
			setsSuppliersForCastomers[i] = append(setsSuppliersForCastomers[i], supplier)
		}
	}

	return &resourceAllocationTask{
		countSuppliers:                   countSuppliers,
		countCastomers:                   countCastomers,
		countTacts:                       countTacts,
		volumeGoodsFromSuppliers:         volumeGoodsFromSuppliers,
		volumeGoodsFromSuppliersOnTact:   volumeGoodsFromSuppliersOnTact,
		volumeUsedGoodsByCastomersOnTact: volumeUsedGoodsByCastomersOnTact,
		setsSuppliersForCastomers:        setsSuppliersForCastomers,
	}

}

type Vertex struct {
	istock       bool
	stock        bool
	nextBranches []*Branch
	prevBranches []*Branch

	flowFromPreviousVertex int
}

func (vert *Vertex) Mark(flow int) {
	if flow == MaxInt {
		vert.flowFromPreviousVertex = flow - 1
	}
	vert.flowFromPreviousVertex = flow
}

func (vert *Vertex) Unmark() {
	vert.flowFromPreviousVertex = MaxInt
}

func (vert *Vertex) IsMarked() bool {
	return vert.flowFromPreviousVertex != MaxInt
}

func (vert *Vertex) IsStock() bool {
	return vert.stock
}

func MakeVertex(istock bool, stock bool, nextBranches []*Branch, prevBranches []*Branch) *Vertex {
	return &Vertex{
		istock:                 istock,
		stock:                  stock,
		nextBranches:           nextBranches,
		prevBranches:           prevBranches,
		flowFromPreviousVertex: MaxInt,
	}
}

type Branch struct {
	source      *Vertex
	destination *Vertex

	directBandwidth  int
	reverseBandwidth int
}

func (branch *Branch) Modificate(flow int, direct bool) {
	if direct {
		branch.directBandwidth -= flow
		branch.reverseBandwidth += flow
	} else {
		branch.directBandwidth += flow
		branch.reverseBandwidth -= flow
	}
}

func MakeBranch(source *Vertex, destination *Vertex, directBandwidth int) *Branch {
	return &Branch{
		source:           source,
		destination:      destination,
		directBandwidth:  directBandwidth,
		reverseBandwidth: 0,
	}
}

type Graph struct {
	vertexes [][]*Vertex
}

func connectVertexes(first *Vertex, second *Vertex, directBandwidth int) {
	// setup branch and it's vertexes
	branch := MakeBranch(first, second, directBandwidth)
	first.nextBranches = append(first.nextBranches, branch)
	second.prevBranches = append(second.prevBranches, branch)
}

func MakeGraphFromTask(task *resourceAllocationTask) *Graph {
	// istock level
	istock := MakeVertex(true, false, []*Branch{}, []*Branch{})

	// second level
	volumeGoodsFromSuppliersVertexes := make([]*Vertex, task.countSuppliers)

	// third level
	volumeGoodsFromSuppliersOnTactVertexes := make([]*Vertex, task.countSuppliers*task.countTacts)

	for i := 0; i < task.countSuppliers; i++ {
		vertex := MakeVertex(false, false, []*Branch{}, []*Branch{})
		connectVertexes(istock, vertex, task.volumeGoodsFromSuppliers[i])
		volumeGoodsFromSuppliersVertexes[i] = vertex

		for j := 0; j < task.countTacts; j++ {
			vertex := MakeVertex(false, false, []*Branch{}, []*Branch{})
			connectVertexes(volumeGoodsFromSuppliersVertexes[i], vertex, task.volumeGoodsFromSuppliersOnTact[i][j])
			volumeGoodsFromSuppliersOnTactVertexes[i*task.countTacts+j] = vertex
		}
	}

	// stock level
	stock := MakeVertex(false, true, []*Branch{}, []*Branch{})

	// fourth level
	volumeUsedGoodsByCastomersOnTactVertexes := make([]*Vertex, task.countCastomers*task.countTacts)

	for i := 0; i < task.countCastomers; i++ {
		for j := 0; j < task.countTacts; j++ {
			vertex := MakeVertex(false, false, []*Branch{}, []*Branch{})
			connectVertexes(vertex, stock, task.volumeUsedGoodsByCastomersOnTact[i][j])

			for _, supplier := range task.setsSuppliersForCastomers[i] {
				// customers are numbered from 1 in input data
				thirdLevelVertex := volumeGoodsFromSuppliersOnTactVertexes[(supplier-1)*task.countTacts+j]
				// this maybe special branch without bandwidth limitation (using MaxInt for now)
				connectVertexes(thirdLevelVertex, vertex, MaxInt)
			}

			volumeUsedGoodsByCastomersOnTactVertexes[i*task.countTacts+j] = vertex
		}
	}

	// connect third level with fourth level

	return &Graph{
		vertexes: [][]*Vertex{
			[]*Vertex{istock},
			volumeGoodsFromSuppliersVertexes,
			volumeGoodsFromSuppliersOnTactVertexes,
			volumeUsedGoodsByCastomersOnTactVertexes,
			[]*Vertex{stock},
		},
	}
}

func (task *resourceAllocationTask) Compute() int {
	graph := MakeGraphFromTask(task)
	// modificate graph
	maxFlow := fordFulkerson(graph)
	return maxFlow
}

func findPossibleTransactions(currentVertex *Vertex) []*Branch {
	possibleTransitions := []*Branch{}
	for _, nextBranch := range currentVertex.nextBranches {
		if nextBranch.directBandwidth != 0 && !nextBranch.destination.IsMarked() {
			possibleTransitions = append(possibleTransitions, nextBranch)
		}
	}
	for _, prevBranch := range currentVertex.prevBranches {
		if prevBranch.reverseBandwidth != 0 && !prevBranch.source.IsMarked() {
			possibleTransitions = append(possibleTransitions, prevBranch)
		}
	}
	return possibleTransitions
}

func computeFlow(wayVertexes []*Vertex, wayBranches []*Branch) int {
	minFlow := MaxInt

	// find min flow (exclude istock vertex) and unmark
	for _, vertex := range wayVertexes[1:] {
		if flow := vertex.flowFromPreviousVertex; flow < minFlow {
			minFlow = flow
		}
		vertex.Unmark()
	}

	for i, branch := range wayBranches {
		switch branch.source {
		case wayVertexes[i]:
			// direct branch
			branch.Modificate(minFlow, true)
		default:
			// reverse branch
			branch.Modificate(minFlow, false)
		}

	}

	return minFlow
}

func fordFulkerson(graph *Graph) int {
	maxFlow := 0

	// start from istock
	istockVertex := graph.vertexes[0][0]
	// MaxInt is used as a flag for check mark, so MaxInt - 1 (this is an unreachable value)
	istockVertex.Mark(MaxInt - 1)

	for {
		vertexesForUnmark := []*Vertex{}
		// initialize with istock vertex
		wayVertexes := []*Vertex{istockVertex}
		wayBranches := []*Branch{}

		currentVertex := istockVertex

		for !currentVertex.IsStock() {
			possibleTransitions := findPossibleTransactions(currentVertex)
			if len(possibleTransitions) == 0 {
				if len(wayVertexes) == 1 {
					goto maxFlowFound
				}
				vertexesForUnmark = append(vertexesForUnmark, currentVertex)
				wayVertexes = wayVertexes[:len(wayVertexes)-1]
				wayBranches = wayBranches[:len(wayBranches)-1]
				currentVertex = wayVertexes[len(wayVertexes)-1]
				continue
			}

			maxBandwidth := 0
			chosenBranch := possibleTransitions[0]
			chosenVertex := chosenBranch.destination
			for _, possibleTransition := range possibleTransitions {
				switch possibleTransition.source {
				case currentVertex:
					// direct branch
					if possibleTransition.directBandwidth > maxBandwidth {
						maxBandwidth = possibleTransition.directBandwidth
						chosenBranch = possibleTransition
						chosenVertex = chosenBranch.destination
					}
				default:
					// reversed branch
					if possibleTransition.reverseBandwidth > maxBandwidth {
						maxBandwidth = possibleTransition.reverseBandwidth
						chosenBranch = possibleTransition
						chosenVertex = chosenBranch.source
					}
				}
			}
			wayBranches = append(wayBranches, chosenBranch)
			wayVertexes = append(wayVertexes, chosenVertex)

			currentVertex = chosenVertex
			currentVertex.Mark(maxBandwidth)
		}

		// also do unmark vertexes and modificate bandwidth branches
		maxFlow += computeFlow(wayVertexes, wayBranches)

		for _, vertex := range vertexesForUnmark {
			vertex.Unmark()
		}

	}

maxFlowFound:

	return maxFlow
}
