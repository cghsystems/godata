package record_test

import (
	"encoding/json"
	"time"

	"github.com/cghsystems/godata/config"
	domain "github.com/cghsystems/gosum/record"
	"github.com/fzzy/radix/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Record Suite")
}

func cleanRedis() {
	redisUrl, err := config.RedisUrl()
	Expect(err).NotTo(HaveOccurred())
	c, _ := redis.DialTimeout("tcp", redisUrl, 3*time.Second)
	defer c.Close()
	c.Cmd("select", 0)
	c.Cmd("FLUSHDB")
}

func actualRecords() domain.Records {
	redisUrl, _ := config.RedisUrl()
	redisClient, err := redis.DialTimeout("tcp", redisUrl, 3*time.Second)
	bytes, err := redisClient.Cmd("SMEMBERS", "chris:gold:records").ListBytes()
	Expect(err).NotTo(HaveOccurred())

	records := domain.Records{}
	for x := range bytes {
		var record domain.Record
		err = json.Unmarshal(bytes[x], &record)
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
		}

		records = append(records, record)
	}
	return records
}
