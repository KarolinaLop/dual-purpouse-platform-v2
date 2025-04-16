package web

import (
	"encoding/xml"
	"testing"

	"github.com/KarolinaLop/dp/models"
	"github.com/google/go-cmp/cmp"
)

func TestParseNmapXML(t *testing.T) {

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

	// Calls the function
	got, err := ParseNmapXML("../testdata/scan-result-copy.xml")
	if err != nil {
		t.Fatalf("Could not parse XML: %s", err)
	}

	// Compares result
	if !cmp.Equal(want, got) {
		t.Errorf("parseNampXML() mismatch (-want +got):\n%s", cmp.Diff(want, got))
	}

	// if diff := cmp.Diff(want, got); diff != "" {
	// 	t.Errorf("parseNampXML() mismatch (-want +got):\n%s", diff)
	// }

	printScanResultsAsJson(t, "../testdata/scan-result-copy.xml") // for debugging
}
