package handlera

import (
	"fmt"
	"net/http"
	"triner/internal/dbase"
	"triner/internal/models"
)

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
