package models_test

import (
	. "github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobs", func() {

	Describe("ID function", func() {
		It("returns the name of the Job", func() {
			job := Job{Name: "name"}
			Expect(job.ID()).To(Equal("name"))
		})
	})

	Describe("AppendTask", func() {
		It("appends a new task", func() {
			job := Job{Name: "name", Tasks: make([]string, 0), Status: "status", Triggers: make([]string, 0)}
			job.AppendTask("task")
			Expect(job.Tasks).To(ContainElement("task"))
		})
	})

})
