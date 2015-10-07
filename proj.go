package conllx

import (
	"container/list"
	"sort"
	"strings"
)

// A projectivizer can transform a dependency structure in one that has
// no non-projective edges. Moreover, it can reverse this transformation.
type Projectivizer interface {
	Deprojectivize([]Token) []Token
	Projectivize([]Token) []Token
}

type edge struct {
	from  int
	to    int
	label string
}

type edgeByLength []edge

func (e edgeByLength) Len() int { return len(e) }

func (e edgeByLength) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e edgeByLength) Less(i, j int) bool {
	return max(e[i].from, e[i].to)-min(e[i].from, e[i].to) < max(e[j].from, e[j].to)-min(e[j].from, e[j].to)
}

// Return tokens with lifted edges by ascending depth.
func depthSortedLifted(graph graph, lifted []bool) []int {
	seen := make([]bool, graph.NVertices())
	var liftedSorted []int

	// Initial queue only has the root vertex.
	q := list.New()
	q.PushBack(0)

	for q.Len() != 0 {
		first := q.Front()
		q.Remove(first)

		vertex := first.Value.(int)

		for _, edge := range graph.Edges(vertex) {
			// Only process the vertex if it was not seen before.
			if !seen[edge.To] {
				seen[edge.To] = true
				q.PushBack(edge.To)
				if lifted[edge.To] {
					liftedSorted = append(liftedSorted, edge.To)
				}
			}
		}
	}

	return liftedSorted
}

// Return non-projective arcs, ordered by length.
func nonProjectiveArcs(g graph) []edge {
	var nonProjective []edge

	for i := 0; i < g.NVertices(); i++ {
		paths := bfs(g, i)
		edges := g.Edges(i)

		for _, k := range edges {
			// An edge i -> k is projective, iff:
			//
			// i > j > k or i < j < k, and i ->* j
			for j := min(i, k.To) + 1; j < max(i, k.To); j++ {
				if _, ok := paths[j]; !ok {
					nonProjective = append(nonProjective, edge{i, k.To, k.Label})
					break
				}
			}
		}
	}

	sort.Stable(edgeByLength(nonProjective))

	return nonProjective
}

// Remove head markers from the dependency graph and return slices
// indicating lifted nodes and the corresponding preferred head labels.
func prepareProjectiveGraph(graph graph) ([]bool, []string, bool) {
	lifted := make([]bool, graph.NVertices())
	preferredHeadLabel := make([]string, graph.NVertices())
	hasLiftedEdges := false

	for i := 0; i < graph.NVertices(); i++ {
		for _, dep := range graph.Edges(i) {
			sepIdx := strings.IndexRune(dep.Label, '|')

			if sepIdx == -1 {
				continue
			}

			hasLiftedEdges = true
			lifted[dep.To] = true
			preferredHeadLabel[dep.To] = dep.Label[sepIdx+1:]

			graph.Remove(i, dep)

			dep.Label = dep.Label[:sepIdx]
			graph.AddEdge(i, dep)
		}
	}

	return lifted, preferredHeadLabel, hasLiftedEdges
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
