package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Tts(w http.ResponseWriter, r *http.Request) {
	input := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
		if text, ok := input["text"].(string); ok {
			if speech, err := MicrosoftTtsService(text); err == nil {
				response := map[string]interface{}{"speech": speech}
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

func MicrosoftTtsService(text string) (speech string, err error) {

	return fmt.Sprintf("base64( '%s'.wav )", text), nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", Tts).Methods("POST")
	fmt.Println("Text to speech service running on localhost:3003")
	log.Fatal(http.ListenAndServe(":3003", r))
}
