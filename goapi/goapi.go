package goapi

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var model *genai.GenerativeModel
var ctx context.Context

func init() {
	ctx = context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)

	key, _ := os.ReadFile("key.txt")

	client, err := genai.NewClient(ctx, option.WithAPIKey(string(key)))
	if err != nil {
		log.Fatal(err)
	}

	// For text-only input, use the gemini-pro model
	model = client.GenerativeModel("gemini-pro")
}

func GetReply(sender, msg string) (reply string) {
	resp, err := model.GenerateContent(ctx, genai.Text(msg))
	if err != nil {
		log.Printf("reply error: %v \n", err)
		return "AI挂了，我一会发现了就去修；或者你可以试试重发"
	}

	if len(resp.Candidates) > 0 {
		cand := resp.Candidates[0]
		if cand.Content != nil && len(cand.Content.Parts) > 0 {
			part := cand.Content.Parts[0]
			reply = fmt.Sprintln(part)
		}
	}

	if reply == "" {
		return "AI没回复，不知道为啥，我一会发现了就去喵一喵~"
	}

	reply = strings.Trim(reply, "\n")

	return
}
