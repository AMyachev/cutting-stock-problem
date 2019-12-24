package knapsack

import (
	"testing"
)

var task = MakeKnapsackProblemFromFile("C:\\Users\\Толик\\Desktop\\спец.семинар\\спецсеминар1\\cutting-stock-problem\\example-tasks\\knapsack\\task_3_10_n1000.txt")

func BenchmarkTable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task.TableSolutionByDefault(task.GetCompanyPerfomance())
	}
}
func BenchmarkRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		task.RecursiveSolutionDefaultOrder(task.GetCompanyPerfomance(), true)
	}
}
