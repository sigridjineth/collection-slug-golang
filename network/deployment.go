package network

import (
	"golang.org/x/time/rate"
	"time"
)

type Network string

const (
	ETHEREUM  Network = "ethereum"
	POLYGON   Network = "matic"
	ARBITRUM  Network = "arbitrum"
	AVALANCHE Network = "avalanche"
	OPTIMISM  Network = "optimism"
)

type Deployment struct {
	ContractAddress string
	Network         Network
}

const (
	BASE_COLLECTION_URL  = "https://opensea.io/collection/"
	DEFAULT_REDUNDANCY   = 2
	INITIAL_BACKOFF_TIME = 5000
	BACKOFF_MULTIPLIER   = 1.25
)

var (
	OPENSTORE_DEPLOYMENTS = []Deployment{
		{
			ContractAddress: "0x495f947276749ce646f68ac8c248420045cb7b5e",
			Network:         ETHEREUM,
		},
		{
			ContractAddress: "0x2953399124f0cbb46d2cbacd8a89cf0599974963",
			Network:         POLYGON,
		},
	}
	// Creating a rate limiter with a specific number of tokens per second
	BOTTLENECK_WAYBACK_MACHINE = rate.NewLimiter(rate.Every(time.Second/time.Duration(4)), 4)
)
