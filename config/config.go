package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type service struct {
	Credentials credentials `json:"credentials"`
}

type credentials struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func RedisUrl() (string, error) {
	services, err := getVcapServices()
	if err != nil {
		return "", err
	}

	redis := services["p-redis"]
	if len(redis) < 1 {
		return "", errors.New("Cannot find service with name 'p-redis'")
	}

	credentials := redis[0].Credentials
	host := credentials.Host
	port := credentials.Port
	url := fmt.Sprintf("%v:%v", host, port)

	fmt.Println("Using redis url: " + url)
	return url, nil
}

func getVcapServices() (map[string][]service, error) {
	servicesEnv := os.Getenv("VCAP_SERVICES")
	if servicesEnv == "" {
		return nil, errors.New("Cannot find VCAP_SERVICES environment variable")
	}

	services := map[string][]service{}
	json.Unmarshal([]byte(servicesEnv), &services)
	return services, nil
}
