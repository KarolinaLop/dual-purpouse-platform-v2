package web

import (
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// ScanResult represents the structure of the Nmap scan results.
type ScanResult struct {
	Host            string `json:"host"`
	Ports           string `json:"ports"`
	Services        string `json:"services"`
	OS              string `json:"os"`
	Vulnerabilities string `json:"vulnerabilities"`
	Status          string `json:"status"`
}

var currentCmd *exec.Cmd

// StartScanHandler starts an Nmap scan and returns the results.
func StartScanHandler(c *gin.Context) {
	// Define the target (default to localhost)
	target := "127.0.0.1"

	// Run the Nmap scan
	currentCmd = exec.Command("nmap", "-Pn", "-T4", "-sS", "-sV", "-O", "--script", "vuln", target)
	output, err := currentCmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse the Nmap output (basic parsing for now)
	result := parseNmapOutput(string(output))
	result.Status = "Scan completed"

	// Return the results as JSON
	c.JSON(http.StatusOK, result)
}

// StopScanHandler stops the currently running Nmap scan.
func StopScanHandler(c *gin.Context) {
	if currentCmd != nil && currentCmd.Process != nil {
		err := currentCmd.Process.Kill()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop scan"})
			return
		}
		currentCmd = nil
		c.JSON(http.StatusOK, gin.H{"message": "Scan stopped"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No scan is currently running"})
	}
}

// parseNmapOutput parses the raw Nmap output into a structured format.
func parseNmapOutput(output string) ScanResult {
	var result ScanResult

	// Host
	hostRegex := regexp.MustCompile(`Nmap scan report for (.+)`)
	if matches := hostRegex.FindStringSubmatch(output); len(matches) > 1 {
		result.Host = matches[1]
	}

	// Ports (all lines with open ports)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if matched, _ := regexp.MatchString(`^\d+/tcp\s+open`, line); matched {
			result.Ports += line + "\n"
		}
	}

	// Services (already partially covered)
	serviceInfo := regexp.MustCompile(`Service Info:\s*(.+)`)
	if matches := serviceInfo.FindStringSubmatch(output); len(matches) > 1 {
		result.Services = matches[1]
	}

	// OS detection
	osRegex := regexp.MustCompile(`OS details:\s*(.+)`)
	if matches := osRegex.FindStringSubmatch(output); len(matches) > 1 {
		result.OS = matches[1]
	}

	// Vulnerabilities (lines with "VULNERABLE")
	for _, line := range lines {
		if strings.Contains(line, "VULNERABLE") {
			result.Vulnerabilities += line + "\n"
		}
	}

	return result
}
