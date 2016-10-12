package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
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
		log.Println("on connection")

		so.Join("chat")

		so.On("testcallback", func(msg string) map[string]interface{} {
			m := make(map[string]interface{})
			m["id"] = so.Id()
			m["message"] = msg + msg
			return m
		})

		so.On("new msg", func(data string) {
			log.Println("new msg", so.Id(), data)

			m := make(map[string]interface{})
			m["username"] = so.Id()
			m["message"] = data

			so.BroadcastTo("chat", "new msg", m)
		})

		so.On("disconnection", func() {
		})
	})

	//http.Handle("/socket.io/", server)
	http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://www.app-echo.com")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		server.ServeHTTP(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
