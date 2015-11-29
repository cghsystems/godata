package repository

import (
	"encoding/json"
	"time"

	"github.com/cghsystems/gosum/record"
	"github.com/fzzy/radix/redis"
)

type RecordRepository struct {
	redis *redis.Client
}

func NewRecordRepository(url string) *RecordRepository {
	c, _ := redis.DialTimeout("tcp", url, time.Duration(10)*time.Second)
	c.Cmd("select", 0)
	return &RecordRepository{redis: c}
}

func (c *RecordRepository) BulkInsert(records record.Records) error {
	for _, record := range records {
		if err := c.set(record); err != nil {
			return err
		}
	}

	return nil
}

func (c *RecordRepository) set(record record.Record) error {
	json, err := json.Marshal(record)
	if err != nil {
		return err
	}

	exists, err := c.redis.Cmd("SISMEMBER", "chris:gold:records", json).Bool()
	if err != nil {
		return err
	}

	if !exists {
		r := c.redis.Cmd("SADD", "chris:gold:records", json)
		return r.Err
	}

	return nil
}

func (c *RecordRepository) Close() {
	c.redis.Close()
}
