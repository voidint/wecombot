package wecombot

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

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

// UploadMediaBytes 文件上传
func (bot *Bot) UploadMediaBytes(tpe FileType, f []byte, filename string, fileLength int64) (*UploadedMedia, error) {
	return bot.UploadMedia(tpe, bytes.NewReader(f), filename, fileLength)
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
