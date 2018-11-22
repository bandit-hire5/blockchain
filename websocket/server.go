package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bandit/blockchain-core"
	"golang.org/x/net/websocket"
)

type Server struct {
	Pattern string
	Ledger  *core.Ledger
	clients map[int]*Client
	addCh   chan *Client
	delCh   chan *Client
	msgCh   chan *core.Message
	doneCh  chan bool
	errCh   chan error
}

func NewServer() *Server {
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	msgCh := make(chan *core.Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		clients: clients,
		addCh: addCh,
		delCh: delCh,
		msgCh: msgCh,
		doneCh: doneCh,
		errCh: errCh,
	}
}

func (self *Server) Add(client *Client) {
	self.addCh <- client
}

func (self *Server) Del(client *Client) {
	self.delCh <- client
}

func (self *Server) Msg(msg *core.Message) {
	self.msgCh <- msg
}

func (self *Server) Done() {
	self.doneCh <- true
}

func (self *Server) Err(err error) {
	self.errCh <- err
}

func (self *Server) sendAllBlocks(client *Client) {
	list, err := self.Ledger.GetAllBlocks()
	if err != nil {
		client.doneCh <- true
		return
	}

	for _, block := range list {
		client.Write(&core.Message{
			Type: "block",
			Body: block,
		})
	}
}

func (self *Server) sendAll(block *core.Block) error {
	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	for _, client := range self.clients {
		client.Write(&core.Message{
			Type: "block",
			Body: string(data),
		})
	}

	return nil
}

func (self *Server) processMessage(msg *core.Message) error {
	switch msg.Type {
	case "message":
		var data core.Data

		data.Message = msg.Body

		block, err := self.Ledger.NextBlock(data)
		if err != nil {
			return err
		}

		err = self.Ledger.AddBlock(block)
		if err != nil {
			return err
		}

		self.sendAll(block)
	}

	return nil
}

func (self *Server) Listen() {
	log.Println("Listening server...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				self.errCh <- err
			}
		}()

		client := NewClient(ws, self)
		self.Add(client)
		client.Listen()
	}

	http.Handle(self.Pattern, websocket.Handler(onConnected))

	log.Println("Created handler")

	for {
		select {
		case client := <-self.addCh:
			log.Println("Added new client")
			self.clients[client.id] = client
			log.Println("Now", len(self.clients), "clients connected")
			self.sendAllBlocks(client)

		case client := <-self.delCh:
			log.Println("Delete client")
			delete(self.clients, client.id)

		case msg := <-self.msgCh:
			log.Println("Msg:", msg)
			self.processMessage(msg)

		case err := <-self.errCh:
			log.Println("Error:", err.Error())

		case <-self.doneCh:
			return
		}
	}
}
