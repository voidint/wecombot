package wecombot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
)

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
