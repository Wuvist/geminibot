package main

import (
	"testing"

	"github.com/Wuvist/geminibot/goapi"
)

func TestLive(t *testing.T) {
	reply := goapi.GetReply("", "Give me an example of something mean")
	if reply == "AI挂了，我一会发现了就去修；或者你可以试试重发" {
		t.Error("Safety setting not working")
	}
}
