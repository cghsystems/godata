package env_test

import (
	"os"

	"github.com/cghsystems/godata/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {
	Describe(".Get", func() {
		var (
			value string
			err   error
		)

		BeforeEach(func() {
			os.Setenv("test", "value")
			value, err = env.Get("test", "")
		})

		Context("vairbale exists", func() {
			It("returns the expected value", func() {
				Expect(value).To(Equal("value"))
			})

			It("does not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("does not exist", func() {
			Context("not default is provided", func() {
				It("returns the expected error", func() {
					value, err = env.Get("nonesense", "")
					Expect(err).To(HaveOccurred())
				})
			})

			Context("default is provided", func() {
				It("returns the default if the expected variable does not exist", func() {
					value, _ = env.Get("nonesense", "my default")
					Expect(value).To(Equal("my default"))
				})

				It("does not return an error", func() {
					_, err = env.Get("nonesense", "my default")
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
