package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Alpha(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"text": "What is the melting point of silver?"})
	w.Write(res)
}

func WolframService(text string) (string, error) {
	return "", nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	fmt.Println("Alpha service running on localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", r))
}
