package models

import (
	"encoding/xml"
	"fmt"
	"html"
	"strings"
)

// ScanResult model to hold values of my XML file.
type ScanResult struct {
	XMLName xml.Name `xml:"nmaprun"`
	Hosts   []Host   `xml:"host"`
}

// Host represents the entire <host> node, holds the addresses IP and MAC, hostnames and ports
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
	Ports      []Port       `xml:"port"`
	Extraports []Extraports `xml:"extraports"`
}

func (p Ports) OpenPortsWithServices() string {
	var result []string
	for _, port := range p.Ports {
		if port.State.State == "open" {
			// Default to "N/A" if no service info available
			service := "N/A"
			if port.Service != nil && port.Service.Name != "" {
				service = port.Service.Name
			}
			result = append(result, fmt.Sprintf("%d (%s)", port.PortID, service))
		}
	}
	return strings.Join(result, ", ")
}

func (sr ScanResult) OpenPorts() int {
	var result int
	for _, host := range sr.Hosts {
		for _, port := range host.Ports.Ports {
			if port.State.State == "open" {
				result += 1
			}
		}
	}

	fmt.Printf("Open Ports: %d\n", result)
	return result
}

func (sr ScanResult) ClosedPorts() int {
	var result int
	for _, host := range sr.Hosts {
		for _, port := range host.Ports.Extraports {
			if port.State == "closed" {
				result += port.Count
			}
		}
	}
	fmt.Printf("Closed Ports: %d\n", result)
	return result
}

func (sr ScanResult) FileredPorts() int {
	var result int
	for _, host := range sr.Hosts {
		for _, port := range host.Ports.Extraports {
			if port.State == "filtered" {
				result += port.Count
			}
		}
	}
	fmt.Printf("Filtered Ports: %d\n", result)
	return result
}

type Port struct {
	Protocol string   `xml:"protocol,attr"`
	PortID   int      `xml:"portid,attr"`
	State    State    `xml:"state"`
	Service  *Service `xml:"service"`
}

// Extraports represents the <extraports> element in the XML.
type Extraports struct {
	State   string        `xml:"state,attr"`
	Count   int           `xml:"count,attr"`
	Proto   string        `xml:"proto,attr,omitempty"`
	Ports   string        `xml:"ports,attr,omitempty"`
	Reason  string        `xml:"reason,attr,omitempty"`
	Reasons []ExtraReason `xml:"extrareasons"`
}

// ExtraReason represents the <extrareasons> element inside <extraports>.
type ExtraReason struct {
	Reason string `xml:"reason,attr"`
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
