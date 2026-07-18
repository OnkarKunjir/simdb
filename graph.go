package main

import "fmt"

// Holds information required to construct [HierarchicalNavigableSmallWorld] and [NavigableSmallWorld].
// Each node represents and vector which is inserted into the graph.
type Node struct {
	id         string
	vector     *Vector
	neighbours []map[string]*Node
}

func (node *Node) String() string {
	return fmt.Sprintf("Node{id: %s}", node.id)
}

// Helper struct to store node with score, to be used in [PriorityQueue] 
type QueuedNode struct {
	node  *Node
	score float64
}

func (scoredNode *QueuedNode) String() string {
	return fmt.Sprintf("ScoredNode{id: %s, score: %0.2f}", scoredNode.node.id, scoredNode.score)
}

// Priority queue of [QueuedNode] built using min heap (smallest score is at 0th index)
type PriorityQueue []*QueuedNode

func (queue PriorityQueue) Len() int { return len(queue) }

func (queue PriorityQueue) Less(i, j int) bool {
	return queue[i].score < queue[j].score
}

func (queue PriorityQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue *PriorityQueue) Push(x any) {
	*queue = append(*queue, x.(*QueuedNode))
}

func (queue *PriorityQueue) Pop() any {
	old := *queue
	n := len(old)
	item := old[n-1]
	old[n-1] = nil

	*queue = old[0 : n-1]
	return item
}
