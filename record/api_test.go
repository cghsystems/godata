package record_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/cghsystems/godata/record"
	domain "github.com/cghsystems/gosum/record"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Records Api", func() {
	var (
		endpoint        string
		endpointHandler func(w http.ResponseWriter, r *http.Request)
	)

	BeforeEach(func() {
		endpoint, endpointHandler = record.NewApi().Endpoint()
	})

	It("returns /records endpoint url", func() {
		Expect(endpoint).To(Equal("/records"))
	})

	Describe("POST /data", func() {
		const recordsUrl = "http://localhost:8080/records"
		var testRecords domain.Records

		BeforeEach(func() {
			cleanRedis()

			time := time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
			testRecords = domain.Records{
				domain.Record{TransactionDate: time, Balance: 100},
			}
		})

		It("persists all of the expected records", func() {
			recordAsJson, _ := json.Marshal(testRecords)
			request, _ := http.NewRequest("POST", recordsUrl, strings.NewReader(string(recordAsJson)))

			response := httptest.NewRecorder()
			endpointHandler(response, request)

			Expect(actualRecords()).To(Equal(testRecords))
			Expect(response.Code).To(Equal(http.StatusOK))
		})
	})
})
