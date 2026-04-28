package handle

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"week12/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境允许所有来源，生产环境应限制为指定域名
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan Message
}

type Message struct {
	Type      string    `json:"type"` //chat system join leave image
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Users     []string  `json:"users,omitempty"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var hub = &Hub{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan Message, 256),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".webp": true,
}

func HandleConnections(c *gin.Context) {
	var err error

	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	// 升级HTTP连接为websocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	client := &Client{
		conn:     conn,
		username: username,
		send:     make(chan Message, 256),
	}
	hub.register <- client

	//启动读写goroutine
	go client.writePump()
	client.readPump() //阻塞直到连接关闭
}

func HandleImageUpload(c *gin.Context) {
	// 打印请求信息
	log.Printf("收到图片上传请求:")
	log.Printf("  Method: %s", c.Request.Method)
	log.Printf("  URL: %s", c.Request.URL.String())
	log.Printf("  Content-Type: %s", c.Request.Header.Get("Content-Type"))
	log.Printf("  Content-Length: %s", c.Request.Header.Get("Content-Length"))

	// 尝试解析表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("解析表单失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无法解析表单数据: " + err.Error(),
		})
		return
	}

	// 打印表单字段
	log.Printf("表单字段:")
	for key, values := range c.Request.MultipartForm.Value {
		log.Printf("  %s: %v", key, values)
	}
	for key, files := range c.Request.MultipartForm.File {
		log.Printf("  %s: %d个文件", key, len(files))
		for i, file := range files {
			log.Printf("    文件%d: %s (%d bytes)", i+1, file.Filename, file.Size)
		}
	}

	// 获取用户名（先从表单获取，再从查询参数获取）
	username := c.PostForm("username")
	if username == "" {
		username = c.Query("username")
		log.Printf("从查询参数获取用户名: %s", username)
	} else {
		log.Printf("从表单获取用户名: %s", username)
	}

	if username == "" {
		log.Printf("用户名为空")
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("获取文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择图片文件: " + err.Error()})
		return
	}

	if file.Size <= 0 {
		log.Printf("文件大小为0")
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片文件为空"})
		return
	}

	log.Printf("处理文件: %s, 大小: %d", file.Filename, file.Size)

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageTypes[ext] {
		log.Printf("不支持的文件类型: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":    "不支持的文件类型",
			"allowed":  []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"},
			"received": ext,
		})
		return
	}

	savedPath, err := utils.SaveImages(c, file, username)
	if err != nil {
		log.Printf("保存图片失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败: " + err.Error()})
		return
	}

	urlPath := "/uploads/" + filepath.Base(savedPath)
	log.Printf("图片保存成功: %s", urlPath)

	imgMsg := Message{
		Type:      "image",
		Username:  username,
		ImageURL:  urlPath,
		Text:      file.Filename,
		Timestamp: time.Now(),
	}

	hub.broadcast <- imgMsg

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     urlPath,
		"message": "图片上传成功",
	})
}

func (c *Client) readPump() {
	defer func() {
		// 发送离开通知
		leaveMsg := Message{
			Type:      "system",
			Text:      c.username + " 离开了聊天室",
			Timestamp: time.Now(),
		}
		hub.broadcast <- leaveMsg
		hub.unregister <- c
		c.conn.Close()
	}()

	// 设置读取限制
	c.conn.SetReadLimit(512) //最大信息大小
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.sendOnlineUserList()

	//发送加入通知
	joinMsg := Message{
		Type:      "system",
		Text:      c.username + "加入了聊天室",
		Timestamp: time.Now(),
	}
	hub.broadcast <- joinMsg

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取错误: %v\n", err)
			}
			break
		}
		//补充消息
		msg.Username = c.username
		msg.Timestamp = time.Now()
		if msg.Type == "" {
			msg.Type = "chat"
		}

		//广播消息
		hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second) //心跳间隔
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		//select是尝试从各个管道读取
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			// 发送心跳ping
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendOnlineUserList() {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	usernames := make([]string, 0, len(hub.clients))
	for client := range hub.clients {
		usernames = append(usernames, client.username)
	}
	// 发送用户列表消息（可以选择发送给自己）
	userListMsg := Message{
		Type:      "userlist",
		Text:      "当前在线用户",
		Timestamp: time.Now(),
		Users:     usernames, // 需要添加 Users 字段到 Message 结构体
	}

	// 只发送给当前客户端
	select {
	case c.send <- userListMsg:
	default:
		log.Println("发送用户列表失败，客户端通道已满")
	}
}

func RunHub() {
	for {
		select {
		case client := <-hub.register:
			hub.mu.Lock()
			hub.clients[client] = true
			hub.mu.Unlock()

		case client := <-hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
				log.Println("客户端断开: " + client.username)
			}
			hub.mu.Unlock()

		case msg := <-hub.broadcast:
			var clientsToRemove []*Client

			hub.mu.RLock()
			for client := range hub.clients {
				select {
				case client.send <- msg:
				default:
					//发送失败，标记
					clientsToRemove = append(clientsToRemove, client)
				}
			}
			hub.mu.RUnlock()

			if len(clientsToRemove) > 0 {
				hub.mu.Lock()
				for _, client := range clientsToRemove {
					if _, ok := hub.clients[client]; ok {
						close(client.send)
						delete(hub.clients, client)
						log.Println("移除了客户端: " + client.username)
					}
				}

				hub.mu.Unlock()
			}
		}
	}
}

// GetOnlineCount 获取当前在线人数
func GetOnlineCount() int {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return len(hub.clients)
}
