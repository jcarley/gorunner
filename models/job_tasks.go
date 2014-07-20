package models

type JobTasks struct {
	JobId   int64 `db:"job_id"`
	TaskId  int64 `db:"task_id"`
	Created int64 `db:"created_at"`
	Updated int64 `db:"updated_at"`
}
