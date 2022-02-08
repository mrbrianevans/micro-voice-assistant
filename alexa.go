package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var ServiceRegistry = map[string]string{
	"alpha": "http://localhost:3001/alpha",
	"stt":   "http://localhost:3002/stt",
	"tts":   "http://localhost:3003/tts",
}

const (
	ALPHA = "http://localhost:3001/alpha"
	STT   = "http://localhost:3002/stt"
	TTS   = "http://localhost:3003/tts"
)

func Alexa(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to Alexa")
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

type SpeechBody struct {
	Speech string `json:"speech"`
}
type TextBody struct {
	Text string `json:"text"`
}

func AlphaService(question string) (answer string, err error) {
	reqBody, err := json.Marshal(TextBody{Text: question})
	response, err := http.Post(ServiceRegistry["alpha"], "application/json", bytes.NewReader(reqBody))
	var resBody TextBody
	if err := json.NewDecoder(response.Body).Decode(&resBody); err == nil {
		return resBody.Text, err
	}
	return "", err
}
func TtsService(text string) (speech string, err error) {
	reqBody, err := json.Marshal(TextBody{Text: text})
	response, err := http.Post(ServiceRegistry["tts"], "application/json", bytes.NewReader(reqBody))
	var resBody SpeechBody
	if err := json.NewDecoder(response.Body).Decode(&resBody); err == nil {
		return resBody.Speech, err
	}
	return "", err
}
func SttService(speech string) (text string, err error) {
	reqBody, err := json.Marshal(SpeechBody{Speech: speech})
	if err != nil {
		return "", errors.New("could not marshal JSON")
	}
	response, err := http.Post(ServiceRegistry["stt"], "application/json", bytes.NewReader(reqBody))
	if err != nil || response.StatusCode != http.StatusOK {
		return "", errors.New("could not request STT service on " + ServiceRegistry["stt"])
	}
	var resBody TextBody
	err = json.NewDecoder(response.Body).Decode(&resBody)
	if err != nil {
		return "", errors.New("could not unmarshal JSON")
	}
	return resBody.Text, err
}

func AlexaService(questionSpeech string) (answerSpeech string, err error) {
	//todo: convert this to non-nested style to more easily handle errors.

	// convert input speech to text, using stt
	if questionText, err := SttService(questionSpeech); err == nil {
		// get answer to question using alpha
		if answerText, err := AlphaService(questionText); err == nil {
			// convert text answer to speech using tts
			if answerSpeech, err = TtsService(answerText); err == nil {
				// return speech
				return answerSpeech, nil
			}
		}
	}
	log.Println("Error in Alexa")
	return "", errors.New("Alexa service failed to answer question")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alexa", Alexa).Methods("POST")
	fmt.Println("Alexa service running on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
