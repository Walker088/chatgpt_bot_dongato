package web

import (
	"log"
	"testing"
)

func TestAskQuestion(t *testing.T) {
	var sessionToken = ""
	var cfClearance = ""
	e, err := NewEngine(sessionToken, cfClearance)
	if err != nil {
		log.Fatalf("Couldn't start chatgpt engine: %v", err)
	}
	question := "How to study English?"
	reply, err := e.AskQuestion(874894996, question)
	if err != nil {
		log.Fatalf("Couldn't start chatgpt engine: %v", err)
	}
	log.Println(string(reply))
}
