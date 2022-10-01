package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	name string 			// name of our current user
	request string  // what info is being send (ex: message request)
	context Context // info being sent

	// mode int 			// we are not using this at the moment
}

var savedSocketReaders = make(map[string]*socketReader)

type Context struct {
	// 	chat Chat
	// }
	// type Chat struct {
	// 	content string
	recipient   string
	onlineUsers []string
	// offlineUsers []string
}

var savedSocketReader []*socketReader = make([]*socketReader, 0)

func (app *application) socket(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		// if savedSocketReader == nil {
		// 	savedSocketReader = make([]*socketReader, 0)
		// }

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

		// ptrSocketReader.con.WriteMessage(websocket.TextMessage, []byte("Greetings from golang"))

		_, message, _ := ptrSocketReader.con.ReadMessage()

		var name = string(message)

		// ptrSocketReader.name = name

		fmt.Println("Message retrieved: ", name)

		savedSocketReaders[name] = ptrSocketReader;

		// savedSocketReader = append(savedSocketReader, ptrSocketReader)

		var onlineArr []string

		for key, _ := range savedSocketReaders {
			onlineArr = append(onlineArr, key)
		}

		var names = strings.Join(onlineArr, ", ")

		// var len = len(onlineArr)
		// s1 := strconv.Itoa(len)

		for _, socket := range savedSocketReaders {
			fmt.Println(socket.con.UnderlyingConn())
			socket.con.WriteMessage(websocket.TextMessage, []byte(names))
		}

		// for _, socket := range savedSocketReader {
		// 	socket.con.WriteMessage(websocket.TextMessage, []byte("other"))
		// }

	case "POST":
		app.serverError(w, errors.New("POST METHOD NOT ALLOWED"))
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
