package wecombot

// MarkdownMessage Markdown 类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B
type MarkdownMessage struct {
	// MsgType 必填。消息类型，此时固定为 markdown 。
	MsgType string `json:"msgtype"`
	// 消息内容
	Markdown struct {
		// Content 必填。markdown内容，最长不超过4096个字节，必须是utf8编码。
		Content string `json:"content"`
	} `json:"markdown"`
}

// SendMarkdownMessage 发送 Markdown 消息
func (bot *Bot) SendMarkdownMessage(msg *MarkdownMessage) (err error) {
	return bot.send(msg)
}

// SendMarkdown 发送 Markdown 消息
func (bot *Bot) SendMarkdown(content string) (err error) {
	var msg MarkdownMessage
	msg.MsgType = "markdown"
	msg.Markdown.Content = content

	return bot.SendMarkdownMessage(&msg)
}
