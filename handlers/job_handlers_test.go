package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/jcarley/gorunner/handlers"
	"github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobHandlers", func() {

	BeforeEach(func() {
		dbContext := models.NewDbContext()
		defer dbContext.Dbmap.Db.Close()
		err := dbContext.TruncateTables()

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
			path := "/jobs"
			appContext, dbContext, w := NewAppContext("GET", path, "", nil)
			defer dbContext.Dbmap.Db.Close()

			handlers.ListJobs(appContext)

			var payload []models.Job
			err := json.Unmarshal(w.Body.Bytes(), &payload)
			if err != nil {
				log.Fatal(err)
			}

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).NotTo(BeNil())
			Expect(payload[0].Name).To(Equal("test build"))
			Expect(payload[0].Status).To(Equal("New"))
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

	Describe("AddTaskToJob", func() {
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
		})

		It("adds a task to an existing job", func() {
			job_id := job.Id
			task_id := task.Id

			vars := make(map[string]string)
			vars["job"] = strconv.FormatInt(job_id, 10)

			path := fmt.Sprintf("/jobs/%d/tasks", job_id)

			bodyString := fmt.Sprintf(`{"task":"%s"}`, strconv.FormatInt(task_id, 10))
			appContext, dbContext, w := NewAppContext("POST", path, bodyString, vars)
			defer dbContext.Dbmap.Db.Close()

			handlers.AddTaskToJob(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from job_tasks where job_id = ? and task_id = ?", job_id, task_id)

			Expect(err).NotTo(HaveOccurred())
			Expect(w.Code).To(Equal(201))
			Expect(count).To(Equal(int64(1)))
		})
	})

	Describe("AddTriggerToJob", func() {
		var (
			job     models.Job
			trigger models.Trigger
		)

		BeforeEach(func() {
			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			job = models.Job{Name: "test_job", Status: "New"}
			trigger = models.Trigger{Name: "test_trigger"}
			dbContext.Dbmap.Insert(&job, &trigger)
		})

		It("adds a trigger to an existing job", func() {
			job_id := job.Id
			trigger_id := trigger.Id

			vars := make(map[string]string)
			vars["job"] = strconv.FormatInt(job_id, 10)

			path := fmt.Sprintf("/jobs/%d/triggers", job_id)

			bodyString := fmt.Sprintf(`{"trigger":"%s"}`, strconv.FormatInt(trigger_id, 10))
			appContext, dbContext, w := NewAppContext("POST", path, bodyString, vars)
			defer dbContext.Dbmap.Db.Close()

			handlers.AddTriggerToJob(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from job_triggers where job_id = ? and trigger_id = ?", job_id, trigger_id)

			Expect(err).NotTo(HaveOccurred())
			Expect(w.Code).To(Equal(201))
			Expect(count).To(Equal(int64(1)))
		})
	})

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

			vars := make(map[string]string)
			vars["job"] = strconv.FormatInt(job_id, 10)
			vars["task"] = strconv.FormatInt(task_id, 10)

			path := fmt.Sprintf("/jobs/%d/tasks/%d", job_id, task_id)

			appContext, dbContext, w := NewAppContext("DELETE", path, "", vars)
			defer dbContext.Dbmap.Db.Close()
			handlers.RemoveTaskFromJob(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from job_tasks where job_id = ? and task_id = ?", job_id, task_id)

			Expect(err).NotTo(HaveOccurred())
			Expect(w.Code).To(Equal(200))
			Expect(count).To(Equal(int64(0)))
		})
	})

	Describe("RemoveTriggerFromJob", func() {
		var (
			job     models.Job
			trigger models.Trigger
		)

		BeforeEach(func() {
			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			job = models.Job{Name: "test_job", Status: "New"}
			trigger = models.Trigger{Name: "test_trigger"}
			dbContext.Dbmap.Insert(&job, &trigger)

			job_trigger := &models.JobTrigger{JobId: job.Id, TriggerId: trigger.Id}
			dbContext.Dbmap.Insert(job_trigger)
		})

		It("removes a trigger from a job", func() {
			job_id := job.Id
			trigger_id := trigger.Id

			vars := make(map[string]string)
			vars["job"] = strconv.FormatInt(job_id, 10)
			vars["trigger"] = strconv.FormatInt(trigger_id, 10)

			path := fmt.Sprintf("/jobs/%d/trigger/%d", job_id, trigger_id)

			appContext, dbContext, w := NewAppContext("DELETE", path, "", vars)
			defer dbContext.Dbmap.Db.Close()
			handlers.RemoveTriggerFromJob(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from job_triggers where job_id = ? and trigger_id = ?", job_id, trigger_id)

			Expect(err).NotTo(HaveOccurred())
			Expect(w.Code).To(Equal(200))
			Expect(count).To(Equal(int64(0)))
		})

		It("Removes the trigger from the list of executors", func() {

		})

	})
})
