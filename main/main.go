package main

import (
	"fmt"
	"net/http"

	"github.com/cghsystems/godata/config"
	"github.com/cghsystems/godata/health"
	"github.com/cghsystems/godata/record"
)

type Api interface {
	EndpointHandleFunc() (url string, handleFunc func(string, string))
}

func main() {
	healthApi := health.NewApi(config.RedisUrl())
	recordsApi := record.NewApi()

	fmt.Println("Starting godata server")
	http.HandleFunc(recordsApi.Endpoint())
	http.HandleFunc(healthApi.Endpoint())
	http.ListenAndServe(":8080", nil)
}
