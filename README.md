# SimDB

A vector similarity search engine implemented from scratch in Go, built around the **HNSW (Hierarchical Navigable Small World)** algorithm — the same algorithm powering production vector databases like Weaviate, Pinecone, and pgvector.

---

## Why I built this

I've spent the last few years building production RAG pipelines and vector search systems — using pgvector, evaluating Elasticsearch, and shipping semantic search over 100K+ metadata attributes at enterprise scale. I understood what these systems *did*, but I wanted to understand what they were *built on*.

This project is that exercise. No libraries. No shortcuts. Just the algorithm, the data structures, and Go.

---

## What is HNSW?

Most vector databases need to answer one question fast: *"given this vector, what are the K most similar vectors in my dataset?"*

The naive approach — compute cosine similarity against every vector — is O(n). Fine for 1,000 vectors. Unusable for 1,000,000.

HNSW solves this with a **hierarchical graph structure**:

- Vectors are nodes in a graph. Each node connects to its nearest neighbours.
- The graph has multiple layers. Upper layers are sparse (fast, coarse navigation). Lower layers are dense (slow, precise search).
- On **insert**, a node is assigned a random layer level using exponential decay — most nodes land at layer 0, fewer at layer 1, fewer still at layer 2, and so on.
- On **search**, we enter at the top layer, greedily navigate toward the query vector, then descend layer by layer, getting more precise as we go.

This gives approximate nearest neighbour search in **O(log n)** — orders of magnitude faster than brute force, with recall that's good enough for production use cases.

---

## Implementation

### NSW (`nsw.go`)

The foundation. A single-layer Navigable Small World graph with:

- **Beam search** using a dual priority queue — a min-heap for candidates to explore, a max-heap for current best results
- **Early stopping** — search terminates when the closest unexplored candidate is farther than the worst current result
- **Diversity-based neighbour pruning (RNG heuristic)** — when a node exceeds its max connections, we don't just drop the farthest neighbour. We keep neighbours that "cover different directions" in vector space, preserving graph navigability

### HNSW (`hnsw.go`)

The full hierarchical index:
- Probabilistic layer assignment on insert using exponential decay (`-ln(rand) * mL`)
- Layer-by-layer greedy search during insert to find entry points for each layer
- NSW-based search within each layer
- Final precise search at layer 0

