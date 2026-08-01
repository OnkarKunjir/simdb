# SimDB

A lightweight, high-performance vector search engine written in Go.

SimDB indexes dense vector embeddings using **Hierarchical Navigable Small World (HNSW)** graphs and exposes a simple REST API with an accompanying Python client. It is designed for semantic search, Retrieval-Augmented Generation (RAG), recommendation systems, and other embedding-based applications.

> **Status:** Active development. The project is functional but not yet production-ready.

---

## Features

* 🚀 Fast approximate nearest neighbor (ANN) search using **HNSW**
* 📚 Multiple independent vector stores
* 🔍 Semantic search over text documents
* 🌐 REST API
* 🐍 Python client library
* 🤖 Automatic embedding generation through an OpenAI-compatible embedding endpoint (e.g. Ollama)
* 📄 Document metadata support
* 🧩 Clean modular architecture

---

## Architecture

```
                 +----------------+
                 | Python Client  |
                 +-------+--------+
                         |
                         | HTTP
                         |
                 +-------v--------+
                 |   REST Server  |
                 +-------+--------+
                         |
                 +-------v--------+
                 | Search Engine  |
                 +-------+--------+
                         |
         +---------------+----------------+
         |                                |
+--------v---------+             +--------v--------+
| Document Store   |             |   HNSW Index    |
+------------------+             +-----------------+
                         |
                         |
                 Embedding Provider
          (Ollama / OpenAI Compatible)
```

---

## How it works

When a document is indexed:

1. The text is sent to the configured embedding endpoint.
2. The returned embedding vector is stored.
3. The vector is inserted into the HNSW graph.
4. The original document and metadata are retained.

During search:

1. The query is converted into an embedding.
2. HNSW performs approximate nearest neighbor search.
3. Matching documents are returned, ranked by similarity.

---

## Project Structure

```
simdb/
├── core/
│   ├── graph.go          # Graph data structures
│   ├── vector.go         # Vector operations
│   ├── nsw.go            # Plain NSW implementation
│   ├── hnsw.go           # HNSW implementation
│   └── searchengine.go   # Indexing and search logic
│
├── rest/
│   └── server.go         # REST API
│
├── python/
│   ├── simdb.py          # Python SDK
│   └── example.py
│
├── main.go
└── go.mod
```

---

## Getting Started

### Clone

```bash
git clone https://github.com/OnkarKunjir/simdb.git
cd simdb
```

### Run

```bash
go run .
```

By default the REST server starts on port **8080**.

---

## Python Client

```python
from simdb import SimDB

db = SimDB("http://localhost:8080")

db.create_store("books")

db.index(
    store="books",
    document="Crime and Punishment was written by Fyodor Dostoevsky.",
    metadata={
        "author": "Fyodor Dostoevsky"
    }
)

results = db.search(
    store="books",
    query="Russian novels",
    top_k=5
)

print(results)
```

---

## REST API

### Create Store

```http
POST /stores
```

### List Stores

```http
GET /stores
```

### Index Document

```http
POST /documents
```

Example

```json
{
    "store": "books",
    "document": "Crime and Punishment was written by Fyodor Dostoevsky.",
    "metadata": {
        "author": "Fyodor Dostoevsky"
    }
}
```

### Search

```http
POST /search
```

Example

```json
{
    "store": "books",
    "query": "Russian novels",
    "top_k": 5
}
```

---

## HNSW

SimDB uses the **Hierarchical Navigable Small World (HNSW)** algorithm for approximate nearest neighbor search.

The implementation includes:

* Multi-layer graph construction
* Random level generation
* Greedy graph traversal
* Beam search
* Neighbor pruning heuristic
* Configurable `M`
* Configurable `efConstruction`
* Configurable `efSearch`

A plain **NSW** implementation is also included for comparison and experimentation.

---

## Current Limitations

SimDB is still under active development.

Current limitations include:

* In-memory storage only
* No persistence
* No replication
* No authentication
* No document deletion
* No update operations
* Limited concurrent access guarantees

---

## Roadmap

* [ ] Persistent storage
* [ ] Batch indexing
* [ ] Concurrent indexing/search
* [ ] Pluggable embedding providers
* [ ] Filtered search
* [ ] Hybrid lexical + vector search
* [ ] Index serialization
* [ ] Benchmark suite
* [ ] Performance profiling

---

## Why SimDB?

The project started as an implementation of the HNSW algorithm and evolved into a complete vector search engine with a clean API and developer-friendly architecture.

The goal is to provide a simple, understandable, and extensible codebase for experimenting with modern vector search while remaining practical enough to power real semantic search applications.

---

## Contributing

Contributions, bug reports, and feature requests are welcome.

If you'd like to improve the implementation or add new capabilities, feel free to open an issue or submit a pull request.

---
