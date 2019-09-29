// Branch and Bounds algorithm
package exact_algorithms

import (
	"fmt"
	"csp/task"
)

func BnB(taskInst task.TaskBnB) (task.Vertex, int) {
	vertexes := taskInst.CreateInitVertexSet()

	iterator := vertexes.Iterator()
	iterator.First()
	bestVertex := iterator.Value().(task.Vertex)

	for {
		//stop condition
		if vertexes.Size() == 0 {
			break
		} else if vertexes.Size() == 1 {
			iterator.First()
			// should be last vertex
			lastVertex := iterator.Value().(task.Vertex)
			if taskInst.UpperBound(lastVertex, false) == taskInst.LowerBound(lastVertex) {
				bestVertex = lastVertex
				break
			}
		}

		//Branching
		iterator.First()
		vertexInst := iterator.Value().(task.Vertex)
		vertexes.Remove(vertexInst)

		bestChildVertex, childVertexes := vertexInst.NextVertexes(taskInst)
		for _, childVertex := range childVertexes {
			vertexes.Add(childVertex)
		}
		if taskInst.UpperBound(bestChildVertex, false) <= taskInst.UpperBound(bestVertex, false) {
			bestVertex = bestChildVertex
		}

		//Clipping - use best Vertex
		//Found stage 
		vertexesForClipping := make([]task.Vertex, 0)
		for iterator := vertexes.Iterator(); iterator.Next(); {
			vertexForClipping := iterator.Value().(task.Vertex)
			if taskInst.UpperBound(bestVertex, false) <= taskInst.LowerBound(vertexForClipping) {
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
	taskInst.UpperBound(bestVertex, true)
	return bestVertex, taskInst.CountTraversedVertexes()
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
