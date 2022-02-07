package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Stt(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"text": "What is the melting point of silver?"})
	w.Write(res)
}

func MicrosoftSttService(text string) (string, error) {

	return "", nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", Stt).Methods("POST")
	fmt.Println("Speech to text service running on localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", r))
}
