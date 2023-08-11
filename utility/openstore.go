package utility

import "strings"

type Network string

type Deployment struct {
	ContractAddress string
	Network         Network
}

// OpenStoreDeployments defines the list of deployments for OpenStore
var OpenStoreDeployments = []Deployment{
	{
		ContractAddress: "0x495f947276749ce646f68ac8c248420045cb7b5e",
		Network:         "ethereum",
	},
	{
		ContractAddress: "0x2953399124f0cbb46d2cbacd8a89cf0599974963",
		Network:         "matic",
	},
	// Add other deployments as needed
}

// IsOpenStoreDeployment checks if the given contract address and network is part of the OpenStore deployments
func IsOpenStoreDeployment(contractAddress string, network Network) bool {
	for _, deployment := range OpenStoreDeployments {
		if strings.EqualFold(string(deployment.Network), string(network)) &&
			strings.EqualFold(deployment.ContractAddress, contractAddress) {
			return true
		}
	}
	return false
}
