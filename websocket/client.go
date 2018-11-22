package websocket

import (
	"fmt"
	"io"
	"log"

	"github.com/bandit/blockchain-core"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *core.Message
	doneCh chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan *core.Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, ws, server, ch, doneCh}
}

func (self *Client) Conn() *websocket.Conn {
	return self.ws
}

func (self *Client) Write(msg *core.Message) {
	select {
	case self.ch <- msg:
	default:
		self.server.Del(self)
		err := fmt.Errorf("client %d is disconnected.", self.id)
		self.server.Err(err)
	}
}

func (self *Client) Done() {
	self.doneCh <- true
}

func (self *Client) Listen() {
	go self.listenWrite()
	self.listenRead()
}

func (self *Client) listenWrite() {
	log.Println("Listening write to client")

	for {
		select {
		case msg := <-self.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(self.ws, msg)

		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true

			return
		}
	}
}

func (self *Client) listenRead() {
	log.Println("Listening read from client")

	for {
		select {
		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true
			return

		default:
			var msg core.Message
			err := websocket.JSON.Receive(self.ws, &msg)
			if err == io.EOF {
				self.doneCh <- true
			} else if err != nil {
				self.server.Err(err)
			} else {
				self.server.Msg(&msg)
			}
		}
	}
}
