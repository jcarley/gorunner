package models

import (
	"encoding/json"
	"time"

	"github.com/coopernurse/gorp"
)

type Trigger struct {
	Id       int64  `db:"id" json:"id,omitempty"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Created  int64  `db:"created_at" json:"created_at,omitempty"`
	Updated  int64  `db:"updated_at" json:"updated_at,omitempty"`
	Version  int64  `db:"version" json:"version,omitempty"`
}

func (t Trigger) ID() string {
	return t.Name
}

func (this *Trigger) PreInsert(s gorp.SqlExecutor) error {
	this.Created = time.Now().UnixNano()
	this.Updated = this.Created
	return nil
}

func (this *Trigger) PreUpdate(s gorp.SqlExecutor) error {
	this.Updated = time.Now().UnixNano()
	return nil
}

type TriggerList struct {
	list
}

func (l *TriggerList) Load(read ListReader) {
	bytes := read(l.fileName)
	var triggers []Trigger
	err := json.Unmarshal([]byte(string(bytes)), &triggers)
	if err != nil {
		panic(err)
	}
	l.elements = nil
	for _, trigger := range triggers {
		l.elements = append(l.elements, trigger)
	}
}
