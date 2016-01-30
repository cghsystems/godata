package record_test

import (
	"fmt"
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
		redisUrl         string
		testTime         = time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)

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
		redisUrl, _ = config.RedisUrl()
		recordRepository = record.NewRepository(redisUrl)
		cleanRedis()
	})

	AfterEach(func() {
		recordRepository.Close()
	})

	FDescribe("Get a months records", func() {
		FIt("get some records", func() {
			startDate, _ := time.Parse(time.RFC3339, "2009-01-01T00:00:00+00:00")
			endDate, _ := time.Parse(time.RFC3339, "2009-01-31T00:00:00+00:00")

			results := domain.Records{}
			for _, record := range actualRecords() {
				transactionDate := record.TransactionDate
				if (transactionDate.Equal(startDate) || transactionDate.After(startDate)) &&
					(transactionDate.Equal(endDate) || transactionDate.Before(endDate)) {
					results = append(results, record)
				}
			}
			fmt.Println(len(results))
		})
	})

	Describe("close", func() {
		It("closes the connection", func() {
			recordRepository.Close()
			records := domain.Records{testRecord}
			err := recordRepository.BulkInsert(records)
			Ω(err).Should(HaveOccurred())
		})
	})

	Describe("BulkInsert", func() {
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
