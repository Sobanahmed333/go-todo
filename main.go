package main

import (
	"fmt"
	"github.com/rs/cors"
	"os"
	"time"

	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// struct for every Task
type Task struct {
	gorm.Model
	createdAt time.Time
	Task      string
	Completed bool
}

// struct for Environmental Variables
type EnvVars struct {
	dialect  string
	host     string
	dbPort   string
	user     string
	dbName   string
	password string
}

var db *gorm.DB
var err error

func main() {
	db := ConnectDB()
	defer DisconnectDB(db)
	db.AutoMigrate(&Task{})
	router := RouteHandler()
	handler := Handler(router)
	log.Fatal(http.ListenAndServe(":8081", handler))
}

func ConnectDB() *gorm.DB {
	envVars := InitializeEnvVars()
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		envVars.host, envVars.user, envVars.dbName, envVars.password, envVars.dbPort)
	db, err = gorm.Open(envVars.dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s Successfully connected to database!\n", time.Now())
	}
	return db
}

func DisconnectDB(db *gorm.DB) {
	err := db.Close()
	if err != nil {
		return
	}
}

// Get environmental variables loaded
func InitializeEnvVars() EnvVars {
	return EnvVars{
		dialect:  os.Getenv("DIALECT"),
		host:     os.Getenv("HOST"),
		dbPort:   os.Getenv("DBPORT"),
		user:     os.Getenv("USER"),
		dbName:   os.Getenv("NAME"),
		password: os.Getenv("PASSWORD"),
	}
}

// this will make API routes
func RouteHandler() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tasks/list", GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", GetTask).Methods("GET")
	router.HandleFunc("/task/gen", GenerateTask).Methods("POST")
	router.HandleFunc("/task/del/{id}", DeleteTask).Methods("DELETE")
	return router
}

func Handler(router *mux.Router) http.Handler {
	options := cors.Options{
		// 		AllowedOrigins: []string{},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		Debug:          false,
	}
	handler := cors.New(options).Handler(router)
	return handler
}

// Create Task
func GenerateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)
	createdTask := db.Create(&newTask)
	err = createdTask.Error
	if err != nil {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			return
		}
	} else {
		json.NewEncoder(w).Encode(&newTask)
	}
	fmt.Println(time.Now(), "Create Task")
}

// Get an Array of tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	db.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
	w.WriteHeader(200)
	fmt.Println(time.Now(), "Get tasks")
}

// Get a Specific task based on ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	params := mux.Vars(r)
	db.First(&task, params["id"])
	json.NewEncoder(w).Encode(&task)
	fmt.Printf("%s Get task %s \n", time.Now(), params["id"])
}

// Delete a specific task based on ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	params := mux.Vars(r)
	db.First(&task, params["id"])
	db.Delete(&task)
	var tasks []Task
	db.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
	fmt.Printf("%s Del task %s \n", time.Now(), params["id"])
}
