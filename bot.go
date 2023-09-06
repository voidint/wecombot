package wecombot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Bot 企业微信群机器人
type Bot struct {
	webhookURL string
	key        string

	threadSafe bool
	reqbuf     *bytes.Buffer
	client     *http.Client
}

// NewBot 返回企业微信群机器人实例
func NewBot(key string, opts ...func(*Bot)) *Bot {
	bot := Bot{
		webhookURL: fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key),
		key:        key,
		client:     http.DefaultClient,
	}

	for _, setter := range opts {
		setter(&bot)
	}

	if !bot.threadSafe {
		bot.reqbuf = bytes.NewBuffer(nil)
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

var jsonReqHeader = map[string]string{
	"Content-Type": "application/json",
}

func (bot *Bot) send(msg interface{}) (err error) {
	var reqBody *bytes.Buffer
	if bot.threadSafe {
		reqBody = bytes.NewBuffer(nil)
	} else {
		reqBody = bot.reqbuf
	}
	defer reqBody.Reset()

	if err = json.NewEncoder(reqBody).Encode(msg); err != nil {
		return err
	}

	var resData resData
	if err = bot.doPost(bot.webhookURL, jsonReqHeader, reqBody, &resData); err != nil {
		return err
	}
	return resData.ToError()
}

// isSuccess 返回 http 请求是否成功
func isSuccess(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func (bot *Bot) doPost(url string, reqHeader map[string]string, reqBody io.Reader, resData interface{}) error {
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

	if !isSuccess(res.StatusCode) {
		return fmt.Errorf("server response abnormality: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(resData)
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

// MsgType 消息类型
type MsgType string

const (
	// TextMsgType 文本类型
	TextMsgType MsgType = "text"
	// MarkdownMsgType Markdown 类型
	MarkdownMsgType MsgType = "markdown"
	// ImageMsgType 图片类型
	ImageMsgType MsgType = "image"
	// NewsMsgType 图文类型
	NewsMsgType MsgType = "news"
	// FileMsgType 文件类型
	FileMsgType MsgType = "file"
	// VoiceMsgType 语音类型
	VoiceMsgType MsgType = "voice"
	// TemplateCardMsgType 模板卡片类型
	TemplateCardMsgType MsgType = "template_card"
)

// ExtractKey 返回 webhook 中包含的 key 信息。若 URL 地址无效，或 URL 中不包含 key 参数，则返回空字符串。
func ExtractKey(webhookURL string) (key string) {
	one, err := url.Parse(webhookURL)
	if err != nil {
		return ""
	}
	return one.Query().Get("key")
}
