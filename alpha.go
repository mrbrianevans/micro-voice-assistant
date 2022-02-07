package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Alpha(w http.ResponseWriter, r *http.Request) {
	input := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
		if question, ok := input["text"].(string); ok {
			if answer, err := WolframService(question); err == nil {
				response := map[string]interface{}{"text": answer}
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

func WolframService(question string) (string, error) {
	return "I can't answer real questions yet", nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	fmt.Println("Alpha service running on localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", r))
}
