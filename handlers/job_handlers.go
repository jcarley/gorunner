package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcarley/gorunner/executor"
	"github.com/jcarley/gorunner/models"
)

func ListJobs(appContext *AppContext) {
	appContext.Write([]byte(appContext.Database.GetJobList().Json()))
}

func AddJob(appContext *AppContext) {

	payload := appContext.Unmarshal("name")

	job := &models.Job{Name: payload["name"], Status: "New"}

	err := appContext.Database.AddJob(job)
	if err != nil {
		appContext.Error(err, http.StatusInternalServerError)
		return
	}
	appContext.WriteHeader(201)
}

func GetJob(appContext *AppContext) {
	r := appContext.Request

	vars := mux.Vars(r)
	jobId := vars["job"]

	job, err := appContext.Database.GetJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}

	appContext.Marshal(job)
}

func DeleteJob(appContext *AppContext) {
	r := appContext.Request

	vars := mux.Vars(r)
	jobId := vars["job"]

	database := appContext.Database

	err := database.DeleteJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusInternalServerError)
		return
	}
}

func AddTaskToJob(appContext *AppContext) {
	r := appContext.Request

	vars := mux.Vars(r)
	jobId := vars["job"]

	job, err := appContext.Database.GetJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}
	// j := job.(models.Job)

	payload := appContext.Unmarshal("task")
	job.AppendTask(payload["task"])
	// jobList.Update(j)

	appContext.WriteHeader(201)
}

func RemoveTaskFromJob(appContext *AppContext) {
	r := appContext.Request
	vars := mux.Vars(r)

	job_id, _ := strconv.ParseInt(vars["job"], 10, 64)
	task_id, _ := strconv.ParseInt(vars["task"], 10, 64)

	err := appContext.Database.RemoveTaskFromJob(job_id, task_id)
	if err != nil {
		appContext.Error(err, http.StatusBadRequest)
		return
	}
}

func AddTriggerToJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	payload := unmarshal(r.Body, "trigger", w)

	j.AppendTrigger(payload["trigger"])
	triggerList := models.GetTriggerList()
	t, err := triggerList.Get(payload["trigger"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	executor.AddTrigger(t.(models.Trigger))
	jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTriggerFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobList()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	t := vars["trigger"]
	j.DeleteTrigger(t)
	jobList.Update(j)

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	jobs := jobList.GetJobsWithTrigger(t)

	if len(jobs) == 0 {
		executor.RemoveTrigger(t)
	}
}
