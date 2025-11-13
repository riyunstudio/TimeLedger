package ws

import (
	"akali/app"
	"akali/libs/logs/ws"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	app       *app.App
	server    *http.Server
	router    *gin.Engine
	upgrader  websocket.Upgrader
	clients   map[string]*clientInfo
	clientsMu sync.Mutex
	closeChan chan struct{}
	wg        sync.WaitGroup
}

type clientInfo struct {
	Uuid       string
	Conn       *websocket.Conn
	LastActive time.Time
	IP         string
}

func InitializeServer(app *app.App) *WebSocketServer {
	w := &WebSocketServer{
		app: app,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		router:    gin.New(),
		clients:   make(map[string]*clientInfo),
		closeChan: make(chan struct{}),
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.Discard

	w.router.GET("/healthy", func(c *gin.Context) {
		c.String(http.StatusOK, "Healthy")
	})
	w.router.GET("/ws", w.handleConnection)

	w.server = &http.Server{
		Addr:    ":" + app.Env.WsServerPort,
		Handler: w.router,
	}

	return w
}

// 啟動 WebSocket Server
func (w *WebSocketServer) Start() {
	// 監控心跳
	go w.monitorHeartbeats()

	// 監聽Http服務
	log.Println("Websocket server started")
	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Websocket ListenAndServe error: %v", err))
		}
	}()
}

// 心跳監控
func (w *WebSocketServer) monitorHeartbeats() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			w.clientsMu.Lock()
			for uuid, c := range w.clients {
				if now.Sub(c.LastActive) > 60*time.Second {
					ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_HEART_TIMEOUT).
						SetClientIP(c.IP).
						PrintInfo("Client disconnected due to heartbeat timeout")
					_ = c.Conn.Close()
					delete(w.clients, uuid)
				}
			}
			w.clientsMu.Unlock()
		case <-w.closeChan:
			return
		}
	}
}

// 處理新的連線
func (w *WebSocketServer) handleConnection(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		// 沒有 uuid，直接拒絕
		http.Error(c.Writer, "missing required parameter: uuid", http.StatusBadRequest)
		return
	}

	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_UPGRADE_ERR).
			SetUuid(uuid).SetClientIP(c.ClientIP()).SetError(err).
			PrintError("Upgrade error")
		return
	}

	client := &clientInfo{
		Uuid:       uuid,
		IP:         c.ClientIP(),
		Conn:       conn,
		LastActive: time.Now(),
	}

	w.clientsMu.Lock()
	w.clients[uuid] = client
	w.clientsMu.Unlock()

	ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_CLIENT_CONN).
		SetUuid(uuid).SetClientIP(conn.RemoteAddr().String()).SetExtraInfo(map[string]any{"clientCount": len(w.clients)}).
		PrintInfo("Client connected")

	w.wg.Add(1)
	go w.handleClient(client)
}

func (w *WebSocketServer) handleClient(client *clientInfo) {
	defer func() {
		w.removeClient(client)
		w.wg.Done()
	}()

	conn := client.Conn
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		select {
		case <-w.closeChan:
			// Server 停掉，退出 goroutine
			return
		default:
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 如果是 server 主動關閉連線，ReadMessage 會返回錯誤
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_CLIENT_DIS_CONN).
					SetUuid(client.Uuid).SetClientIP(client.IP).
					PrintInfo("Client disconnected (server closed)")
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_HEART_TIMEOUT).
					SetUuid(client.Uuid).SetClientIP(client.IP).
					PrintInfo("Client heartbeat timeout")
			} else {
				ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_READ_ERR).
					SetUuid(client.Uuid).SetClientIP(client.IP).SetError(err).
					PrintError("Client read message error")
			}
			break
		}

		message := string(msg)

		if message == "ping" {
			client.LastActive = time.Now()
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			_ = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			continue
		}

		// client 發送訊息，只回應該 client
		reply := "Received: " + message
		if err := conn.WriteMessage(websocket.TextMessage, []byte(reply)); err != nil {
			ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_BROADCAST_ERR).
				SetUuid(client.Uuid).SetClientIP(client.IP).SetError(err).
				PrintError("Failed to send message, client removed")
			break
		}
	}
}

// 移除 client
func (w *WebSocketServer) removeClient(client *clientInfo) {
	w.clientsMu.Lock()
	delete(w.clients, client.Uuid)
	w.clientsMu.Unlock()
	_ = client.Conn.Close()

	ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_CLIENT_DIS_CONN).
		SetUuid(client.Uuid).SetClientIP(client.IP).SetExtraInfo(map[string]any{"clientCount": len(w.clients)}).
		PrintInfo("Client disconnected")
}

// 廣播給所有 client
func (w *WebSocketServer) BroadcastAll(message string) {
	w.clientsMu.Lock()
	defer w.clientsMu.Unlock()

	for _, c := range w.clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

// 廣播給指定 UUID 的 client
func (w *WebSocketServer) BroadcastTo(uuids []string, message string) {
	if len(uuids) == 0 {
		return
	}

	w.clientsMu.Lock()
	defer w.clientsMu.Unlock()

	for _, uuid := range uuids {
		if c, ok := w.clients[uuid]; ok {
			_ = c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

func (w *WebSocketServer) Stop() {
	log.Println("Websocket stopping server...")

	// 先關閉 HTTP Server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := w.server.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
		ws.WsLogInit().SetTopic(ws.TOPIC_SRV).SetEvent(ws.EVENT_SRV_SHUTDOWN_ERR).SetError(err).
			PrintError("Failed to shutdown server")
	}

	// 關閉心跳監控 goroutine
	close(w.closeChan)

	// 關閉所有 client 連線
	w.clientsMu.Lock()
	for uuid, c := range w.clients {
		_ = c.Conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"))
		_ = c.Conn.Close()
		delete(w.clients, uuid)
	}
	w.clientsMu.Unlock()

	// 等待 handleClient goroutine 完成
	w.wg.Wait()

	log.Println("Websocket server stopped gracefully")
}
