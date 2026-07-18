package main

import (
	"cmp"
	"container/heap"
	"maps"
	"slices"
	"strings"
)

type NavigableSmallWorld struct {
	nodes []*Node
}

func (nsw *NavigableSmallWorld) String() string {
	var builder strings.Builder
	for _, node := range nsw.nodes {
		builder.WriteString(node.id)
		builder.WriteString(" -> [ ")

		for _, neighbour := range node.neighbours[0] {
			builder.WriteString(neighbour.id)
			builder.WriteString(" ")
		}
		builder.WriteString("]\n")
	}
	return builder.String()
}

// Search k nearest nodes from graph where at max efsearch nodes are explored
func (nsw *NavigableSmallWorld) Search(vector *Vector, k, efsearch int) []*Node {
	if len(nsw.nodes) == 0 {
		return []*Node{}
	}

	result := make(PriorityQueue, 0, k)            // max distnace on top
	candidates := make(PriorityQueue, 0, efsearch) // min distnace on top
	heap.Push(&candidates, &QueuedNode{node: nsw.nodes[0], score: vector.Distance(nsw.nodes[0].vector)})
	visitedNodes := make(map[string]struct{})
	visitedNodes[candidates[0].node.id] = struct{}{}

	for candidates.Len() > 0 {
		queuedNode := heap.Pop(&candidates).(*QueuedNode)
		if result.Len() == k && -result[0].score < queuedNode.score {
			// results are full and current shortest distance is grateer than recorded longest distance
			// no need to search any further
			break
		}

		visitedNodes[queuedNode.node.id] = struct{}{}
		if result.Len() < k {
			queuedNode.score = -queuedNode.score
			heap.Push(&result, queuedNode)
		} else if -result[0].score > queuedNode.score {
			queuedNode.score = -queuedNode.score

			heap.Pop(&result)
			heap.Push(&result, queuedNode)
		}

		for _, neighbour := range queuedNode.node.neighbours[0] {
			if _, ok := visitedNodes[neighbour.id]; ok {
				// already visited node, skip it.
				continue
			}

			score := vector.Distance(neighbour.vector)
			if candidates.Len() < efsearch {
				heap.Push(&candidates, &QueuedNode{neighbour, score})
			} else if candidates[0].score > score {
				heap.Pop(&candidates)
				heap.Push(&candidates, &QueuedNode{neighbour, score})
			}
		}

	}

	searchedNodes := make([]*Node, result.Len())
	index := result.Len() - 1
	for result.Len() > 0 {
		searchedNodes[index] = heap.Pop(&result).(*QueuedNode).node
		index -= 1
	}

	return searchedNodes
}

func pruneNeighbours(node *Node, m int) {
	neighbours := slices.SortedFunc(maps.Values(node.neighbours[0]), func(a, b *Node) int {
		return cmp.Compare(node.vector.Distance(a.vector), node.vector.Distance(b.vector))
	})

	selected := make([]*Node, 0, m)
	selected = append(selected, neighbours[0])
	toPrune := make([]string, 0, m)

	for _, candidate := range neighbours[1:] {
		currentDistance := node.vector.Distance(candidate.vector)

		// if lenght of selected nodes is > m then skip everything it's too far but also skip it node is closer to any of the selected nodes
		skipNode := len(selected) >= m || slices.ContainsFunc(selected, func(selected *Node) bool {
			return candidate.vector.Distance(selected.vector) < currentDistance
		})

		if skipNode {
			toPrune = append(toPrune, candidate.id)
		} else {
			selected = append(selected, candidate)
		}
	}

	for _, nodeId := range toPrune {
		delete(node.neighbours[0][nodeId].neighbours[0], node.id)
		delete(node.neighbours[0], nodeId)
	}
}

// Prunes neighbours of
func (nsw *NavigableSmallWorld) Insert(vector *Vector, m, efconstruct int) {
	neighbours := []map[string]*Node{make(map[string]*Node)}
	toInsert := &Node{id: vector.Id, vector: vector, neighbours: neighbours}
	for _, node := range nsw.Search(vector, m, efconstruct) {
		node.neighbours[0][toInsert.id] = toInsert
		toInsert.neighbours[0][node.id] = node
		if len(node.neighbours) < m {
			// no pruning needed
			continue
		}

		pruneNeighbours(node, m)
	}

	nsw.nodes = append(nsw.nodes, toInsert)
}
