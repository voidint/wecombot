package wecombot

import "io"

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
