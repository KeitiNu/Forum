package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Create a variable to store the page where the client was before action
// (ex. logging in and returning directly to the post)
var back string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type socketReader struct {
	con *websocket.Conn // pointer to the socket
	// name string 			// name of our current user
	// request string  // what info is being send (ex: message request)
	context Context // info being sent
    mu      sync.Mutex

	// mode int 			// we are not using this at the moment
}

var savedSocketReaders = make(map[string]*socketReader)

type Context struct {
	ContextType string
	// 	content string
	// recipient   string
	OnlineUsers []string
	OfflineUser string
	Sender      string
	Message     string
	// offlineUsers []string
}

type RecievedMessage struct {
	MessageType string `json:"messageType"`
	Context     string `json:"context"`
	Recipient   string `json:"recipient"`
}

// var savedSocketReader []*socketReader = make([]*socketReader, 0)

func (app *application) socket(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		var name string

		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
			}
			r.Body.Close()
		}()

		con, _ := upgrader.Upgrade(w, r, nil)
		ptrSocketReader := &socketReader{
			con: con,
		}

		var msg RecievedMessage



		err := ptrSocketReader.con.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
		}


		name = msg.Context

		switch msg.MessageType {
		case "offline":
			removeSocketReader(msg.Context)
			break
		case "online":
			savedSocketReaders[name] = ptrSocketReader
			sendOnlineUserInfo()
			break
		case "typing":
			fmt.Println("typing")
			break
		default:
			break
		}

		// if msg.MessageType == "offline" {
		// 	removeSocketReader(msg.Context)
		// }
		// if msg.MessageType == "online" {

		// 	savedSocketReaders[name] = ptrSocketReader
		// 	sendOnlineUserInfo()
		// }

		// if msg.MessageType == "typing" {
		// 	typingInProgress("sender", "recipient")

		// }

	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
	}
}

func sendOnlineUserInfo() {
	var onlineArr []string

	for key := range savedSocketReaders {
		onlineArr = append(onlineArr, key)
	}

	var context = &Context{
		ContextType: "online",
		OnlineUsers: onlineArr,
	}

	for _, socket := range savedSocketReaders {
		socket.send(context)
	}
}

func removeSocketReader(name string) {
	delete(savedSocketReaders, name)

	var context = &Context{
		ContextType: "offline",
		OfflineUser: name,
	}

	for _, socket := range savedSocketReaders {
		socket.send(context)
	}
}

func sendChatNotification(sender string, recipient string, message string) {
	var context = &Context{
		ContextType: "chat",
		Sender:      sender,
		Message:     message,
	}

	for name, socket := range savedSocketReaders {
		if name == recipient {
			socket.send(context)
		}
	}
}

func typingInProgress(sender string, recipient string) {
	var context = &Context{
		ContextType: "typing",
		Sender:      sender,
	}

	for name, socket := range savedSocketReaders {
		if name == recipient {
			socket.send(context)
		}
	}
}


func (p *socketReader) send(v interface{}) error {
    p.mu.Lock()
    defer p.mu.Unlock()
    return p.con.WriteJSON(v)
}


func fillInfo(recipient string, offset int) {
	// const db = new sqlite3.Database("database.db");WriteJSON
	var sql string = `
    SELECT
        sender_id,
        recipient_id,
        content,
        sent_at
    FROM messages
    WHERE
        sender_id = ${recipient} AND recipient_id = ${user} OR
        recipient_id = ${recipient} AND sender_id = ${user}
    ORDER BY sent_at ASC
    LIMIT 10
    OFFSET ${offset}
    `
	println(sql)
	// db.all(sql, (err, rows) => {
	//     console.log(rows)
	//     if (err) {throw err}
	// })
}
