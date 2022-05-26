package config

import (
	"os"
	"strconv"
)

type Config struct {
	BatchSize int
	BatchSecondInterval int
	PostEndpoint string
}

func GetConfig() (Config) {
	batchSize, err := strconv.Atoi(os.Getenv("BATCH_SIZE")) 
	if err != nil {
		panic("BATCH_SIZE env variable must be an integer")
	}

	batchSecondInterval, err := strconv.Atoi(os.Getenv("BATCH_SECOND_INTERVAL"))
	if err != nil {
		panic("BATCH_SECOND_INTERVAL env variable must be an integer")
	}

	postEndPoint := os.Getenv("POST_ENDPOINT_URL")
	if postEndPoint == "" {
		panic("POST_ENDPOINT_URL env variable not set")
	}

	return Config{
		BatchSize: batchSize,
		BatchSecondInterval: batchSecondInterval,
		PostEndpoint: postEndPoint,
	}
}