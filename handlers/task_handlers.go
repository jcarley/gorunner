package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcarley/gorunner/models"
)

func ListTasks(appContext *AppContext) {
	taskList := appContext.Database.GetTaskList()

	appContext.Header().Set("Content-Type", "application/json")
	appContext.Write([]byte(taskList.Json()))
}

func AddTask(appContext *AppContext) {
	payload := appContext.Unmarshal("name")

	task := models.Task{Name: payload["name"], Script: ""}
	err := appContext.Database.AddTask(&task)
	if err != nil {
		appContext.Error(err, http.StatusBadRequest)
		return
	}
	appContext.WriteHeader(201)
}

func GetTask(appContext *AppContext) {
	vars := appContext.Vars

	task_id, _ := strconv.ParseInt(vars["task"], 10, 64)
	task, err := appContext.Database.GetTask(task_id)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}

	appContext.Marshal(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	payload := unmarshal(r.Body, "script", w)

	t := task.(models.Task)
	t.Script = payload["script"]
	taskList.Update(t)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskList := models.GetTaskList()

	vars := mux.Vars(r)
	task, err := taskList.Get(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	taskList.Delete(task.ID())
}

func ListJobsForTask(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()
	vars := mux.Vars(r)
	jobs := jobList.GetJobsWithTask(vars["task"])
	marshal(jobs, w)
}
