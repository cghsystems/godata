package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cghsystems/godata/health"
	"github.com/cghsystems/godata/repository"
	"github.com/cghsystems/gosum/record"
)

const redisUrl = "local.lattice.cf:6379"

func recordsPostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing /records post request")
	var records record.Records
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&records)

	if err != nil {
		panic(err)
	}

	repository := repository.NewRecordRepository(redisUrl)
	repository.BulkInsert(records)
}

func main() {
	health := health.NewApi(redisUrl)

	fmt.Println("Starting godata server")
	http.HandleFunc("/records", recordsPostHandler)
	http.HandleFunc("/health", health.Status)
	http.ListenAndServe(":8080", nil)
}
