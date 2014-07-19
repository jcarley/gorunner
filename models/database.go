package models

import (
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/coopernurse/gorp"
)

const (
	runsFile     = "runs.json"
	tasksFile    = "tasks.json"
	triggersFile = "triggers.json"
)

var (
	taskList    TaskList
	runList     RunList
	triggerList TriggerList
)

type ListWriter func([]byte, string)
type ListReader func(string) []byte

func InitDatabase() {
	taskList = TaskList{list{elements: make([]elementer, 10), fileName: tasksFile}}
	triggerList = TriggerList{list{elements: make([]elementer, 10), fileName: triggersFile}}
	runList = RunList{list{elements: make([]elementer, 10), fileName: runsFile}}

	taskList.Load(readFile)
	triggerList.Load(readFile)
	runList.Load()
}

type Database struct {
	dbContext     *DbContext
	transactional bool
}

func NewDatabase(context *DbContext) *Database {
	database := &Database{dbContext: context, transactional: false}
	return database
}

func (this *Database) InTransaction() *Database {
	this.transactional = true
	return this
}

func (this *Database) exec(f func(s gorp.SqlExecutor) error) error {
	if this.transactional {
		return this.withTransaction(f)
	} else {
		return this.withOutTransaction(f)
	}

}

func (this *Database) withTransaction(f func(s gorp.SqlExecutor) error) error {
	trans, err := this.dbContext.Dbmap.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if err := f(trans); err != nil {
		return trans.Rollback()
	}

	return trans.Commit()
}

func (this *Database) withOutTransaction(f func(s gorp.SqlExecutor) error) error {
	return f(this.dbContext.Dbmap)
}

func (this *Database) AddJob(job *Job) error {
	return this.exec(func(s gorp.SqlExecutor) error {
		return s.Insert(job)
	})
}

func (this *Database) GetJobList() *JobList {

	var jobs []Job

	var err error
	this.exec(func(s gorp.SqlExecutor) error {
		_, err := s.Select(&jobs, "select * from jobs")
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	jobList := JobList{list{elements: make([]elementer, 0)}}

	for _, job := range jobs {
		jobList.elements = append(jobList.elements, job)
	}

	return &jobList
}

func (this *Database) GetJob(jobId string) (*Job, error) {
	var job Job

	err := this.exec(func(s gorp.SqlExecutor) error {
		return s.SelectOne(&job, "select * from jobs where name=?", jobId)
	})

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func GetJobList() *JobList {
	dbContext := NewDbContext()
	defer dbContext.Dbmap.Db.Close()

	var jobs []Job
	_, err := dbContext.Dbmap.Select(&jobs, "select * from jobs")
	if err != nil {
		log.Fatal(err)
	}

	jobList := JobList{list{elements: make([]elementer, 0)}}

	for _, job := range jobs {
		jobList.elements = append(jobList.elements, job)
	}

	return &jobList
}

func DeleteJob(jobId string) error {
	dbContext := NewDbContext()
	defer dbContext.Dbmap.Db.Close()

	trans, err := dbContext.Dbmap.Begin()
	if err != nil {
		return err
	}

	trans.Exec("delete from jobs where name = ?", jobId)

	return trans.Commit()
}

func GetRunList() *RunList {
	return &runList
}

func GetRunListSorted() *RunList {
	sort.Sort(Reverse{&runList})
	return &runList
}

func GetTaskList() *TaskList {
	return &taskList
}

func GetTriggerList() *TriggerList {
	return &triggerList
}

func writeFile(bytes []byte, filePath string) {
	err := ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []byte {
	_, err := os.Stat(filePath)
	if err != nil {
		println("Couldn't read file, creating fresh:", filePath)
		err = ioutil.WriteFile(filePath, []byte("[]"), 0644)
		if err != nil {
			panic(err)
		}
	}

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return bytes
}

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Interface.Less(j, i)
}
