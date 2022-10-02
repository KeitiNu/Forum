package main

import (
	"errors"
	"log"
	"net/http"

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
			// delete(savedSocketReaders, name)
		}()
		// if savedSocketReader == nil {
		// 	savedSocketReader = make([]*socketReader, 0)
		// }
		con, _ := upgrader.Upgrade(w, r, nil)
		ptrSocketReader := &socketReader{
			con: con,
		}

		var msg RecievedMessage

		// _, message, _ := ptrSocketReader.con.ReadMessage()
		// json.Unmarshal(message, &msg)

		err := ptrSocketReader.con.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
		}
		// // Get struct from string.
		// json.Unmarshal(message, &msg)
		// fmt.Println("Message retrieved: ", msg.Context)

		// ptrSocketReader.con.WriteMessage(websocket.TextMessage, []byte("Greetings from golang"))

		name = msg.Context
		if msg.MessageType == "offline" {
			removeSocketReader(msg.Context)
		} else {
			savedSocketReaders[name] = ptrSocketReader

			sendOnlineUserInfo()
		}

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
		socket.con.WriteJSON(context)
	}
}

func removeSocketReader(name string) {
	delete(savedSocketReaders, name)

	var context = &Context{
		ContextType: "offline",
		OfflineUser: name,
	}

	for _, socket := range savedSocketReaders {
		socket.con.WriteJSON(context)
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
			socket.con.WriteJSON(context)

		}
	}
}

func fillInfo(recipient string, offset int) {
	// const db = new sqlite3.Database("database.db");
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
