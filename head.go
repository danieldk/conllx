package conllx

import (
	"fmt"
	"math"
)

var _ Projectivizer = HeadProjectivizer{}

// A projectivizer using the 'head' marking strategy (Nivre and Nillson,
// 2005).
type HeadProjectivizer struct{}

func (h HeadProjectivizer) Projectivize(sentence []Token) []Token {
	graph := sentenceToDepGraph(sentence)
	lifted := make([]bool, graph.NVertices())

	for {
		npArcs := nonProjectiveArcs(graph)
		if len(npArcs) == 0 {
			break
		}

		h.lift(graph, lifted, npArcs[0])

	}

	projective := make([]Token, len(sentence))
	copy(projective, sentence)
	depGraphToSentence(graph, projective)

	return projective
}

func (p HeadProjectivizer) lift(graph graph, lifted []bool, edge edge) {
	headIdx := edge.from
	depIdx := edge.to
	depRel := edge.label
	newHead := graph.Parent(headIdx)

	graph.Remove(headIdx, edgeTo{
		To:    depIdx,
		Label: depRel,
	})

	if !lifted[depIdx] {
		graph.AddEdge(newHead.To, edgeTo{
			To:    depIdx,
			Label: fmt.Sprintf("%s|%s", depRel, newHead.Label),
		})
	} else {
		graph.AddEdge(newHead.To, edgeTo{
			To:    depIdx,
			Label: depRel,
		})
	}

	lifted[depIdx] = true
}

func (p HeadProjectivizer) Deprojectivize(sentence []Token) []Token {
	graph := sentenceToDepGraph(sentence)
	lifted, headLabels, hasLifted := prepareProjectiveGraph(graph)
	sortedLifted := depthSortedLifted(graph, lifted)

	if !hasLifted {
		return sentence
	}

	for {
		prevLen := len(sortedLifted)
		sortedLifted = p.deprojectivizeNext(graph, sortedLifted, headLabels)
		if len(sortedLifted) == prevLen {
			break
		}
	}

	pSentence := make([]Token, len(sentence))
	copy(pSentence, sentence)
	depGraphToSentence(graph, pSentence)

	return pSentence
}

func (p HeadProjectivizer) deprojectivizeNext(graph graph, sortedLifted []int, headLabels []string) []int {
	for idx, liftedVertex := range sortedLifted {
		preferredHeadRel := headLabels[liftedVertex]
		curHead := graph.Parent(liftedVertex)

		if newHead, found := p.searchCandidate(graph, curHead.To, liftedVertex, preferredHeadRel); found {
			deprel := curHead.Label

			graph.Remove(curHead.To, edgeTo{
				To:    liftedVertex,
				Label: deprel,
			})

			graph.AddEdge(newHead, edgeTo{
				To:    liftedVertex,
				Label: deprel,
			})

			return append(sortedLifted[:idx], sortedLifted[idx+1:]...)
		}
	}

	return sortedLifted
}

func (HeadProjectivizer) searchCandidate(graph graph, start int, skip int, headRel string) (int, bool) {
	var depth []int
	for _, dep := range graph.Edges(start) {
		depth = append(depth, dep.To)
	}

	// Search breadth first, by depth
	for len(depth) != 0 {
		candidates := make([]int, 0)
		nextDepth := make([]int, 0)
		for _, vertex := range depth {
			// Attaching to self or a dependent leads to cycles.
			if vertex == skip {
				continue
			}

			if graph.Parent(vertex).Label == headRel {
				candidates = append(candidates, vertex)
			}

			for _, dep := range graph.Edges(vertex) {
				nextDepth = append(nextDepth, dep.To)
			}
		}

		if len(candidates) == 0 {
			depth = nextDepth
			continue
		}

		smallestDist := int(math.MaxInt32)
		smallestCand := candidates[0]
		for _, candidate := range candidates {
			dist := max(candidate, skip) - min(candidate, skip)
			if dist < smallestDist {
				smallestDist = dist
				smallestCand = candidate
			}
		}

		return smallestCand, true
	}

	return 0, false
}
