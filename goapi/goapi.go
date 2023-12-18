package goapi

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func GetReply(sender, msg string) (reply string) {
	err := os.WriteFile("req.txt", []byte(msg), 0644)
	if err != nil {
		log.Printf("request error: %v \n", err)
		return "机器人挂了，我一会发现了就去修；或者你可以试试重发"
	}

	command := "python"
	args := []string{"call_gemini.py"}

	// Create a new command object
	cmd := exec.Command(command, args...)

	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("reply error: %v \n %v", err, string(output))
		return "AI挂了，我一会发现了就去修；或者你可以试试重发"
	}

	if output == nil {
		return "AI没回复，不知道为啥，我一会发现了就去喵一喵~"
	}

	reply = strings.TrimSpace(string(output))
	reply = strings.Trim(reply, "\n")

	return
}
