package models

import (
	"encoding/json"
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
