package ws

import (
	"akali/app"
	"akali/libs/logs/ws"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	app        *app.App
	upgrader   websocket.Upgrader
	server     *http.Server
	router     *gin.Engine
	clients    map[*websocket.Conn]*clientInfo
	clientsMu  sync.Mutex
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan []byte

	closeOnce sync.Once
	closeChan chan struct{}
	wg        sync.WaitGroup

	// 心跳設置
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration
}

type clientInfo struct {
	conn       *websocket.Conn
	lastActive time.Time
}

func Initialize(app *app.App) *WebSocketServer {
	ws := &WebSocketServer{
		app: app,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 允許所有來源
			},
		},
		router:            gin.New(),
		clients:           make(map[*websocket.Conn]*clientInfo),
		register:          make(chan *websocket.Conn),
		unregister:        make(chan *websocket.Conn),
		broadcast:         make(chan []byte, 128),
		closeChan:         make(chan struct{}),
		heartbeatInterval: 15 * time.Second,
		heartbeatTimeout:  60 * time.Second,
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.Discard

	// 路由
	ws.router.GET("/ws", ws.handleConnection)

	ws.server = &http.Server{
		Addr:    ":" + app.Env.WsServerPort,
		Handler: ws.router,
	}

	return ws
}

func (w *WebSocketServer) Start() {
	log.Println("Websocket server started")

	// 啟動 hub
	w.wg.Go(func() {
		w.runHub()
	})

	// 啟動 HTTP server
	w.wg.Go(func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Websocket ListenAndServe error: %v", err))
		}
	})
}

// WebSocket 連線
func (w *WebSocketServer) handleConnection(c *gin.Context) {
	conn, err := w.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ws.WsLogInit().
			SetClientIP(c.ClientIP()).SetError(err).
			PrintError("Upgrade failed")
		return
	}

	select {
	case w.register <- conn:
	case <-w.closeChan:
		_ = conn.Close()
		return
	}

	w.wg.Go(func() {
		defer func() {
			select {
			case w.unregister <- conn:
			case <-w.closeChan:
			}
			_ = conn.Close()
		}()

		for {
			select {
			case <-w.closeChan:
				return
			default:
			}

			_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			_, msg, err := conn.ReadMessage()
			if err != nil {
				ws.WsLogInit().
					SetEvent("ClientReadError").SetClientIP(c.ClientIP()).SetError(err).
					PrintError("Client disconnected")
				return
			}

			// 更新最後活動時間
			w.clientsMu.Lock()
			if ci, ok := w.clients[conn]; ok {
				ci.lastActive = time.Now()
			}
			w.clientsMu.Unlock()

			select {
			case w.broadcast <- msg:
			case <-w.closeChan:
				return
			default:
				ws.WsLogInit().
					SetEvent("BroadcastDrop").SetClientIP(c.ClientIP()).SetError(err).SetExtraInfo(string(msg)).
					PrintError("Broadcast channel full, dropping message")
			}
		}
	})
}

// runHub 註冊、廣播、斷線、心跳檢查
func (w *WebSocketServer) runHub() {
	ticker := time.NewTicker(w.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case conn := <-w.register:
			w.clientsMu.Lock()
			w.clients[conn] = &clientInfo{
				conn:       conn,
				lastActive: time.Now(),
			}
			count := len(w.clients)
			w.clientsMu.Unlock()

			ws.WsLogInit().
				SetEvent("ClientConnect").SetClientIP(conn.RemoteAddr().String()).SetExtraInfo(map[string]any{"clientCount": count}).
				PrintInfo("Client connected")

		case conn := <-w.unregister:
			w.clientsMu.Lock()
			delete(w.clients, conn)
			count := len(w.clients)
			w.clientsMu.Unlock()

			ws.WsLogInit().
				SetEvent("ClientDisconnect").SetClientIP(conn.RemoteAddr().String()).SetExtraInfo(map[string]any{"clientCount": count}).
				PrintInfo("Client disconnected")

		case msg := <-w.broadcast:
			w.clientsMu.Lock()
			for _, ci := range w.clients {
				_ = ci.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := ci.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					_ = ci.conn.Close()
					delete(w.clients, ci.conn)

					ws.WsLogInit().
						SetEvent("BroadcastError").SetClientIP(ci.conn.RemoteAddr().String()).SetError(err).
						PrintError("Failed to send message, client removed")

				} else {
					ci.lastActive = time.Now()
				}
			}
			w.clientsMu.Unlock()

		case <-ticker.C: // 心跳檢查
			now := time.Now()
			w.clientsMu.Lock()
			for _, ci := range w.clients {
				// 超時就斷線
				if now.Sub(ci.lastActive) > w.heartbeatTimeout {
					_ = ci.conn.Close()
					delete(w.clients, ci.conn)

					ws.WsLogInit().
						SetEvent("HeartbeatTimeout").SetClientIP(ci.conn.RemoteAddr().String()).
						PrintInfo("Client disconnected due to heartbeat timeout")
					continue
				}

				// 向仍在線的客戶端發送心跳訊息
				if err := ci.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					_ = ci.conn.Close()
					delete(w.clients, ci.conn)

					ws.WsLogInit().
						SetEvent("HeartbeatSendError").SetClientIP(ci.conn.RemoteAddr().String()).SetError(err).
						PrintError("Failed to send heartbeat ping")
				}
			}
			w.clientsMu.Unlock()

		case <-w.closeChan:
			w.clientsMu.Lock()
			for _, ci := range w.clients {
				_ = ci.conn.Close()
				delete(w.clients, ci.conn)
			}
			w.clientsMu.Unlock()

			return
		}
	}
}

// 廣播
func (w *WebSocketServer) Broadcast(msg string) {
	select {
	case <-w.closeChan:
		return
	default:
	}

	select {
	case w.broadcast <- []byte(msg):
	case <-w.closeChan:
	}
}

func (w *WebSocketServer) Stop() {
	w.closeOnce.Do(func() {
		log.Println("Websocket stopping server...")

		// 停止 HTTP server（不再接受新連線）
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := w.server.Shutdown(shutdownCtx); err != nil && err != http.ErrServerClosed {
			log.Printf("Websocket server shutdown error: %v", err)

			ws.WsLogInit().
				SetEvent("ServerShutdown").SetError(fmt.Errorf("server shutdown error: %w", err)).
				PrintError("Server stop error")
		}

		// 通知所有 goroutine 停止（停止收訊息/廣播）
		close(w.closeChan)

		// 等待所有 goroutine 完成
		w.wg.Wait()

		log.Println("Websocket server stopped gracefully")
	})
}
