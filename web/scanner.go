package web

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"slices"

	"github.com/KarolinaLop/dp/data"
	"github.com/KarolinaLop/dp/models"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

// ShowScansListHandler renders the scans list page.
func ShowScansListHandler(c *gin.Context) {
	// call some data package func that loads all scans for the current user from the db
	user, ok := c.Value("user").(models.User) // type assertion and interfaces
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
		return
	}

	scans, err := data.GetAllNmapScans(data.DB, user.ID)
	if err != nil {
		err = fmt.Errorf("failed to retrieve scans: %w", err)
		c.Error(err)
		return
	}

	IP, err := getTheHostIPAddress()
	if err != nil {
		err = fmt.Errorf("failed to retrieve local IP address: %w", err)
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "scans_list.html", gin.H{
		"userName": c.Value("user").(models.User).Name,
		"scans":    scans,
		"localIP":  IP,
	})

}

// StartScanHandler runs an Nmap scan, reads the XML output and saves to DB.
func StartScanHandler(c *gin.Context) {

	log.Println("Scan has started, it may take a while")

	// Generate a unique scan-result-*** file name
	filename := fmt.Sprintf("scan-result-%s.xml", uuid.New().String())

	log.Println("A new scan file has been crated:  ", filename)

	defer os.Remove(filename)

	target, err := getNetworkAddress()

	if err != nil {
		log.Printf("Scan failed: could not retrieve network address: %s", err)
		c.Error(err)
		return
	}

	// Run the command
	// nmap -Pn -T4 -sS -sV -sP -oX new-file.xml -F 192.168.1.0/24
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

func getNetworkAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down or not running
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("failed to get addresses for interface %s: %w", iface.Name, err)
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil { // Check for IPv4
				network := ipNet.IP.Mask(ipNet.Mask)
				return fmt.Sprintf("%s/%d", network.String(), maskToCIDR(ipNet.Mask)), nil
			}
		}
	}

	return "", fmt.Errorf("no active network interface found")
}

func getTheHostIPAddress() (string, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("failed to get local host IP address: %w", err)
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down or not running
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("failed to get the address on an interafce %s: %w", iface.Name, err)
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil { // Check for IPv4
				return ipNet.IP.To4().String(), nil
			}
		}
	}
	return "", errors.New("could not detect IP address")
}

func maskToCIDR(mask net.IPMask) int {
	ones, _ := mask.Size()
	return ones
}

func DeleteScanHandler(c *gin.Context) {
	scanID := c.Param("id")
	if err := data.DeleteScan(data.DB, scanID); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

// ShowScanDetailsHandler displays scan by ID from the database, pasres the XML, and returns it as structured HTML data
func ShowScanDetailsHandler(c *gin.Context) {
	scanID := c.Param("id")

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

	u := c.Value("user")
	user, ok := u.(models.User)
	if !ok {
		err := errors.New("failed to find the user in this context")
		c.Error(err)
		return
	}

	IP, err := getTheHostIPAddress()
	if err != nil {
		err = fmt.Errorf("failed to parse the xml: %w", err)
		c.Error(err)
		return
	}

	// extract IPv4 address from host
	rows := []hostRow{}

	// iterate over the result's hosts and append a new hostRow to rows
	for i, host := range result.Hosts {
		newRow := hostRow{
			Index:     i + 1,
			OpenPorts: host.Ports.OpenPortsWithServices(),
		}

		for _, hostAddress := range host.Addresses {
			switch hostAddress.AddrType {
			case "ipv4":
				newRow.IPv4 = hostAddress.Addr
				if newRow.IPv4 == IP {
					newRow.IsLocalHost = true
				}
			case "mac":
				newRow.MAC = hostAddress.Addr
				//newRow.Vendor = hostAddress.Vendor
				newRow.Vendor = models.CleanVendorName(hostAddress.Vendor)
				if newRow.Vendor == "" {
					newRow.Vendor = "N/A"
				}
			}

		}

		// append it to the rows slice
		rows = append(rows, newRow)
	}

	c.HTML(http.StatusOK, "scan_details.html", gin.H{
		"result":             result,
		"rows":               rows,
		"userName":           user.Name,
		"identifiedServices": identifiedServices(result),
		"numOpenPorts":       result.OpenPorts(),
		"numClosedPorts":     result.ClosedPorts(),
		"numFilteredPorts":   result.FileredPorts(),
		"hostsUp":            len(rows),
	})
}

type ServiceDetails struct {
	Name       string
	Count      int
	Percentage float64
}

func identifiedServices(result models.ScanResult) []ServiceDetails {
	m := make(map[string]ServiceDetails)

	// Loop through all hosts and their ports
	for _, host := range result.Hosts {
		for _, port := range host.Ports.Ports {
			if port.Service == nil {
				continue
			}
			details := m[port.Service.Name]
			details.Count += 1
			m[port.Service.Name] = details
		}
	}

	// get the max count
	maxCount := 0
	for _, details := range m {
		maxCount = max(maxCount, details.Count)
	}

	for key, details := range m {
		details.Percentage = calcPercentage(maxCount, details.Count)

		m[key] = details
	}

	sd := []ServiceDetails{}
	for k, v := range m {
		v.Name = k
		sd = append(sd, v)
	}

	slices.SortFunc(sd, func(a, b ServiceDetails) int {
		switch {
		case a.Count < b.Count:
			return 1
		case a.Count > b.Count:
			return -1
		default:
			return 0
		}
	})
	return sd

}

func calcPercentage(max int, count int) float64 {
	return (float64(count) / float64(max)) * 100.0
}

type hostRow struct {
	Index       int
	IPv4        string
	MAC         string
	Vendor      string
	OpenPorts   string
	IsLocalHost bool
}
