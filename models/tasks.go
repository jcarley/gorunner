package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/coopernurse/gorp"
)

type Task struct {
	Id      int64  `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
	Script  string `db:"script" json:"script"`
	Created int64  `db:"created_at" json:"created_at,omitempty"`
	Updated int64  `db:"updated_at" json:"updated_at,omitempty"`
	Version int64  `db:"version" json:"version,omitempty"`
}

func (t Task) ID() string {
	return t.Name
}

func (this *Task) PreInsert(s gorp.SqlExecutor) error {
	this.Created = time.Now().UnixNano()
	this.Updated = this.Created
	return nil
}

func (this *Task) PreUpdate(s gorp.SqlExecutor) error {
	this.Updated = time.Now().UnixNano()
	return nil
}

func (this *Database) AddTask(task *Task) error {
	return this.transaction(func() error {
		return this.sqlExecutor.Insert(task)
	})
}

func (this *Database) GetTaskList() *TaskList {

	var tasks []Task
	_, err := this.sqlExecutor.Select(&tasks, "select * from tasks")
	if err != nil {
		log.Fatal(err)
	}

	taskList := TaskList{list{elements: make([]elementer, 0)}}

	for _, task := range tasks {
		taskList.elements = append(taskList.elements, task)
	}

	return &taskList
}

type TaskList struct {
	list
}

func (l *TaskList) Load(read ListReader) {
	bytes := read(l.fileName)
	var tasks []Task
	err := json.Unmarshal([]byte(string(bytes)), &tasks)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, task := range tasks {
		l.elements = append(l.elements, task)
	}
}
