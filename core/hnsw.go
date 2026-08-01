package core

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strings"
)

type HierarchicalNavigableSmallWorld struct {
	nodes      map[string]*Node
	entryPoint *Node
}

// Struct to represent vectors searched in the graph
type SearchResult struct {
	Vector *Vector
	Score  float64
}

func (searchResult *SearchResult) String() string {
	return fmt.Sprintf("Document {id: %s, score: %.02f}", searchResult.Vector.Id, searchResult.Score)
}

// Helper struct to sort nodes by search score, it's equivalent to [QueuedNode] right now but usecase is bit different.
type ScoredNodes struct {
	node  *Node
	score float64
}

func (hnsw *HierarchicalNavigableSmallWorld) String() string {
	var builder strings.Builder

	for level := range slices.Backward(hnsw.entryPoint.neighbours) {
		fmt.Fprintf(&builder, "------- Level: %d -------\n", level)

		for _, node := range hnsw.nodes {
			if level >= len(node.neighbours) {
				continue
			}

			builder.WriteString(node.id)
			builder.WriteString(" -> [ ")
			for id := range node.neighbours[level] {
				builder.WriteString(id)
				builder.WriteString(" ")
			}

			builder.WriteString("]\n")
		}
	}
	return builder.String()
}

// searchLevel searches given vector on specified level with provided entry point.
// a beam search similar to [NavigableSmallWorld] is performed at this level at max efsearch nodes are explored.
// nearest k nodes are returned as result
func (hnsw *HierarchicalNavigableSmallWorld) searchLevel(vector *Vector, entryPoint *Node, level, k, efsearch int) []*SearchResult {
	// collect nearest efsearch nodes, used as max heap (fartest node at 0th index)
	result := make(PriorityQueue, 0, efsearch)

	// collects list of nodes to be explored, used as min heap (nearest node at 0th index)
	candidates := make(PriorityQueue, 0, efsearch)
	heap.Push(&candidates, &QueuedNode{entryPoint, vector.Distance(entryPoint.vector)})

	visitedNodes := make(map[string]struct{})
	visitedNodes[candidates[0].node.id] = struct{}{}

	for candidates.Len() > 0 {
		queuedNode := heap.Pop(&candidates).(*QueuedNode)

		// If node being explored is farther than known visited farthest node then we can stop
		if result.Len() == efsearch && -result[0].score < queuedNode.score {
			break
		}

		if result.Len() < efsearch {
			// maybe I can create new queued node here instead of modifying same node but dosent't matter that much
			queuedNode.score = -queuedNode.score
			heap.Push(&result, queuedNode)
		} else if -result[0].score > queuedNode.score {
			queuedNode.score = -queuedNode.score
			heap.Pop(&result)
			heap.Push(&result, queuedNode)
		}

		for _, neighbour := range queuedNode.node.neighbours[level] {
			if _, ok := visitedNodes[neighbour.id]; ok {
				continue
			}

			distance := vector.Distance(neighbour.vector)
			visitedNodes[neighbour.id] = struct{}{}
			heap.Push(&candidates, &QueuedNode{neighbour, distance})
		}
	}

	searchResults := make([]*SearchResult, len(result))
	for index := len(result) - 1; index >= 0; index-- {
		queuedNode := heap.Pop(&result).(*QueuedNode)
		searchResults[index] = &SearchResult{queuedNode.node.vector, -queuedNode.score}
	}

	return searchResults[:min(k, len(searchResults))]
}

// Search given vector in graph and return nearest k vectors as result, while searching at max efsearch nodes are explored at
// each level.
func (hnsw *HierarchicalNavigableSmallWorld) Search(vector *Vector, k, efsearch int) []*SearchResult {
	if hnsw.entryPoint == nil {
		return nil
	}
	entryPoint := hnsw.entryPoint
	var searchResults []*SearchResult
	for currentLevel := range slices.Backward(entryPoint.neighbours) {
		if currentLevel != 0 {
			// just need to pinpoint to correct subgraph fast.
			searchResults = hnsw.searchLevel(vector, entryPoint, currentLevel, 1, 1)
			entryPoint = hnsw.nodes[searchResults[0].Vector.Id]
		} else {
			searchResults = hnsw.searchLevel(vector, entryPoint, currentLevel, k, efsearch)
		}
	}

	return searchResults
}

// prunes neighbours of node form specified level keeping at most m neighbours.
// farthest neighbours are pruned, also to keep good set of neighbours, some candidate neighour which is closer to already
// considered neighbour is also pruned.
func pruneNeighboursFromLevel(node *Node, level, m int) {
	neighbours := make([]*ScoredNodes, 0, len(node.neighbours[level]))
	for _, neighbour := range node.neighbours[level] {
		neighbours = append(neighbours, &ScoredNodes{neighbour, node.vector.Distance(neighbour.vector)})
	}
	sort.Slice(neighbours, func(i, j int) bool {
		return neighbours[i].score < neighbours[j].score
	})

	selected := make([]*ScoredNodes, 0, m)
	selected = append(selected, neighbours[0])
	toPrune := make([]string, 0, m)

	for _, candidate := range neighbours[1:] {
		currentDistance := node.vector.Distance(candidate.node.vector)
		// if lenght of selected nodes is > m then skip everything it's too far but also skip of candidate is closer to any of the selected nodes
		skipNode := len(selected) >= m || slices.ContainsFunc(selected, func(selected *ScoredNodes) bool {
			return candidate.node.vector.Distance(selected.node.vector) < currentDistance
		})

		if skipNode {
			toPrune = append(toPrune, candidate.node.id)
		} else {
			selected = append(selected, candidate)
		}
	}

	for _, nodeId := range toPrune {
		delete(node.neighbours[level][nodeId].neighbours[level], node.id)
		delete(node.neighbours[level], nodeId)
	}
}

// Inserts node at specified level with given entryPoint, node is connected to nearest m other nodes while inserting
func (hnsw *HierarchicalNavigableSmallWorld) insertAtLevel(toInsert *Node, entryPoint *Node, level, m, efconstruct int) *Node {
	documents := hnsw.searchLevel(toInsert.vector, entryPoint, level, m, efconstruct)
	for _, document := range documents {
		node := hnsw.nodes[document.Vector.Id]

		node.neighbours[level][toInsert.id] = toInsert
		toInsert.neighbours[level][node.id] = node
		if len(node.neighbours[level]) <= m {
			continue
		}
		pruneNeighboursFromLevel(node, level, m)
	}
	return hnsw.nodes[documents[0].Vector.Id]
}

// Insert vector into HNSW graph, where each node can have maximum m nodes connected to it at any level, except level-0 where
// each node can have 2*m connected nodes. While inserting node efconstruct nodes are explored.
func (hnsw *HierarchicalNavigableSmallWorld) Insert(vector *Vector, m, efconstruct int) {
	toInsert := &Node{id: vector.Id, vector: vector}
	// if there are no nodes in the graph, no point in calculating levels. just add it to level 0 and call it a day.
	if hnsw.entryPoint == nil {
		hnsw.nodes = make(map[string]*Node)

		toInsert.neighbours = []map[string]*Node{make(map[string]*Node)}
		hnsw.nodes[toInsert.id] = toInsert
		hnsw.entryPoint = toInsert
		return
	}

	// I don't understand exact math here, I just copied it from paper, but overall probability of higher levels is less compared to lower levels
	mL := 1 / math.Log(float64(m))
	level := int(-math.Log(rand.Float64()) * mL)
	for range level + 1 {
		toInsert.neighbours = append(toInsert.neighbours, make(map[string]*Node))
	}

	// Now bit confusing part, if level of vector to be inserted is > hnsw current level then I don't need to perform
	// any action. cause the higher level just going to have single node (newly added one) with no neighbours.
	// Following loop just executes to find entry point for level from where I actully have to insert.
	entryPoint := hnsw.entryPoint
	currentLevel := len(hnsw.entryPoint.neighbours) - 1
	for ; currentLevel > level; currentLevel-- {
		entryPoint = hnsw.nodes[hnsw.searchLevel(vector, entryPoint, currentLevel, 1, 1)[0].Vector.Id]
	}

	for ; currentLevel >= 0; currentLevel-- {
		if currentLevel != 0 {
			entryPoint = hnsw.insertAtLevel(toInsert, entryPoint, currentLevel, m, efconstruct)
		} else {
			entryPoint = hnsw.insertAtLevel(toInsert, entryPoint, currentLevel, 2*m, efconstruct)
		}
	}

	hnsw.nodes[toInsert.id] = toInsert
	// update the entry point
	if level > (len(hnsw.entryPoint.neighbours) - 1) {
		hnsw.entryPoint = toInsert
	}
}
