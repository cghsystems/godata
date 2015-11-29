package repository_test

import (
	"time"

	"github.com/cghsystems/godata/repository"
	"github.com/cghsystems/gosum/record"
	"github.com/fzzy/radix/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	const redisURL = "127.0.0.1:6379"

	var testRecord = record.Record{
		TransactionType:        record.CREDIT,
		SortCode:               "12-34-56",
		AccountNumber:          "123456789",
		TransactionDescription: "A Test Record",
		DebitAmount:            12.12,
		CreditAmount:           0.0,
		Balance:                12.12,
	}

	var (
		recordRepository *repository.RecordRepository
		err              error
	)

	BeforeSuite(func() {
		startRedis()
	})

	cleanRedis := func() {
		c, _ := redis.DialTimeout("tcp", redisURL, time.Duration(10)*time.Second)
		defer c.Close()
		c.Cmd("select", 0)
		c.Cmd("FLUSHDB")
	}

	BeforeEach(func() {
		recordRepository = repository.NewRecordRepository(redisURL)
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
		BeforeEach(func() {
			records := record.Records{testRecord}
			err = recordRepository.BulkInsert(records)
		})

		It("does not return error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
