package utils

import (
	"fmt"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/websocket"
)

var Clients = make(map[string][]*websocket.Conn)
var Mutex = &sync.Mutex{}

func WebsocketAddClient(room string, conn *websocket.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	fmt.Println("New websocket client connected to room:", room)
	Clients[room] = append(Clients[room], conn)
}

func WebsocketRemoveClient(room string, conn *websocket.Conn) {
	Mutex.Lock()
	defer Mutex.Unlock()
	fmt.Println("Websocket client disconnected from room:", room)

	if conns, ok := Clients[room]; ok {
		for i, c := range conns {
			if c == conn {
				Clients[room] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(Clients[room]) == 0 {
			delete(Clients, room)
		}
	}
}

func WebsocketBroadcast(room string, payload any) {
	Mutex.Lock()
	defer Mutex.Unlock()

	conns, ok := Clients[room]
	if !ok {
		return
	}

	data, err := sonic.Marshal(payload)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	for i := 0; i < len(conns); {
		conn := conns[i]
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("write error:", err)
			conn.Close()
			// remove broken conn
			conns = append(conns[:i], conns[i+1:]...)
		} else {
			i++
		}
	}
	Clients[room] = conns
}
