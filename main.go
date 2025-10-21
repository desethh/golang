package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Lead struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type Comments struct {
	Comment string    `json:"comment"`
	User    string    `json:"username"`
	UTime   time.Time `json:"utime"`
}

var db *sql.DB

func postArticle(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT comment, username, utime FROM reviews ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Ошибка получения комментариев: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commentsList []Comments
	for rows.Next() {
		var c Comments
		if err := rows.Scan(&c.Comment, &c.User, &c.UTime); err != nil {
			http.Error(w, "Ошибка чтения комментариев: "+err.Error(), http.StatusInternalServerError)
			return
		}
		commentsList = append(commentsList, c)
	}

	tmpl := template.Must(template.ParseFiles("index.html"))

	err = tmpl.Execute(w, commentsList)
	if err != nil {
		http.Error(w, "Ошибка шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("anons")
	user := r.FormValue("user")
	utime := time.Now()
	if title == "" || user == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_, err := db.Exec("INSERT INTO reviews (comment, username, utime) VALUES (?, ?, ?)", title, user, utime)
	if err != nil {
		http.Error(w, "Ошибка сохранения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func apiLead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var lead Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Новая заявка: %+v\n", lead)

	botToken := "8235650920:AAF-4GqJVEATL5PZucO1OkqFMHy6oX5j1kI"
	chatID := "1937094139"
	msg := fmt.Sprintf(
		"📩 Новая заявка\n\n👤 Имя: %s\n📱 Телефон: %s\n📧 Email: %s\n💬 Сообщение: %s\n🕒 Дата: %s",
		lead.Name, lead.Phone, lead.Email, lead.Message, lead.Date,
	)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := map[string]string{
		"chat_id": chatID,
		"text":    msg,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "telegram error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	var err error
	dsn := "tauren91_itastan:9pV*taaN%baU@tcp(tauren91.beget.tech:3306)/tauren91_itastan?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка при проверке соединения с БД: ", err)
	}
	routing := http.NewServeMux()
	routing.HandleFunc("/", postArticle) // Главная страница
	routing.HandleFunc("POST /save_article", saveArticle)
	routing.HandleFunc("POST /api/lead", apiLead)

	fs := http.FileServer(http.Dir("./static"))
	routing.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", routing))
}
