package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var InsQWRobot *qwRobotConfig

// qwRobotConfig 企业微信机器人客户端
type qwRobotConfig struct {
	webhookKey string
	httpClient *http.Client
	baseURL    string
}

// InitQWRobotClient 创建新的企业微信机器人客户端,详细使用查看 examples/qywx_test.go
func InitQWRobotClient(webhookKey string) {
	InsQWRobot = &qwRobotConfig{
		webhookKey: webhookKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook",
	}
}

// SetTimeout 设置HTTP客户端超时时间
func (c *qwRobotConfig) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetHTTPClient 设置自定义HTTP客户端
func (c *qwRobotConfig) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// QYWXMediaType 媒体文件类型
type QYWXMediaType string

const (
	// QYWXMediaTypeVoice 语音
	QYWXMediaTypeVoice QYWXMediaType = "voice"
	// QYWXMediaTypeFile 文件
	QYWXMediaTypeFile QYWXMediaType = "file"
)

// QYWXUploadResponse 上传媒体文件响应
type QYWXUploadResponse struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// QYWXSendResponse 发送消息响应
type QYWXSendResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// UploadMedia 上传媒体文件,mediaType有file文件和voice语音
func (c *qwRobotConfig) UploadMedia(reader io.Reader, filename string, mediaType QYWXMediaType) (*QYWXUploadResponse, error) {

	// 读取所有数据到内存中，以便进行内容类型检测和获取大小
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}

	// 检测内容类型
	contentType := http.DetectContentType(data)
	dataSize := len(data)

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	// 写入文件内容
	if _, err = part.Write(data); err != nil {
		return nil, fmt.Errorf("写入文件内容失败: %w", err)
	}

	// 添加其他字段
	writer.WriteField("filename", filename)
	writer.WriteField("filelength", fmt.Sprintf("%d", dataSize))
	writer.WriteField("content-type", contentType)

	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("关闭multipart writer失败: %w", err)
	}

	// 构建URL
	url := fmt.Sprintf("%s/upload_media?key=%s&type=%s", c.baseURL, c.webhookKey, mediaType)

	// 创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("上传失败, 状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var uploadResp QYWXUploadResponse
	if err = json.Unmarshal(respBody, &uploadResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if uploadResp.Errcode != 0 {
		return nil, fmt.Errorf("上传失败: %s", uploadResp.Errmsg)
	}

	return &uploadResp, nil
}

// UploadMediaFromFile 从文件路径上传媒体文件的便捷方法
func (c *qwRobotConfig) UploadMediaFromFile(mediaType QYWXMediaType, filePath string) (*QYWXUploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	filename := filepath.Base(filePath)
	return c.UploadMedia(file, filename, mediaType)
}

// QYWXMessage 消息接口
type QYWXMessage interface {
	MessageType() string
}

// QYWXTextMessage 文本消息
type QYWXTextMessage struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content             string   `json:"content"`
		MentionedList       []string `json:"mentioned_list,omitempty"`
		MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
	} `json:"text"`
}

func (t QYWXTextMessage) MessageType() string { return "text" }

// NewQYWXTextMessage 创建文本消息
func NewQYWXTextMessage(content string) *QYWXTextMessage {
	msg := &QYWXTextMessage{Msgtype: "text"}
	msg.Text.Content = content
	return msg
}

// AddMention 添加@用户
func (t *QYWXTextMessage) AddMention(userIDs ...string) *QYWXTextMessage {
	t.Text.MentionedList = append(t.Text.MentionedList, userIDs...)
	return t
}

// AddMentionMobile 添加@手机号
func (t *QYWXTextMessage) AddMentionMobile(mobiles ...string) *QYWXTextMessage {
	t.Text.MentionedMobileList = append(t.Text.MentionedMobileList, mobiles...)
	return t
}

// QYWXMarkdownMessage Markdown消息
type QYWXMarkdownMessage struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

func (m QYWXMarkdownMessage) MessageType() string { return "markdown" }

// NewQYWXMarkdownMessage 创建Markdown消息
func NewQYWXMarkdownMessage(content string) *QYWXMarkdownMessage {
	msg := &QYWXMarkdownMessage{Msgtype: "markdown"}
	msg.Markdown.Content = content
	return msg
}

// QYWXImageMessage 图片消息
type QYWXImageMessage struct {
	Msgtype string `json:"msgtype"`
	Image   struct {
		Base64 string `json:"base64,omitempty"`
		MD5    string `json:"md5,omitempty"`
	} `json:"image"`
}

func (i QYWXImageMessage) MessageType() string { return "image" }

// QYWXFileMessage 文件消息
type QYWXFileMessage struct {
	Msgtype string `json:"msgtype"`
	File    struct {
		MediaID string `json:"media_id"`
	} `json:"file"`
}

func (f QYWXFileMessage) MessageType() string { return "file" }

// NewQYWXFileMessage 创建文件消息
func NewQYWXFileMessage(mediaID string) *QYWXFileMessage {
	msg := &QYWXFileMessage{Msgtype: "file"}
	msg.File.MediaID = mediaID
	return msg
}

// QYWXNewsMessage 图文消息
type QYWXNewsMessage struct {
	Msgtype string `json:"msgtype"`
	News    struct {
		Articles []QYWXArticle `json:"articles"`
	} `json:"news"`
}

func (n QYWXNewsMessage) MessageType() string { return "news" }

// QYWXArticle 图文消息文章
type QYWXArticle struct {
	// 标题
	Title string `json:"title"`
	// 描述
	Description string `json:"description"`
	// 点击跳转的url
	URL string `json:"url"`
	// 聊天框预览的url
	PicURL string `json:"picurl"`
}

// NewQYWXNewsMessage 创建图文消息
func NewQYWXNewsMessage(articles ...QYWXArticle) *QYWXNewsMessage {
	msg := &QYWXNewsMessage{Msgtype: "news"}
	msg.News.Articles = articles
	return msg
}

// AddArticle 添加文章
func (n *QYWXNewsMessage) AddArticle(title, description, url, picURL string) *QYWXNewsMessage {
	article := QYWXArticle{
		Title:       title,
		Description: description,
		URL:         url,
		PicURL:      picURL,
	}
	n.News.Articles = append(n.News.Articles, article)
	return n
}

// SendMessage 发送消息
func (c *qwRobotConfig) SendMessage(msg QYWXMessage) (*QYWXSendResponse, error) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("序列化消息失败: %w", err)
	}

	url := fmt.Sprintf("%s/send?key=%s", c.baseURL, c.webhookKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("发送失败, 状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var sendResp QYWXSendResponse
	if err = json.Unmarshal(respBody, &sendResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if sendResp.Errcode != 0 {
		return nil, fmt.Errorf("发送失败: %s", sendResp.Errmsg)
	}

	return &sendResp, nil
}

// SendText 发送文本消息的便捷方法
func (c *qwRobotConfig) SendText(content string) (*QYWXSendResponse, error) {
	return c.SendMessage(NewQYWXTextMessage(content))
}

// SendMarkdown 发送Markdown消息的便捷方法
func (c *qwRobotConfig) SendMarkdown(content string) (*QYWXSendResponse, error) {
	return c.SendMessage(NewQYWXMarkdownMessage(content))
}

// SendFile 发送文件消息的便捷方法
func (c *qwRobotConfig) SendFile(filePath string, mediaType QYWXMediaType) (*QYWXSendResponse, error) {
	// 上传文件
	uploadResp, err := c.UploadMediaFromFile(mediaType, filePath)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 发送文件消息
	return c.SendMessage(NewQYWXFileMessage(uploadResp.MediaID))
}

// SendFileFromReader 从io.Reader发送文件消息
func (c *qwRobotConfig) SendFileFromReader(mediaType QYWXMediaType, reader io.Reader, filename string) (*QYWXSendResponse, error) {
	// 上传文件
	uploadResp, err := c.UploadMedia(reader, filename, mediaType)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 发送文件消息
	return c.SendMessage(NewQYWXFileMessage(uploadResp.MediaID))
}

// SendNews 发送图文消息的便捷方法
func (c *qwRobotConfig) SendNews(articles ...QYWXArticle) (*QYWXSendResponse, error) {
	return c.SendMessage(NewQYWXNewsMessage(articles...))
}
