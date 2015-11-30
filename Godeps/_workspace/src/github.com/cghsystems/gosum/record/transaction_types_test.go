package record_test

import (
	. "github.com/cghsystems/gosum/record"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransactionTypes", func() {
	Context(".IsDebit", func() {
		It("evaluates to true if record has a transaction type of debit", func() {
			record := Record{TransactionType: DEBIT}
			Ω(record.IsDebit()).Should(BeTrue())
		})

		It("evaluates to false if record does not have a transaction type of debit", func() {
			record := Record{TransactionType: "dave"}
			Ω(record.IsDebit()).Should(BeFalse())
		})
	})

	Context(".IsCredit", func() {
		It("evaluates to true if record has a transaction type of credit", func() {
			record := Record{TransactionType: CREDIT}
			Ω(record.IsCredit()).Should(BeTrue())
		})

		It("evaluates to false if record does not have a transaction type of credit", func() {
			record := Record{TransactionType: "not credit"}
			Ω(record.IsCredit()).Should(BeFalse())
		})
	})
})
