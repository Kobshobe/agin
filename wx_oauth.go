package agin

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
	"time"
)

// oauth websocket connection pool
type ConnectionPool struct {
	Pool   map[string]*Connection
	uuidGenerator SceneID
}

// create connection pool
func NewConnectionPool() ConnectionPool {
	//fmt.Println("create connection Pool")
	return ConnectionPool{
		Pool:   make(map[string]*Connection),
		uuidGenerator: SceneID{lastID: 1000},
	}
}


// create new connection and add to pool
func (cp *ConnectionPool) AddConnection(wsConn *websocket.Conn, mode string) (conn *Connection, err error) {
	if len(cp.Pool) > 80 {
		err = errors.New("too mc")
		return
	}

	id := cp.uuidGenerator.GetUID()
	if G.System.Mode == "test" {
		id = "qq"
		if mode == "admin" {
			id = "ww"
		}
	}

	conn, err = InitConnection(wsConn, id, mode)
	if err != nil {
		return
	}
	cp.Pool[conn.uuid] = conn
	return
}

// get connection from pool
func (cp *ConnectionPool) GetConnection(id string) (conn *Connection, ok bool) {
	conn, ok = cp.Pool[id]
	return
}


// scene uuid generator
type SceneID struct {
	lastID int
	mutex    sync.Mutex
}


// get scene uuid
func (s *SceneID) GetUID() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.lastID ++
	return fmt.Sprintf("%d%d", rand.New(rand.NewSource(time.Now().Unix())).Uint64(), s.lastID)
}


// oauth websocket connection
type Connection struct {
	wsConn          *websocket.Conn
	inChan          chan []byte
	outChan         chan []byte
	allowLoginChan  chan string
	loginResultChan chan string // get msg went allow login
	closeChan       chan byte
	uuid            string
	token           string
	Openid          string
	mode            string // admin ''

	mutex    sync.Mutex
	isClosed bool
}

func InitConnection(wsConn *websocket.Conn, id string, mode string) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:          wsConn,
		inChan:          make(chan []byte, 1),
		outChan:         make(chan []byte, 1),
		allowLoginChan:  make(chan string, 1),
		loginResultChan: make(chan string, 1),
		closeChan:       make(chan byte, 1),
		uuid:            id,
		mode:            mode,
	}

	go conn.readLoop()
	go conn.writeLoop()

	go func() {
		time.Sleep(time.Second * 80)
		conn.Close()
	}()

	return
}

func (conn *Connection) getQrScene() (scene string) {
	if conn.mode == "admin" {
		scene = conn.uuid + "@admin"
	} else {
		scene = conn.uuid
	}
	return
}

func (conn *Connection) SetToken(openid string) (err error) {
	if conn.mode == "admin" {
		conn.token, err = G.WxApp.NewJwtToken(openid, fmt.Sprintf("%s:%s", conn.mode ,G.WxApp.AdminTokenVersion))
	} else {
		conn.token, err = G.WxApp.NewJwtToken(openid, fmt.Sprintf("%s:%s", conn.mode ,G.WxApp.TokenVersion))
	}
	return
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
		if string(data) == "login" {
			err = conn.WriteMessage(data)
		} else if string(data) == "loginOk" {
			conn.loginResultChan <- "loginOk"
		}
	case msg := <-conn.allowLoginChan:
		if msg == "allow" {
			conn.outChan <- []byte(fmt.Sprintf(`{"token":"%s","liveTime":%f,"openid":"%s"}`, conn.token, G.WxApp.GetJwtLife().Hours(), conn.Openid))
		}
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	var (
		buf []byte
	)
	if string(data) == "login" {
		//fmt.Println("write qr")
		if buf, err = G.WxApp.GetQRFromWX(conn.getQrScene()); err != nil {
			err = errors.New("get buf err")
		}
		//fmt.Println("write qr")
		conn.outChan <- buf
	} else if string(data) == "allow" {
		//fmt.Println("login ok")
		conn.outChan <- []byte("login ok...")
	} else {
		conn.outChan <- data
	}
	return
}

func (conn *Connection) Close() {
	_ = conn.wsConn.Close()

	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
		delete(G.WxApp.OAuthConnPool.Pool, conn.uuid)
	}
}

// read msg from websocket
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			//fmt.Println("ReadMessage err: ", err)
			_ = conn.wsConn.Close()
			goto ERR
		}
		select {
		case conn.inChan <- data:
		case <-conn.closeChan:
			goto ERR
		}
		conn.inChan <- data
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)
	for {
		select {
		case data = <-conn.outChan:
			if err = conn.wsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				goto ERR
			}
		case <-conn.closeChan:
			goto ERR
		}
		data = <-conn.outChan

	}
ERR:
	conn.Close()
}
