package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
)

// StartScanHandler starts an Nmap scan and returns the results.
func StartScanHandler(c *gin.Context) {

	// Final Nmap command: sudo nmap -Pn -T4 -sS -sV --open -oX scan-result.xml -F 192.168.1.0/24
	// TODO: make the address configurable or discoverable by the code
	target := "192.168.1.0/24" // Temp my network address

	currentCmd := exec.Command(
		"nmap",            // Run the Nmap scan
		"-Pn",             // Host discovery, disables ping, treats all hosts as online
		"-T4",             // Sets timig for faster scans
		"-sS",             // Stealth SYN scan for open ports
		"-sV",             // Probes ports for srvices and versions
		"--open",          // Lists only open ports
		"-oX",             // Produces XML output
		"scan-result.xml", // Saves output as XML scan-result.xml file
		"-F",              // Fast scan, scans only the most popular ports
		target,
	)

	output, err := currentCmd.CombinedOutput()
	if err != nil {
		log.Printf(">>>>>>>>>>>> exit code: %v\n output: %s", err, string(output))
		c.Error(err)
		return
	}
}

// TODO: parse the output (chekc Go library) - working on my parser - WIP
// TODO: store the parsed output as a scanresult in the database
// data.CreateScan(data.DB)
// parseNmapXML parses the xml file and stores the information in a ScanResult.
func ParseNmapXML(filename string) (models.ScanResult, error) {

	// Open xml file
	xmlFile, err := os.Open(filename)
	// If os.Open returns error then handle it
	if err != nil {
		return models.ScanResult{}, err
	}
	defer xmlFile.Close()

	fmt.Println("Sucessfully Opened scan-result.xml")

	// Read the open XML
	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return models.ScanResult{}, fmt.Errorf("parseNmapXML: could not read xml file: %w", err)
	}

	// Init my ScanResult model
	var result models.ScanResult
	if err := xml.Unmarshal(byteValue, &result); err != nil {
		return models.ScanResult{}, err
	}

	return result, nil
}

func printScanResultsAsJson(_ *testing.T, filename string) {
	data, err := ParseNmapXML(filename)
	if err != nil {
		fmt.Printf("printScanResultAsJson: %s", err)
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ") // do error check
	if err != nil {
		fmt.Printf("Couldn't print JSON output: %s", err)
	}
	fmt.Println(string(jsonData))
}

// TODO: Nice have
// list scans -> list of previous scans that user did
// show scan details -> parse scan output from database and show results to the user
