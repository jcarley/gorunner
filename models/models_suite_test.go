package models_test

import (
	"github.com/jcarley/gorunner/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	dbContext := models.NewDbContext()
	err := dbContext.Migrate()
	Expect(err).NotTo(HaveOccurred())
})
