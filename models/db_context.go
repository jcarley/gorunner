package models

import (
	"database/sql"
	"io"
	"log"
	"os"
	"path"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

type DbContext struct {
	Dbmap *gorp.DbMap
}

func NewDbContext() *DbContext {
	db, err := sql.Open("mysql", "gorunner-admin:letmein123@/gorunner")
	checkError(err, "Opening a connection failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	logFile, err := logfileLocation()
	checkError(err, "")

	dbmap.TraceOn("[gorp]", log.New(logFile, "gorunner", log.Lmicroseconds))
	dbmap.AddTableWithName(Job{}, "jobs").SetKeys(true, "Id")
	dbmap.AddTableWithName(Task{}, "tasks").SetKeys(true, "Id")
	dbmap.AddTableWithName(JobTask{}, "job_tasks").SetKeys(false, "JobId", "TaskId")

	return &DbContext{Dbmap: dbmap}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func logfileLocation() (io.Writer, error) {
	dir, err := os.Getwd()
	checkError(err, "")

	logFilename := path.Join(dir, "../logs/gorp.log")
	logFile, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err, "")

	return logFile, nil
}

func (this *DbContext) Migrate() {
	err := this.Dbmap.DropTablesIfExists()
	checkError(err, "Droping tables failed")

	err = this.Dbmap.CreateTablesIfNotExists()
	checkError(err, "Create tables failed")
}
