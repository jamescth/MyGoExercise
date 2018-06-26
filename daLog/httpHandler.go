package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func GetBugOutput(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	cfg, err := readConf(*conf)
	if err != nil {
		log.Printf("%+v\n", err)
		w.Write([]byte(fmt.Sprintf("Error readConf %+v\n", err)))
		os.Exit(1)
	}

	bugid := "/cores/bug_" + params["id"]
	if err := SearchPath(bugid, cfg); err != nil {
		log.Printf("Error SearchPath %s %+v\n", bugid, err)
		w.Write([]byte(fmt.Sprintf("Error SearchPath %s %+v\n", bugid, err)))
		// os.Exit(1)
		return
	}
	if err := cfg.Output(w); err != nil {
		log.Printf("Error Output %+v\n", err)
		w.Write([]byte(fmt.Sprintf("Error Output %+v\n", err)))
		os.Exit(1)
	}

	log.Printf("report %s\n", bugid)
}
