package main

import (
	"example/es_golang/internal/pkg/storage/elasticsearch"
	"example/es_golang/internal/post"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	elastic, err := elasticsearch.New([]string{"http://es-container:9200"})
	if err != nil {
		log.Fatalln(err)
	}

	if err := elastic.CreateIndex("post"); err != nil {
		log.Fatalln(err)
	}

	storage, err := elasticsearch.NewPostStorage(*elastic)
	if err != nil {
		log.Fatalln(err)
	}

	postAPI := post.New(&storage)
	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/posts", postAPI.Create)
	router.HandlerFunc("PATCH", "/api/v1/posts/:id", postAPI.Update)

	log.Fatalln(http.ListenAndServe(":9000", router))
}
