package handlers_test

import (
	"github.com/jcarley/gorunner/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = BeforeSuite(func() {
	dbContext := models.NewDbContext()
	defer dbContext.Dbmap.Db.Close()
	err := dbContext.Dbmap.TruncateTables()

	Expect(err).NotTo(HaveOccurred())
})
