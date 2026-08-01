package main

import (
	"github.com/onkarkunjir/simdb/rest"
)

func main() {
	server := rest.InitServer(8080)
	server.ListenAndServe()
}
