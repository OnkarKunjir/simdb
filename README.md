# SimDB

A vector similarity search database written in Go, built on a from-scratch implementation of the HNSW (Hierarchical Navigable Small World) algorithm. Uses [Ollama](https://ollama.com) to generate embeddings locally.

## How it works

SimDB implements HNSW — a graph-based approximate nearest neighbour algorithm that organises vectors into multiple layers. Upper layers act as highways for fast coarse navigation, and the base layer (layer 0) performs the fine-grained search. This gives O(log n) search complexity with high recall.

The core algorithm is implemented from scratch in Go with no external dependencies. Embeddings are generated via a local Ollama instance, so no external API keys are needed.

## Project structure

```
core/       HNSW implementation — graph, vectors, search engine
server/     REST API server
python/     Python client wrapper for the REST API
```

## Prerequisites

- Go 1.22+
- [Ollama](https://ollama.com) running locally
- An embedding model pulled in Ollama

```bash
ollama pull qwen3-embedding:0.6b
```

## Running the server

```bash
go run .
```

Server starts on port 8080 by default.

## REST API

### List stores

```
GET /list-stores
```

Returns a list of all existing vector stores.

### Create a store

```
POST /create-store
```

```json
{
  "name": "my-store",
  "m": 16,
  "efConstruct": 64,
  "efSearch": 40,
  "url": "http://localhost:11434/api/embed",
  "model": "qwen3-embedding:0.6b"
}
```

### Index documents

```
POST /index-documents
```

```json
{
  "store": "my-store",
  "documents": ["cat", "dog", "elephant"]
}
```

### Search

```
POST /search
```

```json
{
  "store": "my-store",
  "query": "kitten",
  "k": 5
}
```

## Python client

```python
from simdb import SimDB

db = SimDB("http://localhost:8080")
db.create_store("my-store", model="qwen3-embedding:0.6b")
db.index(["cat", "dog", "elephant", "tiger"])
results = db.search("kitten", k=3)
```

## HNSW parameters

| Parameter | Default | Description |
|---|---|---|
| `m` | 16 | Max connections per node per layer. Higher = better recall, more memory. |
| `efConstruct` | 64 | Beam width during index build. Higher = better graph quality, slower build. |
| `efSearch` | 40 | Beam width during search. Higher = better recall, slower queries. |

For datasets above 1000 documents, `efSearch=80` is recommended for ~0.98 recall.

## Running locally without the server

```bash
go run .
```

This starts an interactive search session using the word dataset in `main.go`.

```
Starting indexing
Finished indexing
Search:
kitten
0 cat
1 dog
2 tiger
3 dolphin
4 penguin
```

## Performance

Benchmarked on 1000 documents with `m=16`, `efConstruct=64`, `efSearch=80`:

- Average recall: ~0.98
- Embedding model: qwen3-embedding:0.6b (500-dimensional vectors)
