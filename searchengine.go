package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Document struct {
	Id      string
	Content string
	Vector  *Vector

	// TODO: for now adding just string vs string map as metadat
	Metadata map[string]string
}

type SearchEngine struct {
	hnsw                     HierarchicalNavigableSmallWorld
	m, efconstruct, efsearch int
	documents                map[string]*Document
	url, model               string
}

type EmbeddingsRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddingsResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float64 `json:"embeddings"`
}

func createVector(url, model, id, content string) *Vector {
	request, err := json.Marshal(&EmbeddingsRequest{model, content})
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(request))
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	embeddingsResponse := EmbeddingsResponse{}

	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(responseContent, &embeddingsResponse); err != nil {
		log.Fatal(err)
	}

	return &Vector{Id: id, Values: embeddingsResponse.Embeddings[0]}
}

func CreateSearchEngine(m, efconstruct, efsearch int, url, model string) *SearchEngine {
	return &SearchEngine{
		hnsw:        HierarchicalNavigableSmallWorld{},
		m:           m,
		efconstruct: efconstruct,
		efsearch:    efsearch,
		documents:   make(map[string]*Document),
		url:         url,
		model:       model,
	}
}

func (searchEngine *SearchEngine) Index(content string, metadata map[string]string) {
	id := strconv.Itoa(len(searchEngine.documents))
	document := &Document{
		Id:       id,
		Content:  content,
		Metadata: metadata,
		Vector:   createVector(searchEngine.url, searchEngine.model, id, content),
	}

	searchEngine.documents[document.Id] = document
	searchEngine.hnsw.Insert(document.Vector, searchEngine.m, searchEngine.efconstruct)
}

func (searchEngine *SearchEngine) Search(content string, k int) []*Document {
	documents := make([]*Document, 0, k)
	vector := createVector(searchEngine.url, searchEngine.model, "search-temp-id", content)
	for _, searchResult := range searchEngine.hnsw.Search(vector, k, searchEngine.efsearch) {
		documents = append(documents, searchEngine.documents[searchResult.Vector.Id])
	}
	return documents
}
