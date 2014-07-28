package models

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/coopernurse/gorp"
)

type Job struct {
	Id       int64    `db:"id" json:"id,omitempty"`
	Name     string   `db:"name" json:"name,omitempty"`
	Tasks    []string `db:"-" json:"tasks,omitempty"`
	Status   string   `db:"status" json:"status,omitempty"`
	Triggers []string `db:"-" json:"triggers,omitempty"`
	Created  int64    `db:"created_at" json:"created_at,omitempty"`
	Updated  int64    `db:"updated_at" json:"updated_at,omitempty"`
	Version  int64    `db:"version" json:"version,omitempty"`
}

func (this *Job) PreInsert(s gorp.SqlExecutor) error {
	this.Created = time.Now().UnixNano()
	this.Updated = this.Created
	return nil
}

func (this *Job) PreUpdate(s gorp.SqlExecutor) error {
	this.Updated = time.Now().UnixNano()
	return nil
}

func (j Job) ID() string {
	return j.Name
}

func (this *Database) AddJob(job *Job) error {
	return this.transaction(func() error {
		return this.sqlExecutor.Insert(job)
	})
}

func (this *Database) GetJobList() *JobList {
	var jobs []Job

	_, err := this.sqlExecutor.Select(&jobs, "select * from jobs")
	if err != nil {
		log.Fatal(err)
	}

	jobList := JobList{list{elements: make([]elementer, 0)}}

	for _, job := range jobs {
		jobList.elements = append(jobList.elements, job)
	}

	return &jobList
}

func (this *Database) GetJob(jobId int64) (*Job, error) {
	var job Job
	if err := this.sqlExecutor.SelectOne(&job, "select * from jobs where id=?", jobId); err != nil {
		return nil, err
	}

	return &job, nil
}

func (this *Database) DeleteJob(jobId string) error {
	return this.transaction(func() error {
		if _, err := this.sqlExecutor.Exec("delete from jobs where name = ?", jobId); err != nil {
			return err
		}
		return nil
	})
}

func (this *Database) RemoveTaskFromJob(job_id, task_id int64) error {
	return this.transaction(func() error {
		jobTask := &JobTask{JobId: job_id, TaskId: task_id}
		if _, err := this.sqlExecutor.Delete(jobTask); err != nil {
			return err
		}
		return nil
	})
}

func (this *Database) AppendTask(job_id, task_id int64) error {
	return this.transaction(func() error {
		jobTask := &JobTask{JobId: job_id, TaskId: task_id}
		if err := this.sqlExecutor.Insert(jobTask); err != nil {
			return err
		}
		return nil
	})
}

func (this *Database) AppendTrigger(job_id, trigger_id int64) error {
	return this.transaction(func() error {
		jobTrigger := &JobTrigger{JobId: job_id, TriggerId: trigger_id}
		if err := this.sqlExecutor.Insert(jobTrigger); err != nil {
			return err
		}
		return nil
	})
}

func (j *Job) DeleteTask(taskPosition int) error {
	i := taskPosition
	j.Tasks = j.Tasks[:i+copy(j.Tasks[i:], j.Tasks[i+1:])]
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
