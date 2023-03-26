package main

import (
	"encoding/json"
	"encoding/xml"
	"os/exec"
)

type SystemProfilerBluetooth struct {
	SPBluetoothDataType []struct {
		DeviceConnected []map[string]struct {
			DeviceAddress string `json:"device_address"`
		} `json:"device_connected"`
	} `json:"SPBluetoothDataType"`
}

func scanBluetoothDevices() ([]string, error) {
	devices := []string{}

	cmd := exec.Command("system_profiler", "-json", "SPBluetoothDataType")
	output, err := cmd.Output()
	if err != nil {
		return devices, err
	}

	var spBluetooth SystemProfilerBluetooth
	err = json.Unmarshal(output, &spBluetooth)
	if err != nil {
		return devices, err
	}

	for _, deviceConnected := range spBluetooth.SPBluetoothDataType[0].DeviceConnected {
		for _, device := range deviceConnected {
			devices = append(devices, device.DeviceAddress)
		}
	}

	return devices, nil
}

type WifiScanResult struct {
	Networks []struct {
		IE string `xml:"IE"`
	} `xml:"array>dict"`
}

func scanWifiNetworks() ([]string, error) {
	networks := []string{}

	cmd := exec.Command("/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport", "-s", "--xml")
	output, err := cmd.Output()
	if err != nil {
		return networks, err
	}

	var wifiScanResult WifiScanResult
	err = xml.Unmarshal(output, &wifiScanResult)
	if err != nil {
		return networks, err
	}

	for _, network := range wifiScanResult.Networks {
		networks = append(networks, network.IE)
	}

	return networks, nil
}