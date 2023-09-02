package wecombot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// Bot 企业微信群机器人
type Bot struct {
	webhookURL string
	key        string

	threadSafe     bool
	reqbuf, resbuf *bytes.Buffer
	client         *http.Client
	marshal        func(v any) ([]byte, error)
	unmarshal      func(data []byte, v any) error
}

// NewBot 返回企业微信群机器人实例
func NewBot(key string, opts ...func(*Bot)) *Bot {
	bot := Bot{
		webhookURL: fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key),
		key:        key,
		client:     http.DefaultClient,
		marshal:    json.Marshal,
		unmarshal:  json.Unmarshal,
	}

	for _, setter := range opts {
		setter(&bot)
	}

	if !bot.threadSafe {
		bot.reqbuf = bytes.NewBuffer(nil)
		bot.resbuf = bytes.NewBuffer(nil)
	}

	return &bot
}

// WithThreadSafe 设置线程安全模式
func WithThreadSafe() func(*Bot) {
	return func(bot *Bot) {
		bot.threadSafe = true
	}
}

// WithHttpClient 设置 HTTP 客户端
func WithHttpClient(c *http.Client) func(*Bot) {
	return func(bot *Bot) {
		bot.client = c
	}
}

// WithMarshal 设置序列化函数实现
func WithMarshal(f func(v any) ([]byte, error)) func(*Bot) {
	return func(bot *Bot) {
		bot.marshal = f
	}
}

// WithUnmarshal 设置反序列化函数实现
func WithUnmarshal(f func(data []byte, v any) error) func(*Bot) {
	return func(bot *Bot) {
		bot.unmarshal = f
	}
}

var jsonReqHeader = map[string]string{
	"Content-Type": "application/json",
}

func (bot *Bot) send(msg any) error {
	data, err := bot.marshal(msg)
	if err != nil {
		return err
	}

	var reqBody *bytes.Buffer
	if bot.threadSafe {
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bot.reqbuf
		reqBody.Write(data)
	}
	defer reqBody.Reset()

	var resData resData
	if err = bot.doPost(bot.webhookURL, jsonReqHeader, reqBody, &resData); err != nil {
		return err
	}
	return resData.ToError()
}

func (bot *Bot) doPost(url string, reqHeader map[string]string, reqBody io.Reader, resData any) error {
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return err
	}
	for k, v := range reqHeader {
		req.Header.Set(k, v)
	}

	res, err := bot.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var resbuf *bytes.Buffer
	if bot.threadSafe {
		resbuf = new(bytes.Buffer)
	} else {
		resbuf = bot.resbuf
	}
	defer resbuf.Reset()

	if _, err = io.Copy(resbuf, res.Body); err != nil { // TODO 优化这次多余的拷贝
		return err
	}

	return bot.unmarshal(resbuf.Bytes(), resData)
}

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

// ImageMessage 图片类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
type ImageMessage struct {
	// MsgType 必填。消息类型，此时固定为 image 。
	MsgType string `json:"msgtype"`
	// 消息内容
	Image struct {
		// Base64 必填。图片内容的base64编码
		Base64 string `json:"base64"`
		// Md5 必填。图片内容（base64编码前）的md5值
		Md5 string `json:"md5"`
	} `json:"image"`
}

// SendImageMessage 发送图片消息
func (bot *Bot) SendImageMessage(msg *ImageMessage) (err error) {
	return bot.send(msg)
}

// SendImage 发送图片消息
func (bot *Bot) SendImage(f io.Reader) (err error) {
	var buf *bytes.Buffer
	if bot.threadSafe {
		buf = bytes.NewBuffer(nil)
	} else {
		buf = bot.reqbuf
	}

	tee := io.TeeReader(f, buf)
	h := md5.New()
	if _, err = io.Copy(h, tee); err != nil {
		buf.Reset()
		return err
	}

	var msg ImageMessage
	msg.MsgType = "image"
	msg.Image.Md5 = fmt.Sprintf("%x", h.Sum(nil))
	msg.Image.Base64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	buf.Reset() // 确保先清除

	return bot.SendImageMessage(&msg)
}

// NewsMessage 图文类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E6%96%87%E7%B1%BB%E5%9E%8B
type NewsMessage struct {
	// MsgType 必填。消息类型，此时固定为 news 。
	MsgType string `json:"msgtype"`
	// 消息内容
	News struct {
		Articles []*Article `json:"articles"`
	} `json:"news"`
}

// Article 图文
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
}

// SendNewsMessage 发送图文消息
func (bot *Bot) SendNewsMessage(msg *NewsMessage) (err error) {
	return bot.send(msg)
}

// SendNews 发送图文消息
func (bot *Bot) SendNews(articles ...*Article) (err error) {
	var msg NewsMessage
	msg.MsgType = "news"
	msg.News.Articles = articles

	return bot.SendNewsMessage(&msg)
}

// FileMessage 文件类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E4%BB%B6%E7%B1%BB%E5%9E%8B
type FileMessage struct {
	// MsgType 必填。消息类型，此时固定为 file 。
	MsgType string `json:"msgtype"`
	// 消息内容
	File struct {
		// MediaID 必填。文件id，通过文件上传接口获取。
		MediaID string `json:"media_id"`
	} `json:"file"`
}

// SendFileMessage 发送文件消息
func (bot *Bot) SendFileMessage(msg *FileMessage) (err error) {
	return bot.send(msg)
}

// SendFile 发送文件
func (bot *Bot) SendFile(f io.Reader, filename string, fileLength int64) (err error) {
	ret, err := bot.UploadMedia(NormalFile, f, filename, fileLength)
	if err != nil {
		return err
	}

	var msg FileMessage
	msg.MsgType = "file"
	msg.File.MediaID = ret.MediaID

	return bot.SendFileMessage(&msg)
}

// VoiceMessage 语音类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E8%AF%AD%E9%9F%B3%E7%B1%BB%E5%9E%8B
type VoiceMessage struct {
	// MsgType 必填。语音类型，此时固定为 voice 。
	MsgType string `json:"msgtype"`
	// 消息内容
	Voice struct {
		// MediaID 必填。语音文件id，通过文件上传接口获取。
		MediaID string `json:"media_id"`
	} `json:"voice"`
}

// SendVoiceMessage 发送语音消息
func (bot *Bot) SendVoiceMessage(msg *VoiceMessage) (err error) {
	return bot.send(msg)
}

// SendVoice 发送语音
func (bot *Bot) SendVoice(f io.Reader, filename string, fileLength int64) (err error) {
	ret, err := bot.UploadMedia(VoiceFile, f, filename, fileLength)
	if err != nil {
		return err
	}

	var msg VoiceMessage
	msg.MsgType = "voice"
	msg.Voice.MediaID = ret.MediaID

	return bot.SendVoiceMessage(&msg)
}

// FileType 文件类型
type FileType string

const (
	// NormalFile 普通文件
	NormalFile FileType = "file"
	// VoiceFile 语音文件
	VoiceFile FileType = "voice"
)

func (bot *Bot) getUploadMediaURL(tpe FileType) string {
	return fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=%s", bot.key, string(tpe))
}

// UploadMedia 文件上传。详见 https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0%E6%8E%A5%E5%8F%A3
func (bot *Bot) UploadMedia(tpe FileType, f io.Reader, filename string, fileLength int64) (*UploadedMedia, error) {
	var reqBody *bytes.Buffer
	if bot.threadSafe {
		reqBody = bytes.NewBuffer(nil)
	} else {
		reqBody = bot.reqbuf
	}
	defer reqBody.Reset()

	writer := multipart.NewWriter(reqBody)

	h := make(textproto.MIMEHeader, 2)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media"; filename="%s"; filelength=%d`, filename, fileLength))
	h.Set("Content-Type", "application/octet-stream")

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, f); err != nil {
		return nil, err
	}
	writer.Close() // finishes the multipart message and writes the trailing boundary end line to the output.

	var resData UploadedMedia
	if err = bot.doPost(bot.getUploadMediaURL(tpe), map[string]string{"Content-Type": writer.FormDataContentType()}, reqBody, &resData); err != nil {
		return nil, err
	}
	if err = resData.ToError(); err != nil {
		return nil, err
	}
	return &resData, nil
}

type resData struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (rd *resData) ToError() error {
	if rd.ErrCode == 0 {
		return nil
	}
	return NewResError(rd.ErrCode, rd.ErrMsg)
}

// UploadedMedia 上传媒体文件结果
type UploadedMedia struct {
	resData
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func (um *UploadedMedia) ToError() error {
	if um.ErrCode == 0 {
		return nil
	}
	return NewResError(um.ErrCode, um.ErrMsg)
}

// ResError 响应错误
type ResError struct {
	errCode int
	errMsg  string
}

// NewResError 返回响应错误实例
func NewResError(code int, msg string) error {
	return &ResError{
		errCode: code,
		errMsg:  msg,
	}
}

// ErrCode 返回错误码
func (e *ResError) ErrCode() int {
	return e.errCode
}

// ErrMsg 返回错误消息内容
func (e *ResError) ErrMsg() string {
	return e.errMsg
}

// Error 返回文本形式的错误描述
func (e *ResError) Error() string {
	return fmt.Sprintf("[%d]%s", e.errCode, e.errMsg)
}
