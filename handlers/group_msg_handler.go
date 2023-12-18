package handlers

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %v Text Msg : %v", group.NickName, msg.Content)

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 替换掉@文本，然后向GPT发起请求
	replaceText := "@" + sender.Self().NickName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	err = os.WriteFile("req.txt", []byte(requestText), 0644)
	if err != nil {
		log.Printf("request error: %v \n", err)
		msg.ReplyText("机器人挂了，我一会发现了就去修；或者你可以试试重发")
		return err
	}

	command := "python"
	args := []string{"call_gemini.py"}

	// Create a new command object
	cmd := exec.Command(command, args...)

	var output []byte
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("reply error: %v \n %v", err, string(output))
		msg.ReplyText("AI挂了，我一会发现了就去修；或者你可以试试重发")
		return err
	}

	if output == nil {
		msg.ReplyText("机器人没回复，我一会发现了就去修。")
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}

	// 回复@我的用户
	reply := strings.TrimSpace(string(output))
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + reply
	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}
