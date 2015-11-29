package repository_test

import (
	"encoding/json"
	"time"

	"github.com/cghsystems/godata/repository"
	"github.com/cghsystems/gosum/record"
	"github.com/fzzy/radix/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	const (
		redisUrl = "127.0.0.1:6379"
	)

	var (
		testTime = time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)

		testRecord = record.Record{
			TransactionDate:        testTime,
			TransactionType:        record.CREDIT,
			SortCode:               "12-34-56",
			AccountNumber:          "123456789",
			TransactionDescription: "A Test Record",
			DebitAmount:            12.12,
			CreditAmount:           0.0,
			Balance:                12.12,
		}
	)

	var (
		recordRepository *repository.RecordRepository
		err              error
	)

	BeforeSuite(func() {
		startRedis()
	})

	cleanRedis := func() {
		c, _ := redis.DialTimeout("tcp", redisUrl, time.Duration(10)*time.Second)
		defer c.Close()
		c.Cmd("select", 0)
		c.Cmd("FLUSHDB")
	}

	BeforeEach(func() {
		recordRepository = repository.NewRecordRepository(redisUrl)
		cleanRedis()
	})

	AfterEach(func() {
		recordRepository.Close()
	})

	Context("close", func() {
		It("closes the connection", func() {
			recordRepository.Close()
			records := record.Records{testRecord}
			err := recordRepository.BulkInsert(records)
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("BulkInsert", func() {
		var testRecords = record.Records{testRecord}

		var actualRecords = func() record.Records {
			redisClient, err := redis.DialTimeout("tcp", redisUrl, time.Duration(10)*time.Second)
			bytes, err := redisClient.Cmd("smembers", "chris:gold:records").ListBytes()
			if err != nil {
				Expect(err).NotTo(HaveOccurred())
			}

			records := record.Records{}
			for x := range bytes {
				var record record.Record
				err = json.Unmarshal(bytes[x], &record)
				if err != nil {
					Expect(err).NotTo(HaveOccurred())
				}

				records = append(records, record)
			}
			return records
		}

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

				expectedRecords := record.Records{testRecord}
				Expect(actualRecords()).To(Equal(expectedRecords))
			})
		})
	})
})
