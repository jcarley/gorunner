package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Install(r *mux.Router) {
	r.HandleFunc("/", App)
	r.HandleFunc("/ws", WsHandler)

	AppRoute(r, "/jobs", ListJobs).Methods("GET")
	AppRoute(r, "/jobs", AddJob).Methods("POST")
	AppRoute(r, "/jobs/{job}", GetJob).Methods("GET")
	AppRoute(r, "/jobs/{job}", DeleteJob).Methods("DELETE")
	AppRoute(r, "/jobs/{job}/tasks", AddTaskToJob).Methods("POST")
	AppRoute(r, "/jobs/{job}/tasks/{task}", RemoveTaskFromJob).Methods("DELETE")
	AppRoute(r, "/jobs/{job}/triggers", AddTriggerToJob).Methods("POST")
	AppRoute(r, "/jobs/{job}/triggers/{trigger}", RemoveTriggerFromJob).Methods("DELETE")

	AppRoute(r, "/tasks", ListTasks).Methods("GET")
	AppRoute(r, "/tasks", AddTask).Methods("POST")
	r.HandleFunc("/tasks/{task}", GetTask).Methods("GET")
	r.HandleFunc("/tasks/{task}", UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{task}", DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{task}/jobs", ListJobsForTask).Methods("GET")

	r.HandleFunc("/runs", ListRuns).Methods("GET")
	r.HandleFunc("/runs", AddRun).Methods("POST")
	r.HandleFunc("/runs/{run}", GetRun).Methods("GET")

	r.HandleFunc("/triggers", ListTriggers).Methods("GET")
	r.HandleFunc("/triggers", AddTrigger).Methods("POST")
	r.HandleFunc("/triggers/{trigger}", GetTrigger).Methods("GET")
	r.HandleFunc("/triggers/{trigger}", UpdateTrigger).Methods("PUT")
	r.HandleFunc("/triggers/{trigger}", DeleteTrigger).Methods("DELETE")
	r.HandleFunc("/triggers/{trigger}/jobs", ListJobsForTrigger).Methods("GET")

	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/")))
}
