package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Alexa(w http.ResponseWriter, r *http.Request) {
	input := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
		if speech, ok := input["speech"].(string); ok {
			if answer, err := AlexaService(speech); err == nil {
				response := map[string]interface{}{"speech": answer}
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

func AlphaService(text string) (string, error) {

	return "", nil
}
func TtsService(text string) (string, error) {

	return "", nil
}
func SttService(text string) (string, error) {

	return "", nil
}

func AlexaService(questionSpeech string) (answerSpeech string, err error) {
	return "base64Encoded (answer.wav)", nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	fmt.Println("Alexa service running on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
