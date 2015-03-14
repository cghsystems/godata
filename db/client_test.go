package db_test

import (
	"github.com/cghsystems/godata/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	const redisURL = "127.0.0.1:6379"

	var client *db.Client

	BeforeEach(func() {
		client = db.New(redisURL)
	})

	AfterEach(func() {
		client.Close()
	})

	Context("close", func() {
		It("closes the connection", func() {
			client.Close()
			err := client.Set("key", "value")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context(".Set", func() {
		var err error

		BeforeEach(func() {
			err = client.Set("key", "value")
		})

		It("does not return error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("sets the expected data", func() {
			data := client.Get("key")
			Ω("value").Should(Equal(data))
		})
	})
})
