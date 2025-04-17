package web

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/KarolinaLop/dp/models"
	"github.com/google/go-cmp/cmp"
)

func TestParseNmapXML(t *testing.T) {
	// Read the test XML (like your real handler)
	xmlBytes, err := os.ReadFile("../testdata/scan-result-copy.xml")
	if err != nil {
		t.Fatalf("Failed to read test XML file: %v", err)
	}

	var got models.ScanResult
	if err := xml.Unmarshal(xmlBytes, &got); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Sets up the correct result of the operation
	want := models.ScanResult{
		XMLName: xml.Name{Local: "nmaprun"},
		Hosts: []models.Host{
			{
				Addresses: []models.Address{
					{Addr: "192.168.1.110", AddrType: "ipv4"},
					{Addr: "DC:A6:32:E1:3D:67", AddrType: "mac", Vendor: "Raspberry Pi Trading"},
				},
				Ports: models.Ports{
					Ports: []models.Port{
						{
							Protocol: "tcp",
							PortID:   111,
							State:    models.State{State: "open"},
							Service: &models.Service{
								Name:      "rpcbind",
								Version:   "2-4",
								ExtraInfo: "RPC #100000",
								Method:    "probed",
								Conf:      "10",
							},
						},
					},
				},
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Parsed XML does not match expected structure (-want +got):\n%s", diff)
	}
}
