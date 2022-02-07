# Micro Voice Assistant
A voice assisstant built in Go using a microservices architecture.

## Speech to text and text to speech
[Microsoft Speech Service](https://docs.microsoft.com/en-gb/azure/cognitive-services/speech-service/)
is used to convert user input speech to text and then convert the output text back to speech.

## Answer engine
[Wolfram Alpha](https://products.wolframalpha.com/api/) is used to get answers to questions.

## Internal microservices
All microservices respond to `POST` requests. 
They take JSON input and respond with JSON.

Example request format:
```http request
POST /endpoint
Content-Type: application/json

{"hello": "world"}
```
Example HTTP requests can be found in `./tests.http`.

### Alpha
Input example:
```json
{ "text": "What is the melting point of silver?" }
```
Output example:
```json
{ "text": "961.78 degrees Celsius" }
```

### Speech-to-text
Input example:
```json
{ "speech": "base64( wav )" }
```
Output example:
```json
{ "text": "words" }
```

### Text-to-speech
Input example:
```json
{ "text": "words" }
```
Output example:
```json
{ "speech": "base64( wav )" }
```

### Alexa
Input example:
```json
{ "speech": "base64( wav )" }
```
Output example:
```json
{ "speech": "base64( wav )" }
```