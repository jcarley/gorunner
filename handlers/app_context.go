package handlers

import (
	"net/http"

	"github.com/jcarley/gorunner/models"
)

type AppContext struct {
	DbContext *models.DbContext
	Response  http.ResponseWriter
	Request   *http.Request
}

func AppContextHandlerFunc(handlerFunc func(appContext *AppContext)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbContext := models.NewDbContext()
		defer dbContext.Dbmap.Db.Close()

		appContext := AppContext{
			DbContext: dbContext,
			Response:  w,
			Request:   r,
		}
		handlerFunc(&appContext)
	})

}
