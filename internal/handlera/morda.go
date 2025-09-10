package handlera

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"triner/internal/dbase"
	"triner/internal/models"

	"github.com/gorilla/mux"
)

type TrinoBaseStruct struct {
	DB *sql.DB
}

func DBPinger(rwr http.ResponseWriter, req *http.Request) {

	err := dbase.Ping(req.Context(), models.DSN)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rwr, `{"Error":"%v"}`, err)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	fmt.Fprintf(rwr, `{"status":"StatusOK"}`)
}

func (db *TrinoBaseStruct) TrinoPinger(rwr http.ResponseWriter, req *http.Request) {

	status := http.StatusOK

	err := db.DB.PingContext(req.Context())
	if err != nil {
		status = http.StatusInternalServerError
	}
	rwr.WriteHeader(status)
	ret := struct {
		Name   string
		Status int
		Err    error
	}{Name: "Ping", Status: status, Err: err}
	json.NewEncoder(rwr).Encode(ret)

	// Быстрая проверка
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

}

func (db *TrinoBaseStruct) AddNameHandler(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(req)
	name := vars["name"]

	status := http.StatusOK

	err := dbase.AddNameToTable(req.Context(), db.DB, name)
	if err != nil {
		status = http.StatusInternalServerError
	}
	rwr.WriteHeader(status)
	ret := struct {
		Name   string
		Status int
		Err    error
	}{Name: name, Status: status, Err: err}
	json.NewEncoder(rwr).Encode(ret)

	models.Logger.Info("Name added", name, ret)
}

func (db *TrinoBaseStruct) GetNamesHandler(rwr http.ResponseWriter, req *http.Request) {

	status := http.StatusOK

	nms, err := dbase.GetNamesFromTable(req.Context(), db.DB)
	if err != nil {
		status = http.StatusInternalServerError
		rwr.WriteHeader(status)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	rwr.WriteHeader(status)

	json.NewEncoder(rwr).Encode(nms)
	models.Logger.Info("Zapros")
}
