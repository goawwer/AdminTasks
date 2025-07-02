package api

import (
	"encoding/json"
	"net/http"

	"github.com/goawwer/admintasks/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *API) PostTask(w http.ResponseWriter, req *http.Request) {
	a.InitHeaders(w)
	a.logger.Info("Post Task POST /api/tasks")
	var task models.Task

	err := json.NewDecoder(req.Body).Decode(&task)

	if err != nil {
		a.logger.Info("Invalid JSON")
		a.InvalidJSON(w)
		return
	}

	defer req.Body.Close()

	t, err := a.storage.TaskRepository().Create(req.Context(), &task)
	if err != nil {
		a.logger.Info("Troubles while creating new user", err)
		a.DatabaseTrouble(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (a *API) GetAllTasks(w http.ResponseWriter, req *http.Request) {
	a.InitHeaders(w)
	a.logger.Info("Get all Tasks GET /api/tasks")

	tasks, err := a.storage.TaskRepository().GetAll(req.Context())
	if err != nil {
		a.logger.Info("Error while Tasks.GetAll")
		a.DatabaseTrouble(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		a.logger.Error("failed to write response", err)
	}
}

func (a *API) GetTask(w http.ResponseWriter, req *http.Request) {
	a.InitHeaders(w)
	idParam := mux.Vars(req)["id"]

	id, err := uuid.Parse(idParam)
	if err != nil {
		a.logger.Info("Invalid UUID in {id} param", err)
		a.InappropriateID(w)
		return
	}

	task, ok, err := a.storage.TaskRepository().GetByID(req.Context(), id)
	if err != nil {
		a.logger.Info("Troubles while accessing tasks database with that id:", err)
		a.DatabaseTrouble(w)
		return
	}

	if !ok {
		a.logger.Info("There is no task with that id")
		a.NotFound(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (a *API) UpdateTask(w http.ResponseWriter, req *http.Request) {
	a.InitHeaders(w)
	a.logger.Info("Update task PUT /api/task/{id}")

	var updatedTask models.Task
	err := json.NewDecoder(req.Body).Decode(&updatedTask)
	if err != nil {
		a.logger.Info("Invalid JSON")
		a.InvalidJSON(w)
		return
	}

	defer req.Body.Close()

	idParam := mux.Vars(req)["id"]

	id, err := uuid.Parse(idParam)
	if err != nil {
		a.logger.Info("Invalid UUID in {id} param", err)
		a.InappropriateID(w)
		return
	}
	task, err := a.storage.TaskRepository().Update(req.Context(), id, &updatedTask)
	if err != nil {
		a.logger.Info("Troubles while updating task with that id", err)
		a.DatabaseTrouble(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (a *API) DeleteTask(w http.ResponseWriter, req *http.Request) {
	a.InitHeaders(w)
	a.logger.Info("Delete task DELETE /api/task/{id}")

	idParam := mux.Vars(req)["id"]

	id, err := uuid.Parse(idParam)
	if err != nil {
		a.logger.Info("Invalid UUID in {id} param", err)
		a.InappropriateID(w)
		return
	}

	err = a.storage.TaskRepository().Delete(req.Context(), id)

	if err != nil {
		a.logger.Info("Troubles while deleting task", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
