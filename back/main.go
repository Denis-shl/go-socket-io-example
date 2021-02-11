package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server, _ := socketio.NewServer(nil)

	br := socketio.NewBroadcast()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		server.JoinRoom("/chat", "bcast", s)
		br.Join("bcast", s)
		return nil
	})

	server.OnEvent("/", "chat", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		fmt.Println(s.ID() + "[" + s.Namespace() + "]: " + msg)
		fmt.Println(s.Rooms())
		fmt.Println()

		server.BroadcastToRoom("/", "bcast", "chat", msg)

		br.Send("bcast", "chat", msg)
		return "recv " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	//http.Handle("/", http.FileServer(http.Dir("./../front")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
