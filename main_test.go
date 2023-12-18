package main

import (
	"fmt"
	"testing"

	"github.com/Wuvist/geminibot/goapi"
)

func TestLive(t *testing.T) {
	reply := goapi.GetReply("", "tell me a joke")
	fmt.Println(reply)
}
