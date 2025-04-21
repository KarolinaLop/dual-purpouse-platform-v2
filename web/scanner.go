package web

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/KarolinaLop/dp/data" // Import the data package
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"
)

// StartScanHandler runs an Nmap scan, reads the XML output and saves to DB.
func StartScanHandler(c *gin.Context) {

	log.Println("Scan has started, it may take a while")
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

	// Read the XML output from file
	xmlBytes, err := os.ReadFile("./scan-result.xml") // os.ReadFiles returns a slice of byte -> func ReadFile(name string) ([]byte, error)
	if err != nil {
		err = fmt.Errorf("failed to read scan results from file: %w", err)
		c.Error(err)
		return
	}

	user, ok := c.Value("user").(models.User)
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
	}

	// Save to DB, uses a placeholder the "1" for user_id temporarily
	err = data.StoreNmapScan(data.DB, user.ID, string(xmlBytes)) // converts the xmlBytes slice (from the previously read file) into a string so it can be stored in the database
	if err != nil {
		err = fmt.Errorf("failed to save the scan results to database: %w", err)
		c.Error(err)
		return
	}
}

// Function that retives scan by ID from the database, pasres the XML, and returns it as structured HTML data
func ShowScanDetails(c *gin.Context) {
	scanID := c.Param("id") // GET from route like /scans/:id

	rawXML, err := data.GetNampXMLScanByID(data.DB, scanID)
	if err != nil {
		err = fmt.Errorf("failed to retrive the scan from database: %w", err)
		c.Error(err)
		return
	}

	// Parse into ScanResult struct
	var result models.ScanResult                                   // var declaration of type -> my big struct ScanResults
	if err := xml.Unmarshal([]byte(rawXML), &result); err != nil { // xml.Unmarshall parses xml data into my struct; &result is a pointer to the result
		err = fmt.Errorf("failed to parse the xml: %w", err)
		c.Error(err)
		return
	}
}

// parseNmapXML parses the xml file and stores the information in a ScanResult.
// func ParseNmapXML(filename string) (models.ScanResult, error) {

// 	// Open xml file
// 	xmlFile, err := os.Open(filename)
// 	// If os.Open returns error then handle it
// 	if err != nil {
// 		return models.ScanResult{}, err
// 	}
// 	defer xmlFile.Close()

// 	fmt.Println("Sucessfully Opened scan-result.xml")

// 	// Read the open XML
// 	byteValue, err := io.ReadAll(xmlFile)
// 	if err != nil {
// 		return models.ScanResult{}, fmt.Errorf("parseNmapXML: could not read xml file: %w", err)
// 	}

// 	// Init my ScanResult model
// 	var result models.ScanResult
// 	if err := xml.Unmarshal(byteValue, &result); err != nil {
// 		return models.ScanResult{}, err
// 	}

// 	return result, nil
// }

// TODO: Nice have
// list scans -> list of previous scans that user did
// show scan details -> parse scan output from database and show results to the user
