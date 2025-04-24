package web

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// StartScanHandler runs an Nmap scan, reads the XML output and saves to DB.
func StartScanHandler(c *gin.Context) {

	log.Println("Scan has started, it may take a while")

	// Generate a unique scan-result file name
	filename := fmt.Sprintf("scan-result-%s.xml", uuid.New().String())

	log.Println("A new scan file has been crated:  ", filename)

	defer os.Remove(filename)

	// TODO: make the address configurable or discoverable by the code
	target := "192.168.1.0/24" // Temp my network address

	currentCmd := exec.Command(
		"nmap",   // Run the Nmap scan
		"-Pn",    // Host discovery, disables ping, treats all hosts as online
		"-T4",    // Sets timig for faster scans
		"-sS",    // Stealth SYN scan for open ports
		"-sV",    // Probes ports for srvices and versions
		"--open", // Lists only open ports
		"-oX",    // Produces XML output
		filename, // Saves output as XML file
		"-F",     // Fast scan, scans only the most popular ports
		target,
	)

	// Run the command
	output, err := currentCmd.CombinedOutput()
	if err != nil {
		log.Printf("Scan failed: %v\n Output: %s", err, string(output))
		c.Error(err)
		return
	}

	// Read the XML output from file
	xmlBytes, err := os.ReadFile(filename) // os.ReadFiles returns a slice of byte -> func ReadFile(name string) ([]byte, error)
	if err != nil {
		err = fmt.Errorf("failed to read scan results from file: %w", err)
		c.Error(err)
		return
	}

	// Get user from context
	user, ok := c.Value("user").(models.User)
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
	}

	// Save to DB
	err = data.StoreNmapScan(data.DB, user.ID, string(xmlBytes)) // converts the xmlBytes slice (from the previously read file) into a string so it can be stored in the database
	if err != nil {
		err = fmt.Errorf("failed to save the scan results to database: %w", err)
		c.Error(err)
		return
	}
}

// ShowScanDetails displays scan by ID from the database, pasres the XML, and returns it as structured HTML data
func ShowScanDetails(c *gin.Context) {
	scanID := c.Param("id") // GET from route like /scans/:id/show

	rawXML, err := data.GetNampXMLScanByID(data.DB, scanID)
	if err != nil {
		err = fmt.Errorf("failed to retrive the scan from database: %w", err)
		c.Error(err)
		return
	}

	fmt.Println(rawXML)

	// Parse into ScanResult struct
	var result models.ScanResult                                   // var declaration of type -> my big struct ScanResults
	if err := xml.Unmarshal([]byte(rawXML), &result); err != nil { // xml.Unmarshall parses xml data into my struct; &result is a pointer to the result
		err = fmt.Errorf("failed to parse the xml: %w", err)
		c.Error(err)
		return
	}

	u := c.Value("user")
	user, ok := u.(models.User)
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
		return
	}

	// extract IPv4 address from host
	rows := []hostRow{}

	// iterate over the result's hosts and append a new hostRow to rows
	for _, host := range result.Hosts {
		newRow := hostRow{
			OpenPorts: host.Ports.OpenPortsWithServices(),
		}
		for _, hostAddress := range host.Addresses {
			switch hostAddress.AddrType {
			case "ipv4":
				newRow.IPv4 = hostAddress.Addr
			case "mac":
				newRow.MAC = hostAddress.Addr
				//newRow.Vendor = hostAddress.Vendor
				newRow.Vendor = models.CleanVendorName(hostAddress.Vendor)
				if newRow.Vendor == "" {
					newRow.Vendor = "Unknown"
				}
			}
		}

		// append it to the rows slice
		rows = append(rows, newRow)
	}

	c.HTML(http.StatusOK, "scan_details.html", gin.H{
		"result":   result,
		"rows":     rows,
		"userName": user.Name,
	})
}

type hostRow struct {
	IPv4      string
	MAC       string
	Vendor    string
	OpenPorts string
}
