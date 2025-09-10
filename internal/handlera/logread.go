package handlera

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ReadPostgresLog(rwr http.ResponseWriter, req *http.Request) {
	content, err := os.ReadFile("/plogs/p.log")
	if err != nil {
		rwr.WriteHeader(http.StatusNotFound)
		json.NewEncoder(rwr).Encode(err)
		return
	}
	strs := strings.Split(string(content), "\n")
	out := []string{}
	for _, s := range strs {
		iBegin := strings.Index(s, ">>>>")
		iEnd := strings.Index(s, "KUERY")
		if iBegin == -1 || iEnd == -1 {
			continue
		}
		kusman := s[iBegin : iEnd+5]
		out = append(out, kusman)
		fmt.Fprintf(rwr, "%s-----\n", kusman)
	}

}
