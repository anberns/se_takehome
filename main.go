package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type ReportingPlan struct {
	PlanName string `json:"plan_name"`
}

type NetworkFile struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}

type ReportingStructure struct {
	ReportingPlans []ReportingPlan `json:"reporting_plans"`
	InNetworkFiles []NetworkFile   `json:"in_network_files"`
}

type TableOfContents struct {
	ReportingStructure []ReportingStructure `json:"reporting_structure"`
}

func main() {
	start := time.Now()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]

	fmt.Println("Opening input file...")
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	fmt.Println("Processing input file...")
	gz, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return
	}
	defer gz.Close()

	decoder := json.NewDecoder(gz)

	// Read the JSON tokens until we find the "reporting_structure" array
	var found bool
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			fmt.Println("Error reading JSON token:", err)
			return
		}
		// Check for the start of the "reporting_structure" array
		if key, ok := token.(string); ok && key == "reporting_structure" {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("reporting_structure not found in the JSON")
		return
	}

	// Read the opening bracket of the "reporting_structure" array
	if _, err := decoder.Token(); err != nil {
		fmt.Println("Error reading opening bracket of reporting_structure array:", err)
		return
	}

	// save in map to avoid duplicate urls
	urlMap := map[string]struct{}{}

	// Process each element in the "reporting_structure" array
	fmt.Println("Processing reporting_structure elements...")

	for decoder.More() {
		var rs ReportingStructure
		if err := decoder.Decode(&rs); err != nil {
			fmt.Println("Error decoding reporting structure:", err)
			return
		}

		for _, rp := range rs.ReportingPlans {
			// New York PPO plans have either "PPO NY" or "NY PPO" in the plan name
			if strings.Contains(rp.PlanName, "PPO NY") || strings.Contains(rp.PlanName, "NY PPO") {
				for _, inf := range rs.InNetworkFiles {
					urlMap[inf.Location] = struct{}{}
				}
			}
		}
	}

	// Read the closing bracket of the "reporting_structure" array
	if _, err := decoder.Token(); err != nil {
		fmt.Println("Error reading closing bracket of reporting_structure array:", err)
		return
	}

	fmt.Println("Writing URLs...")

	// Create output file
	outFile, err := os.Create("urls.json")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// Write each unique url to output file
	urls := make([]string, 0, len(urlMap))
	for k := range urlMap {
		urls = append(urls, k)
	}

	urlJSON, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling URLs to JSON:", err)
		return
	}

	if _, err := outFile.Write(urlJSON); err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("URLs successfully written")

	finish := time.Now()
	fmt.Printf("Time taken: %v", finish.Sub(start))
}
