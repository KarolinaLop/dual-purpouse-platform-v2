package models

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
	"time"
)

// ScanResult model to hold values of my XML file.
type ScanResult struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

// Host holds the addresses
type Host struct {
	Addresses []Address `xml:"address"`
	Ports     Ports     `xml:"ports"`
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
	Vendor   string `xml:"vendor,attr,omitempty"`
}

// CleanVendorName cleans up encoded characters from vendor names.
func CleanVendorName(vendor string) string {
	vendor = html.UnescapeString(vendor)             // Unescapes HTML entities
	vendor = strings.ReplaceAll(vendor, "ï¼", ", ") // Replaces weird symbol with comma
	vendor = strings.TrimSpace(vendor)
	return vendor
}

type Ports struct {
	Ports []Port `xml:"port"`
}

func (p Ports) OpenPortsWithServices() string {
	var result []string
	for _, port := range p.Ports {
		if port.State.State == "open" {
			// Default to "unknown" if no service info available
			service := "unknown"
			if port.Service != nil && port.Service.Name != "" {
				service = port.Service.Name
			}
			result = append(result, fmt.Sprintf("%d (%s)", port.PortID, service))
		}
	}
	return strings.Join(result, ", ")
}

type Port struct {
	Protocol string   `xml:"protocol,attr"`
	PortID   int      `xml:"portid,attr"`
	State    State    `xml:"state"`
	Service  *Service `xml:"service"` // optional
}

type State struct {
	State string `xml:"state,attr"`
}

type Service struct {
	Name      string `xml:"name,attr"`
	Product   string `xml:"product,attr,omitempty"`
	Version   string `xml:"version,attr,omitempty"`
	CPE       string `xml:"cpe"`
	Method    string `xml:"method,attr,omitempty"`
	Conf      string `xml:"conf,attr,omitempty"`
	ExtraInfo string `xml:"extrainfo,attr,omitempty"`
}

type Scan struct {
	ID        int
	Timestamp string
}

type PageData struct {
	UserName string
	Scans    []Scan
}

func formatTimestamp(timestamp string) string {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp // Return the original if parsing fails
	}
	return parsedTime.Format("January 2, 2006, 3:04 PM")
}

// Example usage
func main() {
	scans := []Scan{
		{ID: 1, Timestamp: "2025-04-18T14:46:32Z"},
		{ID: 2, Timestamp: "2025-04-19T10:30:00Z"},
	}

	for i, scan := range scans {
		scans[i].Timestamp = formatTimestamp(scan.Timestamp)
	}
}
