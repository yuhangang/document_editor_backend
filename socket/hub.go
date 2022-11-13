package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	nextID     int
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		nextID:     0,
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

func (hub *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.register <- client

	go client.read()
	go client.write()
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())

	// Make new client
	client.id = hub.nextID
	hub.nextID++
	client.color = generateColor()
	hub.clients = append(hub.clients, client)

	// Make list of all users
	users := []User{}
	for _, c := range hub.clients {
		users = append(users, User{ID: c.id, Color: c.color})
	}

	// Notify that a user joined
	hub.send(NewConnected(client.color, users), client)
	hub.broadcast(NewUserJoined(client.id, client.color), client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())

	client.close()

	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}
	// Delete client from list
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	hub.broadcast(NewUserLeft(client.id), nil)
}

func (hub *Hub) onMessage(data []byte, client *Client) {
	log.Println("onMessage: ", string(data))

	kind := gjson.GetBytes(data, "kind").Int()
	if kind == KindStroke {
		var msg Stroke
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	} else if kind == KindClear {
		var msg Clear
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	}
}
