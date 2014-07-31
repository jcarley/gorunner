package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jcarley/gorunner/handlers"
	"github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskHandlers", func() {

	BeforeEach(func() {
		dbContext := models.NewDbContext()
		defer dbContext.Dbmap.Db.Close()
		err := dbContext.TruncateTables()

		Expect(err).NotTo(HaveOccurred())
	})

	Describe("AddTask", func() {
		It("Adds a task to the database", func() {

			path := "/jobs"
			body := `{"name":"test"}`

			appContext, dbContext, w := NewAppContext("POST", path, body, nil)
			defer dbContext.Dbmap.Db.Close()
			handlers.AddTask(appContext)

			var task models.Task

			dbContext.Dbmap.SelectOne(&task, "select * from tasks where name = ?", "test")

			Expect(w.Code).To(Equal(201))
			Expect(task).NotTo(BeNil())
		})
	})

	Describe("GetTask", func() {

		var task models.Task

		BeforeEach(func() {
			task = models.Task{Name: "test task", Script: "script"}
			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()
			database := models.NewDatabase(dbContext)
			database.AddTask(&task)
		})

		It("Returns a task with the given Id", func() {
			path := fmt.Sprintf("/tasks/%d", task.Id)
			vars := make(map[string]string)
			vars["task"] = strconv.FormatInt(task.Id, 10)

			appContext, dbContext, w := NewAppContext("GET", path, "", vars)
			defer dbContext.Dbmap.Db.Close()

			handlers.GetTask(appContext)

			count, err := dbContext.Dbmap.SelectInt("select count(*) from tasks where id = ?", task.Id)

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(1)))
			Expect(w.Code).To(Equal(200))
		})
	})

	Describe("UpdateTask", func() {

		var task *models.Task

		BeforeEach(func() {
			task = &models.Task{Name: "test task", Script: "script"}

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			database := models.NewDatabase(dbContext)

			database.AddTask(task)
		})

		It("updates the script attribute of a task", func() {
			path := fmt.Sprintf("/tasks/%d", task.Id)
			vars := make(map[string]string)
			vars["task"] = strconv.FormatInt(task.Id, 10)
			body := `{"script": "new bash script"}`

			appContext, dbContext, w := NewAppContext("PUT", path, body, vars)
			defer dbContext.Dbmap.Db.Close()

			handlers.UpdateTask(appContext)

			task, err := appContext.Database.GetTask(task.Id)

			Expect(err).ToNot(HaveOccurred())
			Expect(task.Script).To(Equal("new bash script"))
			Expect(w.Code).To(Equal(200))
		})
	})

	Describe("ListTasks", func() {

		BeforeEach(func() {
			task := &models.Task{Name: "test task", Script: "script"}

			dbContext := models.NewDbContext()
			defer dbContext.Dbmap.Db.Close()

			database := models.NewDatabase(dbContext)

			database.AddTask(task)
		})

		It("Returns a list of all available tasks", func() {
			path := "/tasks"
			appContext, dbContext, w := NewAppContext("GET", path, "", nil)
			defer dbContext.Dbmap.Db.Close()

			handlers.ListTasks(appContext)

			var payload []models.Task
			err := json.Unmarshal(w.Body.Bytes(), &payload)
			if err != nil {
				log.Fatal(err)
			}

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).NotTo(BeNil())
			Expect(payload[0].Name).To(Equal("test task"))
			Expect(payload[0].Script).To(Equal("script"))
		})

	})

})
