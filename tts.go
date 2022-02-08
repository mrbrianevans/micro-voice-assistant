package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// requires MICROSOFT_KEY environment variable to run
const (
	TTS_URI = "https://uksouth.tts.speech.microsoft.com/cognitiveservices/v1"
)

type SsmlVoice struct {
	XMLName  struct{} `xml:"voice"`
	Language string   `xml:"xml:lang,attr"`
	Gender   string   `xml:"xml:gender,attr"`
	Name     string   `xml:"name,attr"`
	Text     string   `xml:",chardata"`
}

// speech synthesis markup - format of xml request body to Microsoft
type SSML struct {
	XMLName  struct{} `xml:"speak"`
	Version  string   `xml:"version,attr"`
	Language string   `xml:"xml:lang,attr"`
	Voice    SsmlVoice
}

func Tts(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to Tts")
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
	ssml := SSML{
		XMLName:  struct{}{},
		Version:  "1.0",
		Language: "en-US",
		Voice: SsmlVoice{
			XMLName:  struct{}{},
			Language: "en-US",
			Gender:   "Male",
			Name:     "en-US-ChristopherNeural",
			Text:     text,
		},
	}
	xmlBytes, _ := xml.Marshal(ssml)
	if req, err := http.NewRequest("POST", TTS_URI, bytes.NewReader(xmlBytes)); err == nil {
		req.Header.Set("Content-Type", "application/ssml+xml")
		req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
		req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("MICROSOFT_KEY"))
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				if speechBytes, err := ioutil.ReadAll(rsp.Body); err == nil {
					speechBase64 := base64.StdEncoding.EncodeToString(speechBytes)
					return speechBase64, nil
				}
			}
		}
	}
	return "", errors.New("Text to speech error")
}
func main() {
	if _, envSet := os.LookupEnv("MICROSOFT_KEY"); !envSet {
		log.Fatal("Environment variable MICROSOFT_KEY not set. Required to access API")
	}
	r := mux.NewRouter()
	r.HandleFunc("/tts", Tts).Methods("POST")
	fmt.Println("Text to speech service running on localhost:3003")
	log.Fatal(http.ListenAndServe(":3003", r))
}
