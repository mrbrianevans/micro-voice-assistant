package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Alexa(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"text": "What is the melting point of silver?"})
	w.Write(res)
}

func AlphaService(text string) (string, error) {

	return "", nil
}
func TtsService(text string) (string, error) {

	return "", nil
}
func SttService(text string) (string, error) {

	return "", nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	fmt.Println("Alexa service running on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
