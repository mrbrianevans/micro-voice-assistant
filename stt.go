package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// requires MICROSOFT_KEY environment variable to run
const (
	STT_URI = "https://uksouth.stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
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
	client := &http.Client{}
	speechBytes := base64.NewDecoder(&base64.Encoding{}, strings.NewReader(speech))
	if req, err := http.NewRequest("POST", STT_URI, speechBytes); err == nil {
		req.Header.Set("Content-Type", "audio/wav;codecs=audio/pcm;samplerate=16000")
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MICROSOFT_KEY"))
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				if body, err := ioutil.ReadAll(rsp.Body); err == nil {
					return string(body), nil
				}
			}
		}
	}
	return "", err
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stt", Stt).Methods("POST")
	fmt.Println("Speech to text service running on localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", r))
}
