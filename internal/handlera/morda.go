package handlera

import (
	"database/sql"
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

func TrinoPinger(rwr http.ResponseWriter, req *http.Request) {

	dsn := "http://trino@trino:8080?catalog=postgresql&schema=public"
	db, err := sql.Open("trino", dsn)
	if err != nil {
		rwr.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rwr, `{"Error":"%v"}`, err)
		return

	}
	defer db.Close()

	// Быстрая проверка
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	rwr.WriteHeader(http.StatusOK)
	fmt.Fprintf(rwr, `{"trinoPing":"StatusOK"}`)

	// return db.PingContext(ctx)
}
