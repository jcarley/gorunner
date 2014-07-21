package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/jcarley/gorunner/handlers"
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
})

func NewAppContext(method, path string, body io.Reader, vars map[string]string) (appContext *handlers.AppContext, dbContext *models.DbContext, w *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, path, body)
	w = httptest.NewRecorder()

	dbContext = models.NewDbContext()
	database := models.NewDatabase(dbContext)
	appContext = &handlers.AppContext{Request: req, Response: w, Database: database, Vars: vars}
	return
}
