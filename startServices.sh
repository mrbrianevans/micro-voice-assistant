#!/bin/bash
nohup go run stt.go &
nohup go run tts.go &
nohup go run alpha.go &
nohup go run alexa.go &