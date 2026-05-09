package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

func main() {
	wsContainer := restful.NewContainer()
	log.Printf("vBlog Core starting on :8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
