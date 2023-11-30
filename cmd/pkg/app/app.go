package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrdr = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var cli []websocket.Conn

func Program() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrdr.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		cli = append(cli, *conn)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))
			for _, client := range cli {
				if err = client.WriteMessage(msgType, msg); err != nil {
					return
				}
			}
		}
	})
}
