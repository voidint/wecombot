package wecombot

import "io"

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
