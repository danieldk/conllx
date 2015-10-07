package conllx

import (
	"container/list"
	"sort"
)

type paths map[int]int

type edgeTo struct {
	To    int
	Label string
}

type graph struct {
	edges        [][]edgeTo
	reverseEdges []edgeTo
}

func newGraph(nVertices int) graph {
	g := graph{
		edges:        make([][]edgeTo, nVertices),
		reverseEdges: make([]edgeTo, nVertices),
	}

	for i := 0; i < nVertices; i++ {
		g.edges[i] = make([]edgeTo, 0, 3)
	}

	return g
}

func (g graph) AddEdge(token int, e edgeTo) {
	edges := g.edges[token]
	idx := sort.Search(len(edges), func(i int) bool { return edges[i].To > e.To })

	if idx < len(edges) && edges[idx] == e {
		return
	}

	g.edges[token] = insert(edges, idx, e)
	g.reverseEdges[e.To] = edgeTo{
		To:    token,
		Label: e.Label,
	}
}

// Return edges going out of the given vertex.
func (g graph) Edges(vertex int) []edgeTo {
	return g.edges[vertex]
}

func (g graph) Parent(vertex int) edgeTo {
	return g.reverseEdges[vertex]
}

func (g graph) NVertices() int {
	return len(g.edges)
}

func (g graph) Remove(token int, e edgeTo) {
	edges := g.edges[token]
	idx := sort.Search(len(edges), func(i int) bool { return edges[i].To >= e.To })
	if idx < len(edges) && edges[idx] == e {
		g.edges[token] = append(edges[:idx], edges[idx+1:]...)
		g.reverseEdges[e.To] = edgeTo{}
	}
}

func sentenceToDepGraph(sentence []Token) graph {
	depGraph := newGraph(len(sentence) + 1)

	for idx, token := range sentence {
		if head, ok := token.Head(); ok {
			if headRel, ok := token.HeadRel(); ok {
				dependent := edgeTo{
					To:    idx + 1,
					Label: headRel,
				}
				depGraph.AddEdge(int(head), dependent)
			}
		}
	}

	return depGraph
}

// Modifies the sentence in-place to add the dependency relations in the graph
func depGraphToSentence(g graph, sentence []Token) {
	for from := 0; from < g.NVertices(); from++ {
		for _, edge := range g.Edges(from) {
			sentence[edge.To-1].SetHeadRel(edge.Label)
			sentence[edge.To-1].SetHead(uint(from))
		}
	}
}

// Do a breadth-first search of the dependency graph, starting at
// 'start'.
func bfs(g graph, start int) paths {
	reachable := make(paths)

	// Initial queue only has the start vertex.
	q := list.New()
	q.PushBack(start)

	for q.Len() != 0 {
		first := q.Front()
		q.Remove(first)

		vertex := first.Value.(int)

		for _, edge := range g.Edges(vertex) {
			// Only process the vertex if it was not seen before.
			if _, ok := reachable[edge.To]; !ok {
				reachable[edge.To] = vertex
				q.PushBack(edge.To)
			}
		}
	}

	return reachable
}

func insert(slice []edgeTo, index int, value edgeTo) []edgeTo {
	slice = append(slice, edgeTo{})
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}
