package wecombot

// NewsMessage 图文类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E6%96%87%E7%B1%BB%E5%9E%8B
type NewsMessage struct {
	// MsgType 必填。消息类型，此时固定为 news 。
	MsgType MsgType `json:"msgtype"`
	// 消息内容
	News struct {
		Articles []*Article `json:"articles"`
	} `json:"news"`
}

// Article 图文
type Article struct {
	// Title 标题，不超过128个字节，超过会自动截断。
	Title string `json:"title"`
	// Description 描述，不超过512个字节，超过会自动截断。
	Description *string `json:"description"`
	// URL 点击后跳转的链接。
	URL string `json:"url"`
	// PicURL 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150。
	PicURL *string `json:"picurl"`
}

// SendNewsMessage 发送图文消息
func (bot *Bot) SendNewsMessage(msg *NewsMessage) (err error) {
	msg.MsgType = NewsMsgType
	return bot.send(msg)
}

// SendNews 发送图文消息
func (bot *Bot) SendNews(articles ...*Article) (err error) {
	var msg NewsMessage
	msg.News.Articles = articles

	return bot.SendNewsMessage(&msg)
}
