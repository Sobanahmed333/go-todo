package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Router() *mux.Router{
	router := mux.NewRouter()
	router.HandleFunc("/tasks/list", GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", GetTask).Methods("GET")
	router.HandleFunc("/task/gen", GenerateTask).Methods("POST")
	router.HandleFunc("/task/del/{id}", DeleteTask).Methods("DELETE")
	return router
}

func TestGetTasks(t *testing.T) {
	envVars := InitializeEnvVars()
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		envVars.host, envVars.user, envVars.dbName, envVars.password, envVars.dbPort)
	db := ConnectDB(envVars.dialect, dbURI)
	defer DisconnectDB(db)
	request, _ := http.NewRequest("GET", "/tasks/list", nil)
	response := httptest.NewRecorder()
	router := Router()
	cors.Default().Handler(router).ServeHTTP(response, request)
	fmt.Println("response -->", reflect.TypeOf(response.Body))
}

func TestGetTask(t *testing.T) {
	envVars := InitializeEnvVars()
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
		envVars.host, envVars.user, envVars.dbName, envVars.password, envVars.dbPort)
	db := ConnectDB(envVars.dialect, dbURI)
	defer DisconnectDB(db)
	request, _ := http.NewRequest("GET", "/task/1", nil)
	response := httptest.NewRecorder()
	router := Router()
	cors.Default().Handler(router).ServeHTTP(response, request)
	fmt.Println("response -->", response.Body)
}

