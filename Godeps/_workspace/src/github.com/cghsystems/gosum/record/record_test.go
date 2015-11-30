package record_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cghsystems/gosum/record"
)

var _ = Describe("Record", func() {
	var TestRecord = Record{
		TransactionType:        CREDIT,
		SortCode:               "12-34-56",
		AccountNumber:          "123456789",
		TransactionDescription: "A Test Record",
		DebitAmount:            12.12,
		CreditAmount:           0.0,
		Balance:                12.12,
	}

	Context("Valid Record", func() {
		Context(".ID", func() {
			It("is the expected ID", func() {
				expectedID := "31323334353637383931322e3132303030312d30312d30312030303a30303a3030202b3030303020555443d41d8cd98f00b204e9800998ecf8427e"
				Ω(expectedID).Should(Equal(TestRecord.ID()))
			})
		})
	})

	Context("NilRecord", func() {
		It("evaluates NilRecord to true", func() {
			record := NilRecord()
			Ω(record.IsNilRecord()).Should(BeTrue())
		})

		It("evaluates empty record to true", func() {
			record := Record{}
			Ω(record.IsNilRecord()).Should(BeTrue())
		})

		It("evaluates a record with data to false", func() {
			record := Record{TransactionDescription: "Not a NilRecord"}
			Ω(record.IsNilRecord()).Should(BeFalse())
		})
	})
})
