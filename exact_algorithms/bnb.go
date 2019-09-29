// Branch and Bounds algorithm
package exact_algorithms

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)


type Vertex interface {
	UpperBound() int
	LowerBound() int
	NextVertexes() (Vertex, []Vertex)
}

type vertex struct {
	permutation []int
	upperBound int
	lowerBound int
}

func MakeVertex(permutation []int) Vertex {
	return &vertex{
		permutation: permutation,
		upperBound: 0,
		lowerBound: 0,
	}
}

func (vert *vertex) UpperBound() int {
	panic("not implemented")
}

func (vert *vertex) LowerBound() int {
	panic("not implemented")
}

func (vert *vertex) NextVertexes() (Vertex, []Vertex) {
	panic("not implemented ")
}


type TaskBnB interface {
	CreateInitVertexSet() *treeset.Set

}

type task struct {
	//include zero vertex
	countVertex int
	deliveryTimes [][]int

	//except zero vertex
	directiveTimes []int
}

func (taskInst *task) CreateInitVertexSet() *treeset.Set {
	// Strategy of choose vertex - create order (should be param)
	comparator := func(a, b interface{}) int {
		aUB := a.(Vertex).UpperBound()
		bUB := b.(Vertex).UpperBound()
		if aUB < bUB {
			return -1
		} else if aUB == bUB {
			return 0
		} else {
			return 1
		}
	}

	vertex := MakeVertex([]int{0})
	return treeset.NewWith(comparator)
}

func (taskInst *task) Branching(set *treeset.Set) *treeset.Set {
	panic("not implemented")
}

func MakeTask() TaskBnB {
	var deliveryTimes [][]int
	var directiveTimes []int

	deliveryTimes = make([][]int, 5)
	deliveryTimes[0] = []int{0,3,4,2,5}
	deliveryTimes[1] = []int{3,0,6,3,9}
	deliveryTimes[2] = []int{4,6,0,2,1}
	deliveryTimes[3] = []int{2,3,2,0,5}
	deliveryTimes[4] = []int{5,9,1,5,0}

	directiveTimes = []int{5,6,7,8,9}
	return &task{
		countVertex: len(deliveryTimes),
		deliveryTimes: deliveryTimes,
		directiveTimes: directiveTimes, 
	}
}


func BnB(task TaskBnB) int {
	vertexes := task.CreateInitVertexSet()
	bestVertex := vertexes.Values()[0].(Vertex)

	for {
		//stop condition
		if vertexes.Size() == 1 {
			vertexes := vertexes.Values()
			// should be last vertex
			lastVertex := vertexes[0].(Vertex)
			if lastVertex.UpperBound() == lastVertex.LowerBound() {
				bestVertex = lastVertex
				break
			}
		}

		//Branching
		iterator := vertexes.Iterator()
		vertexInst := iterator.Value().(Vertex)
		vertexes.Remove(vertexInst)

		bestChildVertex, childVertexes := vertexInst.NextVertexes()
		for _, childVertex := range childVertexes {
			vertexes.Add(childVertex)
		}
		if bestChildVertex.UpperBound() < bestVertex.UpperBound() {
			bestVertex = bestChildVertex
		}

		//Clipping - use best Vertex
		//Found stage 
		vertexesForClipping := make([]Vertex, 0)
		for iterator := vertexes.Iterator(); iterator.Next(); {
			vertexForClipping := iterator.Value().(Vertex)
			if bestVertex.UpperBound() <= vertexForClipping.LowerBound() {
				vertexesForClipping = append(vertexesForClipping, vertexForClipping)
			} else {
				break
			}
		}
		//Remove stage
		for _, vertexForClipping := range vertexesForClipping {
			vertexes.Remove(vertexForClipping)
		}
	}

	fmt.Println(bestVertex)
}

/*set := treeset.NewWithIntComparator() // empty (keys are of type int)
set.Add(1)                            // 1
set.Add(2, 2, 3, 4, 5)                // 1, 2, 3, 4, 5 (in order, duplicates ignored)
set.Remove(4)                         // 1, 2, 3, 5 (in order)
set.Remove(2, 3)                      // 1, 5 (in order)
set.Contains(1)                       // true
set.Contains(1, 5)                    // true
set.Contains(1, 6)                    // false
_ = set.Values()                      // []int{1,5} (in order)
set.Clear()                           // empty
set.Empty()                           // true
set.Size()                            // 0*/
