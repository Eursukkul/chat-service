package handle

import (

	"github.com/gofiber/websocket/v2"
)

var (
	curr *websocket.Conn
	hubs *Hubs
)

type Hubs struct {
	hubs map[string]*Hub
	run  chan *Hub
	stop chan *Hub
}
type Hub struct {
	chat string

	current chan *websocket.Conn

	clients map[*websocket.Conn]bool

	broadcast chan []byte

	register chan *websocket.Conn

	unregister chan *websocket.Conn

	running    chan bool
}

func getHubRun() *Hubs {
	return &Hubs{
		hubs: make(map[string]*Hub),
		run:  make(chan *Hub),
		stop: make(chan *Hub),
	}
}

func HubRunner() {
	hubs = getHubRun()
	for {
		select {
		case hub := <-hubs.run:
			println("Starting hub", hub)
			go hub.Run()

		case hub := <-hubs.stop:
			println("Stopping hub", hub)
			hub.running <- false
			delete(hubs.hubs, hub.chat)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case connection := <-h.register:
			h.clients[connection] = true

		case message := <-h.broadcast:
			for connection := range h.clients {
				if curr != connection {
					if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						h.unregister <- connection
						connection.WriteMessage(websocket.CloseMessage, []byte("Error Occured!"))
						connection.Close()
					}
				}
			}

		case connection := <-h.unregister:
			delete(h.clients, connection)
			if len(h.clients) == 0 {
				println("initiating to stop hub")
				hubs.stop <- h
			}

		case curr := <-h.running:
			if !curr {
				return
			}

		case connection := <-h.current:
			curr = connection

		}
	}
}

func newHub(chatName string) *Hub {
	return &Hub{
		chat:       chatName,
		current:    make(chan *websocket.Conn),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
		running:    make(chan bool),
	}
}

func getCurrHub(chat string) *Hub {
	if hub, ok := hubs.hubs[chat]; ok {
		return hub
	} else {
		hub := newHub(chat)
		hubs.hubs[chat] = hub
		hubs.run <- hub
		return hub
	}
}