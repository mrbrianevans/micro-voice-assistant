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
	TTS_URI = "https://uksouth.stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
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
	client := &http.Client{}
	if req, err := http.NewRequest("POST", STT_URI, strings.NewReader(text)); err == nil {
		req.Header.Set("Content-Type", "application/ssml+xml")
		req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MICROSOFT_KEY"))
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				if speechBytes, err := ioutil.ReadAll(base64.NewDecoder(&base64.Encoding{}, rsp.Body)); err == nil {
					return string(speechBytes), nil
				}
			}
		}
	}
	return "", err
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tts", Tts).Methods("POST")
	fmt.Println("Text to speech service running on localhost:3003")
	log.Fatal(http.ListenAndServe(":3003", r))
}
