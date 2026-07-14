package main

import (
	"container/heap"
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


		for _, neighbour := range node.neighbours {
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
		// TODO: return correct results from here
		return []*Node{}
	}

	result := make(PriorityQueue, 0, k)
	candidates := make(PriorityQueue, 0, efsearch)
	heap.Push(&candidates, &QueuedNode{node: nsw.nodes[0], score: vector.Distance(nsw.nodes[0].vector)})

	visitedNodes := make(map[string]struct{})
	visitedNodes[candidates[0].node.id] = struct{}{}

	for candidates.Len() > 0 {
		queuedNode := heap.Pop(&candidates).(*QueuedNode)
		if result.Len() == k && -result[0].score < queuedNode.score {
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

		for _, neighbour := range queuedNode.node.neighbours {
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

func (nsw *NavigableSmallWorld) Insert(vector *Vector, m, efconstruct int) {
	toInsert := &Node{id: vector.id, vector: vector}
	for _, node := range nsw.Search(vector, m, efconstruct) {
		node.neighbours = append(node.neighbours, toInsert)
		toInsert.neighbours = append(toInsert.neighbours, node)
	}

	// TODO: add pruning logic
	nsw.nodes = append(nsw.nodes, toInsert)
}
