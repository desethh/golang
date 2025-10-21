package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var (
	db       *sql.DB
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	clients = make(map[*websocket.Conn]bool)
)

type Message struct {
	Nick string `json:"nick"`
	Text string `json:"text"`
}

func main() {
	var err error
	db, err = sql.Open("mysql", "tauren91_itastan:9pV*taaN%baU@tcp(tauren91.beget.tech:3306)/tauren91_itastan?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}
	defer ws.Close()

	log.Println("Клиент подключился")
	clients[ws] = true

	sendHistory(ws)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Клиент отключился:", err)
			delete(clients, ws)
			break
		}

		_, err = db.Exec("INSERT INTO messages (nick, text) VALUES (?, ?)", msg.Nick, msg.Text)
		if err != nil {
			log.Println("Ошибка сохранения в БД:", err)
		}

		for c := range clients {
			c.WriteJSON(msg)
		}
	}
}

func sendHistory(ws *websocket.Conn) {
	rows, err := db.Query("SELECT nick, text FROM messages ORDER BY id ASC LIMIT 50")
	if err != nil {
		log.Println("Ошибка чтения истории:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		rows.Scan(&msg.Nick, &msg.Text)
		ws.WriteJSON(msg)
	}
}
