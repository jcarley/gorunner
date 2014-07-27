package handlers_test

import (
	"encoding/json"
	"log"

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
