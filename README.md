# wechatbot

项目基于[openwechat](https://github.com/eatmoreapple/openwechat)
开发 ###目前实现了以下功能

- 群聊@回复
- 私聊回复
- 默认使用`gemini-pro`模型（支持上下文聊天记录，默认 5 分钟超时）
- 需要读图时调用`gemini-pro-vision`模型（似乎不支持上下文）

# 安装使用

```
# 获取项目
git clone https://github.com/Wuvist/geminibot.git

# 进入项目目录
cd geminibot

# 启动项目
go build
./geminibot

启动前需编辑`key.txt`输入google.generativeai的api key
或者通过`API_KEY`的环境变量传递
```
