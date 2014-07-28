package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

type JobTrigger struct {
	JobId     int64 `db:"job_id"`
	TriggerId int64 `db:"trigger_id"`
	Created   int64 `db:"created_at"`
	Updated   int64 `db:"updated_at"`
}

func (this *JobTrigger) PreInsert(s gorp.SqlExecutor) error {
	this.Created = time.Now().UnixNano()
	this.Updated = this.Created
	return nil
}

func (this *JobTrigger) PreUpdate(s gorp.SqlExecutor) error {
	this.Updated = time.Now().UnixNano()
	return nil
}
