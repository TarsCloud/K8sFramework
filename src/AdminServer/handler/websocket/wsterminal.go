package websocket

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second

	// EndOfTransmission end
	EndOfTransmission = "\u0004"
)

// TerminalMessage is the messaging protocol between ShellController and TerminalSession.
type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// HTTP -> WebSocket协议升级
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 5
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// 全局ID：排查bash开启的进程未正常退出的bug
type terminalSessionID struct {
	tsID int32
	lock sync.Mutex
}

var tsGlobalID *terminalSessionID

func init() {
	tsGlobalID = &terminalSessionID{}
}
func addTerminalSessionID() {
	tsGlobalID.lock.Lock()
	defer tsGlobalID.lock.Unlock()
	if tsGlobalID.tsID == math.MaxInt32 {
		tsGlobalID.tsID = 0
	}
	tsGlobalID.tsID += 1
}
func getTerminalSessionID() int32 {
	tsGlobalID.lock.Lock()
	defer tsGlobalID.lock.Unlock()
	return tsGlobalID.tsID
}

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
	Done()
	GetPtyID() int32
}

// WebSocket Terminal Session对象
type TerminalSession struct {
	wsID     int32
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

func NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	addTerminalSessionID()

	session := &TerminalSession{
		wsID:     getTerminalSessionID(),
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}

	//log.Printf("succ. to new wsID:%d, query:%s", session.wsID, r.URL.Query().Encode())
	return session, nil
}

// Done done, must call Done() before connection close, or Next() would not exits.
func (ts *TerminalSession) Done() {
	close(ts.doneChan)
	//log.Printf("succ. to done wsID:%d", ts.wsID)
}

func (ts *TerminalSession) GetPtyID() int32 {
	return ts.wsID
}

// Next called in a loop from remotecommand as long as the process is running
func (ts *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-ts.sizeChan:
		return &size
	case <-ts.doneChan:
		return nil
	}
}

// Read called in a loop from remotecommand as long as the process is running
func (ts *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := ts.wsConn.ReadMessage()
	if err != nil {
		log.Printf("wsID:%d read message err: %v", ts.wsID, err)
		return copy(p, EndOfTransmission), err
	}
	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		log.Printf("wsID:%d read parse message err: %v", ts.wsID, err)
		return copy(p, EndOfTransmission), err
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		ts.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		log.Printf("wsID:%d unknown message type '%s'", ts.wsID, msg.Operation)
		// return 0, nil
		return copy(p, EndOfTransmission), fmt.Errorf("wsID:%d unknown message type '%s'", ts.wsID, msg.Operation)
	}
}

// Write called from remotecommand whenever there is any output
func (ts *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(TerminalMessage{
		Operation: "stdout",
		Data:      string(p),
	})
	if err != nil {
		log.Printf("wsID:%d write parse message err: %v", ts.wsID, err)
		return 0, err
	}
	if err := ts.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("wsID:%d write message err: %v", ts.wsID, err)
		return 0, err
	}
	return len(p), nil
}

// Close close session
func (ts *TerminalSession) Close() error {
	//log.Printf("succ. to close wsID:%d", ts.wsID)
	return ts.wsConn.Close()
}
