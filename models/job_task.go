package models

import (
	"time"

	"github.com/coopernurse/gorp"
)

type JobTask struct {
	JobId   int64 `db:"job_id"`
	TaskId  int64 `db:"task_id"`
	Created int64 `db:"created_at"`
	Updated int64 `db:"updated_at"`
}

func (this *JobTask) PreInsert(s gorp.SqlExecutor) error {
	this.Created = time.Now().UnixNano()
	this.Updated = this.Created
	return nil
}

func (this *JobTask) PreUpdate(s gorp.SqlExecutor) error {
	this.Updated = time.Now().UnixNano()
	return nil
}
