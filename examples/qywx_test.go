package examples

import (
	"bytes"
	"fmt"
	"github.com/loveyu233/gb"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	// 点击机器人后显示的Webhook连接的参数中的key的值,例如:https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=就是这个值
	client := gb.InitQywxClient("就是这个值")
	client.SendText("xxxx")

	client.SendMessage(gb.NewQYWXTextMessage("xxxx").AddMention("@all"))

	client.SendMarkdown("# 一号> ## 二号> ## 三号")

	file, err := os.ReadFile("log.txt")
	if err != nil {
		panic(err)
	}
	media, err := client.UploadMedia(bytes.NewBuffer(file), "abc.txt", gb.QYWXMediaTypeFile)
	if err != nil {
		t.Log(err)
		return
	}
	message, err := client.SendMessage(gb.NewQYWXFileMessage(media.MediaID))
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Printf("%+v", message)

	news, err := client.SendNews(gb.QYWXArticle{
		Title:       "title",
		Description: "desc",
		URL:         "点击跳转的url",
		PicURL:      "聊天框显示的url",
	})
	if err != nil {
		t.Log(err)
		return
	}
	fmt.Printf("%+v", news)
}
