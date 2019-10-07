// Branch and Bounds algorithm
package exact_algorithms

import (
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

	taskInst.UpperBound(bestVertex, true)
	return bestVertex, taskInst.CountTraversedVertexes()
}
