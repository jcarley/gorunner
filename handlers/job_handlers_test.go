package handlers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/jcarley/gorunner/handlers"
	"github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JobHandlers", func() {

	Describe("ListJobs function", func() {

		BeforeEach(func() {
			job := &models.Job{Name: "test build", Status: "New"}
			models.AddJob(job)
		})

		It("returns a json array of all jobs", func() {

			req, err := http.NewRequest("GET", "/jobs", nil)
			if err != nil {
				log.Fatal(err)
			}

			w := httptest.NewRecorder()
			handlers.ListJobs(w, req)

			var payload []models.Job
			err = json.Unmarshal(w.Body.Bytes(), &payload)
			if err != nil {
				log.Fatal(err)
			}

			Expect(w.Code).To(Equal(200))
			Expect(w.Body.String()).NotTo(BeNil())
			Expect(payload[0].Name).To(Equal("test build"))
			Expect(payload[0].Status).To(Equal("New"))

			// fmt.Printf("%v+", payload)
			// fmt.Printf("%d - %s", w.Code, w.Body.String())
		})
	})

})
