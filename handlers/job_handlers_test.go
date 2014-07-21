package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jcarley/gorunner/handlers"
	"github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobHandlers", func() {

	BeforeEach(func() {
		dbContext := models.NewDbContext()
		defer dbContext.Dbmap.Db.Close()
		err := dbContext.Dbmap.TruncateTables()

		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ListJobs function", func() {

		BeforeEach(func() {
			job := &models.Job{Name: "test build", Status: "New"}

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			database := models.NewDatabase(dbContext)

			database.AddJob(job)
		})

		It("returns a json array of all jobs", func() {

			req, err := http.NewRequest("GET", "/jobs", nil)
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()
			database := models.NewDatabase(dbContext)

			appContext := &handlers.AppContext{Request: req, Response: w, Database: database}
			handlers.ListJobs(appContext)

			var payload []models.Job
			err = json.Unmarshal(w.Body.Bytes(), &payload)
			if err != nil {
				log.Fatal(err)
			}

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).NotTo(BeNil())
			Expect(payload[0].Name).To(Equal("test build"))
			Expect(payload[0].Status).To(Equal("New"))

			// fmt.Printf("%d - %s", w.Code, w.Body.String())
		})
	})

	Describe("AddJob function", func() {
		It("returns a http status code of 201", func() {

			body := bytes.NewBufferString(`{"name": "test job name"}`)

			req, err := http.NewRequest("POST", "/jobs", body)
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()
			database := models.NewDatabase(dbContext)

			appContext := &handlers.AppContext{Request: req, Response: w, Database: database}
			handlers.AddJob(appContext)

			Expect(w.Code).To(Equal(201))
		})

		It("adds a job", func() {
			body := bytes.NewBufferString(`{"name": "test job name"}`)

			req, err := http.NewRequest("POST", "/jobs", body)
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()
			database := models.NewDatabase(dbContext)

			appContext := &handlers.AppContext{Request: req, Response: w, Database: database}
			handlers.AddJob(appContext)
			dbContext.Dbmap.Db.Close()

			var job models.Job

			dbContext.Dbmap.SelectOne(&job, "select * from jobs where name = ?", "test job name")

			Expect(job).NotTo(BeNil())
		})
	})

	// r.HandleFunc("/jobs/{job}/tasks/{task}", RemoveTaskFromJob).Methods("DELETE")

	Describe("RemoveTaskFromJob", func() {
		var (
			job  models.Job
			task models.Task
		)

		BeforeEach(func() {
			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			job = models.Job{Name: "test_job", Status: "New"}
			task = models.Task{Name: "test_task"}
			dbContext.Dbmap.Insert(&job, &task)

			job_task := &models.JobTask{JobId: job.Id, TaskId: task.Id}
			dbContext.Dbmap.Insert(job_task)
		})

		It("removes a task from a job", func() {
			job_id := job.Id
			task_id := task.Id

			path := fmt.Sprintf("/jobs/%d/tasks/%d", job_id, task_id)
			fmt.Println(path)
			req, err := http.NewRequest("DELETE", path, nil)
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()
			database := models.NewDatabase(dbContext)

			appContext := &handlers.AppContext{Request: req, Response: w, Database: database}

			handlers.RemoveTaskFromJob(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from job_tasks where job_id = ? and task_id = ?", job_id, task_id)

			Expect(err).NotTo(HaveOccurred())
			Expect(w.Code).To(Equal(200))
			Expect(count).To(Equal(0))
		})
	})

})
