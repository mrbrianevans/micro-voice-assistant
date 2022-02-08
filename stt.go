package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// requires MICROSOFT_KEY environment variable to run
const (
	STT_URI = "https://uksouth.stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US&format=simple"
)

func Stt(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to STT")
	input := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err == nil {
		if speech, ok := input["speech"].(string); ok {
			if speechBytes, err := base64.StdEncoding.DecodeString(speech); err == nil {
				if text, err := MicrosoftSttService(speechBytes); err == nil {
					response := map[string]interface{}{"text": text}
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				} else { // failed to convert speech to text
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else { // base64 encoding is incorrect
				w.WriteHeader(http.StatusBadRequest)
			}
		} else { // json request body doesn't contain "speech" key
			w.WriteHeader(http.StatusBadRequest)
		}
	} else { // request body is not valid json
		w.WriteHeader(http.StatusBadRequest)
	}
}

type MicrosoftRecognitionResponse struct {
	RecognitionStatus string
	DisplayText       string
	Offset            int
	Duration          int
}

func MicrosoftSttService(speech []byte) (text string, err error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", STT_URI, bytes.NewReader(speech)); err == nil {
		req.Header.Set("Content-Type", "audio/wav;codecs=audio/pcm;samplerate=16000")
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MICROSOFT_KEY"))
		req.Header.Set("Accept", "application/json")
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				var response MicrosoftRecognitionResponse
				if err := json.NewDecoder(rsp.Body).Decode(&response); err == nil {
					return response.DisplayText, nil
				}
			}
		}
	}
	return "", errors.New("error while converting speech to text")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", Stt).Methods("POST")
	fmt.Println("Speech to text service running on localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", r))
}
