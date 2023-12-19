package handlers

import (
	"io"
	"log"

	"github.com/Wuvist/geminibot/goapi"
	"github.com/eatmoreapple/openwechat"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// // 好友申请
	// if msg.IsFriendAdd() {
	// 	if config.LoadConfig().AutoPass {
	// 		_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
	// 		if err != nil {
	// 			log.Fatalf("add friend agree error : %v", err)
	// 			return
	// 		}
	// 	}
	// }

	// 私聊
	handlers[UserHandler].handle(msg)
}

func handleIfPicture(msg *openwechat.Message) error {
	if msg.IsPicture() {
		response, err := msg.GetPicture()
		if err != nil {
			log.Printf("get picture error :%v \n", err)
			return err
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("read picture error :%v \n", err)
			return err
		}
		response.Body.Close()
		goapi.SetPicutre(body)
	}
	return nil
}
