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
	sqlExecutor   gorp.SqlExecutor
	transactional bool
}

func NewDatabase(context *DbContext) *Database {
	database := &Database{dbContext: context, sqlExecutor: context.Dbmap, transactional: false}
	return database
}

func (this *Database) transaction(f func() error) error {

	oldSqlExecutor := this.sqlExecutor

	trans, err := this.dbContext.Dbmap.Begin()
	if err != nil {
		log.Fatal(err)
	}

	this.sqlExecutor = trans
	defer func() {
		this.sqlExecutor = oldSqlExecutor
	}()

	if err := f(); err != nil {
		return trans.Rollback()
	}

	return trans.Commit()
}

func (this *Database) AddJob(job *Job) error {
	return this.transaction(func() error {
		return this.sqlExecutor.Insert(job)
	})
}

func (this *Database) GetJobList() *JobList {
	var jobs []Job

	_, err := this.sqlExecutor.Select(&jobs, "select * from jobs")
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
	if err := this.sqlExecutor.SelectOne(&job, "select * from jobs where name=?", jobId); err != nil {
		return nil, err
	}

	return &job, nil
}

func (this *Database) DeleteJob(jobId string) error {
	return this.transaction(func() error {
		if _, err := this.sqlExecutor.Exec("delete from jobs where name = ?", jobId); err != nil {
			return err
		}
		return nil
	})
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
