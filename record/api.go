package record

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cghsystems/godata/config"
	domain "github.com/cghsystems/gosum/record"
)

type api struct {
}

func NewApi() *api {
	return &api{}
}

func (api *api) Endpoint() (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/records", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing /records post request")
		var records domain.Records
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&records)

		if err != nil {
			panic(err)
		}

		redisUrl, err := config.RedisUrl()
		if err != nil {
			panic(err)
		}

		repository := NewRepository(redisUrl)
		repository.BulkInsert(records)
	}
}
