package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/googollee/go-socket.io"
)

var (
	Conn = make(map[int]socketio.Conn, 0)
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Println("ERROR")
	}
	br := socketio.NewBroadcast()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())

		server.ClearRoom("/", s.ID())
		ok := server.JoinRoom("/", "chat", s)
		if !ok {
			log.Println("Ошибка присоединения к комнате")
		}
		//
		//id, err := strconv.Atoi(s.ID())
		//if err != nil {
		//	panic(err)
		//}
		//Conn[id] = s


		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		br.Send("/", "chat", "HELLO")
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
