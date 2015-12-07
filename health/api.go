package health

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fzzy/radix/redis"
)

type Health struct {
	redisUrl string
}

type MetaData struct {
	HttpStatus   int    `json:"http_status"`
	ErrorMessage string `json:"error_message"`
}

type ResponseMessage struct {
	MetaData        MetaData `json:"metadata"`
	RedisConnection bool     `json:"redis_connection"`
}

func NewApi(redisUrl string) Health {
	return Health{redisUrl}
}

func (health Health) Endpoint() (string, func(http.ResponseWriter, *http.Request)) {
	return "/health", health.status
}

func (health Health) status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	health.writeSuccessResponse(w)
}

func (health Health) writeSuccessResponse(w http.ResponseWriter) {
	httpStatus := http.StatusOK
	redisStatus := health.redisConnection()

	if !redisStatus {
		httpStatus = http.StatusBadGateway
	}

	responseMessage := &ResponseMessage{
		MetaData: MetaData{
			HttpStatus: httpStatus,
		},
		RedisConnection: redisStatus,
	}
	recordsJSON, _ := json.Marshal(responseMessage)
	w.Write(recordsJSON)
}

func (health Health) redisConnection() bool {
	client, err := redis.DialTimeout("tcp", health.redisUrl, 1*time.Second)

	if err != nil {
		fmt.Println("ERROR ", err.Error())
		return false
	}

	defer client.Close()
	status, _ := client.Cmd("echo", "ping").Str()
	return status == "ping"
}
