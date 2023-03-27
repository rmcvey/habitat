package main

import "fmt"

func main() {
	// Scan for Bluetooth devices and WiFi networks
	btDevices, _ := scanBluetoothDevices()
	wifiNetworks, _ := scanWifiNetworks()

	// Update trust data based on the scanned devices and networks
	db := initDatabase()
	defer db.Close()

	habitat := Habitat{}
	habitat.determineHabitatTrust(db, btDevices, wifiNetworks)

	fmt.Printf("%f\n", habitat.TrustScore)
}
