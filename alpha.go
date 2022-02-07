package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var WolframKey = os.Getenv("WOLFRAM_KEY") // uses environment variable to avoid exposing API key
var WolframUri = "http://api.wolframalpha.com/v1/result?appid=" + WolframKey

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

func WolframService(question string) (answer string, err error) {
	client := &http.Client{}
	uri, _ := url.Parse(WolframUri)
	query := uri.Query()
	query.Set("i", question)
	uri.RawQuery = query.Encode()
	fmt.Println("Requesting URL:", uri.String(), query.Encode())
	if req, err := http.NewRequest("GET", uri.String(), nil); err == nil {
		if rsp, err := client.Do(req); err == nil {
			if rsp.StatusCode == http.StatusOK {
				if answer, err := ioutil.ReadAll(rsp.Body); err == nil {
					return string(answer), nil
				}
			}
		}
	}
	return "", errors.New("WolframService")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	fmt.Println("Alpha service running on localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", r))
}
