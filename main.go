package main

import (
	"fmt"
)

var words = []string{
	// animals
	"cat", "dog", "elephant", "tiger", "whale", "dolphin", "penguin", "eagle", "cobra", "gorilla",
	// fruits
	"apple", "banana", "mango", "grape", "orange", "strawberry", "watermelon", "pineapple", "peach", "cherry",
	// vegetables
	"carrot", "broccoli", "spinach", "tomato", "potato", "onion", "garlic", "cucumber", "pepper", "mushroom",
	// countries
	"india", "japan", "brazil", "germany", "canada", "australia", "france", "nigeria", "argentina", "thailand",
	// cities
	"mumbai", "tokyo", "london", "paris", "newyork", "berlin", "sydney", "cairo", "toronto", "bangkok",
	// sports
	"cricket", "football", "tennis", "basketball", "swimming", "boxing", "cycling", "golf", "rugby", "volleyball",
	// tech
	"computer", "keyboard", "monitor", "network", "server", "database", "algorithm", "compiler", "processor", "memory",
	// programming languages
	"golang", "python", "javascript", "rust", "java", "typescript", "kotlin", "swift", "haskell", "ruby",
	// emotions
	"happy", "sad", "angry", "fearful", "surprised", "disgusted", "anxious", "excited", "bored", "content",
	// weather
	"sunny", "rainy", "cloudy", "stormy", "windy", "snowy", "foggy", "humid", "freezing", "drought",
}

func main() {
	const M = 16
	const efConstruct = M * 4
	const efsearch = 40
	const url = "http://localhost:11434/api/embed"
	const model = "qwen3-embedding:0.6b"

	searchEngine := CreateSearchEngine(M, efConstruct, efsearch, url, model)

	fmt.Println("Starting indexing")
	for _, word := range words {
		searchEngine.Index(word, nil)
	}
	fmt.Println("Finished indexing")


	var searchString string
	for {
		fmt.Println("Search: ")
		fmt.Scan(&searchString)

		for index, document := range searchEngine.Search(searchString, 5) {
			fmt.Println(index, document.Content)
		}
		fmt.Println("----------------------------")
	}

}
