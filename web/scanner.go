package web

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"os/exec"

	"github.com/KarolinaLop/dp/data" // Import the data package
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
)

// StartScanHandler runs an Nmap scan, reads the XML output and saves to DB.
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
		log.Printf("Scan failed: %v\n Output: %s", err, string(output))
		c.Error(err)
		return
	}
	// file is created

	// Read the XML output from file
	xmlBytes, err := os.ReadFile("../testdata/scan-result.xml")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read scan result"})
		return
	}

	// Save to DB, uses a placeholder the "1" for user_id temporarily
	err = data.StoreNmapScan(data.DB, 1, target, string(xmlBytes))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to store scan result"})
		return
	}
	c.JSON(200, gin.H{"message": "Scan completed and saved"})
}

// Function that retives scan by ID, pasres the XML, and returns it as structured JSON data
func GetParsedScanResult(c *gin.Context) {
	scanID := c.Param("id") // GET from route like /scans/:id

	// Laod form DB
	rawXML, err := data.GetNampXMLScanByID(data.DB, scanID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not load scan result"})
		return
	}

	// Parse into ScanResult struct
	var result models.ScanResult
	if err := xml.Unmarshal([]byte(rawXML), &result); err != nil {
		c.JSON(500, gin.H{"error": "Could not load scan result"})
		return
	}

	// Convert to JSON
	jsonResult, _ := json.MarshalIndent(result, "", " ")
	c.Data(200, "application/json", jsonResult)
}

// TODO: Nice have
// list scans -> list of previous scans that user did
// show scan details -> parse scan output from database and show results to the user
