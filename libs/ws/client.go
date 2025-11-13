package ws

import (
	"akali/libs/logs/ws"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	URL          string
	conn         *websocket.Conn
	connMu       sync.Mutex
	closeChan    chan struct{}
	wg           sync.WaitGroup
	PingInterval time.Duration
	OnMessage    func(message string)
	OnError      func(err error)
}

// 初始化 client
func InitializeClient(url string) *WebSocketClient {
	return &WebSocketClient{
		URL:          url,
		closeChan:    make(chan struct{}),
		PingInterval: 20 * time.Second, // 預設 20 秒發一次 ping
	}
}

// 啟動 client
func (w *WebSocketClient) Start() error {
	var err error
	w.conn, _, err = websocket.DefaultDialer.Dial(w.URL, nil)
	if err != nil {
		return err
	}

	// 接收訊息
	w.wg.Add(1)
	go w.readLoop()

	// 發送心跳
	w.wg.Add(1)
	go w.pingLoop()

	return nil
}

// 接收訊息
func (w *WebSocketClient) readLoop() {
	defer w.wg.Done()
	for {
		select {
		case <-w.closeChan:
			return
		default:
		}

		w.connMu.Lock()
		conn := w.conn
		w.connMu.Unlock()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			if w.OnError != nil {
				w.OnError(err)
			} else {
				ws.WsLogInit().SetTopic(ws.TOPIC_CLI).SetEvent(ws.EVENT_CLI_READ_ERR).SetError(err).PrintError("Read message error")
			}
			return
		}

		message := string(msg)
		switch message {
		case "pong":
			// 心跳回覆
			// log.Println("Received pong")
		default:
			if w.OnMessage != nil {
				w.OnMessage(message)
			} else {
				// log.Println("Received message:", message)
			}
		}
	}
}

// 發送心跳
func (w *WebSocketClient) pingLoop() {
	defer w.wg.Done()
	ticker := time.NewTicker(w.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.Send("ping")
			if err != nil {
				if w.OnError != nil {
					w.OnError(err)
				} else {
					ws.WsLogInit().SetTopic(ws.TOPIC_CLI).SetEvent(ws.EVENT_CLI_PING_ERR).SetError(err).PrintError("Ping error")
				}
				return
			}
			// log.Println("Sent ping")
		case <-w.closeChan:
			return
		}
	}
}

// 發送訊息
func (w *WebSocketClient) Send(msg string) error {
	w.connMu.Lock()
	defer w.connMu.Unlock()

	if w.conn == nil {
		return fmt.Errorf("Send fail, ws connection is nil")
	}
	return w.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

// 停止 client
func (w *WebSocketClient) Stop() {
	close(w.closeChan)
	w.connMu.Lock()
	if w.conn != nil {
		_ = w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
		_ = w.conn.Close()
		w.conn = nil
	}
	w.connMu.Unlock()
	w.wg.Wait()
	log.Println("WebSocket client stopped gracefully")
}
