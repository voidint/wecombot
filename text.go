package wecombot

// TextMessage 文本类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E6%9C%AC%E7%B1%BB%E5%9E%8B
type TextMessage struct {
	// MsgType 必填。消息类型，此时固定为：text。
	MsgType string `json:"msgtype"`
	Text    struct {
		// Content 必填。文本内容，最长不超过2048个字节，必须是utf8编码。
		Content string `json:"content"`
		// MentionedList 可选。userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list。
		MentionedList []string `json:"mentioned_list"`
		// MentionedMobileList 可选。手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人。
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

// WithMentionedList 设置被提醒的 userid 列表
func WithMentionedList(list []string) func(*TextMessage) {
	return func(msg *TextMessage) {
		msg.Text.MentionedList = list
	}
}

// WithMentionedMobileList 设置被提醒的手机号列表
func WithMentionedMobileList(list []string) func(*TextMessage) {
	return func(msg *TextMessage) {
		msg.Text.MentionedMobileList = list
	}
}

// SendTextMessage 发送文本消息
func (bot *Bot) SendTextMessage(msg *TextMessage) error {
	return bot.send(msg)
}

// SendText 发送文本消息
func (bot *Bot) SendText(content string, opts ...func(*TextMessage)) (err error) {
	msg := TextMessage{
		MsgType: "text",
	}
	for _, setter := range opts {
		setter(&msg)
	}

	return bot.SendTextMessage(&msg)
}
