package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	gorm.Model
	Id      int    `json:"id,string"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Text    string `json:"text"`
	Deleted bool   `gorm:"default:false" json:"-"` // Add this line
}

type Credentials struct {
	gorm.Model
	Login    string `json:"login"`
	Password string `json:"password"`
}

var db *gorm.DB

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

const (
	dsn = "host=host.docker.internal user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var entries []Entry
	// Exclude deleted entries by adding a condition to the query
	err := db.Model(&Entry{}).Where("deleted = ?", false).Order("id desc").Limit(3).Find(&entries).Error

	if err != nil {
		fmt.Println(err)
	}
	t, _ := template.ParseFiles("./templates/homepage.html")
	t.Execute(w, entries)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./templates/signIn.html")
	t.Execute(w, nil)
}

func checkAdminHandler(w http.ResponseWriter, r *http.Request) {
	input := decodeCredentials(w, r)
	correct := savedCredentials()
	var errorMessage string

	if input.Login != correct.Login {
		errorMessage = "Wrong login"
	} else if input.Password != correct.Password {
		errorMessage = "Wrong password"
	}

	if errorMessage != "" {
		logout(w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, errorMessage)))
		return
	}

	fmt.Println("Login and password are correct")

	session, _ := store.Get(r, "session-name")
	session.Values["role"] = "Admin"
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"redirectTo": "/new"}`))
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	// Set MaxAge to -1 to delete the session
	session.Options = &sessions.Options{MaxAge: -1}
	// Save the session to send the updated session cookie to the client
	session.Save(r, w)
}

func decodeCredentials(w http.ResponseWriter, r *http.Request) (c Credentials) {
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return c
}

func savedCredentials() (c Credentials) {
	data, err := os.ReadFile("./admin_credentials.txt")
	if err != nil {
		log.Println("Cannot read file:", err)
	}
	fmt.Println("Came from front: ", c)
	c.Login, c.Password, _ = strings.Cut(string(data), "\n")

	_, c.Login, _ = strings.Cut(string(c.Login), ": ")
	_, c.Password, _ = strings.Cut(string(c.Password), ": ")

	return c
}

func isAuthorised(w http.ResponseWriter, r *http.Request) (ans bool) {
	session, _ := store.Get(r, "session-name")
	role, ok := session.Values["role"].(string)
	if !ok || role != "Admin" {
		t, _ := template.ParseFiles("./templates/unauthorised.html")
		t.Execute(w, nil)
		return false
	} else {
		return true
	}
}

func newDataHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthorised(w, r) {
		t, _ := template.ParseFiles("./templates/newArticle.html")
		t.Execute(w, nil)
	}
}

func saveEntryHandler(res http.ResponseWriter, req *http.Request) {
	var p Entry
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&Entry{
		Text:    p.Text,
		Title:   p.Title,
		Preview: p.Preview,
	})

	if result.Error != nil {
		log.Fatalf("Error while inserting entry into database: %v", result.Error)
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write([]byte("{\"code\": 200, \"message\": \"its okay\"}"))
}

func connectDB() error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return err
	}
	err = db.AutoMigrate(&Entry{})
	if err != nil {
		fmt.Println("Error at the AutoMigrate")
		return err
	}
	return err
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var article Entry
	// Exclude deleted entries by adding a condition to the query
	result := db.First(&article, "id = ? AND deleted = ?", id, false)
	if result.Error != nil {
		http.Error(w, "Article not found or has been deleted", http.StatusNotFound)
		return
	}
	fmt.Println("articleHandler: ", article)
	t, _ := template.ParseFiles("./templates/article.html")
	t.Execute(w, article)
}

func EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthorised(w, r) {
		return
	}
	id := r.URL.Query().Get("id")
	var article Entry
	result := db.First(&article, id)
	if result.Error != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("./templates/editArticle.html")
	t.Execute(w, article)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deletePostHandler")
	if !isAuthorised(w, r) {
		return
	}
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	var article Entry
	result := db.First(&article, id)
	if result.Error != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}
	t, _ := template.ParseFiles("./templates/deleteArticle.html")
	t.Execute(w, article)
}

func softDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Ensure the user is authorized
	if !isAuthorised(w, r) {
		return
	}
	// Parse the request body to get the ID of the entry to be soft deleted
	var request struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Convert the ID from string to int
	id, err := strconv.Atoi(request.ID)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Perform the soft delete operation
	result := db.Model(&Entry{}).Where("id = ?", id).Update("deleted", true)
	if result.Error != nil {
		http.Error(w, "Failed to soft delete entry", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Entry soft deleted successfully!"))
}

func saveChangesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !isAuthorised(w, r) {
		return
	}
	var article Entry
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	result := db.Model(&Entry{}).Where("id = ?", article.Id).Updates(article)
	if result.Error != nil {
		http.Error(w, "Failed to update article", http.StatusInternalServerError)
		fmt.Println("Failed to update article")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Article updated successfully!"))
}

func main() {

	mux := http.NewServeMux()

	connectDB()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/admin/", adminHandler)
	mux.HandleFunc("/save-data", saveEntryHandler)
	mux.HandleFunc("/check", checkAdminHandler)
	mux.HandleFunc("/new", newDataHandler)
	mux.HandleFunc("/article", articleHandler)
	mux.HandleFunc("/edit", EditArticleHandler)
	mux.HandleFunc("/save-changes", saveChangesHandler)
	mux.HandleFunc("/deletePost", deletePostHandler)
	mux.HandleFunc("/softDelete", softDeleteHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
