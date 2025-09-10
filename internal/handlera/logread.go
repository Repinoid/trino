package handlera

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"triner/internal/models"
)

func ReadPostgresLog(rwr http.ResponseWriter, req *http.Request) {
	content, err := os.ReadFile("/plogs/p.log")
	if err != nil {
		rwr.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	strs := strings.Split(string(content), "\n")
	out := []string{}
	pref := ">>>>>LOG:"
	prefDetail := ">>>>>DETAIL:"
	suff := "/*KUERY"

	for _, s := range strs {
		iBegin := strings.Index(s, pref)
		iBeginDetail := strings.Index(s, prefDetail)
		iEnd := strings.Index(s, suff)
		if (iBegin == -1 || iEnd == -1) && iBeginDetail == -1 {
			continue
		}
		kusman := ""
		if iBeginDetail == -1 {
			kusman = strings.TrimSpace(s[iBegin+len(pref) : iEnd])
		} else {
			kusman = strings.TrimSpace(s[iBeginDetail+len(prefDetail):])
		}
		out = append(out, kusman)
	}

	if len(out) == 0 {
		fmt.Fprint(rwr, "No Trino queries yet \n\n")
	}

	for _, s := range out {
		fmt.Fprintf(rwr, "%s\n\n", s)
	}
	models.Logger.Info("ReadPostgresLog", "queries", len(out))

}

func DeleteLogFile(rwr http.ResponseWriter, req *http.Request) {
	err := os.Truncate("/plogs/p.log", 0)
	if err != nil {
		rwr.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	rwr.WriteHeader(http.StatusOK)
	ret := struct {
		FName string
		Mess  string
	}{FName: "/plogs/p.log", Mess: "Log File Truncated"}
	json.NewEncoder(rwr).Encode(ret)

	models.Logger.Info("Truncate log file")

}
