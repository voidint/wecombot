# wecombot
[![PkgGoDev](https://pkg.go.dev/badge/github.com/voidint/wecombot)](https://pkg.go.dev/github.com/voidint/wecombot?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/voidint/wecombot)](https://goreportcard.com/report/github.com/voidint/wecombot)
[![codebeat badge](https://codebeat.co/badges/f6d30cce-7c65-4d72-a698-40fbc32eda9d)](https://codebeat.co/projects/github-com-voidint-wecombot-main)

wecombot 是[企业微信群机器人](https://developer.work.weixin.qq.com/document/path/91770)的一个 go 语言 sdk，API友好，无任何第三方依赖。

## 安装
```sehll
$ go get github.com/voidint/wecombot@latest
```


## 用法
### 发送[文本](https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E6%9C%AC%E7%B1%BB%E5%9E%8B)类型消息
```go
package main

import "github.com/voidint/wecombot"

func main() {
	bot := wecombot.NewBot("YOUR_KEY")
	bot.SendText("hello 世界！") // 最长不超过2048个字节
	bot.SendText(
		"你好 world!",
		wecombot.WithMentionedMobileList("186xxxx1234", "170xxxx9876"), // 发送文字内容并@某些手机用户
	)
}
```

### 发送[Markdown](https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B)类型消息
```go
package main

import "github.com/voidint/wecombot"

var md = `
# 成绩单
**姓名：**张三
**数学：**<font color="info">97</font>
**语文：**<font color="comment">72</font>
**英语：**<font color="warning">61</font>
`

func main() {
	wecombot.NewBot("YOUR_KEY").SendMarkdown(md) // 最长不超过4096个字节
}
```

**注意**：企业微信群消息[仅支持有限的 markdown 语法](https://developer.work.weixin.qq.com/document/path/91770#markdown%E7%B1%BB%E5%9E%8B)


### 发送[图片](https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E7%89%87%E7%B1%BB%E5%9E%8B)类型消息
```go
package main

import (
	"log"
	"os"

	"github.com/voidint/wecombot"
)

func main() {
	f, err := os.ReadFile("logo.jpg") // 图片最大不能超过2M，支持JPG,PNG格式。
	if err != nil {
		log.Fatal(err)
	}

	wecombot.NewBot("YOUR_KEY").SendImage(f)
}
```

### 发送[图文](https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E6%96%87%E7%B1%BB%E5%9E%8B)类型消息
```go
package main

import (
	"github.com/voidint/wecombot"
)

func takePointer(s string) *string {
	return &s
}

func main() {
	wecombot.NewBot("YOUR_KEY").SendNews(&wecombot.Article{
		Title:       "中秋节礼品领取",
		Description: takePointer("今年中秋节公司有豪礼相送"), // 可选字段使用指针类型
		URL:         "www.qq.com",
		PicURL:      takePointer("http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png"),
	})
}
```

### 发送[文件](https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E4%BB%B6%E7%B1%BB%E5%9E%8B)类型消息

```go
package main

import (
	"log"
	"os"

	"github.com/voidint/wecombot"
)

func main() {
	data, err := os.ReadFile("学生成绩单.xlsx") // 文件大小不超过20M
	if err != nil {
		log.Fatal(err)
	}

	wecombot.NewBot("YOUR_KEY").SendFile(data, "学生成绩单.xlsx")
}
```


### 发送[语音](https://developer.work.weixin.qq.com/document/path/91770#%E8%AF%AD%E9%9F%B3%E7%B1%BB%E5%9E%8B)类型消息

```go
package main

import (
	"log"
	"os"

	"github.com/voidint/wecombot"
)

func main() {
	data, err := os.ReadFile("生日祝福.amr") // 文件大小不超过2M，播放长度不超过60s，仅支持AMR格式。
	if err != nil {
		log.Fatal(err)
	}

	wecombot.NewBot("YOUR_KEY").SendVoice(data, "生日祝福.amr")
}
```

### 发送[文本通知模板卡片](https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E6%9C%AC%E9%80%9A%E7%9F%A5%E6%A8%A1%E7%89%88%E5%8D%A1%E7%89%87)类型消息

```go
package main

import (
	"github.com/voidint/wecombot"
)

func takePointer(s string) *string {
	return &s
}

func main() {
	var msg wecombot.TextNoticeTemplateCardMessage
	msg.TemplateCard.MainTitle.Title = takePointer("欢迎使用企业微信")
	msg.TemplateCard.MainTitle.Desc = takePointer("您的好友正在邀请您加入企业微信")
	msg.TemplateCard.CardAction.Type = 1
	msg.TemplateCard.CardAction.URL = takePointer("https://work.weixin.qq.com/?from=openApi")
	// 参数较多，详见文档。

	wecombot.NewBot("YOUR_KEY").SendTextNoticeTemplateCardMessage(&msg)
}
```

### 发送[图文展示模板卡片](https://developer.work.weixin.qq.com/document/path/91770#%E5%9B%BE%E6%96%87%E5%B1%95%E7%A4%BA%E6%A8%A1%E7%89%88%E5%8D%A1%E7%89%87)类型消息

```go
package main

import (
	"github.com/voidint/wecombot"
)

func takePointer(s string) *string {
	return &s
}

func main() {
	var msg wecombot.NewsNoticeTemplateCardMessage
	msg.TemplateCard.MainTitle.Title = takePointer("欢迎使用企业微信")
	msg.TemplateCard.MainTitle.Desc = takePointer("您的好友正在邀请您加入企业微信")
	msg.TemplateCard.CardImage.URL = *takePointer("https://wework.qpic.cn/wwpic/354393_4zpkKXd7SrGMvfg_1629280616/0")
	msg.TemplateCard.CardAction.Type = 1
	msg.TemplateCard.CardAction.URL = takePointer("https://work.weixin.qq.com/?from=openApi")
	// 参数较多，详见文档。

	wecombot.NewBot("YOUR_KEY").SendNewsNoticeTemplateCardMessage(&msg)
}
```

### [文件上传](https://developer.work.weixin.qq.com/document/path/91770#%E6%96%87%E4%BB%B6%E4%B8%8A%E4%BC%A0%E6%8E%A5%E5%8F%A3)

```go
package main

import (
	"log"
	"os"

	"github.com/voidint/wecombot"
)

func main() {
	data, err := os.ReadFile("学生成绩单.xlsx") // 文件大小不超过20M
	if err != nil {
		log.Fatal(err)
	}

	mediaRes, err := wecombot.NewBot("YOUR_KEY").UploadMedia(wecombot.NormalFile, data, "学生成绩单.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("media id: %s", mediaRes.MediaID)
}
```