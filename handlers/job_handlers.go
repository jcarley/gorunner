package handlers

import (
	"net/http"
	"strconv"

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
	vars := appContext.Vars
	jobId, _ := strconv.ParseInt(vars["job"], 10, 64)

	job, err := appContext.Database.GetJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}

	appContext.Marshal(job)
}

func DeleteJob(appContext *AppContext) {
	vars := appContext.Vars
	jobId := vars["job"]

	database := appContext.Database

	err := database.DeleteJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusInternalServerError)
		return
	}
}

func AddTaskToJob(appContext *AppContext) {
	vars := appContext.Vars

	jobId, _ := strconv.ParseInt(vars["job"], 10, 64)

	_, err := appContext.Database.GetJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}

	payload := appContext.Unmarshal("task")
	taskId, _ := strconv.ParseInt(payload["task"], 10, 64)
	appContext.Database.AppendTask(jobId, taskId)

	appContext.WriteHeader(201)
}

func RemoveTaskFromJob(appContext *AppContext) {
	vars := appContext.Vars

	job_id, _ := strconv.ParseInt(vars["job"], 10, 64)
	task_id, _ := strconv.ParseInt(vars["task"], 10, 64)

	err := appContext.Database.RemoveTaskFromJob(job_id, task_id)
	if err != nil {
		appContext.Error(err, http.StatusBadRequest)
		return
	}
}

func AddTriggerToJob(appContext *AppContext) {
	vars := appContext.Vars

	jobId, _ := strconv.ParseInt(vars["job"], 10, 64)

	_, err := appContext.Database.GetJob(jobId)
	if err != nil {
		appContext.Error(err, http.StatusNotFound)
		return
	}

	payload := appContext.Unmarshal("trigger")
	triggerId, _ := strconv.ParseInt(payload["trigger"], 10, 64)
	appContext.Database.AppendTrigger(jobId, triggerId)

	appContext.WriteHeader(201)
}

func RemoveTriggerFromJob(appContext *AppContext) {
	vars := appContext.Vars

	job_id, _ := strconv.ParseInt(vars["job"], 10, 64)
	trigger_id, _ := strconv.ParseInt(vars["trigger"], 10, 64)

	err := appContext.Database.RemoveTriggerFromJob(job_id, trigger_id)
	if err != nil {
		appContext.Error(err, http.StatusBadRequest)
		return
	}

	// If Trigger is no longer attached to any Jobs, remove it from Cron to save cycles
	// jobs, err := appContext.Database.GetJobsWithTrigger(trigger_id)

	// if len(jobs) == 0 {
	// executor.RemoveTrigger(t)
	// }
}
