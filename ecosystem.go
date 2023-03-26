package main

import (
	"github.com/willf/bloom"

	"encoding/json"
)

const (
	bloomFilterSize        = 1000
	bloomFilterHashFunctions = 5
	maxExpectedDevices = 20
)

type Ecosystem struct {
	M float64
	K float64
	B string
}

func createEcosystem(btDevices []string, wifiNetworks []string) Ecosystem {
	filter := bloom.New(bloomFilterSize, bloomFilterHashFunctions)

	for _, device := range btDevices {
		filter.AddString(device)
	}
	for _, network := range wifiNetworks {
		filter.AddString(network)
	}

	filterData, _ := filter.MarshalJSON()
	filterJSON := make(map[string]interface{})
	json.Unmarshal(filterData, &filterJSON)

	return Ecosystem{
		M: filterJSON["m"].(float64),
		K: filterJSON["k"].(float64),
		B: filterJSON["b"].(string),
	}
}
