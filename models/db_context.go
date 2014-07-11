package models

import (
	"database/sql"
	"log"

	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
)

type DbContext struct {
	Dbmap *gorp.DbMap
}

func NewDbContext() DbContext {
	db, err := sql.Open("mysql", "gorunner-admin:letmein123@/gorunner")
	checkError(err, "Create tables failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// dbmap.AddTableWithName(Bucket{}, "buckets").SetKeys(true, "Id")

	return DbContext{Dbmap: dbmap}
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func (this DbContext) Migrate() {
	err := this.Dbmap.DropTablesIfExists()
	checkError(err, "Droping tables failed")

	err = this.Dbmap.CreateTablesIfNotExists()
	checkError(err, "Create tables failed")
}
