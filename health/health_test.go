package health_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/cghsystems/godata/health"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {
	var (
		redisUrl        string
		redisConnection bool
		httpStatusCode  int
	)

	JustBeforeEach(func() {
		healthApi := health.NewApi(redisUrl)
		response := httptest.NewRecorder()

		request, err := http.NewRequest("GET", "/health", nil)
		Ω(err).ShouldNot(HaveOccurred())
		healthApi.Status(response, request)

		body, err := ioutil.ReadAll(response.Body)
		Ω(err).ShouldNot(HaveOccurred())
		var responseMessage health.ResponseMessage
		json.Unmarshal(body, &responseMessage)

		httpStatusCode = responseMessage.MetaData.HttpStatus
		redisConnection = responseMessage.RedisConnection
	})

	Context("Redis is running", func() {
		BeforeEach(func() {
			redisUrl = "127.0.0.1:6379"
		})

		It("has an an http status code of 200", func() {
			Expect(httpStatusCode).To(Equal(200))
		})

		It("reports that Redis is running", func() {
			Expect(redisConnection).To(BeTrue())
		})
	})

	Context("Redis is not running", func() {
		BeforeEach(func() {
			redisUrl = "nonesense.0.0.1:6379"
		})

		It("has an an http status code of Bad Gateway (502)", func() {
			Expect(httpStatusCode).To(Equal(502))
		})

		It("reports that Redis is not running", func() {
			Expect(redisConnection).To(BeFalse())
		})
	})
})
