package record_test

import (
	"time"

	"github.com/cghsystems/godata/config"
	"github.com/cghsystems/godata/record"
	domain "github.com/cghsystems/gosum/record"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	var (
		recordRepository *record.Repository
		err              error
		testTime         = time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)
		redisUrl         = config.RedisUrl()

		testRecord = domain.Record{
			TransactionDate:        testTime,
			TransactionType:        domain.CREDIT,
			SortCode:               "12-34-56",
			AccountNumber:          "123456789",
			TransactionDescription: "A Test Record",
			DebitAmount:            12.12,
			CreditAmount:           0.0,
			Balance:                12.12,
		}
	)

	BeforeEach(func() {
		recordRepository = record.NewRepository(redisUrl)
		cleanRedis()
	})

	AfterEach(func() {
		recordRepository.Close()
	})

	Context("close", func() {
		It("closes the connection", func() {
			recordRepository.Close()
			records := domain.Records{testRecord}
			err := recordRepository.BulkInsert(records)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("BulkInsert", func() {
		var testRecords = domain.Records{testRecord}

		BeforeEach(func() {
			err = recordRepository.BulkInsert(testRecords)
		})

		Context("All being well", func() {
			It("does not return error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("persists the expected records", func() {
				Expect(actualRecords()).To(Equal(testRecords))
			})

			It("does not persist duplicates", func() {
				testRecords = append(testRecords, testRecord)
				recordRepository.BulkInsert(testRecords)

				expectedRecords := domain.Records{testRecord}
				Expect(actualRecords()).To(Equal(expectedRecords))
			})
		})
	})
})
