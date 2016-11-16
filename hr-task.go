package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"fmt"
)

func main() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	server.On("connection", func(so socketio.Socket) {
		log.Println(so.Id(), "on connection")

		so.Join("chat")

		so.On("testcallback", func(msg string) map[string]interface{} {
			log.Println(so.Id(), msg)
			m := make(map[string]interface{})
			m["id"] = so.Id()
			m["message"] = msg + msg
			return m
		})

		so.On("log", func(typename, data string) {
			log.Println(so.Id(), typename, data)

			m := make(map[string]interface{})
			m["username"] = so.Id()
			m["type"] = typename
			m["message"] = data

			so.BroadcastTo("chat", "log", m)
		})

		so.On("disconnection", func() {
		})
	})

	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		server.ServeHTTP(w, r)
	})

	http.HandleFunc("/send/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			msg := r.URL.Query().Get("msg")
			fmt.Println(msg)
			server.BroadcastTo("chat", "task u", msg)
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./assets")))
	log.Println("Serving at localhost:5001...")
	log.Fatal(http.ListenAndServe(":5001", nil))
}
