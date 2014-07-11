package models

import (
	"encoding/json"
	"errors"
)

type Job struct {
	Id       int64    `json:"id", db:"id"`
	Name     string   `json:"name"`
	Tasks    []string `json:"tasks"`
	Status   string   `json:"status"`
	Triggers []string `json:"triggers"`
	Created  int64    `db:"created_at"`
	Updated  int64    `db:"updated_at"`
	Version  int64    `db:"version"`
}

func (j Job) ID() string {
	return j.Name
}

func (j *Job) AppendTask(task string) {
	j.Tasks = append(j.Tasks, task)
}

func (j *Job) DeleteTask(taskPosition int) error {
	i := taskPosition
	j.Tasks = j.Tasks[:i+copy(j.Tasks[i:], j.Tasks[i+1:])]
	return nil
}

func (j *Job) AppendTrigger(trigger string) error {
	for _, name := range j.Triggers {
		if name == trigger {
			return errors.New("Trigger already on job")
		}
	}
	j.Triggers = append(j.Triggers, trigger)
	return nil
}

func (j *Job) DeleteTrigger(trigger string) error {
	for i, name := range j.Triggers {
		if name == trigger {
			j.Triggers = j.Triggers[:i+copy(j.Triggers[i:], j.Triggers[i+1:])]
			return nil
		}
	}
	return errors.New("Trigger not found")
}

type JobList struct {
	list
}

func (l *JobList) Load() {
	bytes := readFile(l.fileName)
	var jobs []Job
	err := json.Unmarshal([]byte(string(bytes)), &jobs)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, job := range jobs {
		l.elements = append(l.elements, job)
	}
}

func (l *JobList) GetJobsWithTrigger(triggerName string) (jobs []Job) {
	jobs = make([]Job, 0)
	for _, e := range l.elements {
		job := e.(Job)
		for _, trigger := range job.Triggers {
			if trigger == triggerName {
				jobs = append(jobs, job)
			}
		}
	}
	return
}

func (l *JobList) GetJobsWithTask(taskName string) (jobs []Job) {
	jobs = make([]Job, 0)
	for _, e := range l.elements {
		job := e.(Job)
		for _, task := range job.Tasks {
			if task == taskName {
				jobs = append(jobs, job)
			}
		}
	}
	return
}
