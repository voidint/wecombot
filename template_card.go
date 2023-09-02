package wecombot

// TextNoticeTemplateCardMessage 文本通知模版卡片类型消息
type TextNoticeTemplateCardMessage struct {
	// MsgType 必填。消息类型，此时的消息类型固定为 template_card 。
	MsgType MsgType `json:"msgtype"`
	// TemplateCard 模板卡片
	TemplateCard struct {
		// CardType 模版卡片的模版类型，文本通知模版卡片的类型为 text_notice 。
		CardType CardType `json:"card_type"`
		// Source 可选。卡片来源样式信息。
		Source *Source `json:"source"`
		// MainTitle 模版卡片的主要内容，包括一级标题和标题辅助信息。
		MainTitle MainTitle `json:"main_title"`
		// EmphasisContent 关键数据样式
		EmphasisContent *EmphasisContent `json:"emphasis_content"`
		// QuoteArea 引用文献样式，建议不与关键数据共用。
		QuoteArea *QuoteArea `json:"quote_area"`
		// SubTitleText 二级普通文本，建议不超过112个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写。
		SubTitleText *string `json:"sub_title_text"`
		// HorizontalContentList 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6。
		HorizontalContentList []*HorizontalContent `json:"horizontal_content_list"`
		// JumpList 跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3。
		JumpList []*Jump `json:"jump_list"`
		// CardAction 整体卡片的点击跳转事件，text_notice模版卡片中该字段为必填项。
		CardAction CardAction `json:"card_action"`
	} `json:"template_card"`
}

// NewsNoticeTemplateCardMessage 图文展示模版卡片类型消息
type NewsNoticeTemplateCardMessage struct {
	// MsgType 必填。消息类型，此时的消息类型固定为 template_card 。
	MsgType MsgType `json:"msgtype"`
	// TemplateCard 模板卡片
	TemplateCard struct {
		// CardType 模版卡片的模版类型，图文展示模版卡片的类型为 news_notice 。
		CardType CardType `json:"card_type"`
		// Source 可选。卡片来源样式信息。
		Source *Source `json:"source"`
		// MainTitle 模版卡片的主要内容，包括一级标题和标题辅助信息。
		MainTitle MainTitle `json:"main_title"`
		// CardImage 图片样式
		CardImage CardImage `json:"card_image"`
		// ImageTextArea 左图右文样式
		ImageTextArea *ImageTextArea `json:"image_text_area"`
		// QuoteArea 引用文献样式，建议不与关键数据共用。
		QuoteArea *QuoteArea `json:"quote_area"`
		// VerticalContentList 卡片二级垂直内容，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过4。
		VerticalContentList []*VerticalContent `json:"vertical_content_list"`
		// HorizontalContentList 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6。
		HorizontalContentList []*HorizontalContent `json:"horizontal_content_list"`
		// JumpList 跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3。
		JumpList []*Jump `json:"jump_list"`
		// CardAction 整体卡片的点击跳转事件，text_notice模版卡片中该字段为必填项。
		CardAction CardAction `json:"card_action"`
	} `json:"template_card"`
}

// Source 卡片来源样式信息
type Source struct {
	// IconURL 来源图片的url
	IconURL *string `json:"icon_url"`
	// Desc 来源图片的描述，建议不超过13个字。
	Desc *string `json:"desc"`
	// DescColor 来源文字的颜色，目前支持：0(默认) 灰色，1 黑色，2 红色，3 绿色。
	DescColor *uint8 `json:"desc_color"`
}

// MainTitle 模版卡片的主要内容，包括一级标题和标题辅助信息。
type MainTitle struct {
	// Title 一级标题，建议不超过26个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写
	Title *string `json:"title"`
	// Desc 标题辅助信息，建议不超过30个字
	Desc *string `json:"desc"`
}

// EmphasisContent 关键数据样式
type EmphasisContent struct {
	// Title 关键数据样式的数据内容，建议不超过10个字
	Title *string `json:"title"`
	// Desc 关键数据样式的数据描述内容，建议不超过15个字
	Desc *string `json:"desc"`
}

// QuoteArea 引用文献样式，建议不与关键数据共用。
type QuoteArea struct {
	// Type 引用文献样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序。
	Type *uint8 `json:"type"`
	// URL 点击跳转的url，type是1时必填。
	URL *string `json:"url"`
	// AppID 点击跳转的小程序的appid，type是2时必填。
	AppID *string `json:"appid"`
	// PagePath 点击跳转的小程序的pagepath，type是2时选填。
	PagePath *string `json:"pagepath"`
	// Title 引用文献样式的标题
	Title *string `json:"title"`
	// QuoteText 引用文献样式的引用文案
	QuoteText *string `json:"quote_text"`
}

// HorizontalContent 二级标题+文本列表
type HorizontalContent struct {
	// Type 模版卡片的二级标题信息内容支持的类型，1是url，2是文件附件，3代表点击跳转成员详情。
	Type *uint8 `json:"type"`
	// Keyname 二级标题，建议不超过5个字。
	KeyName string `json:"keyname"`
	// Value 二级文本，如果type是2，该字段代表文件名称（要包含文件类型），建议不超过26个字。
	Value *string `json:"value"`
	// URL 链接跳转的url，type是1时必填。
	URL *string `json:"url"`
	// MediaID 附件的media_id，type是2时必填。
	MediaID *string `json:"media_id"`
	// UserID 成员详情的userid，type是3时必填。
	UserID *string `json:"userid"`
}

// Jump 跳转指引样式
type Jump struct {
	// Type 跳转链接类型，0或不填代表不是链接，1 代表跳转url，2 代表跳转小程序。
	Type *uint8 `json:"type"`
	// Title 跳转链接样式的文案内容，建议不超过13个字。
	Title string `json:"title"`
	// URL 跳转链接的url，type是1时必填。
	URL *string `json:"url"`
	// AppID 跳转链接的小程序的appid，type是2时必填。
	AppID *string `json:"appid"`
	// PagePath 跳转链接的小程序的pagepath，type是2时选填。
	PagePath *string `json:"pagepath"`
}

// 整体卡片的点击跳转事件
type CardAction struct {
	// Type 卡片跳转类型，1 代表跳转url，2 代表打开小程序。
	Type uint8 `json:"type"`
	// URL 跳转事件的url，card_action.type是1时必填。
	URL *string `json:"url"`
	// AppID 跳转事件的小程序的appid，card_action.type是2时必填。
	AppID *string `json:"appid"`
	// PagePath 跳转事件的小程序的pagepath，type是2时选填。
	PagePath *string `json:"pagepath"`
}

// 图片样式
type CardImage struct {
	// URL 图片的url
	URL string `json:"url"`
	// AspectRatio 图片的宽高比，宽高比要小于2.25，大于1.3，不填该参数默认1.3
	AspectRatio *float32 `json:"aspect_ratio"`
}

// ImageTextArea 左图右文样式
type ImageTextArea struct {
	// Type 左图右文样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序。
	Type *uint8 `json:"type"`
	// URL 点击跳转的url，type是1时必填。
	URL *string `json:"url"`
	// AppID 点击跳转的小程序的appid，必须是与当前应用关联的小程序，type是2时必填。
	AppID *string `json:"appid"`
	// PagePath 点击跳转的小程序的pagepath，type是2时选填。
	PagePath *string `json:"pagepath"`
	// Title 左图右文样式的标题
	Title *string `json:"title"`
	// Desc 左图右文样式的描述
	Desc *string `json:"desc"`
	// ImageURL 左图右文样式的图片url
	ImageURL string `json:"image_url"`
}

// VerticalContent 卡片二级垂直内容
type VerticalContent struct {
	// Title 卡片二级标题，建议不超过26个字。
	Title string `json:"title"`
	// Desc 二级普通文本，建议不超过112个字。
	Desc *string `json:"desc"`
}

// CardType 模版卡片类型
type CardType string

const (
	// TextNoticeCardType 文本通知模板卡片类型
	TextNoticeCardType CardType = "text_notice"
	// NewsNoticeCardType 图文展示模板卡片类型
	NewsNoticeCardType CardType = "news_notice"
)

// SendTextNoticeTemplateCardMessage 发送文本通知模板卡片类型消息
func (bot *Bot) SendTextNoticeTemplateCardMessage(msg *TextNoticeTemplateCardMessage) error {
	msg.MsgType = TemplateCardMsgType
	msg.TemplateCard.CardType = TextNoticeCardType
	return bot.send(msg)
}

// SendNewsNoticeTemplateCardMessage 发送图文展示模板卡片类型消息
func (bot *Bot) SendNewsNoticeTemplateCardMessage(msg *NewsNoticeTemplateCardMessage) error {
	msg.MsgType = TemplateCardMsgType
	msg.TemplateCard.CardType = NewsNoticeCardType
	return bot.send(msg)
}
