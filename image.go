package wecombot

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// ImageMessage 图片类型消息。详见 https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B
type ImageMessage struct {
	// MsgType 必填。消息类型，此时固定为 image 。
	MsgType MsgType `json:"msgtype"`
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
	msg.MsgType = ImageMsgType
	return bot.send(msg)
}

// SendImage 发送图片消息
func (bot *Bot) SendImage(f io.Reader) (err error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return bot.SendImageBytes(data)
}

// SendImageBytes 发送图片消息
func (bot *Bot) SendImageBytes(img []byte) (err error) {
	sum := md5.Sum(img)

	var msg ImageMessage
	msg.Image.Md5 = hex.EncodeToString(sum[:])
	msg.Image.Base64 = base64.StdEncoding.EncodeToString(img)

	return bot.SendImageMessage(&msg)
}
