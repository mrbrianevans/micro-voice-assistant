package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Stt(w http.ResponseWriter, r *http.Request) {
	input := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
		if speech, ok := input["speech"].(string); ok {
			if text, err := MicrosoftSttService(speech); err == nil {
				response := map[string]interface{}{"text": text}
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func MicrosoftSttService(speech string) (text string, err error) {

	return " this is what you asked ", nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", Stt).Methods("POST")
	fmt.Println("Speech to text service running on localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", r))
}
