package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcarley/gorunner/models"
)

type AppRouteFunc func(appContext *AppContext)

func (f AppRouteFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dbContext := models.NewDbContext()
	defer dbContext.Dbmap.Db.Close()

	database := models.NewDatabase(dbContext)

	vars := mux.Vars(r)

	appContext := AppContext{
		Database: database,
		Response: w,
		Request:  r,
		Vars:     vars,
	}
	f(&appContext)
}

func AppRoute(r *mux.Router, path string, f func(appContext *AppContext)) *mux.Route {
	return r.HandleFunc(path, AppContextHandlerFunc(f))
}

func AppContextHandlerFunc(f func(appContext *AppContext)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbContext := models.NewDbContext()
		defer dbContext.Dbmap.Db.Close()

		database := models.NewDatabase(dbContext)

		vars := mux.Vars(r)

		appContext := AppContext{
			Database: database,
			Response: w,
			Request:  r,
			Vars:     vars,
		}
		f(&appContext)
	})

}

type AppContext struct {
	Database *models.Database
	Response http.ResponseWriter
	Request  *http.Request
	Vars     map[string]string
}

func (this *AppContext) Marshal(item interface{}) {
	w := this.Response

	bytes, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func (this *AppContext) Unmarshal(k string) (payload map[string]string) {
	r := this.Request.Body
	w := this.Response

	data, err := ioutil.ReadAll(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload[k] == "" {
		http.Error(w, "Please provide a '"+k+"'", http.StatusBadRequest)
		return
	}

	return
}

func (this *AppContext) Header() http.Header {
	return this.Response.Header()
}

func (this *AppContext) Write(data []byte) (int, error) {
	return this.Response.Write(data)
}

func (this *AppContext) WriteHeader(status int) {
	this.Response.WriteHeader(status)
}

func (this *AppContext) Error(err error, code int) {
	http.Error(this, err.Error(), code)
}

func marshal(item interface{}, w http.ResponseWriter) {
	bytes, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(bytes)
}

func unmarshal(r io.Reader, k string, w http.ResponseWriter) (payload map[string]string) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload[k] == "" {
		http.Error(w, "Please provide a '"+k+"'", http.StatusBadRequest)
		return
	}

	return
}
