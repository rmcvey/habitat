package main

import (
	"regexp"
	"strings"
)

func parseWindowsWifiOutput(output string) ([]string, error) {
	bssids := []string{}
	lines := strings.Split(output, "\n")

	bssidRegex := regexp.MustCompile(`^\s*BSSID \d+\s+: (.*)$`)

	for _, line := range lines {
		bssidMatch := bssidRegex.FindStringSubmatch(line)
		if len(bssidMatch) == 2 {
			bssids = append(bssids, bssidMatch[1])
		}
	}

	return bssids, nil
}

func scanWifiNetworks() ([]string, error) {
	cmd := exec.Command("netsh", "wlan", "show", "networks", "mode=Bssid")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseWindowsWifiOutput(string(output))
}
