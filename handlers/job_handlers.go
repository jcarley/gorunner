package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/jcarley/gorunner/executor"
	"github.com/jcarley/gorunner/models"
)

func withTransaction(dbContext *models.DbContext, handler func(trans *gorp.Transaction) error) error {

	trans, err := dbContext.Dbmap.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if err := handler(trans); err != nil {
		return trans.Rollback()
	}

	return trans.Commit()
}

func ListJobs(appContext *AppContext) {
	appContext.Response.Write([]byte(models.GetJobList(appContext.DbContext).Json()))
}

func AddJob(appContext *AppContext) {

	dbContext := appContext.DbContext

	payload := unmarshal(appContext.Request.Body, "name", appContext.Response)

	job := &models.Job{Name: payload["name"], Status: "New"}

	err := withTransaction(dbContext, func(trans *gorp.Transaction) error {
		return models.AddJob(job, trans)
	})

	if err != nil {
		http.Error(appContext.Response, err.Error(), http.StatusInternalServerError)
		return
	}
	appContext.Response.WriteHeader(201)
}

func GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["job"]

	job, err := models.GetJob(jobId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	marshal(job, w)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["job"]

	err := models.DeleteJob(jobId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddTaskToJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["job"]

	job, err := models.GetJob(jobId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// j := job.(models.Job)

	payload := unmarshal(r.Body, "task", w)
	job.AppendTask(payload["task"])
	// jobList.Update(j)

	w.WriteHeader(201)
}

func RemoveTaskFromJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobListOld()

	vars := mux.Vars(r)
	job, err := jobList.Get(vars["job"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	j := job.(models.Job)

	taskPosition, err := strconv.Atoi(vars["task"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j.DeleteTask(taskPosition)
	jobList.Update(j)
}

func AddTriggerToJob(w http.ResponseWriter, r *http.Request) {
	jobList := models.GetJobListOld()

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
	jobList := models.GetJobListOld()

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
