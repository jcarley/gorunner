package models

import (
	"database/sql"
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
	checkError(err, "Open a connection failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	logFilename := path.Join(dir, "../logs/gorp.log")
	logFile, err := os.OpenFile(logFilename, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	dbmap.TraceOn("[gorp]", log.New(logFile, "gorunner", log.Lmicroseconds))
	dbmap.AddTableWithName(Job{}, "jobs").SetKeys(true, "Id")

	return &DbContext{Dbmap: dbmap}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func (this *DbContext) Migrate() error {
	err := this.Dbmap.DropTablesIfExists()
	if err != nil {
		return err
	}
	// checkError(err, "Droping tables failed")

	err = this.Dbmap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	// checkError(err, "Create tables failed")
	return nil
}
