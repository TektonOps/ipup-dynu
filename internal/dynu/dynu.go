package dynu

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

// Configuration settings for the Dynu API
const (
	updateURL = "https://api.dynu.com/nic/update"
)

// UpdateDNSRecordRequest contains the necessary information for updating a DNS record
type UpdateDNSRecordRequest struct {
	Username    string
	Password    string
	Hostname    string
	Location    string
	IPv4Address string
	EnableGroup bool
}

// ResponseMessages maps Dynu API response codes to human-readable messages
var ResponseMessages = map[string]string{
	"unknown":     "Invalid request made to the API server. Check the parameters for validity.",
	"good":        "Action processed successfully.",
	"badauth":     "Failed authentication. Check all parameters including authentication credentials.",
	"servererror": "Server encountered an error. Consider retrying the request.",
	"nochg":       "IP address unchanged. No update was necessary.",
	"notfqdn":     "The hostname is not a valid fully qualified hostname.",
	"numhost":     "Too many hostnames specified. A maximum of 20 is allowed.",
	"abuse":       "Update process failed due to abusive behavior.",
	"nohost":      "Hostname/username not found in the system.",
	"911":         "Update temporarily halted due to scheduled maintenance. Suspend updates for 10 minutes.",
	"dnserr":      "Server-side DNS error. Retry the update process.",
}

// NewDNSRequest initializes a new DNS update request
func NewDNSRequest(username, password, hostname, location, ipv4Address string, enableGroup bool, logger *slog.Logger) (*UpdateDNSRecordRequest, error) {

	if username == "" || password == "" {
		return nil, fmt.Errorf("username or password not set")
	}

	logger.Info("DNS request initialized", "hostname", hostname, "group", location, "enableGroup", enableGroup)

	return &UpdateDNSRecordRequest{
		Username:    username,
		Password:    password,
		Hostname:    hostname,
		Location:    location,
		IPv4Address: ipv4Address,
		EnableGroup: enableGroup,
	}, nil
}

// ConstructURL builds the update URL with the necessary query parameters
func (req *UpdateDNSRecordRequest) ConstructURL(logger *slog.Logger) (string, error) {
	u, err := url.Parse(updateURL)
	if err != nil {
		logger.Error("Error parsing URL", "error", err)
		return "", fmt.Errorf("error parsing URL: %v", err)
	}

	query := u.Query()
	if req.EnableGroup {
		if req.Location == "" {
			return "", fmt.Errorf("group setting enabled but no group provided")
		}
		query.Set("location", req.Location)
	} else {
		query.Set("hostname", req.Hostname)
	}
	query.Set("myip", req.IPv4Address)
	u.RawQuery = query.Encode()

	logger.Debug("Constructed URL", "url", u.String())

	return u.String(), nil
}

// SendRequest sends the DNS update request to the Dynu API
func (req *UpdateDNSRecordRequest) SendRequest(logger *slog.Logger) (*http.Response, error) {
	updateURLWithParams, err := req.ConstructURL(logger)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("GET", updateURLWithParams, nil)
	if err != nil {
		logger.Error("Error creating HTTP request", "error", err)
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	httpReq.SetBasicAuth(req.Username, req.Password)
	client := &http.Client{}
	logger.Debug("Sending DNS update request", "url", updateURLWithParams)
	return client.Do(httpReq)
}

// UpdateDNSRecord updates the DNS record with the provided details
func UpdateDNSRecord(req *UpdateDNSRecordRequest, logger *slog.Logger) (string, error) {
	resp, err := req.SendRequest(logger)
	if err != nil {
		return "", err
	}

	return HandleResponse(resp, logger)
}

// HandleResponse processes the HTTP response from the Dynu API and logs everything returned
func HandleResponse(resp *http.Response, logger *slog.Logger) (string, error) {
	defer resp.Body.Close()

	// Log the status code
	logger.Debug("HTTP Response Status", "status_code", resp.StatusCode)

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body", "error", err)
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	responseText := string(body)
	logger.Debug("HTTP Response Body", "body", responseText)

	// Handle each response code
	codes := strings.Split(responseText, "\r\n")
	for _, code := range codes {
		if message, exists := ResponseMessages[code]; exists {
			logger.Debug("Response Code", "code", code, "message", message)

			if code == "badauth" {
				return "", fmt.Errorf("authentication failed: %s", message)
			}
		}
	}

	return responseText, nil
}
