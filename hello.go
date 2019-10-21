package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", HelloWorld)
	log.Fatal(
		http.ListenAndServe(":8080", nil),
	)
}
