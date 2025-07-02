package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goawwer/admintasks/store"
	"github.com/sirupsen/logrus"
)

var (
	prefix string = "/api"
)

func (a *API) ConfigureLogger() error {
	loggerLvl, err := logrus.ParseLevel(a.config.LoggerLevel)
	if err != nil {
		return fmt.Errorf("failed to configure logger: %w", err)
	}
	a.logger.SetLevel(loggerLvl)
	return nil
}

func (a *API) ConfigureRouter() {
	a.router.HandleFunc(prefix+"/tasks", a.GetAllTasks).Methods("GET")
	a.router.HandleFunc(prefix+"/tasks", a.PostTask).Methods("POST")
	a.router.HandleFunc(prefix+"/tasks/{id}", a.GetTask).Methods("GET")
	a.router.HandleFunc(prefix+"/tasks/{id}", a.UpdateTask).Methods("PUT")
	a.router.HandleFunc(prefix+"/tasks/{id}", a.DeleteTask).Methods("DELETE")
}

func (a *API) ConfigureStorage() error {
	storage := store.New(a.config)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	a.storage = storage
	return nil
}

func (a *API) InitHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

type MessageResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func (a *API) InvalidJSON(w http.ResponseWriter) {
	msg := MessageResponse{
		StatusCode: 400,
		Message:    "Provided json is invalid",
		IsError:    true,
	}
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(msg)
}

func (a *API) DatabaseTrouble(w http.ResponseWriter) {
	msg := MessageResponse{
		StatusCode: 500,
		Message:    "We have some problems to accessing database. Try again later",
		IsError:    true,
	}
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(msg)
}

func (a *API) NotFound(w http.ResponseWriter) {
	msg := MessageResponse{
		StatusCode: 404,
		Message:    "Not found",
		IsError:    true,
	}
	w.WriteHeader(404)
	json.NewEncoder(w).Encode(msg)
}

func (a *API) InappropriateID(w http.ResponseWriter) {
	msg := MessageResponse{
		StatusCode: 400,
		Message:    "Inappropriate id value",
		IsError:    true,
	}
	w.WriteHeader(400)
	json.NewEncoder(w).Encode(msg)
}
