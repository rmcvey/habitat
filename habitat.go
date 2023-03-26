package main

import (
	"database/sql"
	"log"
	"math"
	"time"
)

type Habitat struct {
	ID         int64
	Key        string
	Ecosystem  Ecosystem
	TrustScore float64
	db         *sql.DB
}

func (h *Habitat) determineHabitatTrust(db *sql.DB, btDevices []string, wifiNetworks []string) {
	h.db = db
	h.Ecosystem = createEcosystem(btDevices, wifiNetworks)
	h.Key = h.Ecosystem.B

	h.getHabitatID()

	if h.ID == -1 {
		h.insertNewHabitat()
	}

	h.insertVisitTime()

	h.calculateTrustScore(btDevices, wifiNetworks)

	h.updateTrustScore()
}

func (h *Habitat) getHabitatID() {
	err := h.db.QueryRow("SELECT id FROM habitats WHERE name = ?", h.Key).Scan(&h.ID)

	if err == sql.ErrNoRows {
		h.ID = -1
	} else if err != nil {
		log.Fatal(err)
	}
}

func (h *Habitat) insertNewHabitat() {
	res, err := h.db.Exec("INSERT INTO habitats (name, ecosystem_m, ecosystem_k, ecosystem_b, trust_score) VALUES (?, ?, ?, ?, ?)",
		h.Key, h.Ecosystem.M, h.Ecosystem.K, h.Ecosystem.B, 1.0)
	if err != nil {
		log.Fatal(err)
	}

	h.ID, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Habitat) insertVisitTime() {
	_, err := h.db.Exec("INSERT INTO visit_times (habitat_id, visit_time) VALUES (?, ?)", h.ID, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Habitat) calculateTrustScore(btDevices []string, wifiNetworks []string) {
	rows, err := h.db.Query("SELECT visit_time FROM visit_times WHERE habitat_id = ?", h.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Calculate the trust score using a time decay approach
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
	h.TrustScore = ((weightDevices * deviceTrust) + ((1 - weightDevices) * visitTrust)) / maxVisitTrust * 9 + 1
}

func (h *Habitat) updateTrustScore() {
	_, err := h.db.Exec("UPDATE habitats SET trust_score = ? WHERE id = ?", h.TrustScore, h.ID)
	if err != nil {
		log.Fatal(err)
	}
}
