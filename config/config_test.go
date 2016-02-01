package config_test

import (
	"os"

	"github.com/cghsystems/godata/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	const vcapServices = `{
  "p-redis": [
   {
    "credentials": {
     "host": "192.168.9.147",
     "password": "password",
     "port": 6379
    },
    "label": "p-redis",
    "name": "godata",
    "plan": "dedicated-vm",
    "tags": [
     "pivotal",
     "redis"
    ]
   }
  ]
}`

	var (
		url string
		err error
	)

	Context("Valid VCAP_SERVICES", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", vcapServices)
			url, err = config.RedisUrl()
		})

		It("returns the godata url from vcap services", func() {
			Expect(url).To(Equal("192.168.9.147:6379"))
		})
		It("does not return  error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Invalid VCAP_SERVICES", func() {
		const vcapServices = `{ "some-redis": [ ] }`
		Context("No p-redis provided", func() {
			BeforeEach(func() {
				os.Setenv("VCAP_SERVICES", vcapServices)
				url, err = config.RedisUrl()
			})

			It("returns error", func() {
				Expect(err).To(MatchError("Cannot find service with name 'p-redis'"))
			})

			It("returns no url", func() {
				Expect(url).To(BeEmpty())
			})
		})
	})

	Context("No VCAP_SERVICES", func() {
		BeforeEach(func() {
			os.Unsetenv("VCAP_SERVICES")
			url, err = config.RedisUrl()
		})

		It("returns error", func() {
			Expect(err).To(MatchError("Cannot find VCAP_SERVICES environment variable"))
		})

		It("returns no url", func() {
			Expect(url).To(BeEmpty())
		})
	})
})
