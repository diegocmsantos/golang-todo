package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type TodoItemModel struct {
	Id          int `gorm:"primary_key"`
	Description string
	Completed   bool
}

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is OK")
	w.Header().Set("Content-type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting Todolist API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch, http.MethodOptions},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
	todo := &TodoItemModel{
		Description: description,
		Completed:   false,
	}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if the TodoItem exist in DB
	err := GetItemByID(id)
	if err == false {
		w.Header().Set("Content-type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		todo.Completed = completed
		db.Save(&todo)
		w.Header().Set("Content-type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if the TodoItem exist in DB
	err := GetItemByID(id)
	if err == false {
		w.Header().Set("Content-type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-type", "application/json")
		io.WriteString(w, `{deleted}: true`)
	}
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting completed TodoItems.")
	completedTodoItems := GetTodoItems(true)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting incomplete TodoItems.")
	incompleteItems := GetTodoItems(false)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(incompleteItems)
}

func GetTodoItems(completed bool) interface{} {
	var todos []TodoItemModel
	return db.Where("completed = ?", completed).Find(&todos).Value
}

func GetItemByID(id int) bool {
	todo := &TodoItemModel{}
	result := db.First(&todo, id)
	if result.Error != nil {
		log.Warn("TodoItem not found in database")
		return false
	}
	return true
}
