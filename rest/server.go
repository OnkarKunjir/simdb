package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/onkarkunjir/simdb/core"
)

type RestServer struct {
	searchEngines map[string]*core.SearchEngine
	port          int
}

type CreateStoreRequest struct {
	Name        string `json:"name"`
	M           int    `json:"m"`
	EFConstruct int    `json:"efConstruct"`
	EFsearch    int    `json:"efSearch"`
	Url         string `json:"url"`
	Model       string `json:"model"`
}

type IndexDocumentRequest struct {
	Name     string            `json:"name"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

type SearchRequest struct {
	Name         string `json:"name"`
	SearchString string `json:"searchString"`
	K            int    `json:"k"`
}

type SearchResult struct {
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
}

func InitServer(port int) *RestServer {
	server := &RestServer{
		searchEngines: make(map[string]*core.SearchEngine),
		port:          port,
	}

	http.HandleFunc("/list-stores", server.listStores)
	http.HandleFunc("/create-store", server.createStore)
	http.HandleFunc("/index-document", server.indexDocument)
	http.HandleFunc("/search", server.search)
	http.HandleFunc("/search/{name}/{searchString}/{k}", server.search)

	return server
}

func (server *RestServer) ListenAndServe() {
	log.Println("Listening on port: ", server.port)
	http.ListenAndServe(":"+strconv.Itoa(server.port), nil)
}

func (server *RestServer) listStores(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if len(server.searchEngines) == 0 {
		fmt.Fprintf(writer, "[]")
		return
	}
	keys := slices.Collect(maps.Keys(server.searchEngines))
	response, err := json.Marshal(keys)

	if err != nil {
		log.Println("failed to marshal response:", err)
		http.Error(writer, "Failed to serialize available stores", http.StatusInternalServerError)
		return
	}
	writer.Write(response)
}

func (server *RestServer) createStore(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if http.MethodPost != request.Method {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	createStoreRequest := CreateStoreRequest{}
	err := json.NewDecoder(request.Body).Decode(&createStoreRequest)
	if err != nil {
		log.Println("Failed to unmarshal request body", err)
		http.Error(writer, "Failed to unmarshal request body", http.StatusInternalServerError)
		return
	}

	if _, contains := server.searchEngines[createStoreRequest.Name]; contains {
		log.Println("Store with name: \"" + createStoreRequest.Name + "\" already exists")
		http.Error(writer, "Store with name: \""+createStoreRequest.Name+"\" already exists", http.StatusConflict)
		return
	}

	server.searchEngines[createStoreRequest.Name] = core.CreateSearchEngine(
		createStoreRequest.M,
		createStoreRequest.EFConstruct,
		createStoreRequest.EFsearch,
		createStoreRequest.Url,
		createStoreRequest.Model,
	)

	log.Println("Adding new store with name: \"" + createStoreRequest.Name + "\"")
	json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}

func (server *RestServer) indexDocument(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if http.MethodPost != request.Method {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	indexDocumentRequest := IndexDocumentRequest{}
	err := json.NewDecoder(request.Body).Decode(&indexDocumentRequest)
	if err != nil {
		log.Println("Failed to unmarshal request body", err)
		http.Error(writer, "Failed to unmarshal request body", http.StatusInternalServerError)
		return
	}

	name := indexDocumentRequest.Name
	if searchEngine, ok := server.searchEngines[name]; ok {
		id := searchEngine.Index(indexDocumentRequest.Content, indexDocumentRequest.Metadata)
		log.Println("Indexed docuemnt: " + id + " into store: " + name)

		json.NewEncoder(writer).Encode(map[string]string{"id": id})
		return
	}

	log.Println("Store with " + name + " does not exists")
	http.Error(writer, "Store with "+name+" does not exists", http.StatusNotFound)
}

func (server *RestServer) search(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var name, searchString string
	var k int
	if http.MethodPost == request.Method {
		searchRequest := SearchRequest{}
		err := json.NewDecoder(request.Body).Decode(&searchRequest)
		if err != nil {
			log.Println("Failed to unmarshal request body", err)
			http.Error(writer, "Failed to unmarshal request body", http.StatusInternalServerError)
			return
		}
		name = searchRequest.Name
		searchString = searchRequest.SearchString
		k = searchRequest.K
	} else if http.MethodGet == request.Method {
		name = request.PathValue("name")
		searchString = request.PathValue("searchString")
		k, _ = strconv.Atoi(request.PathValue("k"))
	} else {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if searchEngine, ok := server.searchEngines[name]; ok {
		documents := searchEngine.Search(searchString, k)
		results := make([]*SearchResult, 0)
		for _, document := range documents {
			results = append(results, &SearchResult{Content: document.Content, Metadata: document.Metadata})
		}

		json.NewEncoder(writer).Encode(results)
		return
	}

	log.Println("Store with " + name + " does not exists")
	http.Error(writer, "Store with "+name+" does not exists", http.StatusNotFound)
}
