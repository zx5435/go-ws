package main

import (
	//	"fmt"
	"log"
	"net/http"

	"go-socket.io"
	//	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)

	numUsers := 0

	if err != nil {
		log.Fatal(err)
	}
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	server.On("connection", func(so socketio.Socket) {
		username := ""
		log.Println("on connection")

		so.Join("chat")
		so.On("testcallback", func(msg string) map[string]interface{} {
			m := make(map[string]interface{})
			m["id"] = so.Id()
			m["message"] = msg + msg
			return m
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

		so.On("new message", func(data string) {
			log.Println("new message", username, data)

			m := make(map[string]interface{})
			m["username"] = username
			m["message"] = data

			so.BroadcastTo("chat", "new message", m)
		})

		so.On("add user", func(username2 string) {
			log.Println("add_user", username)
			if username == "" {
				username = username2
				numUsers = numUsers + 1

				l := make(map[string]interface{})
				l["numUsers"] = numUsers
				so.Emit("login", l)

				m := make(map[string]interface{})
				m["username"] = username2
				m["numUsers"] = numUsers
				so.BroadcastTo("chat", "user joined", m)
				log.Println("add_user", username)
			}
		})

		so.On("stop typing", func(_ string) {
			log.Println("stop typing", username)

			m := make(map[string]interface{})
			m["username"] = username
			so.BroadcastTo("chat", "stop typing", m)
		})

		so.On("typing", func(_ string) {
			log.Println("typing", username)

			m := make(map[string]interface{})
			m["username"] = username
			so.BroadcastTo("chat", "typing", m)
		})
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
