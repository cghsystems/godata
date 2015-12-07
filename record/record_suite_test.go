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
	c, _ := redis.DialTimeout("tcp", config.RedisUrl(), 3*time.Second)
	defer c.Close()
	c.Cmd("select", 0)
	c.Cmd("FLUSHDB")
}

func actualRecords() domain.Records {
	redisClient, err := redis.DialTimeout("tcp", config.RedisUrl(), 3*time.Second)
	bytes, err := redisClient.Cmd("smembers", "chris:gold:records").ListBytes()
	if err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

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
