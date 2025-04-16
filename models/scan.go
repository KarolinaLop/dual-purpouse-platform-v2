package models

import "encoding/xml"

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

type Ports struct {
	Ports []Port `xml:"port"`
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
