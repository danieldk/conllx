package conllx

import (
	"reflect"
	"testing"
)

var testGraph graph

var pathsFrom0 = paths{
	1: 0,
	2: 1,
	3: 0,
}

var pathsFrom1 = paths{
	2: 1,
	3: 2,
}

func init() {
	testGraph = newGraph(4)
	testGraph.AddEdge(0, edgeTo{1, ""})
	testGraph.AddEdge(1, edgeTo{2, ""})
	testGraph.AddEdge(2, edgeTo{3, ""})
	testGraph.AddEdge(0, edgeTo{3, ""})
}

func TestBfs(t *testing.T) {
	testEquals(t, pathsFrom0, bfs(testGraph, 0))
	testEquals(t, pathsFrom1, bfs(testGraph, 1))
}

func testEquals(t *testing.T, check, candidate interface{}) {
	if !reflect.DeepEqual(check, candidate) {
		t.Errorf("Expected:\n%v\nGot:\n%v", check, candidate)
	}
}
