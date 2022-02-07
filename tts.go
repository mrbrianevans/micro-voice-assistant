package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Tts(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"text": "What is the melting point of silver?"})
	w.Write(res)
}

func MicrosoftTtsService(text string) (string, error) {

	return "", nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", Tts).Methods("POST")
	fmt.Println("Text to speech service running on localhost:3003")
	log.Fatal(http.ListenAndServe(":3003", r))
}
