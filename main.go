package main

import (
	//"github.com/gdamore/tcell/v2"
	//"github.com/rivo/tview"
	//"github.com/go-git/go-git/v5"
	"fmt"
	"io"

	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

// Work on menu
// Make an app that gets the user requested golang github package from github
// and downloads and imports it automatically

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// Example of user signing up for orderbooks
// to recieve real-time live data
func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		fmt.Println(payload)
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Printf("new incoming connection from client: %s", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) realUpdate(ws *websocket.Conn) {
	var sum int
	for {
		fmt.Printf("money: %d", sum)
		s := strconv.Itoa(sum)
		ws.Write([]byte(s))
		sum += 2
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	// A buffer in temporary storage to transfer data
	// In this case a sequence of bytes or a string
	buf := make([]byte, 1024)
	for {
		// Read the length of the buffer
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		// Store bytes up to the length (Now whole buffer)
		msg := buf[:n]

		// Displays message to everyone connected
		// to Websocket connection endpoint
		s.broadcast(msg)
	}

}

// Wrapper function which disables Origin header checking
// This means that soemthing like Postman will work
// Without Golang restricting certain URL's
// such as ours, http://localhost:3000
func createWebsocketServer(handler websocket.Handler) websocket.Server {
	return websocket.Server{
		Handshake: func(c *websocket.Config, req *http.Request) error {
			return nil
		},

		Handler: handler,
	}
}

// Main websocket part
// We essentially iterate through every
// websocket connection and display the message
// sent by the current connection to every other connection
// This is how real-time chat messaging works
func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		// Go routine func ensures a client does not
		// have to wait long to recieve message
		//  because of other clients delay
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write error:", err)
			}
		}(ws)
	}
}

func main() {

	// Start server and handle at /ws endpoint
	server := NewServer()
	http.Handle("/ws", createWebsocketServer(websocket.Handler(server.handleWS)))
	http.Handle("/orderbookfeed", createWebsocketServer(websocket.Handler(server.handleWSOrderbook)))
	http.Handle("/updates", createWebsocketServer(websocket.Handler(server.realUpdate)))
	http.ListenAndServe(":3000", nil)

}
