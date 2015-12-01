package main_test

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/cghsystems/gosum/record"
	"github.com/fzzy/radix/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("gosum", func() {
	const (
		redisUrl = "local.lattice.cf:6379"
	)

	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		binPath, err := gexec.Build("github.com/cghsystems/godata/main/")
		Expect(err).ToNot(HaveOccurred())
		cmd := exec.Command(binPath)

		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		time.Sleep(3 * time.Second)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		session.Terminate()
		session.Wait()
		Eventually(session).Should(gexec.Exit())
	})

	Describe("GET /health", func() {
		It("returns the all clear", func() {
			response, err := http.Get("http://localhost:8080/health")
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))
		})
	})

	Describe("POST /data", func() {
		const recordsUrl = "http://localhost:8080/records"
		var testRecords record.Records

		BeforeEach(func() {
			c, _ := redis.DialTimeout("tcp", redisUrl, time.Duration(10)*time.Second)
			defer c.Close()
			c.Cmd("select", 0)
			c.Cmd("FLUSHDB")
		})

		BeforeEach(func() {
			time := time.Date(2015, 1, 1, 12, 0, 0, 0, time.UTC)

			testRecords = record.Records{
				record.Record{TransactionDate: time, Balance: 100},
			}
		})

		It("persists all of the expected records", func() {
			postRecords := func() bool {
				recordAsJson, _ := json.Marshal(testRecords)
				req, _ := http.NewRequest("POST", recordsUrl, strings.NewReader(string(recordAsJson)))

				client := &http.Client{}
				_, err := client.Do(req)
				return err == nil
			}

			actualRecords := func() record.Records {
				redisClient, err := redis.DialTimeout("tcp", redisUrl, time.Duration(10)*time.Second)
				if err != nil {
					Expect(err).NotTo(HaveOccurred())
				}

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

			Eventually(postRecords).Should(BeTrue())
			Eventually(actualRecords).Should(Equal(testRecords))
		})
	})
})
