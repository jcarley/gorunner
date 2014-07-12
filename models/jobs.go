package models

import (
	"encoding/json"
	"errors"
)

type Job struct {
	Id       int64    `db:"id" json:"id,omitempty"`
	Name     string   `db:"name" json:"name,omitempty"`
	Tasks    []string `db:"tasks" json:"tasks,omitempty"`
	Status   string   `db:"status" json:"status,omitempty"`
	Triggers []string `db:"triggers" json:"triggers,omitempty"`
	Created  int64    `db:"created_at" json:"created_at,omitempty"`
	Updated  int64    `db:"updated_at" json:"updated_at,omitempty"`
	Version  int64    `db:"version" json:"version,omitempty"`
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
