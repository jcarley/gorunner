package models_test

import (
	. "github.com/jcarley/gorunner/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jobs", func() {

	BeforeEach(func() {
		dbContext := NewDbContext()
		defer dbContext.Dbmap.Db.Close()
		err := dbContext.TruncateTables()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("ID function", func() {
		It("returns the name of the Job", func() {
			job := Job{Name: "name"}
			Expect(job.ID()).To(Equal("name"))
		})
	})

})
