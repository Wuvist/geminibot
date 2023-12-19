package goapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var model, vModel *genai.GenerativeModel
var ctx context.Context

type session struct {
	cs       *genai.ChatSession
	lastChat time.Time
}

func (s *session) HasExpired() bool {
	return time.Since(s.lastChat).Minutes() > 5
}

func (s *session) Update() {
	s.lastChat = time.Now()
}

var sessions = make(map[string]*session)

func init() {
	ctx = context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)

	key := os.Getenv("API_KEY")
	if key == "" {
		keyTxt, err := os.ReadFile("key.txt")
		if err != nil {
			log.Fatal(err)
		}
		key = string(keyTxt)
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatal(err)
	}

	// For text-only input, use the gemini-pro model
	model = client.GenerativeModel("gemini-pro")
	model.SetCandidateCount(1)
	model.SetTemperature(0.1)

	// Not sure why only three Categories are allowed
	// vModel.SafetySettings = []*genai.SafetySetting{
	// 	{
	// 		Category:  genai.HarmCategoryUnspecified,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryDerogatory,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryToxicity,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryViolence,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategorySexual,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryMedical,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryDangerous,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryHarassment,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryHateSpeech,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategorySexuallyExplicit,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// 	{
	// 		Category:  genai.HarmCategoryDangerousContent,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// }

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
	}

	vModel = client.GenerativeModel("gemini-pro-vision")
	vModel.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
	}
	vModel.SetCandidateCount(1)
	vModel.SetTemperature(1)
}

var picData []byte

func SetPicutre(pic []byte) {
	picData = pic
}

func GetReply(sender, msg string) (reply string) {
	if msg == "重来" {
		delete(sessions, sender)
		return "好的~对话记录已经清除~"
	}

	needVmodel := false
	if strings.Index(msg, "[图片]") > 0 && picData != nil {
		i := strings.LastIndex(msg, " - - -")
		msg = msg[i+1:]
		needVmodel = true
	}

	var err error
	var ses *session
	var resp *genai.GenerateContentResponse

	if needVmodel {
		// sessions[sender] = ses
		// Somehow gemini-pro-vision close chat after first msg
		// So, remove from sessions to ensure new session next time
		delete(sessions, sender)

		mime := http.DetectContentType(picData)
		log.Printf("mime: %s \n", mime)
		if strings.Contains(mime, "png") {
			resp, err = vModel.GenerateContent(ctx, genai.Text(msg), genai.ImageData("png", picData))
		} else {
			resp, err = vModel.GenerateContent(ctx, genai.Text(msg), genai.ImageData("jpeg", picData))
		}

		picData = nil
	} else {
		ses = sessions[sender]
		if ses == nil || ses.HasExpired() {
			ses = &session{model.StartChat(), time.Now()}
			sessions[sender] = ses
		} else {
			ses.Update()
		}
		resp, err = ses.cs.SendMessage(ctx, genai.Text(msg))
	}

	if err != nil {
		log.Printf("reply error: %v \n", err)
		u, _ := json.Marshal(resp)
		log.Printf("resp json: %s \n", string(u))
		delete(sessions, sender)
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
		delete(sessions, sender)
		return "AI没回复，不知道为啥，我一会发现了就去喵一喵~"
	}

	reply = strings.Trim(reply, "\n")
	reply = strings.ReplaceAll(reply, "**", "")

	return
}
