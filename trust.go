package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	"github.com/willf/bloom"

	"encoding/json"
	"time"
)

const (
	bloomFilterSize        = 1000
	bloomFilterHashFunctions = 5
	maxExpectedDevices = 20
)

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

func updateTrustData(db *sql.DB, btDevices []string, wifiNetworks []string) {
	ecosystem := createEcosystem(btDevices, wifiNetworks)
	habitatKey := ecosystem.B

	var habitatID int64
	err := db.QueryRow("SELECT id FROM habitats WHERE name = ?", habitatKey).Scan(&habitatID)

	if err == sql.ErrNoRows {
		// Habitat not found, insert a new habitat
		res, err := db.Exec("INSERT INTO habitats (name, ecosystem_m, ecosystem_k, ecosystem_b, trust_score) VALUES (?, ?, ?, ?, ?)",
			habitatKey, ecosystem.M, ecosystem.K, ecosystem.B, 1.0)
		if err != nil {
			log.Fatal(err)
		}

		habitatID, err = res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	// Add visit time to the visit_times table
	_, err = db.Exec("INSERT INTO visit_times (habitat_id, visit_time) VALUES (?, ?)", habitatID, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the trust score using a time decay approach
	rows, err := db.Query("SELECT visit_time FROM visit_times WHERE habitat_id = ?", habitatID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	now := time.Now()
	decayFactor := 0.5
	visitTrust := 0.0
	visitCount := 0
	for rows.Next() {
		var visitTime time.Time
		err = rows.Scan(&visitTime)
		if err != nil {
			log.Fatal(err)
		}
		timeDiff := now.Sub(visitTime)
		visitTrust += math.Pow(decayFactor, timeDiff.Hours()/24.0)
		visitCount++
	}

	// Calculate the maximum possible trust score for the current number of visits
	maxVisitTrust := 0.0
	for i := 0; i < visitCount; i++ {
		maxVisitTrust += math.Pow(decayFactor, float64(i))
	}
	maxVisitTrust *= 10.0 / (1.0 - decayFactor)

	// Normalize the trust score to the 1.0-10.0 range
	weightDevices := 0.7
	deviceTrust := float64(len(btDevices) + len(wifiNetworks)) / float64(maxExpectedDevices)
	trustScore := ((weightDevices * deviceTrust) + ((1 - weightDevices) * visitTrust)) / maxVisitTrust * 9 + 1

	// Update the trust score in the habitats table
	_, err = db.Exec("UPDATE habitats SET trust_score = ? WHERE id = ?", trustScore, habitatID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(trustScore)
}

